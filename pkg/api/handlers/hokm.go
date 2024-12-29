package handlers

import (
    "context"
    "encoding/json"
    "errors"
    "fmt"
    "github.com/alirezadp10/hokm/internal/util/crypto"
    "github.com/alirezadp10/hokm/internal/util/my_bool"
    "github.com/alirezadp10/hokm/internal/util/my_slice"
    "github.com/alirezadp10/hokm/internal/util/trans"
    "github.com/alirezadp10/hokm/pkg/api/request"
    "github.com/alirezadp10/hokm/pkg/api/transformer"
    "github.com/alirezadp10/hokm/pkg/api/validator"
    "github.com/alirezadp10/hokm/pkg/repository"
    "github.com/alirezadp10/hokm/pkg/service"
    "github.com/google/uuid"
    "github.com/labstack/echo/v4"
    "github.com/redis/rueidis"
    "net/http"
    "strconv"
    "strings"
    "time"
)

type HokmHandler struct {
    GameService    service.GameService
    CardsService   service.CardsService
    PlayersService service.PlayersService
    PointsService  service.PointsService
    RedisService   service.RedisService
}

func NewHokmHandler(gameService *service.GameService, cardsService *service.CardsService, pointsService *service.PointsService, playersService *service.PlayersService, redisService *service.RedisService) *HokmHandler {
    return &HokmHandler{
        GameService:    *gameService,
        CardsService:   *cardsService,
        PointsService:  *pointsService,
        PlayersService: *playersService,
        RedisService:   *redisService,
    }
}

func (h *HokmHandler) CreateGame(c echo.Context) error {
    username := c.Get("username").(string)

    if err := validator.CreateGameValidator(h.GameService, validator.CreateGameValidatorData{
        Username: username,
    }); err != nil {
        return c.JSON(err.StatusCode, map[string]interface{}{"message": err.Message, "details": err.Details})
    }

    gameID := uuid.New().String()
    distributedCards := h.CardsService.DistributeCards()
    kingCards, king := h.PlayersService.ChooseFirstKing()

    go h.GameService.Matchmaking(c.Request().Context(), username, gameID, distributedCards, kingCards, king)

    err := h.RedisService.Subscribe(c.Request().Context(), "game_creation", func(msg rueidis.PubSubMessage) {
        message := strings.Split(msg.Message, "|")
        players := strings.Split(message[0], ",")
        if my_slice.Has(players, username) {
            gameID = message[1]
            _, err := h.PlayersService.PlayersRepo.AddPlayerToGame(username, gameID)
            if err != nil {
                fmt.Println(err)
            }
            h.RedisService.Unsubscribe(c.Request().Context(), "game_creation")
        }
    })

    if err != nil {
        if errors.Is(err, context.Canceled) {
            h.GameService.RemovePlayerFromWaitingList(username)
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

func (h *HokmHandler) GetGameInformation(c echo.Context) error {
    username := c.Get("username").(string)
    gameID := c.Param("gameID")

    if err := validator.GetGameInformationValidator(h.GameService, validator.GetGameInformationValidatorData{
        Username: username,
        GameID:   gameID,
    }); err != nil {
        return c.JSON(err.StatusCode, map[string]interface{}{"message": err.Message, "details": err.Details})
    }

    gameInformation, err := h.GameService.GameRepo.GetGameInformation(c.Request().Context(), gameID)

    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": trans.Get("Something went wrong. Please try again later.")})
    }

    players := strings.Split(gameInformation["players"].(string), ",")

    uIndex := my_slice.GetIndex(username, players)

    return c.JSON(http.StatusOK, transformer.GameInformationTransformer(h.PlayersService, h.PointsService, h.CardsService, transformer.GameInformationTransformerData{
        GameInformation: gameInformation,
        Players:         players,
        UIndex:          uIndex,
    }))
}

func (h *HokmHandler) ChooseTrump(c echo.Context) error {
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
        return c.JSON(err.StatusCode, map[string]interface{}{"message": err.Message, "details": err.Details})
    }

    lastMoveTimestamp := strconv.FormatInt(time.Now().Unix(), 10)

    err = h.CardsService.SetTrump(c.Request().Context(), gameID, requestBody.Trump, strconv.Itoa(uIndex), lastMoveTimestamp)

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

func (h *HokmHandler) GetCards(c echo.Context) error {
    username := c.Get("username").(string)
    gameID := c.Param("gameID")

    if err := validator.GetCardsValidator(h.GameService, validator.GetCardsValidatorData{
        Username: username,
        GameID:   gameID,
    }); err != nil {
        return c.JSON(err.StatusCode, map[string]interface{}{"message": err.Message, "details": err.Details})
    }

    gameInformation, err := h.GameService.GameRepo.GetGameInformation(c.Request().Context(), gameID)

    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": trans.Get("Something went wrong. Please try again later.")})
    }

    trump := gameInformation["trump"].(string)

    if gameInformation["trump"].(string) == "" {
        err := h.RedisService.Subscribe(c.Request().Context(), "choosing_trump", func(msg rueidis.PubSubMessage) {
            messages := strings.Split(msg.Message, ",")
            messageId := my_slice.HasLike(messages, func(s string) bool {
                return strings.Contains(s, gameID+"|")
            })
            if messageId != -1 {
                data := strings.Split(messages[messageId], "|")
                trump = data[1]
                h.RedisService.Unsubscribe(c.Request().Context(), "choosing_trump")
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

func (h *HokmHandler) PlaceCard(c echo.Context) error {
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
        return c.JSON(err.StatusCode, map[string]interface{}{"message": err.Message, "details": err.Details})
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

    params := repository.PlaceCardParams{
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

    if err = h.CardsService.CardsRepo.PlaceCard(c.Request().Context(), params); err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]interface{}{
            "message": trans.Get("Something went wrong, Please try again later."),
        })
    }

    return c.JSON(http.StatusOK, transformer.PlaceCardTransformer(h.PlayersService, h.PointsService, h.CardsService, transformer.PlaceCardTransformerData{
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

func (h *HokmHandler) initializeGameState(gameInformation map[string]interface{}) map[string]interface{} {
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

func (h *HokmHandler) determineLeadSuit(card string, currentLeadSuit string) string {
    if currentLeadSuit == "" {
        return h.CardsService.GetCardSuit(card)
    }
    return currentLeadSuit
}

func (h *HokmHandler) updateWinnersAndPoints(gameState *map[string]interface{}, cardsWinner string) {
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

func (h *HokmHandler) startNewRound(gameState *map[string]interface{}) {
    (*gameState)["cards"] = h.CardsService.DistributeCards()
    (*gameState)["isItNewRound"] = true
    if (*gameState)["wasKingChanged"].(bool) {
        (*gameState)["trump"] = ""
    }
}

func (h *HokmHandler) GetUpdate(c echo.Context) error {
    username := c.Get("username").(string)
    gameID := c.Param("gameID")

    if err := validator.GetUpdateValidator(h.GameService, validator.GetUpdateValidatorData{
        Username: username,
        GameID:   gameID,
    }); err != nil {
        return c.JSON(err.StatusCode, map[string]interface{}{"message": err.Message, "details": err.Details})
    }

    gameInformation, err := h.GameService.GameRepo.GetGameInformation(c.Request().Context(), gameID)

    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": trans.Get("Something went wrong. Please try again later.")})
    }

    players := strings.Split(gameInformation["players"].(string), ",")

    uIndex := my_slice.GetIndex(username, players)

    var player int
    var card string

    err = h.RedisService.Subscribe(c.Request().Context(), "placing_card", func(msg rueidis.PubSubMessage) {
        messages := strings.Split(msg.Message, "|")
        if messages[0] == gameID {
            gameInformation, _ = h.GameService.GameRepo.GetGameInformation(c.Request().Context(), messages[0])
            player, _ = strconv.Atoi(messages[1])
            card = messages[2]
            h.RedisService.Unsubscribe(c.Request().Context(), "placing_card")
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

    return c.JSON(http.StatusOK, transformer.GetUpdateTransformer(h.PointsService, h.PlayersService, transformer.GetUpdateTransformerData{
        GameInformation: gameInformation,
        UIndex:          uIndex,
        PlayerIndex:     player,
        Card:            card,
    }))
}

func (h *HokmHandler) GetSplashPage(c echo.Context) error {
    return c.Render(200, "splash.html", nil)
}

func (h *HokmHandler) GetMenuPage(c echo.Context) error {
    var user request.User

    err := json.Unmarshal([]byte(c.QueryParam("user")), &user)

    if err != nil {
        fmt.Println("Error unmarshalling user JSON:", err)
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "unauthorized"})
    }

    encryptedUsername, err := crypto.Encrypt(user.Username)
    if err != nil {
        fmt.Println("Error in encryption:", err)
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "unauthorized"})
    }

    chatInstance, _ := strconv.ParseInt(c.QueryParam("chat_instance"), 10, 64)

    _, _ = h.PlayersService.PlayersRepo.SavePlayer(user, chatInstance)

    return c.Render(200, "menu.html", map[string]interface{}{
        "userReferenceKey": encryptedUsername,
    })
}

func (h *HokmHandler) GetGamePage(c echo.Context) error {
    return c.Render(200, "game.html", map[string]interface{}{
        "userReferenceKey": c.QueryParam("user_id"),
        "gameID":           c.QueryParam("game_id"),
    })
}
