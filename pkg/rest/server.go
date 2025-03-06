package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/IBM/sarama"
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

func NewServer(port string, brokers []string) (*Server, error) {

	config := sarama.NewConfig()
	config.Version = sarama.V2_8_0_0

	client, err := sarama.NewClient(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka client: %v", err)
	}

	server := &Server{
		kafkaClient: client,
		startTime:   time.Now(),
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/status", server.handleStatus)
	mux.HandleFunc("/api/v1/topics", server.handleTopics)
	mux.HandleFunc("/api/v1/brokers", server.handleBrokers)

	server.httpServer = &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	return server, nil
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	brokers := s.kafkaClient.Brokers()
	kafkaStatus := "disconnected"
	if len(brokers) > 0 {
		kafkaStatus = "connected"
	}

	clusterInfo := make(map[string]interface{})

	clusterInfo["api_version"] = "1.0.0"
	clusterInfo["kafka_version"] = s.kafkaClient.Config().Version.String()

	brokerDetails := make([]map[string]interface{}, 0)
	for _, broker := range brokers {

		connected, err := broker.Connected()
		if err != nil {
			connected = false
		}

		brokerInfo := map[string]interface{}{
			"id":        broker.ID(),
			"addr":      broker.Addr(),
			"connected": connected,
		}
		brokerDetails = append(brokerDetails, brokerInfo)
	}
	clusterInfo["broker_details"] = brokerDetails

	response := Response{
		Status:  "ok",
		Message: "OpenKommander REST API is running",
		Data: map[string]interface{}{
			"kafka_status":   kafkaStatus,
			"brokers_count":  len(brokers),
			"cluster_info":   clusterInfo,
			"uptime_seconds": time.Since(s.startTime).Seconds(),
		},
	}

	sendJSON(w, http.StatusOK, response)
}

func (s *Server) handleTopics(w http.ResponseWriter, r *http.Request) {
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

func (s *Server) listTopics(w http.ResponseWriter, r *http.Request) {

	if refreshable, ok := s.kafkaClient.(interface{ RefreshMetadata() error }); ok {
		refreshable.RefreshMetadata()
	}

	topics, err := s.kafkaClient.Topics()
	if err != nil {
		sendError(w, "Failed to list topics", err)
		return
	}

	response := Response{
		Status: "ok",
		Data:   topics,
	}
	sendJSON(w, http.StatusOK, response)
}

func (s *Server) createTopic(w http.ResponseWriter, r *http.Request) {
	var req TopicRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Invalid request body", err)
		return
	}

	brokerList := make([]string, 0)
	for _, broker := range s.kafkaClient.Brokers() {
		brokerList = append(brokerList, broker.Addr())
	}

	admin, err := sarama.NewClusterAdmin(brokerList, s.kafkaClient.Config())
	if err != nil {
		sendError(w, "Failed to create admin client", err)
		return
	}
	defer admin.Close()

	err = admin.CreateTopic(req.Name, &sarama.TopicDetail{
		NumPartitions:     req.Partitions,
		ReplicationFactor: req.ReplicationFactor,
	}, false)
	if err != nil {
		sendError(w, "Failed to create topic", err)
		return
	}

	response := Response{
		Status:  "ok",
		Message: fmt.Sprintf("Topic '%s' created successfully", req.Name),
	}
	sendJSON(w, http.StatusCreated, response)
}

func (s *Server) deleteTopic(w http.ResponseWriter, r *http.Request) {
	topicName := r.URL.Query().Get("name")
	if topicName == "" {
		http.Error(w, "Topic name is required", http.StatusBadRequest)
		return
	}

	brokerList := make([]string, 0)
	for _, broker := range s.kafkaClient.Brokers() {
		brokerList = append(brokerList, broker.Addr())
	}

	admin, err := sarama.NewClusterAdmin(brokerList, s.kafkaClient.Config())
	if err != nil {
		sendError(w, "Failed to create admin client", err)
		return
	}
	defer admin.Close()

	err = admin.DeleteTopic(topicName)
	if err != nil {
		sendError(w, "Failed to delete topic", err)
		return
	}

	response := Response{
		Status:  "ok",
		Message: fmt.Sprintf("Topic '%s' deleted successfully", topicName),
	}
	sendJSON(w, http.StatusOK, response)
}

func (s *Server) handleBrokers(w http.ResponseWriter, r *http.Request) {
	brokers := s.kafkaClient.Brokers()
	brokerList := make([]map[string]interface{}, 0)

	for _, broker := range brokers {

		connected, err := broker.Connected()
		if err != nil {

			connected = false
		}

		brokerInfo := map[string]interface{}{
			"id":        broker.ID(),
			"addr":      broker.Addr(),
			"connected": connected,
			"rack":      broker.Rack(),
		}
		brokerList = append(brokerList, brokerInfo)
	}

	response := Response{
		Status: "ok",
		Data:   brokerList,
	}
	sendJSON(w, http.StatusOK, response)
}

func sendJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func sendError(w http.ResponseWriter, message string, err error) {
	response := Response{
		Status:  "error",
		Message: fmt.Sprintf("%s: %v", message, err),
	}
	sendJSON(w, http.StatusInternalServerError, response)
}

func (s *Server) Start() error {
	fmt.Printf("REST server starting on port %s\n", s.httpServer.Addr)
	fmt.Printf("Connected to Kafka brokers: %v\n", s.kafkaClient.Brokers())
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	if err := s.kafkaClient.Close(); err != nil {
		return fmt.Errorf("failed to close Kafka client: %v", err)
	}
	return s.httpServer.Shutdown(ctx)
}
