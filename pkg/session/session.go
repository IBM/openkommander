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
	brokers         []string
	client          sarama.Client
	adminClient     sarama.ClusterAdmin
	isAuthenticated bool
	version         sarama.KafkaVersion
}

type SessionData struct {
	Brokers         []string `json:"brokers"`
	IsAuthenticated bool     `json:"isAuthenticated"`
	Version         string   `json:"version"`
}

func (s *session) Info() string {
	return fmt.Sprintf("Brokers: %v, Authenticated: %v, Version: %v", s.brokers, s.isAuthenticated, s.version)
}

func (s *session) Connect(ctx context.Context) (sarama.Client, error) {
	if s.client != nil {
		return s.client, nil
	}
	client, err := cluster.NewCluster(s.brokers, s.version).Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("error connecting to cluster: %w", err)
	}
	adminClient, err := cluster.NewCluster(s.brokers, s.version).ConnectAdmin(ctx)
	if err != nil {
		return nil, fmt.Errorf("error connecting to cluster as admin: %w", err)
	}
	s.client = client
	s.adminClient = adminClient
	s.isAuthenticated = true
	return client, nil
}

func (s *session) Disconnect() {
	if s.client != nil {
		s.client.Close()
	}
	s.client = nil
	s.adminClient = nil
	s.isAuthenticated = false
	s.version = constants.SaramaKafkaVersion
	fmt.Println("Logged out successfully!")
}

func (s *session) IsAuthenticated() bool {
	return s.isAuthenticated
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

func (s *session) GetAdminClient() (sarama.ClusterAdmin, error) {
	if s.adminClient != nil {
		return s.adminClient, nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	adminClient, err := cluster.NewCluster(s.brokers, s.version).ConnectAdmin(ctx)
	if err != nil {
		return nil, err
	}
	s.adminClient = adminClient
	return adminClient, nil
}

func (s *session) GetBrokers() []string {
	return s.brokers
}

func GetCurrentSession() *session {
	return currentSession
}

var currentSession *session

func createDefaultSession() error {
	file, err := os.Create(constants.OpenKommanderConfigFilename)
	if err != nil {
		fmt.Println("Error creating session file:", err)
		return err
	}
	defer file.Close()

	sessionData := SessionData{Brokers: []string{}, IsAuthenticated: false, Version: constants.SaramaKafkaVersion.String()}
	return json.NewEncoder(file).Encode(sessionData)
}

func saveSession() error {
	file, err := os.Create(constants.OpenKommanderConfigFilename)
	if err != nil {
		fmt.Println("Error creating session file:", err)
		return err
	}
	defer file.Close()

	sessionData := SessionData{Brokers: currentSession.brokers, IsAuthenticated: currentSession.isAuthenticated, Version: currentSession.version.String()}
	err = json.NewEncoder(file).Encode(sessionData)
	if err != nil {
		fmt.Println("Error encoding session data:", err)
		return err
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
	defer file.Close()

	var data SessionData
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	if err != nil {
		fmt.Println("Error decoding session data:", err)
		return err
	}

	currentSession.brokers = data.Brokers
	currentSession.isAuthenticated = data.IsAuthenticated
	currentSession.version, _ = sarama.ParseKafkaVersion(data.Version)
	return nil
}

func init() {
	currentSession = &session{
		brokers:         []string{},
		isAuthenticated: false,
		client:          nil,
		adminClient:     nil,
		version:         constants.SaramaKafkaVersion,
	}

	err := loadSession()
	if err != nil {
		fmt.Println("Error loading session:", err)
	}
}

func Login() {
	if currentSession.IsAuthenticated() {
		fmt.Println("Already logged in.")
		return
	}

	versionReader := bufio.NewReader(os.Stdin)
	fmt.Printf("Enter kafka version [%s]: ", constants.KafkaVersion)

	version, _ := versionReader.ReadString('\n')
	version = strings.TrimSpace(version)
	if version == "" {
		version = constants.KafkaVersion
	}
	currentSession.version, _ = sarama.ParseKafkaVersion(version)

	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Enter broker address [%s]: ", constants.KafkaBroker)

	broker, _ := reader.ReadString('\n')
	broker = strings.TrimSpace(broker)
	if broker == "" {
		broker = constants.KafkaBroker
	}

	currentSession.brokers = []string{broker}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := currentSession.Connect(ctx)

	if client != nil && err == nil {
		fmt.Println("Logged in successfully!")
		fmt.Printf("Kafka Version [%s]\n", currentSession.version)
		err = saveSession()
		if err != nil {
			fmt.Println("Error saving session:", err)
		}
	} else {
		fmt.Printf("Error connecting to cluster: %v\n", err)
	}
}

func Logout() {
	if !currentSession.IsAuthenticated() {
		fmt.Println("No active session.")
		return
	}

	currentSession.Disconnect()
	err := saveSession()
	if err != nil {
		fmt.Println("Error saving session:", err)
	}
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
