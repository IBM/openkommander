package session

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/IBM/openkommander/pkg/cluster"
	"github.com/IBM/openkommander/pkg/constants"
	"github.com/IBM/openkommander/pkg/logger"
	"github.com/IBM/sarama"
)

type Session interface {
	Info() string
	Connect(ctx context.Context) (sarama.Client, error)
	Disconnect()
	GetClient() (sarama.Client, error)
	GetAdminClient() (sarama.ClusterAdmin, error)
	IsAuthenticated() bool
}

type session struct {
	clusters      []ClusterConnection
	activeCluster string
	client        sarama.Client
	adminClient   sarama.ClusterAdmin
}

type ClusterConnection struct {
	Name            string   `json:"name"`
	Brokers         []string `json:"brokers"`
	Version         string   `json:"version"`
	IsAuthenticated bool     `json:"isAuthenticated"`
}

type SessionData struct {
	Clusters      []ClusterConnection `json:"clusters"`
	ActiveCluster string              `json:"activeCluster"`
}

func (s *session) Info() string {
	if s.activeCluster == "" {
		return "No active cluster selected"
	}

	for _, cluster := range s.clusters {
		if cluster.Name == s.activeCluster {
			return fmt.Sprintf("Active Cluster: %s, Brokers: %v, Authenticated: %v, Version: %v",
				cluster.Name, cluster.Brokers, cluster.IsAuthenticated, cluster.Version)
		}
	}
	return "Active cluster not found"
}

func (s *session) getActiveCluster() *ClusterConnection {
	idx := s.getActiveClusterIndex()
	if idx >= 0 {
		return &s.clusters[idx]
	}
	return nil
}

func (s *session) getActiveClusterIndex() int {
	for i := range s.clusters {
		if s.clusters[i].Name == s.activeCluster {
			return i
		}
	}
	return -1
}

func (s *session) Connect(ctx context.Context) (sarama.Client, error) {
	if s.client != nil {
		return s.client, nil
	}

	activeCluster := s.getActiveCluster()
	if activeCluster == nil {
		return nil, fmt.Errorf("no active cluster selected")
	}

	version, err := sarama.ParseKafkaVersion(activeCluster.Version)
	if err != nil {
		return nil, fmt.Errorf("invalid kafka version: %w", err)
	}
	client, err := cluster.NewCluster(activeCluster.Brokers, version).Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("error connecting to cluster: %w", err)
	}
	adminClient, err := cluster.NewCluster(activeCluster.Brokers, version).ConnectAdmin(ctx)
	if err != nil {
		return nil, fmt.Errorf("error connecting to cluster as admin: %w", err)
	}
	s.client = client
	s.adminClient = adminClient
	index := s.getActiveClusterIndex()
	if index >= 0 {
		s.clusters[index].IsAuthenticated = true
	}
	return client, nil
}

func (s *session) Disconnect() {
	if s.client != nil {
		if err := s.client.Close(); err != nil {
			logger.Error("Error closing client", "error", err)
		}
	}
	s.client = nil
	s.adminClient = nil

	// Mark active cluster as disconnected
	index := s.getActiveClusterIndex()
	if index >= 0 {
		s.clusters[index].IsAuthenticated = false
	}
	fmt.Println("Logged out successfully!")
}

func (s *session) IsAuthenticated() bool {
	activeCluster := s.getActiveCluster()
	return activeCluster != nil && activeCluster.IsAuthenticated
}

func (s *session) GetAdminClient() (sarama.ClusterAdmin, error) {
	if s.adminClient != nil {
		return s.adminClient, nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	activeCluster := s.getActiveCluster()
	if activeCluster == nil {
		return nil, fmt.Errorf("no active cluster selected")
	}

	version, err := sarama.ParseKafkaVersion(activeCluster.Version)
	if err != nil {
		return nil, fmt.Errorf("invalid kafka version: %w", err)
	}
	adminClient, err := cluster.NewCluster(activeCluster.Brokers, version).ConnectAdmin(ctx)
	if err != nil {
		return nil, err
	}
	s.adminClient = adminClient
	return adminClient, nil
}

func (s *session) GetClient() (sarama.Client, error) {
	if s.client != nil {
		return s.client, nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := s.Connect(ctx)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (s *session) GetBrokers() []string {
	activeCluster := s.getActiveCluster()
	if activeCluster == nil {
		return []string{}
	}
	return activeCluster.Brokers
}

func GetCurrentSession() *session {
	return currentSession
}

var currentSession *session

func createDefaultSession() error {
	err := os.MkdirAll(constants.OpenKommanderFolder, 0755)
	if err != nil {
		return fmt.Errorf("error creating directory %s: %w", constants.OpenKommanderFolder, err)
	}

	file, err := os.OpenFile(constants.OpenKommanderConfigFilename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("error creating session file %s: %w", constants.OpenKommanderConfigFilename, err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			logger.Error("Error closing file", "error", err)
		}
	}()

	sessionData := SessionData{Clusters: []ClusterConnection{}, ActiveCluster: ""}
	return json.NewEncoder(file).Encode(sessionData)
}

func saveSession() error {
	err := os.MkdirAll(constants.OpenKommanderFolder, 0755)
	if err != nil {
		return fmt.Errorf("error creating directory %s: %w", constants.OpenKommanderFolder, err)
	}

	file, err := os.OpenFile(constants.OpenKommanderConfigFilename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("error creating session file %s: %w", constants.OpenKommanderConfigFilename, err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			logger.Error("Error closing session file", "error", err)
		}
	}()

	sessionData := SessionData{
		Clusters:      currentSession.clusters,
		ActiveCluster: currentSession.activeCluster,
	}
	err = json.NewEncoder(file).Encode(sessionData)
	if err != nil {
		return fmt.Errorf("error encoding session data: %w", err)
	}
	return nil
}

func loadSession() error {
	file, err := os.Open(constants.OpenKommanderConfigFilename)
	if err != nil {
		err = createDefaultSession()
		if err != nil {
			fmt.Println("Error creating session file:", err)
			return err
		}

		file, err = os.Open(constants.OpenKommanderConfigFilename)
		if err != nil {
			fmt.Println("Error opening session file:", err)
			return err
		}
	}
	defer func() {
		if err := file.Close(); err != nil {
			logger.Error("Error closing session file", "error", err)
		}
	}()

	var data SessionData
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	if err != nil {
		fmt.Println("Error decoding session data:", err)
		return err
	}

	currentSession.clusters = data.Clusters
	currentSession.activeCluster = data.ActiveCluster
	return nil
}

func init() {
	currentSession = &session{
		clusters:      []ClusterConnection{},
		activeCluster: "",
		client:        nil,
		adminClient:   nil,
	}

	err := loadSession()
	if err != nil {
		logger.Error("Error loading session", "error", err)
	}
}

func Login() {
	versionReader := bufio.NewReader(os.Stdin)
	fmt.Printf("Enter kafka version [%s]: ", constants.KafkaVersion)

	version, _ := versionReader.ReadString('\n')
	version = strings.TrimSpace(version)
	if version == "" {
		version = constants.KafkaVersion
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("(Any broker address from the cluster - other brokers will be auto-discovered)")
	fmt.Printf("Enter broker address [%s]: ", constants.KafkaBroker)

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		input = constants.KafkaBroker
	}

	// Create temporary cluster connection for testing
	tempCluster := ClusterConnection{
		Brokers:         []string{input},
		Version:         version,
		IsAuthenticated: false,
	}

	fmt.Printf("Connecting to cluster via: %s\n", input)

	// Test connection
	kafkaVersion, err := sarama.ParseKafkaVersion(version)
	if err != nil {
		logger.Error("Invalid Kafka version string", "version", version, "error", err)
		fmt.Println("Failed to parse Kafka version. Please enter a valid version string (e.g., 2.1.0.0).")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := cluster.NewCluster(tempCluster.Brokers, kafkaVersion).Connect(ctx)

	if client != nil && err == nil {
		fmt.Println("Logged in successfully!")
		fmt.Printf("Kafka Version [%s]\n", version)

		discoveredBrokers := discoverBrokers(client)

		// Update cluster with discovered brokers
		tempCluster.Brokers = discoveredBrokers
		tempCluster.IsAuthenticated = true

		// Get cluster name
		fmt.Print("Enter a name for this cluster connection: ")
		nameInput, err := reader.ReadString('\n')
		if err != nil {
			logger.Error("Error reading cluster name input", "error", err)
			return
		}
		nameInput = strings.TrimSpace(nameInput)
		if nameInput == "" {
			nameInput = fmt.Sprintf("cluster-%d", len(currentSession.clusters)+1)
		}
		tempCluster.Name = nameInput

		// Check if cluster with this name already exists
		for i, existingCluster := range currentSession.clusters {
			if existingCluster.Name == nameInput {
				currentSession.clusters[i] = tempCluster
				currentSession.activeCluster = nameInput
				fmt.Printf("Updated existing cluster connection: %s\n", nameInput)

				// Close the test client
				err = client.Close()
				if err != nil {
					logger.Error("Error closing client", "error", err)
				}

				// Save session
				err = saveSession()
				if err != nil {
					logger.Error("Error saving session", "error", err)
				}
				return
			}
		}

		// Add new cluster
		currentSession.clusters = append(currentSession.clusters, tempCluster)
		currentSession.activeCluster = nameInput
		fmt.Printf("Added new cluster connection: %s\n", nameInput)

		// Close the test client
		err = client.Close()

		if err != nil {
			logger.Error("Error closing client", "error", err)
		}

		// Save session
		err = saveSession()
		if err != nil {
			logger.Error("Error saving session", "error", err)
		}
	} else {
		logger.Error("Error connecting to cluster", "error", err)
	}
}

func LoginWithParams(brokers []string, version string, clusterName string) (bool, string) {
	// Create temporary cluster connection for testing
	tempCluster := ClusterConnection{
		Brokers:         brokers,
		Version:         version,
		IsAuthenticated: false,
	}

	// Test connection
	kafkaVersion, err := sarama.ParseKafkaVersion(version)
	if err != nil {
		logger.Error("Invalid Kafka version string", "version", version, "error", err)
		return false, "Invalid Kafka version string: " + err.Error()
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := cluster.NewCluster(tempCluster.Brokers, kafkaVersion).Connect(ctx)

	if client != nil && err == nil {
		discoveredBrokers := discoverBrokers(client)

		// Update cluster with discovered brokers
		tempCluster.Brokers = discoveredBrokers
		tempCluster.IsAuthenticated = true
		tempCluster.Name = clusterName

		existing := false

		// Check if cluster with this name already exists
		for i, existingCluster := range currentSession.clusters {
			if existingCluster.Name == clusterName {
				currentSession.clusters[i] = tempCluster
				currentSession.activeCluster = clusterName
				existing = true
			}
		}

		if !existing {
			// Add new cluster
			currentSession.clusters = append(currentSession.clusters, tempCluster)
			currentSession.activeCluster = clusterName
		}

		// Close the test client
		err = client.Close()
		if err != nil {
			logger.Error("Error closing client", "error", err)
		}

		// Save session
		err = saveSession()
		if err != nil {
			logger.Error("Error saving session", "error", err)
			return false, "Error saving session: " + err.Error()
		}
		return true, "Saved cluster connection: " + clusterName
	} else {
		logger.Error("Error connecting to cluster", "error", err)
		return false, "Error connecting to cluster: " + err.Error()
	}
}

func discoverBrokers(client sarama.Client) []string {
	// Auto-discover and display all brokers in the cluster
	brokers := client.Brokers()
	fmt.Printf("Auto-discovered %d brokers in cluster:\n", len(brokers))

	// Update cluster with all discovered brokers
	discoveredBrokers := make([]string, 0, len(brokers))
	for _, broker := range brokers {
		fmt.Printf("  - Broker %d: %s\n", broker.ID(), broker.Addr())
		discoveredBrokers = append(discoveredBrokers, broker.Addr())
	}

	return discoveredBrokers
}

func Logout(clusterName string) bool {
	if clusterName == "" {
		if currentSession.activeCluster == "" {
			fmt.Println("No active cluster session.")
			return false
		}
		clusterName = currentSession.activeCluster
	}

	// Find and remove the cluster
	for i, cluster := range currentSession.clusters {
		if cluster.Name == clusterName {
			currentSession.clusters = append(currentSession.clusters[:i], currentSession.clusters[i+1:]...)

			// If this was the active cluster, clear it
			if currentSession.activeCluster == clusterName {
				currentSession.activeCluster = ""
				// Disconnect current client if any
				if currentSession.client != nil {
					err := currentSession.client.Close()
					if err != nil {
						logger.Error("Error closing client", "error", err)
					}
					currentSession.client = nil
				}
				if currentSession.adminClient != nil {
					err := currentSession.adminClient.Close()
					if err != nil {
						logger.Error("Error closing admin client", "error", err)
					}
					currentSession.adminClient = nil
				}
			}

			// Save session
			err := saveSession()
			if err != nil {
				fmt.Println("Error saving session:", err)
				return false
			}

			fmt.Printf("Logged out from cluster: %s\n", clusterName)
			return true
		}
	}

	fmt.Printf("Cluster '%s' not found in active sessions.\n", clusterName)
	return false
}

func GetClusterConnections() []ClusterConnection {
	return currentSession.clusters
}

func GetActiveClusterName() string {
	return currentSession.activeCluster
}

func ListClusters() {
	if len(currentSession.clusters) == 0 {
		fmt.Println("No cluster connections found.")
		return
	}

	fmt.Println("Available cluster connections:")
	for i, cluster := range currentSession.clusters {
		status := "Disconnected"
		if cluster.IsAuthenticated {
			status = "Connected"
		}

		active := ""
		if cluster.Name == currentSession.activeCluster {
			active = " (ACTIVE)"
		}

		fmt.Printf("%d. %s - %s - %d brokers%s\n",
			i+1, cluster.Name, status, len(cluster.Brokers), active)

		for j, broker := range cluster.Brokers {
			fmt.Printf("   Broker %d: %s\n", j+1, broker)
		}
	}
}

func cleanupClients() {
	if currentSession.client != nil {
		err := currentSession.client.Close()
		if err != nil {
			logger.Error("Error closing client", "error", err)
		}
		currentSession.client = nil
	}
	if currentSession.adminClient != nil {
		err := currentSession.adminClient.Close()
		if err != nil {
			logger.Error("Error closing admin client", "error", err)
		}
		currentSession.adminClient = nil
	}
}

func SelectCluster(clusterName string) {
	for _, cluster := range currentSession.clusters {
		if cluster.Name == clusterName {
			cleanupClients()

			currentSession.activeCluster = clusterName
			fmt.Printf("Selected cluster: %s\n", clusterName)

			// Save session
			err := saveSession()
			if err != nil {
				logger.Error("Error saving session", "error", err)
			}
			return
		}
	}

	fmt.Printf("Cluster '%s' not found. Use 'ok cluster list' to see available clusters.\n", clusterName)
}

func DisplaySession() {
	if currentSession == nil {
		fmt.Println("No active session.")
		return
	}

	if currentSession.IsAuthenticated() {
		fmt.Println("Current session:", currentSession.Info())
	} else {
		fmt.Println("No active session.")
	}
}

func InitAPI() error {
	err := loadSession()
	if err != nil {
		logger.Error("Error loading session", "error", err)
		return err
	}
	return nil
}

func GetClusterByName(clusterName string) *ClusterConnection {
	for i := range currentSession.clusters {
		if currentSession.clusters[i].Name == clusterName {
			return &currentSession.clusters[i]
		}
	}
	return nil
}
