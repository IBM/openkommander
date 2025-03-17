package cli

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

func startRESTServer(cmd cobraCmd, args cobraArgs) {
	
}
