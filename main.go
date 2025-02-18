package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"openkommander/pkg/session"
)

func main() {
	reader := bufio.NewScanner(os.Stdin)
	fmt.Println("CLI App: commands available: login, logout, session, exit")
	for {
		fmt.Print("> ")
		if !reader.Scan() {
			break
		}
		cmd := strings.TrimSpace(reader.Text())
		switch cmd {
		case "login":
			session.Login()
		case "logout":
			session.Logout()
		case "session":
			session.DisplaySession()
		case "exit":
			fmt.Println("Exiting application.")
			return
		default:
			fmt.Println("Unknown command. Available commands: login, logout, session, exit")
		}
	}
}
