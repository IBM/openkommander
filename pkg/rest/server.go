package rest

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

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

	// Status endpoint for API supports GET only
	router.HandleFunc("/api/v1/status", wrapWithLogging(s.handleStatus))

	// Topics endpoint supports GET, POST, DELETE
	router.HandleFunc("/api/v1/topics", wrapWithLogging(s.handleTopics))
	router.HandleFunc("/api/v1/topics/{name}", wrapWithLogging(s.handleTopics))

	// Brokers endpoint supports GET, POST
	router.HandleFunc("/api/v1/brokers", wrapWithLogging(s.handleBrokers))

	// Metrics/messages/minute endpoint supports GET only
	router.HandleFunc("/api/v1/metrics/messages/minute", wrapWithLogging(s.handleMessagesPerMinute))

	// // Health endpoint supports GET only
	// router.HandleFunc("/api/v1/{broker}/health", wrapWithLogging(s.handleHealth))

	// Clusters endpoint supports GET only
	router.HandleFunc("/api/v1/clusters", wrapWithLogging(s.handleClusters))

	router.HandleFunc("/api/v1/cluster/{name}", wrapWithLogging(s.handleClusterByName))

	// Cluster metadata endpoint supports GET only
	router.HandleFunc("/api/v1/cluster/metadata", wrapWithLogging(s.handleClusterMetadata))

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

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	if !enforceMethod(w, r, []string{http.MethodGet}) {
		return
	}

	switch r.Method {
	case http.MethodGet:
		HandleStatus(w, r, s.startTime)
	}
}

func (s *Server) handleTopics(w http.ResponseWriter, r *http.Request) {
	if !enforceMethod(w, r, []string{http.MethodGet, http.MethodPost, http.MethodDelete}) {
		return
	}

	switch r.Method {
	case http.MethodGet:
		ListTopics(w, r)
	case http.MethodPost:
		CreateTopic(w, r)
	case http.MethodDelete:
		DeleteTopic(w, r)
	}
}

func (s *Server) handleBrokers(w http.ResponseWriter, r *http.Request) {
	if !enforceMethod(w, r, []string{http.MethodGet, http.MethodPost}) {
		return
	}

	switch r.Method {
	case http.MethodGet:
		HandleGetBrokers(w, r)
	case http.MethodPost:
		HandleCreateBroker(w, r)
	}
}

func (s *Server) handleMessagesPerMinute(w http.ResponseWriter, r *http.Request) {
	if !enforceMethod(w, r, []string{http.MethodGet}) {
		return
	}

	switch r.Method {
	case http.MethodGet:
		HandleMessagesPerMinute(w, r)
	}
}

func (s *Server) handleClusters(w http.ResponseWriter, r *http.Request) {
	if !enforceMethod(w, r, []string{http.MethodGet, http.MethodPost}) {
		return
	}

	switch r.Method {
	case http.MethodGet:
		HandleListClusters(w, r)
	case http.MethodPost:
		HandleLoginCluster(w, r)
	}
}

func (s *Server) handleClusterByName(w http.ResponseWriter, r *http.Request) {
	if !enforceMethod(w, r, []string{http.MethodGet, http.MethodPost, http.MethodDelete}) {
		return
	}

	clusterName := r.PathValue("name")

	if clusterName == "" {
		logger.Warn("Cluster name is required in path", "path", r.URL.Path)
		SendError(w, "Cluster name is required in path", nil)
		return
	}

	switch r.Method {
	case http.MethodGet:
		HandleGetClusterByName(w, r, clusterName)
	case http.MethodPost:
		HandleSelectCluster(w, r, clusterName)
	case http.MethodDelete:
		HandleDeleteClusterByName(w, r, clusterName)
	}
}

func (s *Server) handleClusterMetadata(w http.ResponseWriter, r *http.Request) {
	if !enforceMethod(w, r, []string{http.MethodGet}) {
		return
	}

	switch r.Method {
	case http.MethodGet:
		HandleClusterMetadata(w, r)
	}
}
