package telegram

import (
    "context"
    "fmt"
    "github.com/alirezadp10/hokm/internal/util/trans"
    "gopkg.in/telebot.v4"
    "log"
    "os"
    "time"
)

func Start(ctx context.Context) {
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
                Text:   trans.Get("start the game"),
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
