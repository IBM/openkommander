package cli

import (
	"fmt"
	"strings"

	"github.com/IBM/openkommander/pkg/session"
)

type ClusterCommandList struct{}

func (ClusterCommandList) GetParentCommand() *OkParentCmd {
	return &OkParentCmd{
		Use:   "cluster <command>",
		Short: "Cluster management commands",
	}
}

func (m ClusterCommandList) GetCommands() []*OkCmd {
	return []*OkCmd{
		{ // List cluster connections
			Use:   "list",
			Short: "List all cluster connections",
			Run:   listClusterConnections,
		},
		{ // Select cluster
			Use:   "select <cluster-name>",
			Short: "Select active cluster",
			Run:   selectCluster,
		},
	}
}

func (ClusterCommandList) GetSubcommands() []CommandList {
	return nil
}

func listClusterConnections(cmd cobraCmd, args cobraArgs) {
	clusters := session.GetClusterConnections()
	activeCluster := session.GetActiveClusterName()

	if len(clusters) == 0 {
		fmt.Println("No cluster connections found.")
		return
	}

	// Prepare table headers and rows
	connectionHeaders := []string{"Name", "Status", "Brokers", "Version", "Active"}
	connectionRows := [][]interface{}{}

	for _, cluster := range clusters {
		status := "Disconnected"
		if cluster.IsAuthenticated {
			status = "Connected"
		}

		active := "No"
		if cluster.Name == activeCluster {
			active = "Yes"
		}

		brokersStr := strings.Join(cluster.Brokers, "\n")

		connectionRows = append(connectionRows, []interface{}{
			cluster.Name,
			status,
			brokersStr,
			cluster.Version,
			active,
		})
	}

	RenderTable("Clusters:", connectionHeaders, connectionRows)
}

func selectCluster(cmd cobraCmd, args cobraArgs) {
	if len(args) == 0 {
		fmt.Println("Usage: ok cluster select <cluster-name>")
		return
	}

	clusterName := args[0]
	session.SelectCluster(clusterName)
}
