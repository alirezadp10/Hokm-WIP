package main

import (
	"github.com/joho/godotenv"
	"gopkg.in/telebot.v4"
	"log"
	"os"
)

func main() {
	_ = godotenv.Load()
	bot, err := telebot.NewBot(telebot.Settings{
		Token: os.Getenv("TELEGRAM_BOT_TOKEN"),
		Poller: &telebot.Webhook{
			Listen: "0.0.0.0:8443", // Address to listen for incoming webhook requests
			Endpoint: &telebot.WebhookEndpoint{
				PublicURL: "https://8bea-46-100-55-166.ngrok-free.app", // Replace with your actual public URL
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// Define a command handler
	bot.Handle("/start", func(c telebot.Context) error {
		return c.Send("Webhook is now active!")
	})

	// Start the bot
	bot.Start()
}
