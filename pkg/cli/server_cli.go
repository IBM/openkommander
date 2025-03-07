package cli

import "github.com/IBM/openkommander/pkg/functions"

type ServerCommands struct {
	Server *OkParentCmd

	Start *OkCmd
	Stop  *OkCmd
}

func GetServerCommands() *ServerCommands {
	serverCommands := &ServerCommands{
		Server: &OkParentCmd{
			Use:   "server",
			Short: "REST server commands",
		},
		Start: &OkCmd{
			Use:   "start",
			Short: "Start the REST server",
			Run:   startRESTServer,
		},
		Stop: &OkCmd{
			Use:   "stop",
			Short: "Stop the REST server",
			Run:   stopRESTServer,
		},
	}

	return serverCommands
}

func startRESTServer(cmd cobraCmd, args cobraArgs) {
	functions.StartRESTServer()
}

func stopRESTServer(cmd cobraCmd, args cobraArgs) {
	functions.StopRESTServer()
}
