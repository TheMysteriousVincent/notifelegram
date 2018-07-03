package commands

import (
	"fmt"
	"testing"

	tba "github.com/TheMysteriousVincent/telegram-bot-api"
)

func TestExec(t *testing.T) {
	cmds := &tba.Commands{
		tba.Command{
			Name:      "test",
			Arguments: []string{"subcmd1", "subcmd2", "testVar1", "testVar2"},
		},
	}

	cmdHndl := NewCommandHandler()
	cmdHndl.AddCommand("test").AddSubCommand("subcmd1").AddSubCommand("subcmd2").HandlerFunc(func(msg *tba.Message, vars []string) error {
		fmt.Println(vars)
		return nil
	})
	ce := cmdHndl.NewCommandExecutor(&tba.Message{}, cmds)
	for ce.Next() {
		err := ce.Execute()
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}
