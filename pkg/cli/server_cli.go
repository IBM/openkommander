package cli

import (
	"fmt"
	"strconv"

	"github.com/IBM/openkommander/pkg/logger"
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
		},
	}
}

func (ServerCommandList) GetSubcommands() []CommandList {
	return nil
}

func startRESTServer(cmd *cobra.Command, args []string) {
	port, _ := cmd.Flags().GetString("port")

	if port == "" {
		fmt.Print("Enter port for the REST server: ")
		if _, err := fmt.Scanln(&port); err != nil {
			fmt.Println("Error reading port number:", err)
			return
		}
	}

	if err := validatePort(port); err != nil {
		logger.Error("Invalid port", "error", err)
		return
	}

	rest.StartRESTServer(port)
}

func validatePort(port string) error {
	if port == "" || port == "0" {
		return fmt.Errorf("port number is required and cannot be 0")
	}

	portNumber, err := strconv.Atoi(port)
	if err != nil {
		return fmt.Errorf("invalid port number: %s", port)
	}

	if portNumber < 1 || portNumber > 65535 {
		return fmt.Errorf("port number must be between 1 and 65535")
	}

	return nil
}
