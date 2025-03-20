package consumers

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"openkommander/lib/factory"
	"openkommander/lib/utils"
)

func NewCommand(clientFactory *factory.ClientFactory) *cobra.Command {
	consumersCmd := &cobra.Command{
		Use:   "consumers",
		Short: "Manage Kafka consumer groups",
		Long:  `List consumer groups, view consumer group details, and consume messages from topics`,
	}
	
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all consumer groups",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := clientFactory.CreateClientFromFlags(cmd)
			utils.HandleCLIError(err, "Failed to connect to Kafka")
			defer client.Close()
			
			groups, err := client.GetConsumerGroups()
			utils.HandleCLIError(err, "Failed to list consumer groups")
			
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "GROUP ID\tMEMBERS\tTOPICS\tTOTAL LAG\tSTATE")
			fmt.Fprintln(w, "--------\t-------\t------\t---------\t-----")
			
			for _, group := range groups {
				fmt.Fprintf(w, "%s\t%d\t%d\t%d\t%s\n",
					group.GroupID,
					group.Members,
					group.Topics,
					group.Lag,
					group.State)
			}
			w.Flush()
		},
	}
	
	describeCmd := &cobra.Command{
		Use:   "describe [group-id]",
		Short: "Describe a consumer group",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := clientFactory.CreateClientFromFlags(cmd)
			utils.HandleCLIError(err, "Failed to connect to Kafka")
			defer client.Close()
			
			groupID := args[0]
			
			group, err := client.GetConsumerGroup(groupID)
			utils.HandleCLIError(err, "Failed to get consumer group info")
			
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			
			fmt.Fprintf(w, "Group ID:\t%s\n", group.GroupID)
			fmt.Fprintf(w, "State:\t%s\n", group.State)
			fmt.Fprintf(w, "Members:\t%d\n", group.Members)
			fmt.Fprintf(w, "Topics:\t%d\n", group.Topics)
			fmt.Fprintf(w, "Total Lag:\t%d\n", group.Lag)
			fmt.Fprintf(w, "Coordinator Broker ID:\t%d\n", group.Coordinator)
			
			w.Flush()
			
			if len(group.TopicLags) > 0 {
				fmt.Println("\nLag by Topic-Partition:")
				
				lagWriter := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
				
				fmt.Fprintln(lagWriter, "TOPIC\tPARTITION\tLAG")
				fmt.Fprintln(lagWriter, "-----\t---------\t---")
				
				topicMap := make(map[string][]struct{Partition int32; Lag int64})
				
				for _, lag := range group.TopicLags {
					topicMap[lag.Topic] = append(topicMap[lag.Topic], struct{Partition int32; Lag int64}{lag.Partition, lag.Lag})
				}
				
				topics := make([]string, 0, len(topicMap))
				for topic := range topicMap {
					topics = append(topics, topic)
				}
				sort.Strings(topics)
				
				for _, topic := range topics {
					lags := topicMap[topic]
					
					sort.Slice(lags, func(i, j int) bool {
						return lags[i].Partition < lags[j].Partition
					})
					
					for _, lag := range lags {
						fmt.Fprintf(lagWriter, "%s\t%d\t%d\n", topic, lag.Partition, lag.Lag)
					}
				}
				
				lagWriter.Flush()
			}
		},
	}
	
	consumersCmd.AddCommand(listCmd, describeCmd)
	
	return consumersCmd
}
