package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/playnet-public/flagenv"

	"github.com/urfave/negroni"

	tba "github.com/TheMysteriousVincent/telegram-bot-api"
)

const (
	webhookPath = "/v1/webhooks/telegram/"
)

var (
	apiKey      = flagenv.String("telegram-bot-api-key", "", "This is your telegram bot-api-token")
	webhookHost = flagenv.String("webhook-host", "localhost", "The webhook host")
	webhookPort = flagenv.Int("webhook-port", 88, "The webhook port")
	host        = flagenv.String("host", "0.0.0.0", "The host of the http server")
	port        = flagenv.Int("port", 88, "The port of the http server")
	certFile    = flagenv.String("cert-file", "", "The certfile to establish a secure connection")
	keyFile     = flagenv.String("key-file", "", "The keyfile to establish a secure connection")
)

func main() {
	flagenv.Parse()

	bot, err := tba.NewBotAPI(*apiKey)
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
	log.Println(r)
	update, err := tba.GetWebhookUpdate(r)
	if err != nil {
		w.Write([]byte("false"))
	}

	log.Println(update)
}
