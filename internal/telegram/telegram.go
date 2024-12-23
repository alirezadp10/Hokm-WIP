package telegram

import (
    "context"
    "fmt"
    "github.com/alirezadp10/hokm/internal/database/sqlite"
    "github.com/alirezadp10/hokm/internal/utils/crypto"
    "github.com/alirezadp10/hokm/internal/utils/trans"
    "gopkg.in/telebot.v4"
    "gorm.io/gorm"
    "log"
    "net/url"
    "os"
    "time"
)

func Start(ctx context.Context, db *gorm.DB) {
    bot, err := telebot.NewBot(telebot.Settings{
        Token:  os.Getenv("TELEGRAM_BOT_TOKEN"),
        Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
    })

    if err != nil {
        log.Fatal(err)
    }

    bot.Handle("/start", func(context telebot.Context) error {
        return startHandler(context, db)
    })

    fmt.Println("Bot is running...")
    go bot.Start()
    <-ctx.Done()
    fmt.Println("Shutting down bot...")
    bot.Stop()
}

func startHandler(c telebot.Context, db *gorm.DB) error {
    player, err := sqlite.SavePlayer(db, c.Sender(), c.Chat().ID)
    if err != nil {
        log.Fatalf("couldn't save: %v", err)
    }

    encrypted, _ := crypto.Encrypt(player.Username)

    return c.Send("Let's Play", &telebot.ReplyMarkup{
        InlineKeyboard: [][]telebot.InlineButton{{{
            Text:   trans.Get("start the game"),
            WebApp: &telebot.WebApp{URL: os.Getenv("APP_URL") + "?user_id=" + url.QueryEscape(encrypted)},
        }}},
    })
}
