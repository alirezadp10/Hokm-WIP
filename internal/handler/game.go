package handler

import (
    "context"
    "errors"
    "fmt"
    "github.com/alirezadp10/hokm/internal/database/redis"
    "github.com/alirezadp10/hokm/internal/database/sqlite"
    "github.com/alirezadp10/hokm/internal/hokm"
    "github.com/alirezadp10/hokm/internal/utils/my_slice"
    "github.com/alirezadp10/hokm/internal/utils/trans"
    "github.com/google/uuid"
    "github.com/labstack/echo/v4"
    "github.com/redis/rueidis"
    "net/http"
    "strconv"
    "strings"
    "time"
)

func (h *Handler) GetGameId(c echo.Context) error {
    username := c.Get("username").(string)

    var gameId string

    if gid, ok := sqlite.DoesPlayerHaveAnActiveGame(h.sqliteConnection, username); ok {
        return c.JSON(http.StatusForbidden, map[string]interface{}{
            "message": trans.Get("You have already an active game."),
            "gameId":  *gid,
        })
    }

    gameId = uuid.New().String()
    go hokm.Matchmaking(c.Request().Context(), h.redisConnection, username, gameId)

    err := redis.Subscribe(c.Request().Context(), h.redisConnection, "game_creation", func(msg rueidis.PubSubMessage) {
        message := strings.Split(msg.Message, "|")
        players := strings.Split(message[0], ",")
        if my_slice.Has(players, username) {
            gameId = message[1]
            _, err := sqlite.AddPlayerToGame(h.sqliteConnection, username, gameId)
            if err != nil {
                fmt.Println(err)
            }
            redis.Unsubscribe(c.Request().Context(), h.redisConnection, "game_creation")
            //cancel()
        }
    })

    if err != nil {
        if errors.Is(c.Request().Context().Err(), context.DeadlineExceeded) {
            return c.JSON(http.StatusRequestTimeout, map[string]interface{}{
                "message": trans.Get("No body have found. please try again later."),
                "gameId":  nil,
            })
        }
        return c.JSON(http.StatusInternalServerError, map[string]interface{}{
            "message": trans.Get("Something went wrong. please try again later."),
            "gameId":  nil,
        })
    }

    return c.JSON(http.StatusOK, map[string]interface{}{
        "message": trans.Get("Game has been made."),
        "gameId":  gameId,
    })
}

func (h *Handler) GetGameData(c echo.Context) error {
    username := c.Get("username").(string)

    gameId := c.Param("gameId")

    if sqlite.HasGameFinished(h.sqliteConnection, gameId) {
        return c.JSON(http.StatusForbidden, map[string]interface{}{
            "message": trans.Get("The game has already finished."),
        })
    }

    if !sqlite.DoesPlayerBelongsToThisGame(h.sqliteConnection, username, gameId) {
        return c.JSON(http.StatusForbidden, map[string]interface{}{
            "message": trans.Get("It's not your game."),
        })
    }

    gameInformation := redis.GetGameInformation(c.Request().Context(), h.redisConnection, gameId)

    players := strings.Split(gameInformation["players"].(string), ",")

    uIndex := my_slice.GetIndex(username, players)

    response := map[string]interface{}{
        "players":              hokm.GetPlayersWithDirections(players, uIndex),
        "points":               hokm.GetPoints(gameInformation["points"].(string), uIndex),
        "centerCards":          hokm.GetCenterCards(gameInformation["center_cards"].(string), uIndex),
        "turn":                 hokm.GetTurn(gameInformation["turn"].(string), uIndex),
        "king":                 hokm.GetKing(gameInformation["king"].(string), uIndex),
        "kingCards":            hokm.GetKingCards(gameInformation["king_cards"].(string)),
        "timeRemained":         hokm.GetTimeRemained(gameInformation["last_move_timestamp"].(string)),
        "hasKingCardsFinished": gameInformation["has_king_cards_finished"].(string),
        "trump":                gameInformation["trump"],
    }

    if response["hasKingCardsFinished"] == "false" {
        response["playerCards"] = hokm.GetPlayerCards(gameInformation["cards"].(string), uIndex)
    }

    return c.JSON(http.StatusOK, response)
}

func (h *Handler) ChooseTrump(c echo.Context) error {
    username := c.Get("username").(string)
    gameId := c.Param("gameId")

    var requestBody struct {
        Trump string `json:"trump"`
    }

    if err := c.Bind(&requestBody); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": trans.Get("Invalid JSON.")})
    }

    if !my_slice.Has([]string{"H", "D", "C", "S"}, requestBody.Trump) {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": trans.Get("Invalid trump.")})
    }

    if !sqlite.DoesPlayerBelongsToThisGame(h.sqliteConnection, username, gameId) {
        return c.JSON(http.StatusForbidden, map[string]interface{}{
            "message": trans.Get("It's not your game."),
        })
    }

    gameInformation := redis.GetGameInformation(c.Request().Context(), h.redisConnection, gameId)

    players := strings.Split(gameInformation["players"].(string), ",")

    uIndex := my_slice.GetIndex(username, players)

    if gameInformation["king"].(string) != strconv.Itoa(uIndex) {
        return c.JSON(http.StatusForbidden, map[string]interface{}{
            "message": trans.Get("You're not king in this round."),
        })
    }

    if gameInformation["has_king_cards_finished"].(string) == "true" {
        return c.JSON(http.StatusForbidden, map[string]interface{}{
            "message": trans.Get("You're not allowed to choose a trump at the moment."),
        })
    }

    lastMoveTimestamp := strconv.FormatInt(time.Now().Unix(), 10)

    err := redis.SetTrump(c.Request().Context(), h.redisConnection, gameId, requestBody.Trump, strconv.Itoa(uIndex), lastMoveTimestamp)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]interface{}{
            "message": trans.Get("Something went wrong, Please try again later."),
        })
    }

    return c.JSON(http.StatusOK, map[string]interface{}{
        "trump": requestBody.Trump,
        "cards": hokm.GetPlayerCards(gameInformation["cards"].(string), uIndex),
    })
}

func (h *Handler) GetYourCards(c echo.Context) error {
    username := c.Get("username").(string)
    gameId := c.Param("gameId")

    if !sqlite.DoesPlayerBelongsToThisGame(h.sqliteConnection, username, gameId) {
        return c.JSON(http.StatusForbidden, map[string]interface{}{
            "message": trans.Get("It's not your game."),
        })
    }

    gameInformation := redis.GetGameInformation(c.Request().Context(), h.redisConnection, gameId)

    trump := gameInformation["trump"].(string)

    if gameInformation["trump"].(string) == "" {
        err := redis.Subscribe(c.Request().Context(), h.redisConnection, "choosing_trump", func(msg rueidis.PubSubMessage) {
            messages := strings.Split(msg.Message, ",")
            messageId := my_slice.HasLike(messages, func(s string) bool {
                return strings.Contains(s, gameId+"|")
            })
            if messageId != -1 {
                data := strings.Split(messages[messageId], "|")
                trump = data[1]
                redis.Unsubscribe(c.Request().Context(), h.redisConnection, "choosing_trump")
                //cancel()
            }
        })

        if err != nil {
            if errors.Is(c.Request().Context().Err(), context.DeadlineExceeded) {
                return c.JSON(http.StatusRequestTimeout, map[string]interface{}{
                    "message": trans.Get("Something went wrong, Please try again later."),
                })
            }
            return c.JSON(http.StatusInternalServerError, map[string]interface{}{
                "message": trans.Get("Something went wrong, Please try again later."),
            })
        }
    }

    players := strings.Split(gameInformation["players"].(string), ",")

    uIndex := my_slice.GetIndex(username, players)

    return c.JSON(http.StatusOK, map[string]interface{}{
        "cards": hokm.GetPlayerCards(gameInformation["cards"].(string), uIndex),
        "trump": trump,
    })
}

func (h *Handler) PlaceCard(c echo.Context) error {
    username := c.Get("username").(string)
    gameId := c.Param("gameId")

    var requestBody struct {
        Card string `json:"card"`
    }

    if err := c.Bind(&requestBody); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
    }

    if !sqlite.DoesPlayerBelongsToThisGame(h.sqliteConnection, username, gameId) {
        return c.JSON(http.StatusForbidden, map[string]interface{}{
            "message": trans.Get("It's not your game."),
        })
    }

    gameInformation := redis.GetGameInformation(c.Request().Context(), h.redisConnection, gameId)

    players := strings.Split(gameInformation["players"].(string), ",")

    uIndex := my_slice.GetIndex(username, players)

    if gameInformation["turn"].(string) != strconv.Itoa(uIndex) {
        return c.JSON(http.StatusForbidden, map[string]interface{}{
            "message": trans.Get("It's not your turn."),
        })
    }

    if gameInformation["trump"].(string) != string(requestBody.Card[len(requestBody.Card)-1]) {
        for _, card := range hokm.GetPlayerCards(gameInformation["cards"].(string), uIndex) {
            if gameInformation["trump"].(string) == card[len(card)-1] {
                return c.JSON(http.StatusForbidden, map[string]interface{}{
                    "message": trans.Get("You're not allowed to select this card."),
                })
            }
        }
    }

    centerCardsList := strings.Split(gameInformation["center_cards"].(string), ",")
    centerCardsList[uIndex] = requestBody.Card
    centerCards := strings.Join(centerCardsList, ",")
    err := redis.PlaceCard(c.Request().Context(), h.redisConnection, uIndex, gameId, requestBody.Card, centerCards)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]interface{}{
            "message": trans.Get("Something went wrong, Please try again later."),
        })
    }

    if len(strings.Split(gameInformation["center_card"].(string), ",")) == 4 {
        // winner
    }

    return c.JSON(http.StatusOK, map[string]interface{}{
        "players":      hokm.GetPlayersWithDirections(players, uIndex),
        "points":       hokm.GetPoints(gameInformation["points"].(string), uIndex),
        "centerCards":  hokm.GetCenterCards(gameInformation["center_cards"].(string), uIndex),
        "turn":         hokm.GetTurn(gameInformation["turn"].(string), uIndex),
        "king":         hokm.GetKing(gameInformation["king"].(string), uIndex),
        "timeRemained": hokm.GetTimeRemained(gameInformation["last_move_timestamp"].(string)),
        "playerCards":  hokm.GetPlayerCards(gameInformation["cards"].(string), uIndex),

        "whoHasWonTheCards": hokm.GetDirection(gameInformation["who_has_won_the_cards"].(int), uIndex),
        "whoHasWonTheRound": hokm.GetDirection(gameInformation["who_has_won_the_round"].(int), uIndex),
        "whoHasWonTheGame":  hokm.GetDirection(gameInformation["who_has_won_the_game"].(int), uIndex),
        "wasKingChanged":    gameInformation["was_king_changed"].(string),
        "trump":             gameInformation["trump"],
    })
}

func (h *Handler) GetUpdate(c echo.Context) error {
    username := c.Get("username").(string)

    var player int
    var card string
    var gameInformation map[string]interface{}
    gameId := c.QueryParam("gameId")

    if !sqlite.DoesPlayerBelongsToThisGame(h.sqliteConnection, username, gameId) {
        return c.JSON(http.StatusForbidden, map[string]interface{}{
            "message": trans.Get("It's not your game."),
        })
    }

    err := redis.Subscribe(c.Request().Context(), h.redisConnection, "placing_card", func(msg rueidis.PubSubMessage) {
        messages := strings.Split(msg.Message, ",")
        messageId := my_slice.HasLike(messages, func(s string) bool {
            return strings.Contains(s, gameId+"|")
        })
        if messageId != -1 {
            data := strings.Split(messages[messageId], "|")
            gameInformation = redis.GetGameInformation(c.Request().Context(), h.redisConnection, data[0])
            player, _ = strconv.Atoi(data[1])
            card = data[2]
            redis.Unsubscribe(c.Request().Context(), h.redisConnection, "placing_card")
            //cancel()
        }
    })

    if err != nil {
        if errors.Is(c.Request().Context().Err(), context.DeadlineExceeded) {
            return c.JSON(http.StatusRequestTimeout, map[string]interface{}{
                "message": trans.Get("Something went wrong, Please try again later."),
            })
        }
        return c.JSON(http.StatusInternalServerError, map[string]interface{}{
            "message": trans.Get("Something went wrong, Please try again later."),
        })
    }

    players := strings.Split(gameInformation["players"].(string), ",")

    uIndex := my_slice.GetIndex(username, players)

    return c.JSON(http.StatusOK, map[string]interface{}{
        "lastMove": map[string]string{
            "from": hokm.GetDirection(player, uIndex),
            "card": card,
        },
        "points":            hokm.GetPoints(gameInformation["points"].(string), uIndex),
        "centerCards":       hokm.GetCenterCards(gameInformation["center_cards"].(string), uIndex),
        "turn":              hokm.GetTurn(gameInformation["turn"].(string), uIndex),
        "king":              hokm.GetKing(gameInformation["king"].(string), uIndex),
        "timeRemained":      hokm.GetTimeRemained(gameInformation["last_move_timestamp"].(string)),
        "whoHasWonTheCards": hokm.GetDirection(gameInformation["who_has_won_the_cards"].(int), uIndex),
        "whoHasWonTheRound": hokm.GetDirection(gameInformation["who_has_won_the_round"].(int), uIndex),
        "whoHasWonTheGame":  hokm.GetDirection(gameInformation["who_has_won_the_game"].(int), uIndex),
        "wasKingChanged":    gameInformation["was_king_changed"].(string),
    })
}
