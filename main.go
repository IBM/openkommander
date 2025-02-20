package main

import (
	"fmt"

	"github.com/IBM/openkommander/pkg/commands"
)

func main() {
	fmt.Println("CLI App: type exit to quit")
	commands.DisplayCommands()
	commands.RunCLI()
}
