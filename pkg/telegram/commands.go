package telegram

import (
	"strings"

	"gopkg.in/telegram-bot-api.v4"
)

type CommandHandle func(*tgbotapi.Message)

type CommandHandler struct {
	commands             map[string]CommandHandle
	DefaultCommandHandle CommandHandle
}

func NewCommandHandler() *CommandHandler {
	return &CommandHandler{
		commands: map[string]CommandHandle{},
	}
}

func (ch *CommandHandler) AddCommand(name string, hndl CommandHandle) {
	ch.commands[strings.ToLower(name)] = hndl
}

func (ch *CommandHandler) Execute(msg *tgbotapi.Message) {
	if hndl, ok := ch.commands[strings.ToLower(msg.Command())]; ok {
		hndl(msg)
	} else {
		ch.DefaultCommandHandle(msg)
	}
}
