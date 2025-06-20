package session

import (
	"bufio"
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strconv"
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
	isSecureKafka   bool
	username        string
	password        string
}

func (s *session) Info() string {
	return fmt.Sprintf("Brokers: %v, Authenticated: %v, Version: %v, Secure: %v", s.brokers, s.isAuthenticated, s.version, s.isSecureKafka)
}

func (s *session) Connect(ctx context.Context) (sarama.Client, error) {
	if s.client != nil {
		return s.client, nil
	}
	client, err := cluster.NewCluster(getSaramaConfig(s), s.brokers).Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("error connecting to cluster: %w", err)
	}
	adminClient, err := cluster.NewCluster(getSaramaConfig(s), s.brokers).ConnectAdmin(ctx)
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
		if err := s.client.Close(); err != nil {
			fmt.Printf("Error closing client: %v\n", err)
		}
	}
	s.client = nil
	s.adminClient = nil
	s.isAuthenticated = false
	s.version = constants.SaramaKafkaVersion
	s.isSecureKafka = false
	s.username = ""
	s.password = ""
	s.brokers = []string{}
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
	adminClient, err := cluster.NewCluster(getSaramaConfig(s), s.brokers).ConnectAdmin(ctx)
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

func saveSession() error {
	envMap := make(map[string]string)
	envMap["OK_BROKERS"] = strings.Join(currentSession.brokers, ",")
	envMap["OK_IS_AUTHENTICATED"] = strconv.FormatBool(currentSession.isAuthenticated)
	envMap["OK_VERSION"] = currentSession.version.String()
	envMap["OK_IS_SECURE_KAFKA"] = strconv.FormatBool(currentSession.isSecureKafka)
	envMap["OK_USERNAME"] = currentSession.username
	envMap["OK_PASSWORD"] = currentSession.password

	err := godotenv.Write(envMap, ".env")
	if err != nil {
		fmt.Println("Error encoding session data:", err)
		return err
	}
	return nil
}

func loadSession() error {
	currentSession.brokers = strings.Split(os.Getenv("OK_BROKERS"), ",")
	currentSession.isAuthenticated, _ = strconv.ParseBool(os.Getenv("OK_IS_AUTHENTICATED"))
	currentSession.version, _ = sarama.ParseKafkaVersion(os.Getenv("OK_VERSION"))
	currentSession.isSecureKafka, _ = strconv.ParseBool(os.Getenv("OK_IS_SECURE_KAFKA"))
	currentSession.username = os.Getenv("OK_USERNAME")
	currentSession.password = os.Getenv("OK_PASSWORD")
	return nil
}

func init() {
	currentSession = &session{
		brokers:         []string{},
		isAuthenticated: false,
		client:          nil,
		adminClient:     nil,
		version:         constants.SaramaKafkaVersion,
		isSecureKafka:   false,
		username:        "",
		password:        "",
	}
	envReadError := godotenv.Load()
	err := loadSession()
	if err != nil || envReadError != nil {
		fmt.Println("Error loading session:", err)
	}
}

func Login() {
	if currentSession.IsAuthenticated() {
		fmt.Println("Already logged in.")
		return
	}
	version := readUserInput("Enter kafka version [" + constants.KafkaVersion + "]: ")
	currentSession.version, _ = sarama.ParseKafkaVersion(version)

	currentSession.isSecureKafka = readUserClosedInput("Is your kafka configured with SASL_PLAINTEXT security? (y/n): ")
	if currentSession.isSecureKafka {
		currentSession.username = readUserInput("Enter configured username: ")
		currentSession.password = readUserInput("Enter configured password: ")
	}

	broker := readUserInput("Enter broker address [" + constants.KafkaBroker + "]: ")
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

func readUserInput(inputMessage string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf(inputMessage)

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		fmt.Println("Please enter a valid Input")
		readUserInput(inputMessage)
	}
	return input
}

func readUserClosedInput(inputMessage string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf(inputMessage)

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "y" {
		return true
	} else if input == "n" {
		return false
	} else {
		fmt.Println("Please enter a valid Input")
		readUserClosedInput(inputMessage)
	}
	return false
}

func getSaramaConfig(s *session) *sarama.Config {
	config := sarama.NewConfig()
	config.Version = s.version
	if s.isSecureKafka {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = s.username
		config.Net.SASL.Password = s.password
		config.Net.SASL.Mechanism = sarama.SASLTypePlaintext
	}
	return config
}
