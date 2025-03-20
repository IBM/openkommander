package cli

import (
	"fmt"
	"strings"

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
				NewOkFlag(OkFlagString, "brokers", "b", "Specify the Kafka brokers to connect to"),
			},
			RequiredFlags: []string{"port", "brokers"},
		},
	}
}

func (ServerCommandList) GetSubcommands() []CommandList {
	return nil
}

func startRESTServer(cmd *cobra.Command, args []string) {
	port, _ := cmd.Flags().GetString("port")
	brokerslist, _ := cmd.Flags().GetString("brokers")

	brokers := strings.Split(brokerslist, ",")

	if port == "" {
		fmt.Println("Error: Port is required")
		return
	}
	
	if len(brokers) == 0 {
		fmt.Println("Error: At least one broker is required")
		return
	}

	rest.StartRESTServer(port, brokers)
}
