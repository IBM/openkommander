package utils

import (
	"os"
	"fmt"
	"log"
	"sort"
	"strings"
	"net"
	"strconv"

	"github.com/IBM/sarama"
	"github.com/spf13/cobra"
	"github.com/gin-gonic/gin"

)

func PrintRoutes(router *gin.Engine) {
	routes := router.Routes()
	
	sort.Slice(routes, func(i, j int) bool {
		return routes[i].Path < routes[j].Path
	})
	
	fmt.Println("\n=== REGISTERED API ROUTES ===")
	fmt.Printf("%-7s %-50s %s\n", "METHOD", "PATH", "HANDLER")
	fmt.Println(strings.Repeat("-", 80))
	
	for _, route := range routes {
		handlerName := route.Handler
		if idx := strings.LastIndex(handlerName, "."); idx != -1 {
			handlerName = handlerName[idx+1:]
		}
		fmt.Printf("%-7s %-50s %s\n", route.Method, route.Path, handlerName)
	}
	
	fmt.Println(strings.Repeat("-", 80))
	log.Printf("Server registered %d routes\n", len(routes))
}

func GetSASLMechanism(mechanismName string) sarama.SASLMechanism {
	switch mechanismName {
	case "SCRAM-SHA-256":
		return sarama.SASLTypeSCRAMSHA256
	case "SCRAM-SHA-512":
		return sarama.SASLTypeSCRAMSHA512
	default:
		return sarama.SASLTypePlaintext
	}
}

func SplitHostPort(addr string) (string, int, error) {
	host, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		return "", 0, err
	}
	
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return "", 0, err
	}
	
	return host, port, nil
}

func AddKafkaFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().String("cluster", "", "Use named cluster from config")
	cmd.PersistentFlags().String("config", "", "Path to config file (defaults to $HOME/.config/openkommander.json)")
}

func ReadStdin() ([]byte, error) {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return nil, err
	}

	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return nil, fmt.Errorf("no data provided on stdin (expected to be used in a pipe)")
	}

	return os.ReadFile("/dev/stdin")
}

