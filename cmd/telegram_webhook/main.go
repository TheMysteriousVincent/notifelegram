package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	gitlab "github.com/xanzy/go-gitlab"

	"github.com/playnet-public/flagenv"

	"github.com/urfave/negroni"
	"gopkg.in/telegram-bot-api.v4"
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
	bot              *tgbotapi.BotAPI
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
		w.Write([]byte("false"))
		return
	}

	if update.Message == nil {
		w.Write([]byte("false"))
		return
	}

	if update.Message.IsCommand() {
		var tmpBuf bytes.Buffer
		parsedTemplateNotifyMessage.Execute(&tmpBuf, update.Message.Chat)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, tmpBuf.String())
		msg.ReplyToMessageID = update.Message.MessageID
		msg.ParseMode = *defaultParseMode
		msg.ReplyMarkup = tgbotapi.ForceReply{
			ForceReply: true,
		}
		bot.Send(msg)
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

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, string(b))
		msg.ParseMode = *defaultParseMode
		bot.Send(msg)
	}

	w.Write([]byte("true"))
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
