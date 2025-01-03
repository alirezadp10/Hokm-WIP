package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"gopkg.in/telebot.v4"
)

func StartTelegram(ctx context.Context) {
	bot, err := telebot.NewBot(telebot.Settings{
		Token:  os.Getenv("TELEGRAM_BOT_TOKEN"),
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
	}

	bot.Handle("/start", func(c telebot.Context) error {
		return c.Send("Let's Play", &telebot.ReplyMarkup{
			InlineKeyboard: [][]telebot.InlineButton{{{
				Text:   os.Getenv("APP_URL"),
				WebApp: &telebot.WebApp{URL: os.Getenv("APP_URL")},
			}}},
		})
	})

	fmt.Println("Bot is running...")
	go bot.Start()
	<-ctx.Done()
	fmt.Println("Shutting down bot...")
	bot.Stop()
}
