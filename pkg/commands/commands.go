package commands

import (
	"bufio"
	"fmt"
	"os"
)

type CommandFunc func()

var registry = map[string]CommandFunc{}

var display = map[string][]string{}

func init() {
	Register("CLI", "help", DisplayCommands)
	Register("CLI", "exit", ExitCommand)
}

func Register(group string, name string, fn CommandFunc) {
	registry[name] = fn

	display[group] = append(display[group], name)
}

func GetCommands() map[string]CommandFunc {
	return registry
}

func DisplayCommands() {
	fmt.Print("Available commands: ")

	for group, commands := range display {
		fmt.Printf("\n  %s:\n", group)
		for _, command := range commands {
			fmt.Printf("    %s\n", command)
		}
	}

	fmt.Println()
}

func RunCLI() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		cmd := scanner.Text()
		if fn, ok := registry[cmd]; ok {
			fn()
		} else {
			fmt.Println("Unknown command.")
		}
	}
}

func ExitCommand() {
	fmt.Println("Exiting application.")
	os.Exit(0)
}
