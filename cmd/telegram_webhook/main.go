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
	gitlabClient     *gitlab.Client
	bot              *tba.BotAPI
)

func main() {
	flagenv.Parse()

	gitlabClient = gitlab.NewClient(nil, *gitlabKey)
	gitlabClient.SetBaseURL("https://gitlab.allgameplay.de/api/v4/")

	var err error
	bot, err = tba.NewBotAPI(*apiKey)
	if err != nil {
		log.Fatal(err)
	}
	bot.Debug = true

	_, err = bot.SetWebhook(
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

func handleTelegramWebhook(w http.ResponseWriter, r *http.Request) {
	update, err := tba.GetWebhookUpdate(r)
	if err != nil {
		w.Write([]byte("false"))
		return
	}

	if update.Message.IsCommand() {
	} else {
		fmt.Println(22)
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
