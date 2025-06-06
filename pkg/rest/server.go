package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IBM/openkommander/pkg/constants"
	"github.com/IBM/sarama"
	"github.com/gorilla/mux"
)

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

func NewServer(port string) (*Server, error) {
	s := &Server{
		kafkaClient: nil,
		startTime:   time.Now(),
	}

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/{broker}/status", s.handleStatus)
	router.HandleFunc("/api/v1/{broker}/topics", s.handleTopics)
	router.HandleFunc("/api/v1/{broker}/brokers", s.handleBrokers)
	router.HandleFunc("/api/v1/{broker}/metrics/messages/minute", s.handleMessagesPerMinute)
	router.HandleFunc("/api/v1/{broker}/health", func(w http.ResponseWriter, r *http.Request) {
		response := Response{
			Status:  "ok",
			Message: "Health check successful",
		}
		sendJSON(w, http.StatusOK, response)
	})

	frontendDir := constants.OpenKommanderFolder + "/frontend"
	fileServer := http.FileServer(http.Dir(frontendDir))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServer))

	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		filePath := frontendDir + r.URL.Path
		if _, err := os.Stat(filePath); err == nil {
			http.ServeFile(w, r, filePath)
			return
		}

		indexPath := frontendDir + "/index.html"
		if _, err := os.Stat(indexPath); err == nil {
			http.ServeFile(w, r, indexPath)
		} else {
			http.NotFound(w, r)
		}
	})

	log.Printf("Serving frontend from %s", frontendDir)

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
	if err := s.kafkaClient.Close(); err != nil {
		return fmt.Errorf("failed to close Kafka client: %v", err)
	}
	return s.httpServer.Shutdown(ctx)
}

func StartRESTServer(port string) {
	s, err := NewServer(port)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	// Set up shutdown handler
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		<-sigCh

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := s.Stop(ctx); err != nil {
			log.Printf("Error during server shutdown: %v", err)
		}
	}()

	log.Printf("REST API server running on port %s...", port)
	if err := s.Start(); err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
}

func createNewClient(w http.ResponseWriter, r *http.Request, s *Server) (status bool, err error) {
	vars := mux.Vars(r)
	broker := vars["broker"]

	fmt.Println("Creating new Kafka client for broker:", broker)

	if broker == "" {
		http.Error(w, "Broker not specified", http.StatusBadRequest)
		return false, fmt.Errorf("broker not specified")
	}

	config := sarama.NewConfig()
	config.Version = sarama.V2_8_0_0
	client, err := sarama.NewClient([]string{broker}, config)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create Kafka client: %v", err), http.StatusInternalServerError)
		return false, fmt.Errorf("failed to create Kafka client: %w", err)
	}

	s.kafkaClient = client
	return true, nil
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	status, err := createNewClient(w, r, s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !status {
		http.Error(w, "Failed to create Kafka client", http.StatusInternalServerError)
		return
	}

	fmt.Println("Handling status request for broker")

	if s.kafkaClient == nil {
		http.Error(w, "Kafka client not initialized", http.StatusInternalServerError)
		return
	}

	brokers := s.kafkaClient.Brokers()
	kafkaStatus := "disconnected"
	if len(brokers) > 0 {
		kafkaStatus = "connected"
	}

	response := Response{
		Status:  "ok",
		Message: "OpenKommander REST API is running",
		Data: map[string]interface{}{
			"kafka_status":   kafkaStatus,
			"brokers_count":  len(brokers),
			"uptime_seconds": time.Since(s.startTime).Seconds(),
		},
	}
	sendJSON(w, http.StatusOK, response)
}

func (s *Server) handleTopics(w http.ResponseWriter, r *http.Request) {
	status, err := createNewClient(w, r, s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !status {
		http.Error(w, "Failed to create Kafka client", http.StatusInternalServerError)
		return
	}

	fmt.Println("Handling status request for broker " + mux.Vars(r)["broker"])

	if s.kafkaClient == nil {
		http.Error(w, "Kafka client not initialized", http.StatusInternalServerError)
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
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleBrokers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		s.createBroker(w, r)
	case http.MethodGet:
		s.getBrokers(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) createBroker(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Status:  "ok",
		Message: "Broker creation is not implemented yet",
	}

	sendJSON(w, http.StatusNotImplemented, response)
}

func (s *Server) getBrokers(w http.ResponseWriter, r *http.Request) {
	brokers := s.kafkaClient.Brokers()
	brokerList := make([]map[string]interface{}, 0)

	for _, broker := range brokers {
		connected, err := broker.Connected()
		if err != nil {
			connected = false
		}

		tlsState, _ := broker.TLSConnectionState()

		brokerInfo := map[string]interface{}{
			"id":        broker.ID(),
			"addr":      broker.Addr(),
			"connected": connected,
			"rack":      broker.Rack(),
			"state":     tlsState,
		}
		brokerList = append(brokerList, brokerInfo)
	}

	sendJSON(w, http.StatusOK, Response{Status: "ok", Data: brokerList})
}

func (s *Server) listTopics(w http.ResponseWriter, r *http.Request) {
	admin, err := sarama.NewClusterAdminFromClient(s.kafkaClient)

	if err != nil {
		sendError(w, "Failed to create admin client", err)
		return
	}

	// defer admin.Close()

	topics, err := admin.ListTopics()

	if err != nil {
		sendError(w, "Failed to list topics", err)
		return
	}

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
	var req TopicRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Invalid request body", err)
		return
	}

	admin, err := sarama.NewClusterAdminFromClient(s.kafkaClient)
	if err != nil {
		sendError(w, "Failed to create admin client", err)
		return
	}
	// defer admin.Close()

	err = admin.CreateTopic(req.Name, &sarama.TopicDetail{
		NumPartitions:     req.Partitions,
		ReplicationFactor: req.ReplicationFactor,
	}, false)
	if err != nil {
		sendError(w, "Failed to create topic", err)
		return
	}

	sendJSON(w, http.StatusCreated, Response{Status: "ok", Message: fmt.Sprintf("Topic '%s' created successfully", req.Name)})
}

func (s *Server) deleteTopic(w http.ResponseWriter, r *http.Request) {
	var req TopicRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Invalid request body", err)
		return
	}

	topicName := req.Name
	if topicName == "" {
		sendError(w, "Topic name is required", nil)
		return
	}

	admin, err := sarama.NewClusterAdminFromClient(s.kafkaClient)
	if err != nil {
		sendError(w, "Failed to create admin client", err)
		return
	}
	// defer admin.Close()

	err = admin.DeleteTopic(topicName)
	if err != nil {
		sendError(w, "Failed to delete topic", err)
		return
	}

	sendJSON(w, http.StatusOK, Response{Status: "ok", Message: fmt.Sprintf("Topic '%s' deleted successfully", topicName)})
}

func sendJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("Failed to encode JSON response: %v", err)
	}
}

func sendError(w http.ResponseWriter, message string, err error) {
	log.Printf("ERROR: %s: %v", message, err)
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
	admin, err := sarama.NewClusterAdminFromClient(s.kafkaClient)

	if err != nil {
		sendError(w, "Failed to create admin client", err)
		return
	}

	// defer admin.Close()

	topics, err := admin.ListTopics()

	if err != nil {
		sendError(w, "Failed to list topics", err)
		return
	}

	// Calculate produced and consumed message counts in the last minute for each topic
	counts := []MessagesLastMinute{}
	now := time.Now()
	totalAllProduced := int64(0)
	totalAllConsumed := int64(0)
	totalAllProducedSec := int64(0)
	totalAllConsumedSec := int64(0)

	for name := range topics {
		// Produced: sum latest offsets across all partitions
		partitions, err := s.kafkaClient.Partitions(name)
		if err != nil {
			continue
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
		if err == nil {
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
			}
		}

		// Store offset history for rolling window
		producedHistory[name] = append(producedHistory[name], OffsetHistoryEntry{Offset: totalProduced, Timestamp: now})
		consumedHistory[name] = append(consumedHistory[name], OffsetHistoryEntry{Offset: totalConsumed, Timestamp: now})

		// Remove entries older than 1 minute
		pruneMinute := now.Add(-1 * time.Minute)
		for len(producedHistory[name]) > 0 && producedHistory[name][0].Timestamp.Before(pruneMinute) {
			producedHistory[name] = producedHistory[name][1:]
		}
		for len(consumedHistory[name]) > 0 && consumedHistory[name][0].Timestamp.Before(pruneMinute) {
			consumedHistory[name] = consumedHistory[name][1:]
		}

		// Remove entries older than 1 second
		pruneSecond := now.Add(-1 * time.Second)
		producedHistorySec := producedHistory[name]
		consumedHistorySec := consumedHistory[name]
		for len(producedHistorySec) > 0 && producedHistorySec[0].Timestamp.Before(pruneSecond) {
			producedHistorySec = producedHistorySec[1:]
		}
		for len(consumedHistorySec) > 0 && consumedHistorySec[0].Timestamp.Before(pruneSecond) {
			consumedHistorySec = consumedHistorySec[1:]
		}

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

	sendJSON(w, http.StatusOK, Response{Status: "ok", Data: counts})
}
