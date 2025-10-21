package cli

import (
	"fmt"

	"github.com/IBM/openkommander/internal/core/commands"
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
		{ // List cluster info
			Use:   "list",
			Short: "List all cluster info",
			Run:   getClusterList,
		},
	}
}

func (ClusterCommandList) GetSubcommands() []CommandList {
	return nil
}

func getClusterList(cmd cobraCmd, args cobraArgs) {
	// Use the command from internal/core/commands
	clusters, failure := commands.ListClusters()
	if failure != nil {
		fmt.Println(failure.Err)
		return
	}

	// Prepare table headers and rows
	clusterHeaders := []string{"Cluster ID", "Address", "Status", "Rack"}
	clusterRows := [][]interface{}{}

	for _, cluster := range clusters {
		clusterRows = append(clusterRows, []interface{}{
			cluster.ID,
			cluster.Address,
			cluster.Status,
			cluster.Rack,
		})
	}

	RenderTable("Available Clusters:", clusterHeaders, clusterRows)
}
