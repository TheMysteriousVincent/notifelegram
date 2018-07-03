package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	gitlab "github.com/xanzy/go-gitlab"

	"github.com/playnet-public/flagenv"

	"github.com/urfave/negroni"

	tba "github.com/TheMysteriousVincent/telegram-bot-api"
)

const (
	webhookPath = "/v1/webhooks/telegram/"
)

var (
	apiKey           = flagenv.String("telegram-bot-api-key", "", "This is your telegram bot-api-token")
	webhookHost      = flagenv.String("webhook-host", "localhost", "The webhook host")
	webhookPort      = flagenv.Int("webhook-port", 88, "The webhook port")
	host             = flagenv.String("host", "0.0.0.0", "The host of the http server")
	port             = flagenv.Int("port", 88, "The port of the http server")
	certFile         = flagenv.String("cert-file", "", "The certfile to establish a secure connection")
	keyFile          = flagenv.String("key-file", "", "The keyfile to establish a secure connection")
	helpFile         = flagenv.String("helpFile", "help.md", "This is the default help page. Will be displayed everytime if a entered command was unknown or the message was no command")
	defaultParseMode = flagenv.String("default-msg-parse-mode", "Markdown", "The parse mode of a message to send (help page and templates). Can be Markdown or HTML as mentioned in https://core.telegram.org/bots/api#sendmessage")
	gitlabKey        = flagenv.String("gitlab-api-key", "", "The GitLab server API key")
	gitlabBaseURL    = flagenv.String("gitlab-base-url", "", "The GitLab API base url")
	helpContent      string
	gitlabClient     *gitlab.Client
	bot              *tba.BotAPI
)

func main() {
	flagenv.Parse()

	gitlabClient = gitlab.NewClient(nil, *gitlabKey)
	gitlabClient.SetBaseURL(*gitlabBaseURL)

	var err error
	bot, err = tba.NewBotAPI(*apiKey)
	if err != nil {
		log.Fatal(err)
	}

	if err := setWebhook(); err != nil {
		log.Fatal(err)
	}

	n := negroni.Classic()
	r := mux.NewRouter()
	r.Methods("POST").Path(webhookPath + bot.Token).Name("TelegramWebhookHandler").HandlerFunc(handleTelegramWebhook)
	n.UseHandler(r)

	err = http.ListenAndServeTLS(fmt.Sprintf("%s:%d", *host, *port), *certFile, *keyFile, n)
	if err != nil {
		log.Fatal(err)
	}
}

func setWebhook() error {
	_, err := bot.SetWebhook(
		tba.NewWebhookWithCert(
			fmt.Sprintf(
				"https://%s:%d%s%s",
				*webhookHost,
				*webhookPort,
				webhookPath,
				bot.Token,
			),
			*certFile,
		),
	)
	return err
}

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


PSQL Table for notifies:
nid(serial primary key), uid(int), notifier (varchar(256)), notfier_value (text)
*/

func handleTelegramWebhook(w http.ResponseWriter, r *http.Request) {
	update, err := tba.GetWebhookUpdate(r)
	if err != nil {
		w.Write([]byte("false"))
		return
	}

	if update.Message == nil {
		w.Write([]byte("false"))
		return
	}

	if update.Message.IsCommand() {
		cmds, _ := update.Message.GetCommands()

		for _, cmd := range *cmds {
			switch cmd.Name {
			case "notify":
				if len(cmd.Arguments) > 0 {
					switch cmd.Arguments[0] {
					case "add":
						if len(cmd.Arguments) > 1 {
							switch cmd.Arguments[1] {
							case "commit":
								handleAddCommitNotify()
							case "issue":
								if len(cmd.Arguments) > 2 {
									switch cmd.Arguments[2] {
									case "mentioned":
										if len(cmd.Arguments) > 3 {
										}
									case "assigned":
										if len(cmd.Arguments) > 3 {

										}
									}
								}
							}
						}
					case "remove":
						if len(cmd.Arguments) > 1 {
							switch cmd.Arguments[1] {
							case "commit":
							case "issue":
								if len(cmd.Arguments) > 2 {
									switch cmd.Arguments[2] {
									case "mentioned":
										if len(cmd.Arguments) > 3 {
										}
									case "assigned":
										if len(cmd.Arguments) > 3 {

										}
									}
								}
							}
						}
					}
				}
			case "version":
			}
		}
	} else {
		f, err := os.Open(*helpFile)
		if err != nil {
			w.Write([]byte("false"))
			return
		}

		b, err := ioutil.ReadAll(f)
		if err != nil {
			w.Write([]byte("false"))
			return
		}

		msg := tba.NewMessage(update.Message.Chat.ID, string(b))
		msg.ParseMode = *defaultParseMode
		bot.Send(msg)
	}
}

func handleAddCommitNotify() error {

}
