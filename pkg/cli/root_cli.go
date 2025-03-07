package cli

import (
	"github.com/IBM/openkommander/pkg/functions"
)

type RootCommands struct {
	Root *OkParentCmd

	Login    *OkCmd
	Logout   *OkCmd
	Session  *OkCmd
	Metadata *OkCmd

	Children *RootChildren
}

type RootChildren struct {
	Server *ServerCommands
	Topic  *TopicCommands
}

func GetRootCommands() *RootCommands {
	rootCommands := &RootCommands{
		Root: &OkParentCmd{
			Use:   "ok",
			Short: "OpenKommander - A CLI tool for Apache Kafka management",
			Long: `OpenKommander is a command line utility for Apache Kafka compatible brokers.
					Complete documentation is available at https://github.com/IBM/openkommander`,
			Aliases: []string{"openkommander", "kommander", "okm"},
		},
		Login: &OkCmd{
			Use:   "login",
			Short: "Connect to a Kafka cluster",
			Run:   login,
		},
		Logout: &OkCmd{
			Use:   "logout",
			Short: "End the current session",
			Run:   logout,
		},
		Session: &OkCmd{
			Use:   "session",
			Short: "Display current session information",
			Run:   getSessionInfo,
		},
		Metadata: &OkCmd{
			Use:   "metadata",
			Short: "Display cluster information",
			Run:   getClusterMetadata,
		},
		Children: &RootChildren{
			Server: GetServerCommands(),
			Topic:  GetTopicCommands(),
		},
	}

	return rootCommands
}

func login(cmd cobraCmd, args cobraArgs) {
	functions.Login()
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
