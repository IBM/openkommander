package topics

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"context"
	"text/tabwriter"


	"github.com/spf13/cobra"
	"github.com/IBM/sarama"
	"openkommander/lib/factory"
	"openkommander/lib/utils"
)

func NewCommand(clientFactory *factory.ClientFactory) *cobra.Command {
	topicsCmd := &cobra.Command{
		Use:   "topics",
		Short: "Manage Kafka topics",
		Long:  `List, create, describe, and delete Kafka topics`,
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all topics",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := clientFactory.CreateClientFromFlags(cmd)
			utils.HandleCLIError(err, "Failed to connect to Kafka")
			defer client.Close()
			
			topics, err := client.GetTopicInfo()
			utils.HandleCLIError(err, "Failed to list topics")

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			
			fmt.Fprintln(w, "TOPIC\tPARTITIONS\tREPLICATION\tCLEANUP POLICY")
			fmt.Fprintln(w, "-----\t----------\t-----------\t--------------")
			
			for _, topic := range topics {
				fmt.Fprintf(w, "%s\t%d\t%d\t%s\n", 
					topic.Name, 
					topic.Partitions, 
					topic.ReplicationFactor,
					topic.CleanupPolicy)
			}
			
			w.Flush()
		},
	}

	createCmd := &cobra.Command{
		Use:   "create [topic-name]",
		Short: "Create a new Kafka topic",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := clientFactory.CreateClientFromFlags(cmd)
			utils.HandleCLIError(err, "Failed to connect to Kafka")
			defer client.Close()

			partitions := utils.MustGetInt32(cmd, "partitions")
			replication := utils.MustGetInt16(cmd, "replication-factor")
			
			err = client.CreateTopic(args[0], partitions, replication)
			utils.HandleCLIError(err, "Failed to create topic")
			fmt.Printf("Topic '%s' created successfully\n", args[0])
		},
	}
	createCmd.Flags().Int32P("partitions", "p", 1, "Number of partitions")
	createCmd.Flags().Int16P("replication-factor", "r", 1, "Replication factor")

	deleteCmd := &cobra.Command{
		Use:   "delete [topic-name]",
		Short: "Delete a Kafka topic",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := clientFactory.CreateClientFromFlags(cmd)
			utils.HandleCLIError(err, "Failed to connect to Kafka")
			defer client.Close()
			
			err = client.DeleteTopic(args[0])
			utils.HandleCLIError(err, "Failed to delete topic")
			fmt.Printf("Topic '%s' deleted successfully\n", args[0])
		},
	}

	describeCmd := &cobra.Command{
		Use:   "describe [topic-name]",
		Short: "Describe a Kafka topic",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := clientFactory.CreateClientFromFlags(cmd)
			utils.HandleCLIError(err, "Failed to connect to Kafka")
			defer client.Close()

			topic := args[0]
			
			topicDetail, err := client.GetTopic(topic)
			utils.HandleCLIError(err, "Failed to describe topic")

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			
			fmt.Fprintf(w, "Topic:\t%s\n", topicDetail.Name)
			fmt.Fprintf(w, "Partitions:\t%d\n", topicDetail.Partitions)
			fmt.Fprintf(w, "Replication Factor:\t%d\n", topicDetail.ReplicationFactor)
			
			w.Flush()
			
			fmt.Println("\nPartition IDs:")
			for i, id := range topicDetail.PartitionIDs {
				if i > 0 && i%10 == 0 {
					fmt.Println() 
				}
				fmt.Printf("%d ", id)
			}
			fmt.Println() 
		},
	}

	consumeCmd := &cobra.Command{
		Use:   "consume [topic]",
		Short: "Consume messages from a Kafka topic",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := clientFactory.CreateClientFromFlags(cmd)
			utils.HandleCLIError(err, "Failed to connect to Kafka")
			defer client.Close()
	
			topic := args[0]
			group, _ := cmd.Flags().GetString("group")
			fromBeginning, _ := cmd.Flags().GetBool("from-beginning")
	
			initialOffset := sarama.OffsetNewest
			if fromBeginning {
				initialOffset = sarama.OffsetOldest
			}
	
			if group == "" {
				fmt.Printf("Consuming messages from topic '%s' (non-commited mode, no user group)\n", topic)
			} else {
				fmt.Printf("Consuming messages from topic '%s' with group '%s'\n", topic, group)
			}
			fmt.Println("Press Ctrl+C to stop...")
	
			ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
			defer cancel()
	
			err = client.ConsumeMessagesWithOptions(ctx, topic, group, initialOffset, func(msg *sarama.ConsumerMessage) error {
				fmt.Printf("Partition: %d | Offset: %d | Key: %s\n", 
					msg.Partition, msg.Offset, string(msg.Key))
				fmt.Printf("Value: %s\n", string(msg.Value))
				fmt.Println("----------------")
				return nil
			})
	
			if err != nil && err != context.Canceled {
				utils.HandleCLIError(err, "Error consuming messages")
			}
		},
	}
	
	consumeCmd.Flags().String("group", "", "Consumer group ID (optional)")
	consumeCmd.Flags().Bool("from-beginning", false, "Consume messages from beginning of the topic")
	

	topicsCmd.AddCommand(listCmd, createCmd, deleteCmd, describeCmd, consumeCmd)
	return topicsCmd
}
