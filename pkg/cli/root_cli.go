package cli

import (
	"github.com/IBM/openkommander/pkg/functions"
)

type RootCommandList struct{}

func (RootCommandList) GetParentCommand() *OkParentCmd {
	return &OkParentCmd{
		Use:   "ok",
		Short: "OpenKommander - A CLI tool for Apache Kafka management",
		Long: `OpenKommander is a command line utility for Apache Kafka compatible brokers.
				Complete documentation is available at https://github.com/IBM/openkommander`,
	}
}

func (RootCommandList) GetCommands() []*OkCmd {
	return []*OkCmd{
		{
			Use:   "login [URL] [flags]",
			Short: "Connect to a Kafka cluster",
			Run:   login,
			Flags: []OkFlag{
				{
					Name:      "username",
					ShortName: "u",
					ValueType: "string",
					Usage:     "Username for cluster",
				},
				{
					Name:      "password",
					ShortName: "p",
					ValueType: "string",
					Usage:     "Password for cluster",
				},
			},
		},
		{
			Use:   "logout",
			Short: "End the current session",
			Run:   logout,
		},
		{
			Use:   "session",
			Short: "Display current session information",
			Run:   getSessionInfo,
		},
		{
			Use:   "metadata",
			Short: "Display cluster information",
			Run:   getClusterMetadata,
		},
	}
}

func (RootCommandList) GetSubcommands() []CommandList {
	return []CommandList{
		&TopicCommandList{},
		&ServerCommandList{},
	}
}

func login(cmd cobraCmd, args cobraArgs) {
	url := cmd.Flags().Arg(0)
	functions.Login(url)
}

func logout(cmd cobraCmd, args cobraArgs) {
	functions.Logout()
}

func getSessionInfo(cmd cobraCmd, args cobraArgs) {
	functions.GetSessionInfo()
}

func getClusterMetadata(cmd cobraCmd, args cobraArgs) {
	functions.GetClusterMetadata()
}
