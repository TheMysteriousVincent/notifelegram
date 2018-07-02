package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/playnet-public/flagenv"

//	"github.com/gorilla/mux"
//	"github.com/urfave/negroni"

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
	certFile = flagenv.String("cert-file", "", "The certfile to establish a secure connection")
	keyFile = flagenv.String("key-file", "", "The keyfile to establish a secure connection")
)

func main() {
	flagenv.Parse()

	bot, err := tba.NewBotAPI(*apiKey)
	if err != nil {
		log.Fatal(err)
	}
	bot.Debug = true
	webhookAddress := fmt.Sprintf(
		"https://%s:%d/%s",
		*webhookHost,
		*webhookPort,
		bot.Token,
	)

	_, err = bot.SetWebhook(
		tba.NewWebhookWithCert(
			webhookAddress,
			*certFile,
		),
	)

	if err != nil {
		log.Fatal(err)
	}

	updates := bot.ListenForWebhook("/" + bot.Token)
	go http.ListenAndServeTLS("0.0.0.0:88", "cert.pem", "key.pem", nil)

	for update := range updates {
		log.Printf("%+v\n", update)
	}
}

func handleTelegramWebhook(w http.ResponseWriter, r *http.Request) {
	log.Println(r)
/*	update, err := tba.GetWebhookUpdate(r)
	if err != nil {
		w.Write([]byte("false"))
	}

	log.Println(update)*/
}
