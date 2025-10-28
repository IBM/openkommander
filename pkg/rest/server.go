package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"slices"
	"strings"
	"syscall"
	"time"

	"github.com/IBM/openkommander/internal/core/commands"
	"github.com/IBM/openkommander/pkg/constants"
	"github.com/IBM/openkommander/pkg/logger"
	"github.com/IBM/sarama"
)

func wrapWithLogging(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		LoggingMiddleware(http.HandlerFunc(fn)).ServeHTTP(w, r)
	}
}

type Server struct {
	httpServer  *http.Server
	kafkaClient sarama.Client
	startTime   time.Time
}

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type TopicRequest struct {
	Name              string `json:"name"`
	Partitions        int32  `json:"partitions"`
	ReplicationFactor int16  `json:"replication_factor"`
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		isAPIRequest := len(r.URL.Path) >= 4 && r.URL.Path[:4] == "/api"
		requestType := "UI"
		if isAPIRequest {
			requestType = "API"
		}
		if isAPIRequest {
			logger.HTTP("API request",
				r.Method,
				r.URL.Path,
				0, // status not available yet
				0, // duration not available yet
				"url", r.URL.String(),
				"remote_addr", r.RemoteAddr,
				"user_agent", r.UserAgent(),
				"content_length", r.ContentLength,
				"host", r.Host,
				"referer", r.Referer(),
			)
		} else {
			logger.Debug("Incoming UI request",
				"type", requestType,
				"method", r.Method,
				"path", r.URL.Path,
				"remote_addr", r.RemoteAddr,
				"user_agent", r.UserAgent(),
			)
		}

		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrapped, r)

		duration := time.Since(start)
		if isAPIRequest {
			logger.HTTP("API request completed",
				r.Method,
				r.URL.Path,
				wrapped.statusCode,
				duration.Milliseconds(),
				"remote_addr", r.RemoteAddr,
			)
		} else {
			logger.Debug("UI request completed",
				"type", requestType,
				"method", r.Method,
				"path", r.URL.Path,
				"status_code", wrapped.statusCode,
				"duration_ms", duration.Milliseconds(),
			)
		}
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func enforceMethod(w http.ResponseWriter, r *http.Request, allowedMethods []string) bool {
	for _, method := range allowedMethods {
		if r.Method == method {
			return true
		}
	}

	// Join allowed methods for header
	w.Header().Set("Allow", strings.Join(allowedMethods, ", "))
	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)

	logger.Info("Method Not Allowed",
		"path", r.URL.Path,
		"method", r.Method,
		"remote_addr", r.RemoteAddr,
	)
	return false
}

func NewServer(port string) (*Server, error) {
	s := &Server{
		kafkaClient: nil,
		startTime:   time.Now(),
	}

	router := http.NewServeMux()

	// Topics endpoint supports GET, POST, DELETE
	router.HandleFunc("/api/v1/{broker}/topics", wrapWithLogging(s.handleTopics))

	// Brokers endpoint supports GET, POST
	router.HandleFunc("/api/v1/{broker}/brokers", wrapWithLogging(s.handleBrokers))

	// Metrics/messages/minute endpoint supports GET only
	router.HandleFunc("/api/v1/{broker}/metrics/messages/minute", wrapWithLogging(s.handleMessagesPerMinute))

	// Status endpoint supports GET only
	router.HandleFunc("/api/v1/{broker}/status", wrapWithLogging(s.handleStatus))

	// Health endpoint supports GET only
	router.HandleFunc("/api/v1/{broker}/health", wrapWithLogging(s.handleHealth))

	// Clusters endpoint supports GET only
	router.HandleFunc("/api/v1/clusters", wrapWithLogging(s.handleClusters))

	// Cluster metadata endpoint supports GET only
	router.HandleFunc("/api/v1/clusters/{clusterId}/metadata", wrapWithLogging(s.handleClusterMetadata))

	frontendDir := constants.OpenKommanderFolder + "/frontend"
	fileServer := http.FileServer(http.Dir(frontendDir))
	router.Handle("/static/", http.StripPrefix("/static/", fileServer))

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		filePath := frontendDir + r.URL.Path

		// Serve static files directly if they exist
		if _, err := os.Stat(filePath); err == nil {
			if !enforceMethod(w, r, []string{http.MethodGet}) {
				return
			}
			http.ServeFile(w, r, filePath)
			return
		}

		// Block unmatched API/static routes
		if strings.HasPrefix(r.URL.Path, "/api") || strings.HasPrefix(r.URL.Path, "/static") {
			http.Error(w, "404 Not Found", http.StatusNotFound)
			return
		}

		// Serve index.html for all other frontend GET requests and let frontend handle routing
		indexPath := frontendDir + "/index.html"
		if _, err := os.Stat(indexPath); err == nil {
			if !enforceMethod(w, r, []string{http.MethodGet}) {
				return
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			http.ServeFile(w, r, indexPath)
			return
		}

		// Fallback in case index.html missing
		http.Error(w, "404 Not Found", http.StatusNotFound)
	})

	logger.Info("Serving frontend", "directory", frontendDir)

	s.httpServer = &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	return s, nil
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	if s.kafkaClient != nil {
		brokers := s.kafkaClient.Brokers()
		if len(brokers) > 0 {
			if err := s.kafkaClient.Close(); err != nil {
				logger.Warn("Failed to close Kafka client during server shutdown", "error", err)
			} else {
				logger.Info("Kafka client closed successfully")
			}
		} else {
			logger.Info("Kafka client was already disconnected")
		}
		s.kafkaClient = nil
	}
	return s.httpServer.Shutdown(ctx)
}

func StartRESTServer(port string) {
	s, err := NewServer(port)
	if err != nil {
		logger.Error("Failed to start server", "error", err)
		os.Exit(1)
	}

	// Set up shutdown handler
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		<-sigCh

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := s.Stop(ctx); err != nil {
			logger.Error("Error during server shutdown", "error", err)
		}
	}()
	logger.Info("REST API server running on port", "port", port)
	if err := s.Start(); err != http.ErrServerClosed {
		logger.Error("Server error", "error", err)
		os.Exit(1)
	}
}

func createNewClient(w http.ResponseWriter, r *http.Request, s *Server) (status bool, err error) {
	broker := r.PathValue("broker")

	if broker == "" {
		broker = "kafka-cluster1:9093"
	}

	logger.Kafka("Creating new Kafka client", broker, "connect", "client_addr", r.RemoteAddr)

	if broker == "" {
		logger.Warn("Broker not specified in request", "url", r.URL.String())
		sendError(w, "Broker not specified", nil)
		return false, fmt.Errorf("broker not specified")
	}
	config := sarama.NewConfig()
	config.Version = constants.SaramaKafkaVersion
	client, err := sarama.NewClient([]string{broker}, config)
	if err != nil {
		logger.Error("Failed to create Kafka client", "broker", broker, "error", err)
		sendError(w, "Failed to create Kafka client", err)
		return false, fmt.Errorf("failed to create Kafka client: %w", err)
	}

	s.kafkaClient = client
	logger.Kafka("Successfully created Kafka client", broker, "connect")
	return true, nil
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		sendJSON(w, http.StatusMethodNotAllowed, Response{
			Status:  "error",
			Message: fmt.Sprintf("Method %s not allowed", r.Method),
		})
		return
	}

	broker := r.PathValue("broker")

	status, err := createNewClient(w, r, s)
	if err != nil {
		logger.Error("Failed to create Kafka client for status check", "broker", broker, "error", err)
		sendError(w, "Failed to create Kafka client", nil)
		return
	}
	if !status {
		logger.Error("Client creation failed for status check", "broker", broker)
		sendError(w, "Failed to create Kafka client", fmt.Errorf("client creation failed"))
		return
	}

	if s.kafkaClient == nil {
		sendError(w, "Kafka client not initialized", fmt.Errorf("kafka client is nil"))
		return
	}
	brokers := s.kafkaClient.Brokers()
	kafkaStatus := "disconnected"
	if len(brokers) > 0 {
		kafkaStatus = "connected"
	}

	uptime := time.Since(s.startTime).Seconds()
	logger.Info("Status check completed", "broker", broker, "kafka_status", kafkaStatus, "brokers_count", len(brokers), "uptime_seconds", uptime)

	response := Response{
		Status:  "ok",
		Message: "OpenKommander REST API is running",
		Data: map[string]interface{}{
			"kafka_status":   kafkaStatus,
			"brokers_count":  len(brokers),
			"uptime_seconds": uptime,
		},
	}
	sendJSON(w, http.StatusOK, response)
}

func (s *Server) handleTopics(w http.ResponseWriter, r *http.Request) {
	broker := r.PathValue("broker")

	status, err := createNewClient(w, r, s)
	if err != nil {
		logger.Error("Failed to create Kafka client for topics operation", "broker", broker, "method", r.Method, "error", err)
		sendError(w, "Failed to create Kafka client", err)
		return
	}
	if !status {
		logger.Error("Client creation failed for topics operation", "broker", broker, "method", r.Method)
		sendError(w, "Failed to create Kafka client", fmt.Errorf("client creation failed"))
		return
	}

	if s.kafkaClient == nil {
		logger.Error("Kafka client not initialized for topics operation", "broker", broker)
		sendError(w, "Kafka client not initialized", fmt.Errorf("kafka client is nil"))
		return
	}

	switch r.Method {
	case http.MethodGet:
		s.listTopics(w, r)
	case http.MethodPost:
		s.createTopic(w, r)
	case http.MethodDelete:
		s.deleteTopic(w, r)
	default:
		logger.Warn("Method not allowed for topics endpoint", "method", r.Method, "broker", broker)
		w.WriteHeader(http.StatusMethodNotAllowed)
		sendJSON(w, http.StatusMethodNotAllowed, Response{
			Status:  "error",
			Message: fmt.Sprintf("Method %s not allowed", r.Method),
		})
	}
}

func (s *Server) handleBrokers(w http.ResponseWriter, r *http.Request) {
	broker := r.PathValue("broker")
	switch r.Method {
	case http.MethodPost:
		s.createBroker(w, r)
	case http.MethodGet:
		s.getBrokers(w, r)
	default:
		logger.Warn("Method not allowed for brokers endpoint", "method", r.Method, "broker", broker)
		w.WriteHeader(http.StatusMethodNotAllowed)
		sendJSON(w, http.StatusMethodNotAllowed, Response{
			Status:  "error",
			Message: fmt.Sprintf("Method %s not allowed", r.Method),
		})
	}
}

func (s *Server) createBroker(w http.ResponseWriter, r *http.Request) {
	_ = r // Not used in this stub implementation
	response := Response{
		Status:  "ok",
		Message: "Broker creation is not implemented yet",
	}

	sendJSON(w, http.StatusNotImplemented, response)
}

func (s *Server) getBrokers(w http.ResponseWriter, r *http.Request) {
	broker := r.PathValue("broker")

	status, err := createNewClient(w, r, s)
	if err != nil {
		logger.Error("Failed to create Kafka client for brokers operation", "broker", broker, "error", err)
		sendError(w, "Failed to create Kafka client", err)
		return
	}
	if !status {
		logger.Error("Client creation failed for brokers operation", "broker", broker)
		sendError(w, "Failed to create Kafka client", fmt.Errorf("client creation failed"))
		return
	}

	if s.kafkaClient == nil {
		logger.Error("Kafka client not initialized for brokers operation", "broker", broker)
		sendError(w, "Kafka client not initialized", fmt.Errorf("kafka client is nil"))
		return
	}

	brokers := s.kafkaClient.Brokers()
	brokerList := make([]map[string]interface{}, 0)

	for _, brokerInfo := range brokers {
		connected, err := brokerInfo.Connected()
		if err != nil {
			connected = false
		}

		tlsState, _ := brokerInfo.TLSConnectionState()

		brokerData := map[string]interface{}{
			"id":        brokerInfo.ID(),
			"addr":      brokerInfo.Addr(),
			"connected": connected,
			"rack":      brokerInfo.Rack(),
			"state":     tlsState,
		}
		brokerList = append(brokerList, brokerData)
	}

	logger.Info("Successfully retrieved brokers", "broker", broker, "broker_count", len(brokerList))
	sendJSON(w, http.StatusOK, Response{Status: "ok", Data: brokerList})
}

func (s *Server) listTopics(w http.ResponseWriter, r *http.Request) {
	broker := r.PathValue("broker")
	admin, err := sarama.NewClusterAdminFromClient(s.kafkaClient)

	if err != nil {
		logger.Error("Failed to create admin client for listing topics", "broker", broker, "error", err)
		sendError(w, "Failed to create admin client", err)
		return
	}
	defer func() {
		if closeErr := admin.Close(); closeErr != nil {
			logger.Warn("Failed to close admin client", "error", closeErr)
		}
	}()

	topics, err := admin.ListTopics()

	if err != nil {
		logger.Error("Failed to list topics from Kafka", "broker", broker, "error", err)
		sendError(w, "Failed to list topics", err)
		return
	}

	logger.Debug("Successfully retrieved topics", "broker", broker, "topic_count", len(topics))

	topicList := make([]map[string]interface{}, 0, len(topics))
	for name, details := range topics {
		replicas := int(details.NumPartitions) * int(details.ReplicationFactor)
		inSyncReplicas := replicas
		topicList = append(topicList, map[string]interface{}{
			"name":               name,
			"partitions":         details.NumPartitions,
			"replication_factor": details.ReplicationFactor,
			"replicas":           replicas,
			"in_sync_replicas":   inSyncReplicas,
		})
	}

	sendJSON(w, http.StatusOK, Response{Status: "ok", Data: topicList})
}

func (s *Server) createTopic(w http.ResponseWriter, r *http.Request) {
	broker := r.PathValue("broker")

	var req TopicRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Invalid request body for topic creation", "broker", broker, "error", err)
		sendError(w, "Invalid request body", err)
		return
	}

	logger.Info("Topic creation request details",
		"broker", broker,
		"topic_name", req.Name,
		"partitions", req.Partitions,
		"replication_factor", req.ReplicationFactor)
	admin, err := sarama.NewClusterAdminFromClient(s.kafkaClient)
	if err != nil {
		logger.Error("Failed to create admin client for topic creation", "broker", broker, "topic_name", req.Name, "error", err)
		sendError(w, "Failed to create admin client", err)
		return
	}
	defer func() {
		if closeErr := admin.Close(); closeErr != nil {
			logger.Warn("Failed to close admin client", "error", closeErr)
		}
	}()
	err = admin.CreateTopic(req.Name, &sarama.TopicDetail{
		NumPartitions:     req.Partitions,
		ReplicationFactor: req.ReplicationFactor,
	}, false)
	if err != nil {
		logger.Error("Failed to create topic in Kafka", "broker", broker, "topic_name", req.Name, "error", err)
		sendError(w, "Failed to create topic", err)
		return
	}

	logger.Info("Topic created successfully", "broker", broker, "topic_name", req.Name, "partitions", req.Partitions, "replication_factor", req.ReplicationFactor)
	sendJSON(w, http.StatusCreated, Response{Status: "ok", Message: fmt.Sprintf("Topic '%s' created successfully", req.Name)})
}

func (s *Server) deleteTopic(w http.ResponseWriter, r *http.Request) {
	broker := r.PathValue("broker")

	var req TopicRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Invalid request body for topic deletion", "broker", broker, "error", err)
		sendError(w, "Invalid request body", err)
		return
	}
	topicName := req.Name
	if topicName == "" {
		logger.Warn("Topic name is required for deletion", "broker", broker)
		sendError(w, "Topic name is required", nil)
		return
	}

	logger.Info("Topic deletion request details", "broker", broker, "topic_name", topicName)
	admin, err := sarama.NewClusterAdminFromClient(s.kafkaClient)
	if err != nil {
		logger.Error("Failed to create admin client for topic deletion", "broker", broker, "topic_name", topicName, "error", err)
		sendError(w, "Failed to create admin client", err)
		return
	}
	defer func() {
		if closeErr := admin.Close(); closeErr != nil {
			logger.Warn("Failed to close admin client", "error", closeErr)
		}
	}()

	err = admin.DeleteTopic(topicName)
	if err != nil {
		logger.Error("Failed to delete topic from Kafka", "broker", broker, "topic_name", topicName, "error", err)
		sendError(w, "Failed to delete topic", err)
		return
	}

	logger.Info("Topic deleted successfully", "broker", broker, "topic_name", topicName)
	sendJSON(w, http.StatusOK, Response{Status: "ok", Message: fmt.Sprintf("Topic '%s' deleted successfully", topicName)})
}

func sendJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		logger.Error("Failed to encode JSON response", "error", err)
	}
}

func sendError(w http.ResponseWriter, message string, err error) {
	logger.Error(message, "error", err)
	sendJSON(w, http.StatusInternalServerError, Response{
		Status:  "error",
		Message: fmt.Sprintf("%s: %v", message, err),
	})
}

// MessagesLastMinute holds the count of produced and consumed messages in the last minute and last second
type MessagesLastMinute struct {
	Topic          string `json:"topic"`
	ProducedCount  int64  `json:"produced_count"`
	ConsumedCount  int64  `json:"consumed_count"`
	ProducedPerSec int64  `json:"produced_per_sec"`
	ConsumedPerSec int64  `json:"consumed_per_sec"`
}

// OffsetHistoryEntry stores offset and timestamp
type OffsetHistoryEntry struct {
	Offset    int64
	Timestamp time.Time
}

// In-memory rolling window for offsets (not production safe)
var producedHistory = make(map[string][]OffsetHistoryEntry)
var consumedHistory = make(map[string][]OffsetHistoryEntry)

// In-memory store for last offsets and timestamps (for demo purposes, not production safe)

// Handler for messages per minute
func (s *Server) handleMessagesPerMinute(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		sendJSON(w, http.StatusMethodNotAllowed, Response{
			Status:  "error",
			Message: fmt.Sprintf("Method %s not allowed", r.Method),
		})
		return
	}

	broker := r.PathValue("broker")

	status, err := createNewClient(w, r, s)
	if err != nil {
		logger.Error("Failed to create Kafka client for messages per minute", "broker", broker, "error", err)
		sendError(w, "Failed to create Kafka client", err)
		return
	}
	if !status {
		logger.Error("Client creation failed for messages per minute", "broker", broker)
		sendError(w, "Failed to create Kafka client", fmt.Errorf("client creation failed"))
		return
	}

	if s.kafkaClient == nil {
		logger.Error("Kafka client not initialized for messages per minute", "broker", broker)
		sendError(w, "Kafka client not initialized", fmt.Errorf("kafka client is nil"))
		return
	}

	admin, err := sarama.NewClusterAdminFromClient(s.kafkaClient)
	if err != nil {
		logger.Error("Failed to create admin client for messages per minute", "broker", broker, "error", err)
		sendError(w, "Failed to create admin client", nil)
		return
	}

	defer func() {
		if closeErr := admin.Close(); closeErr != nil {
			logger.Warn("Failed to close admin client", "error", closeErr)
		}
	}()

	topics, err := admin.ListTopics()

	if err != nil {
		logger.Error("Failed to list topics for messages per minute", "broker", broker, "error", err)
		sendError(w, "Failed to list topics", nil)
		return
	}

	logger.Debug("Processing message metrics", "broker", broker, "topic_count", len(topics))

	// Calculate produced and consumed message counts in the last minute for each topic
	counts := []MessagesLastMinute{}
	now := time.Now()
	totalAllProduced := int64(0)
	totalAllConsumed := int64(0)
	totalAllProducedSec := int64(0)
	totalAllConsumedSec := int64(0)

	for name := range topics { // Produced: sum latest offsets across all partitions
		partitions, err := s.kafkaClient.Partitions(name)
		if err != nil {
			logger.Warn("Failed to get partitions for topic", "broker", broker, "topic", name, "error", err)
		}
		var totalProduced int64 = 0
		for _, partition := range partitions {
			offset, err := s.kafkaClient.GetOffset(name, partition, sarama.OffsetNewest)
			if err == nil {
				totalProduced += offset
			}
		}
		// Consumed: sum committed offsets for all consumer groups
		var totalConsumed int64 = 0
		groups, err := admin.ListConsumerGroups()
		if err != nil {
			logger.Warn("Failed to list consumer groups", "broker", broker, "error", err)
		}

		for group := range groups {
			offsets, err := admin.ListConsumerGroupOffsets(group, map[string][]int32{name: partitions})
			if err == nil && offsets.Blocks != nil {
				for _, partition := range partitions {
					block := offsets.GetBlock(name, partition)
					if block != nil && block.Offset > 0 {
						totalConsumed += block.Offset
					}
				}
			}
		} // Store offset history for rolling window
		producedHistory[name] = append(producedHistory[name], OffsetHistoryEntry{Offset: totalProduced, Timestamp: now})
		consumedHistory[name] = append(consumedHistory[name], OffsetHistoryEntry{Offset: totalConsumed, Timestamp: now})

		// Remove entries older than 1 minute using slices.DeleteFunc
		pruneMinute := now.Add(-1 * time.Minute)
		producedHistory[name] = slices.DeleteFunc(producedHistory[name], func(entry OffsetHistoryEntry) bool {
			return entry.Timestamp.Before(pruneMinute)
		})
		consumedHistory[name] = slices.DeleteFunc(consumedHistory[name], func(entry OffsetHistoryEntry) bool {
			return entry.Timestamp.Before(pruneMinute)
		})

		// Create filtered copies for second-based calculations
		pruneSecond := now.Add(-1 * time.Second)
		producedHistorySec := slices.DeleteFunc(slices.Clone(producedHistory[name]), func(entry OffsetHistoryEntry) bool {
			return entry.Timestamp.Before(pruneSecond)
		})
		consumedHistorySec := slices.DeleteFunc(slices.Clone(consumedHistory[name]), func(entry OffsetHistoryEntry) bool {
			return entry.Timestamp.Before(pruneSecond)
		})

		// Calculate produced/consumed in last minute
		producedCount := int64(0)
		consumedCount := int64(0)
		if len(producedHistory[name]) > 1 {
			producedCount = producedHistory[name][len(producedHistory[name])-1].Offset - producedHistory[name][0].Offset
		}
		if len(consumedHistory[name]) > 1 {
			consumedCount = consumedHistory[name][len(consumedHistory[name])-1].Offset - consumedHistory[name][0].Offset
		}

		// Calculate produced/consumed in last second
		producedPerSec := int64(0)
		consumedPerSec := int64(0)
		if len(producedHistorySec) > 1 {
			producedPerSec = producedHistorySec[len(producedHistorySec)-1].Offset - producedHistorySec[0].Offset
		}
		if len(consumedHistorySec) > 1 {
			consumedPerSec = consumedHistorySec[len(consumedHistorySec)-1].Offset - consumedHistorySec[0].Offset
		}

		totalAllProduced += producedCount
		totalAllConsumed += consumedCount
		totalAllProducedSec += producedPerSec
		totalAllConsumedSec += consumedPerSec

		counts = append(counts, MessagesLastMinute{
			Topic:          name,
			ProducedCount:  producedCount,
			ConsumedCount:  consumedCount,
			ProducedPerSec: producedPerSec,
			ConsumedPerSec: consumedPerSec,
		})
	}

	counts = append(counts, MessagesLastMinute{
		Topic:          "total",
		ProducedCount:  totalAllProduced,
		ConsumedCount:  totalAllConsumed,
		ProducedPerSec: totalAllProducedSec,
		ConsumedPerSec: totalAllConsumedSec,
	})
	logger.Info("Successfully calculated message metrics", "broker", broker, "total_topics", len(counts)-1, "total_produced", totalAllProduced, "total_consumed", totalAllConsumed)
	sendJSON(w, http.StatusOK, Response{Status: "ok", Data: counts})
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		sendJSON(w, http.StatusMethodNotAllowed, Response{
			Status:  "error",
			Message: fmt.Sprintf("Method %s not allowed", r.Method),
		})
		return
	}

	response := Response{
		Status:  "ok",
		Message: "Health check successful",
	}
	sendJSON(w, http.StatusOK, response)
}

// Handler for clusters endpoint
func (s *Server) handleClusters(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		sendJSON(w, http.StatusMethodNotAllowed, Response{
			Status:  "error",
			Message: fmt.Sprintf("Method %s not allowed", r.Method),
		})
		return
	}

	// Use the command from internal/core/commands
	clusters, failure := commands.ListClusters()
	if failure != nil {
		logger.Error("Failed to list clusters", "error", failure.Err)
		sendError(w, "Failed to list clusters", failure.Err)
		return
	}

	logger.Info("Successfully retrieved clusters", "cluster_count", len(clusters))
	sendJSON(w, http.StatusOK, Response{Status: "ok", Data: clusters})
}

// Handler for cluster metadata endpoint
func (s *Server) handleClusterMetadata(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		sendJSON(w, http.StatusMethodNotAllowed, Response{
			Status:  "error",
			Message: fmt.Sprintf("Method %s not allowed", r.Method),
		})
		return
	}

	clusterId := r.PathValue("clusterId")

	status, err := createNewClient(w, r, s)
	if err != nil {
		logger.Error("Failed to create Kafka client for cluster metadata operation", "clusterId", clusterId, "error", err)
		sendError(w, "Failed to create Kafka client", err)
		return
	}
	if !status {
		logger.Error("Client creation failed for cluster metadata operation", "clusterId", clusterId)
		sendError(w, "Failed to create Kafka client", fmt.Errorf("client creation failed"))
		return
	}

	// Use the command from internal/core/commands
	metadata, failure := commands.GetClusterMetadata()
	if failure != nil {
		logger.Error("Failed to get cluster metadata", "clusterId", clusterId, "error", failure.Err)
		sendError(w, "Failed to get cluster metadata", failure.Err)
		return
	}

	logger.Info("Successfully retrieved cluster metadata", "clusterId", clusterId)
	sendJSON(w, http.StatusOK, Response{Status: "ok", Data: metadata})
}
