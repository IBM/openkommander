package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"openkommander/lib/kafka"
	"openkommander/lib/utils"
)

type Server struct {
	config        *Config
	router        *gin.Engine
	client        *kafka.Client
	clusterClients map[string]*kafka.Client
	tracker       *kafka.ConsumerTracker
	httpServer    *http.Server
}

func NewServer(config *Config) (*Server, error) {
	if config.LogLevel == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	var opts []kafka.ClientOption
	if config.SASLEnabled {
		mechanism := utils.GetSASLMechanism(config.SASLMechanism)
		opts = append(opts, kafka.WithSASL(config.SASLUsername, config.SASLPassword, mechanism))
	}
	
	if config.TLSEnabled {
		opts = append(opts, kafka.WithTLS())
	}
	
	client, err := kafka.NewClient(config.Brokers, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka client: %w", err)
	}

	clusterClients := make(map[string]*kafka.Client)
	for _, cluster := range config.Clusters {
		var clusterOpts []kafka.ClientOption
		
		if cluster.SASLEnabled {
			mechanism := utils.GetSASLMechanism(cluster.SASLMechanism)
			clusterOpts = append(clusterOpts, kafka.WithSASL(cluster.SASLUsername, cluster.SASLPassword, mechanism))
		}
		
		if cluster.TLSEnabled {
			clusterOpts = append(clusterOpts, kafka.WithTLS())
		}
		
		clusterClient, err := kafka.NewClient(cluster.Brokers, clusterOpts...)
		if err != nil {
			log.Printf("Warning: Failed to create client for cluster '%s': %v", cluster.Name, err)
			continue
		}
		
		clusterClients[cluster.Name] = clusterClient
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(RequestLogger())
	router.Use(CORS())

	consumerTracker := kafka.NewConsumerTracker()

	return &Server{
		config:         config,
		router:         router,
		client:         client,
		clusterClients: clusterClients,
		tracker:        consumerTracker,
	}, nil
}

func (s *Server) Start() error {
	if len(s.clusterClients) > 0 {
		log.Printf("Server is configured with %d additional cluster(s)", len(s.clusterClients))
		for name := range s.clusterClients {
			log.Printf("- Cluster '%s' is available", name)
		}
		
		SetupRoutesWithClusters(s.router, s.client, s.clusterClients, s.tracker)
	} else {
		SetupRoutes(s.router, s.client, s.tracker)
	}

	utils.PrintRoutes(s.router)

	s.httpServer = &http.Server{
		Addr:    ":" + s.config.Port,
		Handler: s.router,
	}

	go func() {
		log.Printf("Server starting on port %s", s.config.Port)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("Shutting down server...")

	s.tracker.StopAllConsumers()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	if err := s.client.Close(); err != nil {
		log.Printf("Warning: failed to close primary Kafka client: %v", err)
	}

	for name, client := range s.clusterClients {
		if err := client.Close(); err != nil {
			log.Printf("Warning: failed to close Kafka client for cluster '%s': %v", name, err)
		}
	}

	log.Println("Server gracefully stopped")
	return nil
}

func (s *Server) WaitForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}
}

func (s *Server) GetRouter() *gin.Engine {
	return s.router
}

func (s *Server) GetKafkaClient() *kafka.Client {
	return s.client
}

func (s *Server) GetClusterClient(name string) (*kafka.Client, bool) {
	client, exists := s.clusterClients[name]
	return client, exists
}

func (s *Server) GetConsumerTracker() *kafka.ConsumerTracker {
	return s.tracker
}