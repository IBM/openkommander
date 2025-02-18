package commands

import (
	"openkommander/pkg/session"
)

func init() {
	Register("Session", "login", loginCommand)
	Register("Session", "logout", logoutCommand)
	Register("Session", "session", sessionInfoCommand)
}

func loginCommand() {
	session.Login()
}

func logoutCommand() {
	session.Logout()
}

func sessionInfoCommand() {
	session.DisplaySession()
}
