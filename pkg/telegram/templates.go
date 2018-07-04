package telegram

import (
	"html/template"
	"log"
)

func init() {
	var err error
	ParsedTemplateHelp, err = template.New("TemplateHelp").Parse(templateHelp)
	if err != nil {
		log.Fatal(err)
	}
	ParsedTemplateListMentions, err = template.New("TemplateListMentions").Parse(templateListMentions)
	if err != nil {
		log.Fatal(err)
	}
}

var (
	ParsedTemplateHelp *template.Template
	templateHelp       = `
*Hello {{.UserName}}!*
I see you are trying to work with me.
Let me help you by showing you a list of available commands:

You can get updates on commits:
/enableCommits
/disableCommits

These are some commands to *sneak on somebodies* workflow:
/addMentions <username>
/removeMentions <username>

List your current subscriptions:
/commitsEnabled
/listMentions

Well, you can also display the current version - if you want to:
/version

Isn´t that helpful? Just try it out!
`
)

var (
	ParsedTemplateListMentions *template.Template
	templateListMentions       = `
*Subscribed Mentions:*

{{if not .}}None.{{else}}{{range .}}- Mentions of {{.Username}} added {{.DaysAgo}} day(s) ago{{end}}{{end}}
`
)
