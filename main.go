package main

import (
	"fmt"

	"openkommander/pkg/commands"
	_ "openkommander/pkg/commands"
)

func main() {
	fmt.Println("CLI App: type exit to quit")
	commands.DisplayCommands()
	commands.RunCLI()
}
