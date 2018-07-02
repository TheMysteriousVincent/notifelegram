package commands

import (
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

type CommandRoute struct {
}

type CommandRoutes struct {
}

type CommandHandler struct {
	cmds *tba.Commands
}

func NewCommandHandler(cmds *tba.Commands) *CommandHandler {
}
