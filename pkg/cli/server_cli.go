package cli

import (
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
			},
		}
	}

	func (ServerCommandList) GetSubcommands() []CommandList {
		return nil
	}

	func startRESTServer(cmd *cobra.Command, args []string) {
		port := "8080" 
		brokers := []string{"localhost:9092"} 
	
		if len(args) > 0 {
			port = args[0]
		}
		if len(args) > 1 {
			brokers = args[1:]
		}
	
		rest.StartRESTServer(port, brokers)
	}
