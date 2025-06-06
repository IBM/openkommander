package cli

import (
	"fmt"

	"github.com/IBM/openkommander/pkg/rest"
	"github.com/spf13/cobra"
)

type ServerCommandList struct{}

func (ServerCommandList) GetParentCommand() *OkParentCmd {
	return &OkParentCmd{
		Use:   "server",
		Short: "REST server commands",
	}
}

func (ServerCommandList) GetCommands() []*OkCmd {
	return []*OkCmd{
		{
			Use:   "start",
			Short: "Start the REST server",
			Run:   startRESTServer,
			Flags: []OkFlag{
				NewOkFlag(OkFlagString, "port", "p", "Specify the port for the REST server"),
			},
			RequiredFlags: []string{"port"},
		},
	}
}

func (ServerCommandList) GetSubcommands() []CommandList {
	return nil
}

func startRESTServer(cmd *cobra.Command, args []string) {
	port, _ := cmd.Flags().GetString("port")

	if port == "" {
		fmt.Println("Error: Port is required")
		return
	}

	rest.StartRESTServer(port)
}
