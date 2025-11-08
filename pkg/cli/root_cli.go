package cli

import (
	"fmt"

	"github.com/IBM/openkommander/internal/core/commands"
	"github.com/IBM/openkommander/pkg/session"
	"github.com/IBM/sarama"
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
		{ // Login
			Use:   "login",
			Short: "Connect to a Kafka cluster",
			Run:   login,
			Flags: []OkFlag{
				NewOkFlag(OkFlagString, "username", "u", "username for cluster"),
				NewOkFlag(OkFlagString, "password", "p", "password for cluster"),
			},
		},
		{ // Logout
			Use:   "logout [cluster-name]",
			Short: "End a cluster session",
			Run:   logout,
		},
		{ // Session info
			Use:   "session",
			Short: "Display current session information",
			Run:   getSessionInfo,
		},
		{ // Cluster metadata
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
		&BrokerCommandList{},
		&ProduceCommandList{},
		&ClusterCommandList{},
	}
}

func login(cmd cobraCmd, args cobraArgs) {
	session.Login()
}

func logout(cmd cobraCmd, args cobraArgs) {
	clusterName := ""
	if len(args) > 0 {
		clusterName = args[0]
	}
	session.Logout(clusterName)
}

func getSessionInfo(cmd cobraCmd, args cobraArgs) {
	session.DisplaySession()
}

func getClusterMetadata(cmd cobraCmd, args cobraArgs) {
	client, validateFailure := commands.GetClient()
	if validateFailure != nil {
		fmt.Print(validateFailure.Err)
		return
	}

	brokers := client.Brokers()

	brokerHeaders := []string{"ID", "Address", "Connected"}
	brokerRows := [][]interface{}{}
	for _, b := range brokers {
		connected := "No"
		if err := b.Open(client.Config()); err == nil || err == sarama.ErrAlreadyConnected {
			connected = "Yes"
		}
		brokerRows = append(brokerRows, []interface{}{
			b.ID(),
			b.Addr(),
			connected,
		})
	}
	RenderTable("Cluster Brokers:", brokerHeaders, brokerRows)
}
