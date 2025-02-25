package commands

import (
	"github.com/IBM/openkommander/pkg/session"
)

func loginCommand() {
	session.Login()
}

func logoutCommand() {
	session.Logout()
}

func sessionInfoCommand() {
	session.DisplaySession()
}
