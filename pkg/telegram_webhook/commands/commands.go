package commands

import (
	"strings"

	tba "github.com/TheMysteriousVincent/telegram-bot-api"
)

/*
Currently available commands:
notify:
	add:
		commit
		issue:
			mentioned:
				<username>
			assigned:
				<username>
	remove:
		commit
		issue:
			mentioned:
				<username>
			assigned:
				<username>
	list:
		-> List of all active notifys
version:



*/

//CommandRoutes represents a slice of CommandRoute
type CommandRoutes map[string]*CommandRoute

//CommandFunc represents a single command function (handler)
type CommandFunc func(*CommandRoute) error

//CommandRoute represents single command
type CommandRoute struct {
	subroutes CommandRoutes
	name      string
	hndl      *CommandFunc
}

//AddSubCommand adds a command as subcommand to another command
func (cr *CommandRoute) AddSubCommand(name string) *CommandRoute {
	r := CommandRoute{
		name: name,
	}
	cr.subroutes[strings.ToLower(name)] = &r
	return &r
}

//GetName returns a name of a route
func (cr *CommandRoute) GetName() string {
	return cr.name
}

//HandlerFunc adds a function to
func (cr *CommandRoute) HandlerFunc(hndl CommandFunc) *CommandRoute {
	cr.hndl = &hndl
	return cr
}

//CommandHandler represents a handler function for commands
type CommandHandler struct {
	routes CommandRoutes
}

//NewCommandHandler creates a new instance of CommandHandler
func NewCommandHandler() *CommandHandler {
	return &CommandHandler{
		routes: CommandRoutes{},
	}
}

//AddCommand creates and appends a command route to the CommandHandler
func (ch *CommandHandler) AddCommand(name string) *CommandRoute {
	r := CommandRoute{
		name: name,
	}
	ch.routes[strings.ToLower(name)] = &r
	return &r
}

func (ch *CommandHandler) Exec(cmds *tba.Commands) error {
}
