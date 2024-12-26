package handler

import (
    "context"
    "errors"
    "fmt"
    "github.com/alirezadp10/hokm/internal/database/redis"
    "github.com/alirezadp10/hokm/internal/database/sqlite"
    "github.com/alirezadp10/hokm/internal/transformer"
    "github.com/alirezadp10/hokm/internal/util/my_bool"
    "github.com/alirezadp10/hokm/internal/util/my_slice"
    "github.com/alirezadp10/hokm/internal/util/trans"
    "github.com/alirezadp10/hokm/internal/validator"
    "github.com/google/uuid"
    "github.com/labstack/echo/v4"
    "github.com/redis/rueidis"
    "net/http"
    "strconv"
    "strings"
    "time"
)

func (h *Handler) CreateGame(c echo.Context) error {
    username := c.Get("username").(string)

    if err := validator.CreateGameValidator(h.GameService, validator.CreateGameValidatorData{
        Username: username,
    }); err != nil {
        return c.JSON(err.StatusCode, map[string]interface{}{"message": err.Message})
    }

    gameID := uuid.New().String()
    distributedCards := h.CardsService.DistributeCards()
    kingCards, king := h.PlayersService.ChooseFirstKing()
    go h.GameService.Matchmaking(c.Request().Context(), h.redis, username, gameID, distributedCards, kingCards, king)

    err := redis.Subscribe(c.Request().Context(), h.redis, "game_creation", func(msg rueidis.PubSubMessage) {
        message := strings.Split(msg.Message, "|")
        players := strings.Split(message[0], ",")
        if my_slice.Has(players, username) {
            gameID = message[1]
            _, err := sqlite.AddPlayerToGame(h.sqlite, username, gameID)
            if err != nil {
                fmt.Println(err)
            }
            redis.Unsubscribe(c.Request().Context(), h.redis, "game_creation")
        }
    })

    if err != nil {
        if errors.Is(err, context.Canceled) {
            redis.RemovePlayerList(context.Background(), h.redis, "matchmaking", username)
        }
        if errors.Is(c.Request().Context().Err(), context.DeadlineExceeded) {
            return c.JSON(http.StatusRequestTimeout, map[string]interface{}{
                "message": trans.Get("No body have found. please try again later."),
                "gameID":  nil,
            })
        }
        return c.JSON(http.StatusInternalServerError, map[string]interface{}{
            "message": trans.Get("Something went wrong. please try again later."),
            "gameID":  nil,
        })
    }

    return c.JSON(http.StatusOK, map[string]interface{}{
        "message": trans.Get("Game has been made."),
        "gameID":  gameID,
    })
}

func (h *Handler) GetGameInformation(c echo.Context) error {
    username := c.Get("username").(string)
    gameID := c.Param("gameID")

    if err := validator.GetGameInformationValidator(h.GameService, validator.GetGameInformationValidatorData{
        Username: username,
        GameID:   gameID,
    }); err != nil {
        return c.JSON(err.StatusCode, map[string]interface{}{"message": err.Message})
    }

    gameInformation, err := h.GameService.GameRepo.GetGameInformation(c.Request().Context(), gameID)

    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": trans.Get("Something went wrong. Please try again later.")})
    }

    players := strings.Split(gameInformation["players"].(string), ",")

    uIndex := my_slice.GetIndex(username, players)

    return c.JSON(http.StatusOK, transformer.GameInformationTransformer(h, transformer.GameInformationTransformerData{
        GameInformation: gameInformation,
        Players:         players,
        UIndex:          uIndex,
    }))
}

func (h *Handler) ChooseTrump(c echo.Context) error {
    var requestBody struct {
        Trump string `json:"trump"`
    }

    username := c.Get("username").(string)
    gameID := c.Param("gameID")

    if err := c.Bind(&requestBody); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": trans.Get("Invalid JSON.")})
    }

    gameInformation, err := h.GameService.GameRepo.GetGameInformation(c.Request().Context(), gameID)

    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": trans.Get("Something went wrong. Please try again later.")})
    }

    players := strings.Split(gameInformation["players"].(string), ",")

    uIndex := my_slice.GetIndex(username, players)

    if err := validator.ChooseTrumpValidator(h.GameService, validator.ChooseTrumpValidatorData{
        GameInformation: gameInformation,
        UIndex:          uIndex,
        Trump:           requestBody.Trump,
        GameID:          gameID,
        Username:        username,
    }); err != nil {
        return c.JSON(err.StatusCode, map[string]interface{}{"message": err.Message})
    }

    lastMoveTimestamp := strconv.FormatInt(time.Now().Unix(), 10)

    err = redis.SetTrump(c.Request().Context(), h.redis, gameID, requestBody.Trump, strconv.Itoa(uIndex), lastMoveTimestamp)

    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]interface{}{
            "message": trans.Get("Something went wrong, Please try again later."),
        })
    }

    return c.JSON(http.StatusOK, map[string]interface{}{
        "trump":        requestBody.Trump,
        "cards":        h.CardsService.GetPlayerCards(gameInformation["cards"].(string), uIndex)[1:],
        "timeRemained": h.PlayersService.GetTimeRemained(gameInformation["last_move_timestamp"].(string)),
    })
}

func (h *Handler) GetCards(c echo.Context) error {
    username := c.Get("username").(string)
    gameID := c.Param("gameID")

    if err := validator.GetCardsValidator(h.GameService, validator.GetCardsValidatorData{
        Username: username,
        GameID:   gameID,
    }); err != nil {
        return c.JSON(err.StatusCode, map[string]interface{}{"message": err.Message})
    }

    gameInformation, err := h.GameService.GameRepo.GetGameInformation(c.Request().Context(), gameID)

    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": trans.Get("Something went wrong. Please try again later.")})
    }

    trump := gameInformation["trump"].(string)

    if gameInformation["trump"].(string) == "" {
        err := redis.Subscribe(c.Request().Context(), h.redis, "choosing_trump", func(msg rueidis.PubSubMessage) {
            messages := strings.Split(msg.Message, ",")
            messageId := my_slice.HasLike(messages, func(s string) bool {
                return strings.Contains(s, gameID+"|")
            })
            if messageId != -1 {
                data := strings.Split(messages[messageId], "|")
                trump = data[1]
                redis.Unsubscribe(c.Request().Context(), h.redis, "choosing_trump")
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
        "cards": h.CardsService.GetPlayerCards(gameInformation["cards"].(string), uIndex),
        "trump": trump,
    })
}

func (h *Handler) PlaceCard(c echo.Context) error {
    username := c.Get("username").(string)
    gameID := c.Param("gameID")

    var requestBody struct {
        Card string `json:"card"`
    }

    if err := c.Bind(&requestBody); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"message": trans.Get("Invalid JSON.")})
    }

    gameInformation, err := h.GameService.GameRepo.GetGameInformation(c.Request().Context(), gameID)

    players := strings.Split(gameInformation["players"].(string), ",")

    uIndex := my_slice.GetIndex(username, players)

    gameState := h.initializeGameState(gameInformation)

    leadSuit := h.determineLeadSuit(requestBody.Card, gameState["leadSuit"].(string))

    if err := validator.PlaceCardValidator(h.GameService, h.CardsService, validator.PlaceCardValidatorData{
        GameInformation: gameInformation,
        Username:        username,
        GameID:          gameID,
        UIndex:          uIndex,
        LeadSuit:        leadSuit,
    }); err != nil {
        return c.JSON(err.StatusCode, map[string]interface{}{"message": err.Message})
    }

    centerCards := h.CardsService.UpdateCenterCards(gameInformation["center_cards"].(string), requestBody.Card, uIndex)

    cardsWinner := h.PointsService.FindCardsWinner(centerCards, gameState["trump"].(string), leadSuit)

    if cardsWinner != "" {
        h.updateWinnersAndPoints(&gameState, cardsWinner)
        if gameState["roundWinner"].(string) != "" {
            h.startNewRound(&gameState)
        }
    }

    gameState["lastMoveTimestamp"] = strconv.FormatInt(time.Now().Unix(), 10)
    gameState["cards"] = h.CardsService.UpdateUserCards(gameInformation["cards"].(string), requestBody.Card, uIndex)
    gameState["turn"] = h.PlayersService.GetNewTurn(gameInformation["turn"].(string))

    params := redis.PlaceCardParams{
        GameId:            gameID,
        Card:              requestBody.Card,
        CenterCards:       centerCards,
        LeadSuit:          leadSuit,
        CardsWinner:       cardsWinner,
        Points:            gameState["points"].(string),
        Turn:              gameState["turn"].(string),
        King:              gameState["king"].(string),
        WasKingChanged:    my_bool.ToString(gameState["wasKingChanged"].(bool)),
        LastMoveTimestamp: gameState["lastMoveTimestamp"].(string),
        Trump:             gameState["trump"].(string),
        IsItNewRound:      my_bool.ToString(gameState["isItNewRound"].(bool)),
        Cards:             gameState["cards"].([]string),
        PlayerIndex:       uIndex,
    }

    if err = redis.PlaceCard(c.Request().Context(), h.redis, params); err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]interface{}{
            "message": trans.Get("Something went wrong, Please try again later."),
        })
    }

    return c.JSON(http.StatusOK, transformer.PlaceCardTransformer(h, transformer.PlaceCardTransformerData{
        GameInformation:   gameInformation,
        Players:           players,
        UIndex:            uIndex,
        Points:            gameState["points"].(string),
        CenterCards:       centerCards,
        Turn:              gameState["turn"].(string),
        King:              gameState["king"].(string),
        LastMoveTimestamp: gameState["lastMoveTimestamp"].(string),
        WasKingChanged:    gameState["wasKingChanged"].(bool),
        CardsWinner:       cardsWinner,
        RoundWinner:       gameState["roundWinner"].(string),
        GameWinner:        gameState["gameWinner"].(string),
    }))
}

func (h *Handler) initializeGameState(gameInformation map[string]interface{}) map[string]interface{} {
    return map[string]interface{}{
        "trump":          gameInformation["trump"].(string),
        "king":           gameInformation["king"].(string),
        "points":         gameInformation["points"].(string),
        "leadSuit":       gameInformation["lead_suit"].(string),
        "gameWinner":     "",
        "roundWinner":    "",
        "wasKingChanged": false,
        "isItNewRound":   false,
    }
}

func (h *Handler) determineLeadSuit(card string, currentLeadSuit string) string {
    if currentLeadSuit == "" {
        return h.CardsService.GetCardSuit(card)
    }
    return currentLeadSuit
}

func (h *Handler) updateWinnersAndPoints(gameState *map[string]interface{}, cardsWinner string) {
    points, roundWinner, gameWinner := h.PointsService.UpdatePoints((*gameState)["points"].(string), cardsWinner)
    (*gameState)["points"] = points
    (*gameState)["roundWinner"] = roundWinner
    (*gameState)["gameWinner"] = gameWinner

    if roundWinner == "" {
        (*gameState)["turn"] = cardsWinner
    } else {
        (*gameState)["turn"] = ""
    }

    (*gameState)["king"] = h.PlayersService.GiveKing(roundWinner, (*gameState)["king"].(string))
    (*gameState)["wasKingChanged"] = (*gameState)["king"] == (*gameState)["king"].(string)
    (*gameState)["leadSuit"] = ""
    (*gameState)["centerCards"] = ",,,"
}

func (h *Handler) startNewRound(gameState *map[string]interface{}) {
    (*gameState)["cards"] = h.CardsService.DistributeCards()
    (*gameState)["isItNewRound"] = true
    if (*gameState)["wasKingChanged"].(bool) {
        (*gameState)["trump"] = ""
    }
}

func (h *Handler) GetUpdate(c echo.Context) error {
    username := c.Get("username").(string)
    gameID := c.Param("gameID")

    if err := validator.GetUpdateValidator(h.GameService, validator.GetUpdateValidatorData{
        Username: username,
        GameID:   gameID,
    }); err != nil {
        return c.JSON(err.StatusCode, map[string]interface{}{"message": err.Message})
    }

    gameInformation, err := h.GameService.GameRepo.GetGameInformation(c.Request().Context(), gameID)

    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": trans.Get("Something went wrong. Please try again later.")})
    }

    players := strings.Split(gameInformation["players"].(string), ",")

    uIndex := my_slice.GetIndex(username, players)

    var player int
    var card string

    err = redis.Subscribe(c.Request().Context(), h.redis, "placing_card", func(msg rueidis.PubSubMessage) {
        messages := strings.Split(msg.Message, "|")
        if messages[0] == gameID {
            gameInformation, _ = h.GameService.GameRepo.GetGameInformation(c.Request().Context(), messages[0])
            player, _ = strconv.Atoi(messages[1])
            card = messages[2]
            redis.Unsubscribe(c.Request().Context(), h.redis, "placing_card")
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

    return c.JSON(http.StatusOK, transformer.GetUpdateTransformer(h, transformer.GetUpdateTransformerData{
        GameInformation: gameInformation,
        UIndex:          uIndex,
        PlayerIndex:     player,
        Card:            card,
    }))
}
