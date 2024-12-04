package main

import (
    "fmt"
    "github.com/joho/godotenv"
    "gopkg.in/telebot.v4"
    "log"
    "os"
    "time"
)

func main() {
    _ = godotenv.Load()
    pref := telebot.Settings{
        Token:  os.Getenv("TELEGRAM_BOT_TOKEN"),
        Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
    }
    bot, err := telebot.NewBot(pref)

    if err != nil {
        log.Fatal(err)
    }

    // Define a command handler
    bot.Handle("/start", func(c telebot.Context) error {
        //client, _ := rueidis.NewClient(rueidis.ClientOption{InitAddress: []string{"redis:6379"}})
        //defer client.Close()

        //subClient := client.B().Subscribe().Channel("").Build() // Subscription command
        //
        //message := c.Message()
        //user := message.Sender
        //
        //foo, _ := json.MarshalIndent(user, "", " ")
        // Print message details
        //fmt.Println(string(foo))
        //return c.Send("سلام خوش اومدی، صبر کن بقیه هم بیان خبرت میکنم")
        fmt.Println("Bot is running...")

        return c.Reply("Hi", telebot.ModeHTML)
    })

    fmt.Println("Bot is running...")

    // Start the bot
    bot.Start()
}
