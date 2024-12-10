package main

import (
    "encoding/json"
    "fmt"
    "github.com/joho/godotenv"
    "gopkg.in/telebot.v4"
    "log"
    "net/http"
    "os"
    "time"
)

func main() {
    go startTelegramApp()

    http.Handle("/templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir("./templates"))))

    http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
        if request.Method != "GET" {
            http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }
        http.ServeFile(writer, request, "templates/splash.html")
    })

    http.HandleFunc("/room/123456", func(writer http.ResponseWriter, request *http.Request) {
        if request.Method != "GET" {
            http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }
        writer.Header().Set("Content-Type", "application/json")
        json.NewEncoder(writer).Encode(map[string]interface{}{
            "players": map[string]interface{}{
                "up": map[string]interface{}{
                    "name": "ali",
                },
                "right": map[string]interface{}{
                    "name": "maryam",
                },
                "left": map[string]interface{}{
                    "name": "sara",
                },
                "down": map[string]interface{}{
                    "name": "mammad",
                },
            },
            "points": map[string]interface{}{
                "total":        map[string]interface{}{"right": 3, "down": 2},
                "currentRound": map[string]interface{}{"right": 3, "down": 2},
            },
            "firstKingDeterminationCards": []interface{}{},
            "centerCards": map[string]interface{}{
                "up":    "2S",
                "right": "7S",
            },
            "currentTurn":  "down",
            "timeRemained": 14,
            "yourCards":    []interface{}{"3H", "3H", "3S", "3S", "4C"},
            "judge":        "down",
            "trump":        "heart",
        })
        //json.NewEncoder(writer).Encode(map[string]interface{}{
        //    "players": map[string]interface{}{
        //        "up": map[string]interface{}{
        //            "name": "ali",
        //        },
        //        "right": map[string]interface{}{
        //            "name": "maryam",
        //        },
        //        "left": map[string]interface{}{
        //            "name": "sara",
        //        },
        //        "down": map[string]interface{}{
        //            "name": "mammad",
        //        },
        //    },
        //    "points": map[string]interface{}{
        //        "total":        map[string]interface{}{"right": 3, "down": 2},
        //        "currentRound": map[string]interface{}{"right": 3, "down": 2},
        //    },
        //    "firstKingDeterminationCards": []interface{}{
        //        map[string]interface{}{"direction": "up", "card": "2S"},
        //        map[string]interface{}{"direction": "right", "card": "3C"},
        //        map[string]interface{}{"direction": "down", "card": "3H"},
        //        map[string]interface{}{"direction": "left", "card": "3S"},
        //        map[string]interface{}{"direction": "up", "card": "AC"},
        //    },
        //    "centerCards":  []interface{}{"2H", "3H", "4C"},
        //    "currentTurn":  "right",
        //    "timeRemained": 14,
        //    "yourCards":    []interface{}{"3H", "3H", "3S", "3S", "4C"},
        //"judge":        "up",
        //    "trump":        "heart",
        //})
    })

    http.HandleFunc("/room/123456/cards", func(writer http.ResponseWriter, request *http.Request) {
        if request.Method != "GET" {
            http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }
        writer.Header().Set("Content-Type", "application/json")
        json.NewEncoder(writer).Encode(map[string]interface{}{
            "cards": []interface{}{
                []interface{}{
                    "3C",
                    "3H",
                    "3S",
                    "8S",
                    "9D",
                },
                []interface{}{
                    "AC",
                    "AH",
                    "2S",
                    "6S",
                    "2D",
                },
                []interface{}{
                    "JS",
                    "KH",
                    "QD",
                },
            },
        })
    })

    http.HandleFunc("/room/123456/place", func(writer http.ResponseWriter, request *http.Request) {
        if request.Method != http.MethodPost {
            http.Error(writer, "Method Not Allowed", http.StatusMethodNotAllowed)
            return
        }
        writer.Header().Set("Content-Type", "application/json")
        json.NewEncoder(writer).Encode(map[string]interface{}{
            "points": map[string]interface{}{
                "total":        map[string]interface{}{"right": 4, "down": 2},
                "currentRound": map[string]interface{}{"right": 0, "down": 3},
            },
            "currentTurn":             "down",
            "timeRemained":            14,
            "judge":                   "down",
            "wasGameOvered":           false,
            "wasRoundOvered":          false,
            "whoHasWonTheRound":       "up",
            "whoHasWonTheGame":        nil,
            "wasKingChanged":          false,
            "trumpDeterminationCards": []interface{}{"3H", "3H", "3S", "3S", "4C"},
        })
    })

    fmt.Println("Server is running at 7070")

    err := http.ListenAndServe(":7070", nil)

    if err != nil {
        fmt.Println("Error starting server:", err)
    }
}

func startTelegramApp() {
    _ = godotenv.Load()
    pref := telebot.Settings{
        Token:  os.Getenv("TELEGRAM_BOT_TOKEN"),
        Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
    }
    bot, err := telebot.NewBot(pref)

    if err != nil {
        log.Fatal(err)
    }

    bot.Handle("/start", func(c telebot.Context) error {
        return c.Send("Let's Play", &telebot.ReplyMarkup{
            InlineKeyboard: [][]telebot.InlineButton{
                {
                    {
                        Text:   "Launch App",
                        WebApp: &telebot.WebApp{URL: "https://76b1-46-100-55-166.ngrok-free.app"},
                    },
                },
            },
        })
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
    })

    fmt.Println("Bot is running...")

    bot.Start()
}
