package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/TheMysteriousVincent/notifelegram/pkg/telegram"
	gitlab "github.com/xanzy/go-gitlab"

	"github.com/playnet-public/flagenv"

	"gopkg.in/telegram-bot-api.v4"

	_ "github.com/lib/pq" //needs to be force because its a used database type
)

const (
	webhookPath = "/v1/webhooks/telegram/"
)

var (
	apiKey         = flagenv.String("telegram-bot-api-key", "", "This is your telegram bot-api-token")
	webhookHost    = flagenv.String("webhook-host", "localhost", "The webhook host")
	webhookPort    = flagenv.Int("webhook-port", 88, "The webhook port")
	host           = flagenv.String("host", "0.0.0.0", "The host of the http server")
	port           = flagenv.Int("port", 88, "The port of the http server")
	certFile       = flagenv.String("cert-file", "", "The certfile to establish a secure connection")
	keyFile        = flagenv.String("key-file", "", "The keyfile to establish a secure connection")
	gitlabKey      = flagenv.String("gitlab-api-key", "", "The GitLab server API key")
	gitlabBaseURL  = flagenv.String("gitlab-base-url", "", "The GitLab API base url")
	dbHost         = flagenv.String("dbHost", "localhost", "The default database host")
	dbPort         = flagenv.Int("dbPort", 5432, "The default database port")
	dbUser         = flagenv.String("dbUser", "notifier", "The default database username")
	dbPass         = flagenv.String("dbPass", "1234", "The default database password")
	dbName         = flagenv.String("dbName", "notifier", "The default database name")
	helpContent    string
	gitlabClient   *gitlab.Client
	bot            *tgbotapi.BotAPI
	sqlConnection  *sql.DB
	commandHandler = telegram.NewCommandHandler()
	handler        *telegram.Handler
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

	bot.RemoveWebhook()
	/*if err := setWebhook(); err != nil {
		log.Fatal(err)
	}*/

	sqlConnection, err = sql.Open("postgres", fmt.Sprintf(
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

	handler = telegram.NewHandler(bot, sqlConnection, gitlabClient)

	commandHandler = telegram.NewCommandHandler()
	commandHandler.AddCommand("enableCommits", handler.HandleEnableCommits)
	commandHandler.AddCommand("disableCommits", handler.HandleDisableCommits)
	commandHandler.AddCommand("addMentions", handler.HandleAddMentions)
	commandHandler.AddCommand("removeMentions", handler.HandleRemoveMentions)
	commandHandler.AddCommand("version", handler.HandleVersion)
	commandHandler.AddCommand("commitsEnabled", handler.HandleCommitsEnabled)
	commandHandler.AddCommand("listMentions", handler.HandleListMentions)
	commandHandler.DefaultCommandHandle = handler.HandleHelp

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 20

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}

	for update := range updates {
		handleUpdate(&update)
	}

	/*n := negroni.Classic()
	r := mux.NewRouter()
	r.Methods("POST").Path(webhookPath + bot.Token).Name("TelegramWebhookHandler").HandlerFunc(handleTelegramWebhook)
	n.UseHandler(r)

	err = http.ListenAndServeTLS(fmt.Sprintf("%s:%d", *host, *port), *certFile, *keyFile, n)
	if err != nil {
		log.Fatal(err)
	}*/
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
		log.Fatal(err.Error())
		return
	}

	handleUpdate(update)
}

func handleUpdate(update *tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	if update.Message.IsCommand() {
		commandHandler.Execute(update.Message)
	} else {
		handler.HandleHelp(update.Message)
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
