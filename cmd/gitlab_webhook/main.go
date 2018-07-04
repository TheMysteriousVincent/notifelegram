package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/playnet-public/flagenv"
	"gopkg.in/telegram-bot-api.v4"

	_ "github.com/lib/pq" //needs to be force because its a used database type

	"gopkg.in/go-playground/webhooks.v4"
	"gopkg.in/go-playground/webhooks.v4/gitlab"
)

var (
	apiKey              = flagenv.String("telegram-bot-api-key", "", "This is your telegram bot-api-token")
	host                = flagenv.String("host", "0.0.0.0", "The host of the http server")
	port                = flagenv.Int("port", 88, "The port of the http server")
	gitlabWebhookSecret = flagenv.String("gitlab-webhook-secret", "", "The GitLab server API key")
	dbHost              = flagenv.String("dbHost", "localhost", "The default database host")
	dbPort              = flagenv.Int("dbPort", 5432, "The default database port")
	dbUser              = flagenv.String("dbUser", "notifier", "The default database username")
	dbPass              = flagenv.String("dbPass", "1234", "The default database password")
	dbName              = flagenv.String("dbName", "notifier", "The default database name")
	helpContent         string
	bot                 *tgbotapi.BotAPI
	sqlConnection       *sql.DB
)

func main() {
	flagenv.Parse()

	hook := gitlab.New(&gitlab.Config{
		Secret: *gitlabWebhookSecret,
	})
	hook.RegisterEvents(handleIssues, gitlab.IssuesEvents)
	hook.RegisterEvents(handleComments, gitlab.CommentEvents)
	hook.RegisterEvents(handleConfidentialIssues, gitlab.ConfidentialIssuesEvents)
	hook.RegisterEvents(handlePushes, gitlab.PushEvents)

	var err error
	bot, err = tgbotapi.NewBotAPI(*apiKey)
	if err != nil {
		log.Fatal(err)
	}

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

	err = webhooks.Run(hook, fmt.Sprintf("%s:%d", *host, *port), "/v1/webhooks/gitlab/")
	if err != nil {
		log.Fatal(err.Error())
	}
}

func handleIssues(payload interface{}, header webhooks.Header) {
	pl := payload.(gitlab.IssueEventPayload)

	rows, err := sqlConnection.Query(
		"SELECT chatId FROM mentions WHERE gitlabUsername = $1",
		pl.Assignee.Name,
	)
	if err != nil {
		log.Println(err.Error())
		return
	}

	for rows.Next() {
		var chatId int64
		if err := rows.Scan(&chatId); err != nil {
			log.Println(err.Error())
			return
		}
		if chatId <= 0 {
			continue
		}

		msg := tgbotapi.NewMessage(chatId, "Someone created an issue assigned to you")
		msg.ParseMode = "Markdown"
		bot.Send(msg)
	}
}

func handleComments(payload interface{}, header webhooks.Header) {
	//pl := payload.(gitlab.CommentEventPayload)
}

func handleConfidentialIssues(payload interface{}, header webhooks.Header) {
	//pl := payload.(gitlab.ConfidentialIssueEventPayload)
}

func handlePushes(payload interface{}, header webhooks.Header) {
	//pl := payload.(gitlab.PushEventPayload)
}
