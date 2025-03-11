package cli

import (
	"github.com/spf13/cobra"
)

type cobraCmd = *cobra.Command
type cobraArgs = []string

// Used for command flags
type OkFlag struct {
	Name      string
	ShortName string // Optional
	ValueType string
	Usage     string
}

type OkCmd struct {
	Use                    string
	Short                  string
	Long                   string
	Run                    func(cmd cobraCmd, args cobraArgs)
	Flags                  []OkFlag
	Aliases                []string
	EnforceFlagConstraints func(cmd cobraCmd)
	RequiredFlags          []string
	Args                   []string
}

type OkParentCmd = OkCmd

func Init() cobraCmd {
	return RegisterCommands(RootCommandList{})
}

type CommandList interface {
	GetParentCommand() *OkParentCmd
	GetCommands() []*OkCmd
	GetSubcommands() []CommandList
}

// Go through commands recursively and build tree of commands
func RegisterCommands(commandList CommandList) cobraCmd {
	if commandList == nil {
		return nil
	}

	var parentCommand = cobraCmdFromOkCmd(commandList.GetParentCommand())

	for _, command := range commandList.GetCommands() {
		parentCommand.AddCommand(cobraCmdFromOkCmd(command))
	}

	for _, subcommand := range commandList.GetSubcommands() {
		parentCommand.AddCommand(RegisterCommands(*subcommand))
	}

	return parentCommand
}

func cobraCmdFromOkCmd(command *OkCmd) cobraCmd {
	cmd := &cobra.Command{
		Use:     command.Use,
		Short:   command.Short,
		Long:    command.Long,
		Run:     command.Run,
		Aliases: command.Aliases,
	}

	if command.Flags != nil {
		for _, flag := range command.Flags {
			switch flag.ValueType {
			case "string":
				cmd.Flags().StringP(flag.Name, flag.ShortName, "", flag.Usage)
			case "int":
				cmd.Flags().IntP(flag.Name, flag.ShortName, 0, flag.Usage)
			}
		}
	}

	if len(command.RequiredFlags) > 0 {
		cmd.MarkFlagsRequiredTogether(command.RequiredFlags...)
		cmd.Flags().
	}

	return cmd
}
