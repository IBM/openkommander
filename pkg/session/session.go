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

	"github.com/IBM/sarama"
	"github.com/spf13/viper"
)

type Session interface {
	Info() string
	Connect(ctx context.Context) (sarama.Client, error)
	Disconnect()
	GetClient() sarama.Client
	GetAdminClient() sarama.ClusterAdmin
	IsAuthenticated() bool
}

type session struct {
	brokers         []string
	client          sarama.Client
	adminClient     sarama.ClusterAdmin
	isAuthenticated bool
}

func (s *session) Info() string {
	return fmt.Sprintf("Brokers: %v, Authenticated: %v", s.brokers, s.isAuthenticated)
}

func (s *session) Connect(ctx context.Context) (sarama.Client, error) {
	if s.client != nil {
		return s.client, nil
	}
	client, err := cluster.NewCluster(s.brokers).Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("error connecting to cluster: %w", err)
	}
	adminClient, err := cluster.NewCluster(s.brokers).ConnectAdmin(ctx)
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
	fmt.Println("Logged out successfully!")
}

func (s *session) IsAuthenticated() bool {
	return s.isAuthenticated
}

func (s *session) GetClient() sarama.Client {
	if s.client == nil {
		return nil
	}
	return s.client
}

func (s *session) GetAdminClient() sarama.ClusterAdmin {
	if s.adminClient == nil {
		return nil
	}

	return s.adminClient
}

func GetCurrentSession() *session {
	return currentSession
}

var currentSession *session

func createDefaultSession() error {
	file, err := os.Create(".openkommander_config")
	if err != nil {
		fmt.Println("Error creating session file:", err)
		return err
	}

	defer file.Close()

	sessionData := map[string]interface{}{
		"brokers":         []string{},
		"isAuthenticated": false,
	}

	encoder := json.NewEncoder(file)
	return encoder.Encode(sessionData)
}

func saveSession() error {
	file, err := os.Create(".openkommander_config")
	if err != nil {
		fmt.Println("Error creating session file:", err)
		return err
	}
	defer file.Close()
	sessionData := map[string]interface{}{
		"brokers":         currentSession.brokers,
		"isAuthenticated": currentSession.isAuthenticated,
	}
	encoder := json.NewEncoder(file)
	err = encoder.Encode(sessionData)
	if err != nil {
		fmt.Println("Error encoding session data:", err)
		return err
	}
	return nil
}

func loadSession() error {
	file, err := os.Open(".openkommander_config")
	if err != nil {
		err = createDefaultSession()
		if err != nil {
			fmt.Println("Error creating session file:", err)
			return err
		}

		file, err = os.Open(".openkommander_config")
		if err != nil {
			fmt.Println("Error opening session file:", err)
			return err
		}
	}
	defer file.Close()

	sessionData := map[string]interface{}{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&sessionData)
	if err != nil {
		fmt.Println("Error decoding session data:", err)
		return err
	}

	brokersInterface := sessionData["brokers"].([]interface{})
	brokers := make([]string, len(brokersInterface))
	for i, v := range brokersInterface {
		brokers[i] = v.(string)
	}

	currentSession.brokers = brokers
	currentSession.isAuthenticated = sessionData["isAuthenticated"].(bool)
	return nil
}

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.ReadInConfig()

	currentSession = &session{
		brokers:         []string{},
		isAuthenticated: false,
		client:          nil,
		adminClient:     nil,
	}

	loadSession()
}

func Login() {
	if currentSession.IsAuthenticated() {
		fmt.Println("Already logged in.")
		return
	}

	reader := bufio.NewReader(os.Stdin)

	defaultBroker := viper.GetString("kafka.broker")
	if defaultBroker != "" {
		fmt.Printf("Enter broker address [%s]: ", defaultBroker)
	} else {
		fmt.Print("Enter broker address: ")
	}

	broker, _ := reader.ReadString('\n')
	broker = strings.TrimSpace(broker)
	if broker == "" {
		broker = defaultBroker
	}

	currentSession.brokers = []string{broker}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := currentSession.Connect(ctx)
	if client != nil && err == nil {
		fmt.Println("Logged in successfully!")
		saveSession()
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
	saveSession()
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
