package main

import (
	"encoding/json"
	"fmt"
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
			Listen: "0.0.0.0:8443",
			Endpoint: &telebot.WebhookEndpoint{
				PublicURL: "https://8bb2-194-107-126-16.ngrok-free.app",
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// Define a command handler
	bot.Handle("/start", func(c telebot.Context) error {

		message := c.Message()
		user := message.Sender

		foo, _ := json.MarshalIndent(user, "", " ")
		// Print message details
		fmt.Println(string(foo))
		return c.Send("سلام خوش اومدی، صبر کن بقیه هم بیان خبرت میکنم")
	})

	// Start the bot
	bot.Start()
}
