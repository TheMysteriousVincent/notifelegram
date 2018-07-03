package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	gitlab "github.com/xanzy/go-gitlab"

	"github.com/playnet-public/flagenv"

	"github.com/urfave/negroni"
	"gopkg.in/telegram-bot-api.v4"
)

const (
	version     = "v0.1"
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
	dbHost           = flagenv.String("dbHost", "localhost", "The default database host")
	dbPort           = flagenv.Int("dbPort", 5432, "The default database port")
	dbUser           = flagenv.String("dbUser", "notifier", "The default database username")
	dbPass           = flagenv.String("dbPass", "1234", "The default database password")
	dbName           = flagenv.String("dbName", "notifier", "The default database name")
	helpContent      string
	gitlabClient     *gitlab.Client
	bot              *tgbotapi.BotAPI
	sqlCon           *sql.DB
)

func main() {
	flagenv.Parse()

	gitlabClient = gitlab.NewClient(nil, *gitlabKey)
	gitlabClient.SetBaseURL(*gitlabBaseURL)

	var err error
	bot, err = tgbotapi.NewBotAPI(*apiKey)
	if err != nil {
		log.Fatal(err)
	}

	if err := setWebhook(); err != nil {
		log.Fatal(err)
	}

	if err := parseTemplates(); err != nil {
		log.Fatal(err)
	}

	sqlCon, err = sql.Open("postgres", fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		*dbUser,
		*dbPass,
		*dbHost,
		*dbPort,
		*dbName,
	))
	if err != nil {
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
		tgbotapi.NewWebhookWithCert(
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

func handleTelegramWebhook(w http.ResponseWriter, r *http.Request) {
	update, err := getWebhookUpdate(r)
	if err != nil {
		return
	}

	if update.Message == nil {
		return
	}

	if update.Message.IsCommand() {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		msg.ParseMode = *defaultParseMode
		switch update.Message.Command() {
		case "enableCommits":
			var nid int
			if err := sqlCon.QueryRow(
				"INSERT INTO users (type, uid) SELECT 'commit', $1 WHERE NOT EXISTS (SELECT nid FROM users WHERE type = 'commit' AND uid = $2) RETURNING nid",
				update.Message.From.ID,
				update.Message.From.ID,
			).Scan(&nid); err != nil {
				msg.Text = err.Error()
				break
			}

			if nid > 0 {
				msg.Text = "Commits are already enabled."
				break
			}

			msg.Text = "Enabled commit events."
		case "disableCommits":
			msg.Text = "Disabled commit events."
		case "addIssueMention":
			username := update.Message.CommandArguments()
			msg.Text = fmt.Sprintf("Added issue mentions for user '%s'", username)
		case "removeIssueMention":
			username := update.Message.CommandArguments()
			msg.Text = fmt.Sprintf("Removed issue mentions for user '%s'", username)
		case "list":
		case "version":
			msg.Text = fmt.Sprintf("Version: %s", version)
		default:
			var tmpBuf bytes.Buffer
			err := parsedTemplateNotifyMessage.Execute(&tmpBuf, update.Message.Chat)
			if err != nil {
				msg.Text = err.Error()
				break
			}
			msg.Text = tmpBuf.String()
		}
		bot.Send(msg)
	} else {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		msg.ParseMode = *defaultParseMode
		var tmpBuf bytes.Buffer
		err := parsedTemplateNotifyMessage.Execute(&tmpBuf, update.Message.Chat)
		if err != nil {
			msg.Text = err.Error()
		} else {
			msg.Text = tmpBuf.String()
		}
		bot.Send(msg)
	}
}

func getWebhookUpdate(r *http.Request) (*tgbotapi.Update, error) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var update tgbotapi.Update
	err = json.Unmarshal(b, &update)
	return &update, err
}
