package cli

import (
	"reflect"

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
}

type OkParentCmd struct {
	Use                    string
	Short                  string
	Long                   string
	Flags                  []OkFlag
	Aliases                []string
	EnforceFlagConstraints func(cmd cobraCmd)
}

func Init() cobraCmd {
	return RegisterCommands(GetRootCommands())
}

// Go through commands recursively and build tree of commands
func RegisterCommands(commands any) cobraCmd {
	var rValue = reflect.Indirect(reflect.ValueOf(commands))

	var parentCmd cobraCmd
	var cmd cobraCmd
	var okParentCmd *OkParentCmd
	var okCmd *OkCmd

	var subcommands = make([]*cobra.Command, 0)

	for i := range rValue.NumField() {
		var field = rValue.Field(i)

		switch field.Type() {
		case reflect.TypeOf(&OkParentCmd{}):
			okParentCmd = field.Interface().(*OkParentCmd)
			parentCmd = &cobra.Command{
				Use:   okParentCmd.Use,
				Short: okParentCmd.Short,
			}
		case reflect.TypeOf(&OkCmd{}):
			okCmd = field.Interface().(*OkCmd)
			cmd = &cobra.Command{
				Use:   okCmd.Use,
				Short: okCmd.Short,
				Run:   okCmd.Run,
			}

			if okCmd.Flags != nil {
				for _, flag := range okCmd.Flags {
					switch flag.ValueType {
					case "string":
						cmd.Flags().StringP(flag.Name, flag.ShortName, "", flag.Usage)
					case "int":
						cmd.Flags().IntP(flag.Name, flag.ShortName, 0, flag.Usage)
					}
				}

				if okCmd.EnforceFlagConstraints != nil {
					okCmd.EnforceFlagConstraints(cmd)
				}
			}

			subcommands = append(subcommands, cmd)
		case reflect.TypeOf(&RootChildren{}):
			rootChildren := field.Interface().(*RootChildren)
			parentCmd.AddCommand(RegisterCommands(rootChildren.Server))
			parentCmd.AddCommand(RegisterCommands(rootChildren.Topic))
		}
	}

	for _, cmd := range subcommands {
		parentCmd.AddCommand(cmd)
	}

	return parentCmd
}
