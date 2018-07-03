package commands

import (
	"errors"
	"strings"

	tba "github.com/TheMysteriousVincent/telegram-bot-api"
)

var (
	//ErrNoEndpoint is thrown when a route exists and there is no handler
	ErrNoEndpoint = errors.New("there exist no command handling endpoint")
)

//CommandRoutes represents a slice of CommandRoute
type CommandRoutes map[string]*CommandRoute

//CommandFunc represents a single command function (handler)
type CommandFunc func(*tba.Message, []string) error

//CommandRoute represents single command
type CommandRoute struct {
	subroutes CommandRoutes
	name      string
	hndl      CommandFunc
}

//AddSubCommand adds a command as subcommand to another command
func (cr *CommandRoute) AddSubCommand(name string) *CommandRoute {
	var r *CommandRoute
	if v, ok := cr.subroutes[strings.ToLower(name)]; ok {
		r = v
	} else {
		r = &CommandRoute{
			name:      name,
			subroutes: CommandRoutes{},
		}
		cr.subroutes[strings.ToLower(name)] = r
	}
	return r
}

//GetName returns a name of a route
func (cr *CommandRoute) GetName() string {
	return cr.name
}

//HandlerFunc adds a function to
func (cr *CommandRoute) HandlerFunc(hndl CommandFunc) *CommandRoute {
	cr.hndl = hndl
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
	var r *CommandRoute
	if v, ok := ch.routes[strings.ToLower(name)]; ok {
		r = v
	} else {
		r = &CommandRoute{
			name:      name,
			subroutes: CommandRoutes{},
		}
		ch.routes[strings.ToLower(name)] = r
	}
	return r
}

//CommandExecutor represents an instance to execute all entered commands
type CommandExecutor struct {
	cmds  *tba.Commands
	ch    *CommandHandler
	index int
	msg   *tba.Message
}

//NewCommandExecutor creates a new CommandExecutor instance
func (ch *CommandHandler) NewCommandExecutor(msg *tba.Message, cmds *tba.Commands) *CommandExecutor {
	return &CommandExecutor{
		cmds: cmds,
		ch:   ch,
		msg:  msg,
	}
}

//Next checks if the
func (ce *CommandExecutor) Next() bool {
	return ce.index < len(*ce.cmds)
}

//Execute executes the next command in queue
func (ce *CommandExecutor) Execute() error {
	defer ce.increaseIndex()
	cmd := (*ce.cmds)[ce.index]
	var r *CommandRoute
	for _, r = range ce.ch.routes {
		if cmd.Name == r.name {
			break
		}
	}

	if r == nil {
		return ErrNoEndpoint
	}

	var i = -1
	var arg string
	for i, arg = range cmd.Arguments {
		for k, sR := range r.subroutes {
			if k == arg {
				r = sR
				break
			}
		}
	}

	if r.hndl == nil && i == -1 {
		return ErrNoEndpoint
	}

	return r.hndl(ce.msg, cmd.Arguments[i-1:len(cmd.Arguments)])
}

func (ce *CommandExecutor) increaseIndex() {
	ce.index++
}
