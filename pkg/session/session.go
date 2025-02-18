package session

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Session interface {
	Info() string
}

type session struct {
	cluster string
}

func (s session) Info() string {
	return fmt.Sprintf("Connected to Kafka cluster '%s'", s.cluster)
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
}

func Login() {
	reader := bufio.NewReader(os.Stdin)

	defaultCluster := viper.GetString("kafka.cluster")
	if defaultCluster != "" {
		fmt.Printf("Default Kafka Cluster from config [%s]: ", defaultCluster)
	} else {
		fmt.Print("Enter Kafka Cluster identifier: ")
	}
	cluster, _ := reader.ReadString('\n')
	cluster = strings.TrimSpace(cluster)
	if cluster == "" && defaultCluster != "" {
		cluster = defaultCluster
	}

	currentSession = session{cluster: cluster}
	fmt.Println("Connected successfully!")
}

func Logout() {
	if currentSession == nil {
		fmt.Println("No active session.")
		return
	}

	currentSession = nil
	fmt.Println("Logged out successfully!")
}

func DisplaySession() {
	if currentSession == nil {
		fmt.Println("No active session.")
	} else {
		fmt.Println("Current session:", currentSession.Info())
	}
}
