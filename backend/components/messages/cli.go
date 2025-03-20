package messages

import (
	"fmt"
	"strings"
	"encoding/json"

	"github.com/spf13/cobra"
	"openkommander/lib/factory"
	"openkommander/lib/utils"
)

func NewCommand(clientFactory *factory.ClientFactory) *cobra.Command {
	producerCmd := &cobra.Command{
		Use:   "produce [topic]",
		Short: "Produce a message to a Kafka topic",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := clientFactory.CreateClientFromFlags(cmd)
			utils.HandleCLIError(err, "Failed to connect to Kafka")
			defer client.Close()

			topic := args[0]
			key, _ := cmd.Flags().GetString("key")
			value, _ := cmd.Flags().GetString("value")
			file, _ := cmd.Flags().GetString("file")
			isJSON, _ := cmd.Flags().GetBool("json")

			if file != "" {
				err = client.ProduceMessageFromFile(topic, key, file, isJSON)
			} else if value != "" {
				if isJSON {
					err = client.ProduceMessage(topic, key, json.RawMessage(value))
				} else {
					err = client.ProduceMessage(topic, key, value)
				}
			} else {
				data, err := utils.ReadStdin()
				utils.HandleCLIError(err, "Failed to read from stdin")
				
				content := string(data)
				if isJSON || (len(content) > 0 && strings.TrimSpace(content)[0] == '{') {
					var reader = strings.NewReader(content)
					err = client.ProduceMessageFromReader(topic, key, reader, true)
				} else {
					var reader = strings.NewReader(content)
					err = client.ProduceMessageFromReader(topic, key, reader, false)
				}
			}

			utils.HandleCLIError(err, "Failed to produce message")
			fmt.Println("Message sent successfully")
		},
	}

	producerCmd.Flags().StringP("key", "k", "", "Message key")
	producerCmd.Flags().StringP("value", "v", "", "Message value")
	producerCmd.Flags().StringP("file", "f", "", "Read message from file")
	producerCmd.Flags().BoolP("json", "j", false, "Treat input as JSON")

	return producerCmd
}
