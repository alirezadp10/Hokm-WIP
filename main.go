package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/redis/rueidis"
	"gopkg.in/telebot.v4"
	"log"
	"os"
)

func main() {
	_ = godotenv.Load()
	bot, err := telebot.NewBot(telebot.Settings{
		Token: os.Getenv("TELEGRAM_BOT_TOKEN"),
		Poller: &telebot.Webhook{
			Listen: "0.0.0.0:80", // Address to listen for incoming webhook requests
			Endpoint: &telebot.WebhookEndpoint{
				PublicURL: "https://8bb2-194-107-126-16.ngrok-free.app", // Replace with your actual public URL
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// Define a command handler
	bot.Handle("/start", func(c telebot.Context) error {
		client, _ := rueidis.NewClient(rueidis.ClientOption{InitAddress: []string{"redis:6379"}})
		defer client.Close()
		client.Do(context.Background(), client.B().Set().Key("key").Value("val").Nx().Build())
		return c.Send("Webhook is now active!")
	})

	// Start the bot
	bot.Start()
}
