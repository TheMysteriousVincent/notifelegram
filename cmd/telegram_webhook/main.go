package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/playnet-public/flagenv"

	"github.com/gorilla/mux"

	"github.com/urfave/negroni"

	tba "github.com/TheMysteriousVincent/telegram-bot-api"
)

const (
	webhookPath = "/v1/webhooks/telegram/"
)

var (
	apiKey      = flagenv.String("telegram-bot-api-key", "", "This is your telegram bot-api-token")
	webhookHost = flagenv.String("webhook-host", "localhost", "The webhook host")
	webhookPort = flagenv.Int("webhook-port", 3030, "The webhook port")
	host        = flagenv.String("host", "0.0.0.0", "The host of the http server")
	port        = flagenv.Int("port", 3030, "The port of the http server")
)

func main() {
	flagenv.Parse()

	bot, err := tba.NewBotAPI(*apiKey)
	if err != nil {
		log.Fatal(err)
	}

	_, err = bot.SetWebhook(tba.NewWebhook(fmt.Sprintf("%s:%d/%s", *webhookHost, *webhookPort, webhookPath)))
	if err != nil {
		log.Fatal(err)
	}

	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Fatal(err)
	}

	if info.LastErrorDate != 0 {
		log.Printf("[Telegram callback failed]%s", info.LastErrorMessage)
	}

	n := negroni.Classic()
	r := mux.NewRouter()
	r.Methods("POST").Path("/v1/webhooks/telegram/").Name("TelegramWebhookHandler").HandlerFunc(handleTelegramWebhook)
	n.UseHandler(r)

	http.ListenAndServe("5.9.96.244:3030", n)
}

func handleTelegramWebhook(w http.ResponseWriter, r *http.Request) {
	update, err := tba.GetWebhookUpdate(r)
	if err != nil {
		w.Write([]byte("false"))
	}

	log.Println(update)
}
