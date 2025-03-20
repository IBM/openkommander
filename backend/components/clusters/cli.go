package clusters

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"openkommander/lib/factory"
	"openkommander/lib/utils"
)

func NewCommand(clientFactory *factory.ClientFactory) *cobra.Command {
	clustersCmd := &cobra.Command{
		Use:   "clusters",
		Short: "List configured Kafka clusters",
		Long:  `List configured Kafka clusters from the configuration file`,
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List configured clusters",
		Run: func(cmd *cobra.Command, args []string) {
			configPath, _ := cmd.Flags().GetString("config")
			
			cfg, err := clientFactory.LoadConfig(configPath)
			if err != nil {
				utils.HandleCLIError(err, "Failed to load config")
				return
			}
			
			if len(cfg.Clusters) == 0 {
				fmt.Println("No clusters configured in the config file.")
				return
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			
			fmt.Fprintln(w, "NAME\tBROKERS\tDESCRIPTION")
			fmt.Fprintln(w, "----\t-------\t-----------")
			
			for _, cluster := range cfg.Clusters {
				fmt.Fprintf(w, "%s\t%v\t%s\n", 
					cluster.Name, 
					cluster.Brokers, 
					cluster.Description)
			}
			
			w.Flush()
		},
	}

	clustersCmd.AddCommand(listCmd)
	return clustersCmd
}
