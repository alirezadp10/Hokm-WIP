package handler

import (
    "context"
    "errors"
    "github.com/alirezadp10/hokm/internal/database/redis"
    "github.com/alirezadp10/hokm/internal/database/sqlite"
    "github.com/alirezadp10/hokm/internal/hokm"
    "github.com/alirezadp10/hokm/internal/utils/my_slice"
    "github.com/alirezadp10/hokm/internal/utils/trans"
    "github.com/labstack/echo/v4"
    "github.com/redis/rueidis"
    "net/http"
    "strconv"
    "strings"
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

    go hokm.Matchmaking(h.context, h.redisConnection, username)

    err := redis.Subscribe(h.context, h.redisConnection, "game_creation", func(msg rueidis.PubSubMessage) {
        message := strings.Split(msg.Message, ",")
        players := message[:len(message)-1]
        if my_slice.Has(players, username) {
            _, _ = sqlite.AddPlayerToGame(h.sqliteConnection, username, gameId)
            gameId = message[len(message)-1]
            redis.Unsubscribe(h.context, h.redisConnection, "game_creation")
            //cancel()
        }
    })

    if err != nil {
        if errors.Is(h.context.Err(), context.DeadlineExceeded) {
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

    gameId := c.QueryParam("gameId")

    if !sqlite.DoesPlayerBelongsToThisGame(h.sqliteConnection, username, gameId) {
        return c.JSON(http.StatusForbidden, map[string]interface{}{
            "message": trans.Get("It's not your game."),
        })
    }

    gameInformation := redis.GetGameInformation(h.context, h.redisConnection, gameId)

    players := gameInformation["players"].([]string)

    uIndex := my_slice.GetIndex(username, players)

    return c.JSON(http.StatusOK, map[string]interface{}{
        "players":      hokm.GetPlayersWithDirections(players, uIndex),
        "points":       hokm.GetPoints(gameInformation["points"].(map[string]interface{}), uIndex),
        "centerCards":  hokm.GetCenterCards(gameInformation["center_cards"].(map[int]string), uIndex),
        "turn":         hokm.GetDirection(gameInformation["turn"].(int), uIndex),
        "judge":        hokm.GetDirection(gameInformation["judge"].(int), uIndex),
        "timeRemained": hokm.GetTimeRemained(gameInformation["last_move_timestamp"].(string)),
        "kingsCards":   hokm.GetKingsCards(gameInformation["kings_cards"].([]string), uIndex),
        "trump":        gameInformation["trump"].(string),
        "yourCards":    gameInformation["cards"].(map[int]interface{})[uIndex],
    })
}

func (h *Handler) ChooseTrump(c echo.Context) error {
    username := c.Get("username").(string)

    var requestBody struct {
        GameId string `json:"gameId"`
        Trump  string `json:"trump"`
    }

    if err := c.Bind(&requestBody); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
    }

    if !sqlite.DoesPlayerBelongsToThisGame(h.sqliteConnection, username, requestBody.GameId) {
        return c.JSON(http.StatusUnauthorized, map[string]interface{}{
            "message": trans.Get("It's not your game."),
        })
    }

    gameInformation := redis.GetGameInformation(h.context, h.redisConnection, requestBody.GameId)

    if gameInformation["judge"].(string) != username {
        return c.JSON(http.StatusForbidden, map[string]interface{}{
            "message": trans.Get("You're not judge in this round."),
        })
    }

    if gameInformation["trump"] != nil {
        return c.JSON(http.StatusForbidden, map[string]interface{}{
            "message": trans.Get("You're not allowed to choose a trump at the moment."),
        })
    }

    err := redis.SetTrump(h.context, h.redisConnection, requestBody.GameId, requestBody.Trump)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]interface{}{
            "message": trans.Get("Something went wrong, Please try again later."),
        })
    }

    players := gameInformation["players"].([]string)

    uIndex := my_slice.GetIndex(username, players)

    return c.JSON(http.StatusOK, map[string]interface{}{
        "trump": requestBody.Trump,
        "cards": hokm.GetPlayerCards(gameInformation["cards"].(map[int][]string), uIndex),
    })
}

func (h *Handler) GetYourCards(c echo.Context) error {
    username := c.Get("username").(string)

    var trump string
    var gameInformation map[string]interface{}
    gameId := c.QueryParam("gameId")

    if !sqlite.DoesPlayerBelongsToThisGame(h.sqliteConnection, username, gameId) {
        return c.JSON(http.StatusUnauthorized, map[string]interface{}{
            "message": trans.Get("It's not your game."),
        })
    }

    err := redis.Subscribe(h.context, h.redisConnection, "choosing_trump", func(msg rueidis.PubSubMessage) {
        messages := strings.Split(msg.Message, ",")
        messageId := my_slice.HasLike(messages, func(s string) bool {
            return strings.Contains(s, gameId+"|")
        })
        if messageId != -1 {
            data := strings.Split(messages[messageId], "|")
            gameInformation = redis.GetGameInformation(h.context, h.redisConnection, data[0])
            trump = data[1]
            redis.Unsubscribe(h.context, h.redisConnection, "choosing_trump")
            //cancel()
        }
    })

    if err != nil {
        if errors.Is(h.context.Err(), context.DeadlineExceeded) {
            return c.JSON(http.StatusRequestTimeout, map[string]interface{}{
                "message": trans.Get("Something went wrong, Please try again later."),
            })
        }
        return c.JSON(http.StatusInternalServerError, map[string]interface{}{
            "message": trans.Get("Something went wrong, Please try again later."),
        })
    }

    players := gameInformation["players"].([]string)

    uIndex := my_slice.GetIndex(username, players)

    return c.JSON(http.StatusOK, map[string]interface{}{
        "cards": hokm.GetPlayerCards(gameInformation["cards"].(map[int][]string), uIndex),
        "trump": trump,
    })
}

func (h *Handler) PlaceCard(c echo.Context) error {
    username := c.Get("username").(string)

    var requestBody struct {
        GameId string `json:"gameId"`
        Card   string `json:"card"`
    }

    if err := c.Bind(&requestBody); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
    }

    if !sqlite.DoesPlayerBelongsToThisGame(h.sqliteConnection, username, requestBody.GameId) {
        return c.JSON(http.StatusUnauthorized, map[string]interface{}{
            "message": trans.Get("It's not your game."),
        })
    }

    gameInformation := redis.GetGameInformation(h.context, h.redisConnection, requestBody.GameId)

    players := gameInformation["players"].([]string)

    uIndex := my_slice.GetIndex(username, players)

    err := redis.PlaceCard(h.context, h.redisConnection, uIndex, requestBody.GameId, requestBody.Card, gameInformation["center_cards"].(string))
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]interface{}{
            "message": trans.Get("Something went wrong, Please try again later."),
        })
    }

    return c.JSON(http.StatusOK, map[string]interface{}{
        "points":            hokm.GetPoints(gameInformation["points"].(map[string]interface{}), uIndex),
        "centerCards":       hokm.GetCenterCards(gameInformation["center_cards"].(map[int]string), uIndex),
        "kingsCards":        hokm.GetKingsCards(gameInformation["kings_cards"].([]string), uIndex),
        "turn":              hokm.GetDirection(gameInformation["turn"].(int), uIndex),
        "judge":             hokm.GetDirection(gameInformation["judge"].(int), uIndex),
        "whoHasWonTheCards": hokm.GetDirection(gameInformation["who_has_won_the_cards"].(int), uIndex),
        "whoHasWonTheRound": hokm.GetDirection(gameInformation["who_has_won_the_round"].(int), uIndex),
        "whoHasWonTheGame":  hokm.GetDirection(gameInformation["who_has_won_the_game"].(int), uIndex),
        "timeRemained":      hokm.GetTimeRemained(gameInformation["last_move_timestamp"].(string)),
        "wasKingChanged":    gameInformation["who_king_changed"].(bool),
    })
}

func (h *Handler) GetUpdate(c echo.Context) error {
    username := c.Get("username").(string)

    var player int
    var card string
    var gameInformation map[string]interface{}
    gameId := c.QueryParam("gameId")

    if !sqlite.DoesPlayerBelongsToThisGame(h.sqliteConnection, username, gameId) {
        return c.JSON(http.StatusUnauthorized, map[string]interface{}{
            "message": trans.Get("It's not your game."),
        })
    }

    err := redis.Subscribe(h.context, h.redisConnection, "placing_card", func(msg rueidis.PubSubMessage) {
        messages := strings.Split(msg.Message, ",")
        messageId := my_slice.HasLike(messages, func(s string) bool {
            return strings.Contains(s, gameId+"|")
        })
        if messageId != -1 {
            data := strings.Split(messages[messageId], "|")
            gameInformation = redis.GetGameInformation(h.context, h.redisConnection, data[0])
            player, _ = strconv.Atoi(data[1])
            card = data[2]
            redis.Unsubscribe(h.context, h.redisConnection, "placing_card")
            //cancel()
        }
    })

    if err != nil {
        if errors.Is(h.context.Err(), context.DeadlineExceeded) {
            return c.JSON(http.StatusRequestTimeout, map[string]interface{}{
                "message": trans.Get("Something went wrong, Please try again later."),
            })
        }
        return c.JSON(http.StatusInternalServerError, map[string]interface{}{
            "message": trans.Get("Something went wrong, Please try again later."),
        })
    }

    players := gameInformation["players"].([]string)

    uIndex := my_slice.GetIndex(username, players)

    return c.JSON(http.StatusOK, map[string]interface{}{
        "lastMove": map[string]string{
            "from": hokm.GetDirection(player, uIndex),
            "card": card,
        },
        "centerCards":       hokm.GetCenterCards(gameInformation["center_cards"].(map[int]string), uIndex),
        "points":            hokm.GetPoints(gameInformation["points"].(map[string]interface{}), uIndex),
        "turn":              hokm.GetDirection(gameInformation["turn"].(int), uIndex),
        "judge":             hokm.GetDirection(gameInformation["judge"].(int), uIndex),
        "whoHasWonTheCards": hokm.GetDirection(gameInformation["who_has_won_the_cards"].(int), uIndex),
        "whoHasWonTheRound": hokm.GetDirection(gameInformation["who_has_won_the_round"].(int), uIndex),
        "whoHasWonTheGame":  hokm.GetDirection(gameInformation["who_has_won_the_game"].(int), uIndex),
        "timeRemained":      hokm.GetTimeRemained(gameInformation["last_move_timestamp"].(string)),
        "wasKingChanged":    gameInformation["who_king_changed"].(bool),
    })
}
