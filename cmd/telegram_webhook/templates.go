package main

import "html/template"

var templateNotifyMessage = `
**Hello {{.UserName}}!**
I see you are trying to work with me.
Let me help you by showing you a list of available commands:

These are some commands to *sneak on somebodies* workflow:
/notify
/notify add commit
/notify add issue mentioned <username>
/notify add issue assigned <username>
/notify remove commit
/notify remove issue mentioned <username>
/notify remove issue assigned <username>
/notify list

Well, you can also display the current version - if you want to:
/version

Isn´t that helpful? Just try it out!
`
var parsedTemplateNotifyMessage *template.Template

func parseTemplates() error {
	var err error
	parsedTemplateNotifyMessage, err = template.New("TemplateNotifyMessage").Parse(templateNotifyMessage)
	if err != nil {
		return err
	}

	return nil
}
