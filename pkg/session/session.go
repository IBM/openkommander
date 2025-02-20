package session

import (
	"bufio"
	"context"
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

func (s session) Info() string {
	return fmt.Sprintf("Brokers: %v, Authenticated: %v", s.brokers, s.isAuthenticated)
}

func (s session) Connect(ctx context.Context) (sarama.Client, error) {
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
	currentSession = s
	return client, nil
}

func (s session) Disconnect() {
	if s.client != nil {
		s.client.Close()
	}
	s.client = nil
	s.adminClient = nil
	s.isAuthenticated = false
	currentSession = s
	fmt.Println("Logged out successfully!")
}

func (s session) IsAuthenticated() bool {
	return s.isAuthenticated
}

func (s session) GetClient() sarama.Client {
	if s.client == nil {
		return nil
	}
	return s.client
}

func (s session) GetAdminClient() sarama.ClusterAdmin {
	if s.adminClient == nil {
		return nil
	}

	return s.adminClient
}

func GetCurrentSession() Session {
	return currentSession
}

var currentSession Session

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	err := viper.ReadInConfig()
	if err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		fmt.Println("No config file found, using defaults")
	}

	currentSession = session{
		brokers:         []string{},
		isAuthenticated: false,
		client:          nil,
		adminClient:     nil,
	}
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

	currentSession = session{
		brokers:         []string{broker},
		isAuthenticated: false,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := currentSession.Connect(ctx)
	if client != nil && err == nil {
		fmt.Println("Logged in successfully!")
	} else {
		fmt.Println("Error connecting to cluster.")
	}
}

func Logout() {
	if !currentSession.IsAuthenticated() {
		fmt.Println("No active session.")
		return
	}

	currentSession.Disconnect()
}

func DisplaySession() {
	if currentSession.IsAuthenticated() {
		fmt.Println("Current session:", currentSession.Info())
	} else {
		fmt.Println("No active session.")
	}
}
