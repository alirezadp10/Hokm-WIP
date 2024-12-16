package handler

import (
    "context"
    "errors"
    "github.com/alirezadp10/hokm/internal/database/redis"
    "github.com/alirezadp10/hokm/internal/database/sqlite"
    "github.com/alirezadp10/hokm/internal/helper/crypto"
    "github.com/alirezadp10/hokm/internal/helper/myslice"
    "github.com/alirezadp10/hokm/internal/helper/trans"
    "github.com/alirezadp10/hokm/internal/hokm"
    "github.com/labstack/echo/v4"
    "github.com/redis/rueidis"
    "log"
    "net/http"
    "strings"
    "time"
)

func GetGameData(c echo.Context) error {
    gameId := c.QueryParam("gameId")

    usernameDecrypted, err := crypto.Decrypt([]byte(c.QueryParam("username")))
    if err != nil {
        return c.JSON(http.StatusForbidden, map[string]interface{}{
            "message": trans.Get("Username is not valid."),
        })
    }

    username := string(usernameDecrypted)

    sqliteClient := sqlite.GetNewConnection()

    if !sqlite.DoesPlayerBelongsToThisGame(sqliteClient, username, gameId) {
        return c.JSON(http.StatusForbidden, map[string]interface{}{
            "message": trans.Get("It's not your game."),
        })
    }

    ctx := context.Background()

    redisClient := redis.GetNewConnection()

    gameInformation := redis.GetGameInformation(ctx, redisClient, gameId)

    players := gameInformation["players"].([]string)

    uIndex := myslice.GetIndex(username, players)

    response := map[string]interface{}{
        "players":      hokm.GetPlayersWithDirections(players, uIndex),
        "points":       hokm.GetPoints(gameInformation["points"].(map[string]interface{}), uIndex),
        "centerCards":  hokm.GetCenterCards(gameInformation["center_cards"].(map[string]string), players, uIndex),
        "turn":         hokm.GetDirection(gameInformation["turn"].(string), players, uIndex),
        "judge":        hokm.GetDirection(gameInformation["judge"].(string), players, uIndex),
        "trump":        gameInformation["trump"].(string),
        "timeRemained": 15,
        "yourCards":    gameInformation["cards"].(map[string]interface{})[username],
        "kingsCards":   hokm.GetKingsCards(gameInformation["kings_cards"].([]string), uIndex),
    }

    return c.JSON(http.StatusOK, response)
}

func GetYourCards(c echo.Context) error {
    time.Sleep(2 * time.Second)
    response := map[string]interface{}{
        "trump": "heart",
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
    }
    return c.JSON(http.StatusOK, response)
}

func PlaceCard(c echo.Context) error {
    response := map[string]interface{}{
        "points": map[string]interface{}{
            "total":        map[string]interface{}{"right": 4, "down": 2},
            "currentRound": map[string]interface{}{"right": 0, "down": 3},
        },
        "currentTurn":       "down",
        "timeRemained":      14,
        "judge":             "right",
        "whoHasWonTheCards": "up",
        "whoHasWonTheRound": nil,
        "whoHasWonTheGame":  nil,
        "wasKingChanged":    false,
        //"trumpDeterminationCards": []interface{}{"3H", "3H", "3S", "3S", "4C"},
        "trumpDeterminationCards": nil,
    }
    return c.JSON(http.StatusOK, response)
}

func GetUpdate(c echo.Context) error {
    time.Sleep(2 * time.Second)
    response := map[string]interface{}{
        "lastMove": map[string]interface{}{
            "from": "right",
            "card": "3C",
        },
        "centerCards": map[string]interface{}{"up": "2H", "left": "3H"},
        "points": map[string]interface{}{
            "total":        map[string]interface{}{"right": 4, "down": 2},
            "currentRound": map[string]interface{}{"right": 0, "down": 3},
        },
        "currentTurn":       "down",
        "timeRemained":      14,
        "judge":             "up",
        "whoHasWonTheCards": nil,
        "whoHasWonTheRound": nil,
        "whoHasWonTheGame":  nil,
        "wasKingChanged":    false,
    }
    return c.JSON(http.StatusOK, response)
}

func GetGameId(c echo.Context) error {
    var gameId string
    //TODO fix it
    username := c.QueryParam("username")

    ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
    defer cancel()

    sqliteClient := sqlite.GetNewConnection()

    if gid, ok := sqlite.DoesPlayerHaveAnActiveGame(sqliteClient, username); ok {
        return c.JSON(http.StatusForbidden, map[string]interface{}{
            "message": trans.Get("You have already an active game."),
            "gameId":  *gid,
        })
    }

    redisClient := redis.GetNewConnection()

    go hokm.Matchmaking(ctx, redisClient, username)

    err := redisClient.Receive(context.Background(), redisClient.B().Subscribe().Channel("waiting").Build(), func(msg rueidis.PubSubMessage) {
        message := strings.Split(msg.Message, ",")
        players := message[:len(message)-1]
        if myslice.Has(players, username) {
            _, _ = sqlite.AddPlayerToGame(sqliteClient, username, gameId)
            gameId = message[len(message)-1]
            unsubscribeErr := redisClient.Do(ctx, redisClient.B().Unsubscribe().Channel("waiting").Build()).Error()
            if unsubscribeErr != nil {
                log.Println("Error while unsubscribing:", unsubscribeErr)
            }
            cancel()
        }
    })

    if err != nil {
        log.Println("Error subscribing to channel:", err)
        if errors.Is(ctx.Err(), context.DeadlineExceeded) {
            return c.JSON(http.StatusRequestTimeout, map[string]interface{}{
                "message": "فردی پیدا نشد، بعدا تلاش کنید.",
                "gameId":  nil,
            })
        }
        return c.JSON(http.StatusInternalServerError, map[string]interface{}{
            "message": "مشکلی پیش آمده است، بعدا تلاش کنید.",
            "gameId":  nil,
        })
    }

    return c.JSON(http.StatusCreated, map[string]interface{}{
        "message": trans.Get("Game has been made."),
        "gameId":  gameId,
    })
}
