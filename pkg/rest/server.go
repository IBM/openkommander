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
	if (err != nil) {
		return nil, fmt.Errorf("failed to create Kafka client: %v", err)
	}

	s := &Server{
		kafkaClient: client,
		startTime:   time.Now(),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/status", s.handleStatus)
	mux.HandleFunc("/api/v1/topics", s.handleTopics)
	mux.HandleFunc("/api/v1/brokers", s.handleBrokers)

	s.httpServer = &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	return s, nil
}

func StartRESTServer(port string, brokers []string) {
	s, err := NewServer(port, brokers)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	go s.handleShutdown()
	log.Printf("REST API server running on port %s...", s.httpServer.Addr)
	if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
}

func (s *Server) Stop(ctx context.Context) error {
	if err := s.kafkaClient.Close(); err != nil {
		return fmt.Errorf("failed to close Kafka client: %v", err)
	}
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) handleShutdown() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.Stop(ctx); err != nil {
		log.Printf("Error during server shutdown: %v", err)
	}
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
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

	sendJSON(w, http.StatusOK, Response{Status: "ok", Data: brokerList})
}

func (s *Server) listTopics(w http.ResponseWriter, r *http.Request) {
	topics, err := s.kafkaClient.Topics()
	if err != nil {
		sendError(w, "Failed to list topics", err)
		return
	}
	sendJSON(w, http.StatusOK, Response{Status: "ok", Data: topics})
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
	defer admin.Close()

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
	topicName := r.URL.Query().Get("name")
	if topicName == "" {
		http.Error(w, "Topic name is required", http.StatusBadRequest)
		return
	}

	admin, err := sarama.NewClusterAdminFromClient(s.kafkaClient)
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

	sendJSON(w, http.StatusOK, Response{Status: "ok", Message: fmt.Sprintf("Topic '%s' deleted successfully", topicName)})
}

func sendJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func sendError(w http.ResponseWriter, message string, err error) {
	sendJSON(w, http.StatusInternalServerError, Response{
		Status:  "error",
		Message: fmt.Sprintf("%s: %v", message, err),
	})
}
