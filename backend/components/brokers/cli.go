package brokers

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"openkommander/lib/factory"
	"openkommander/lib/utils"
)

func NewCommand(clientFactory *factory.ClientFactory) *cobra.Command {
	brokerCmd := &cobra.Command{
		Use:   "brokers",
		Short: "List information about Kafka brokers",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := clientFactory.CreateClientFromFlags(cmd)
			utils.HandleCLIError(err, "Failed to connect to Kafka")
			defer client.Close()
			
			brokers, err := client.GetBrokerInfo()
			utils.HandleCLIError(err, "Failed to list brokers")

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			
			fmt.Fprintln(w, "ID\tHOST\tPORT\tLEADER PARTITIONS\tTOTAL PARTITIONS\tIN-SYNC PARTITIONS")
			fmt.Fprintln(w, "--\t----\t----\t-----------------\t----------------\t------------------")
			
			for _, broker := range brokers {
				fmt.Fprintf(w, "%d\t%s\t%d\t%d\t%d\t%d\n", 
					broker.ID, 
					broker.Host, 
					broker.Port,
					broker.PartitionsLeader,
					broker.Partitions,
					broker.InSyncPartitions)
			}
			
			w.Flush()
		},
	}

	return brokerCmd
}
