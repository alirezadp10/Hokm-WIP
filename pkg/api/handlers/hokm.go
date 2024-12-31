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
    "gorm.io/gorm"
    "log"
    "net/http"
    "strconv"
    "strings"
    "time"
)

type HokmHandler struct {
    sqlite         gorm.DB
    redis          rueidis.Client
    gameService    service.GameService
    cardsService   service.CardsService
    playersService service.PlayersService
    pointsService  service.PointsService
}

func NewHokmHandler(sqliteClient *gorm.DB, redisClient *rueidis.Client, gameService *service.GameService, cardsService *service.CardsService, pointsService *service.PointsService, playersService *service.PlayersService) *HokmHandler {
    return &HokmHandler{
        sqlite:         *sqliteClient,
        redis:          *redisClient,
        gameService:    *gameService,
        cardsService:   *cardsService,
        pointsService:  *pointsService,
        playersService: *playersService,
    }
}

func (h *HokmHandler) CreateGame(c echo.Context) error {
    username := c.Get("username").(string)

    if err := validator.CreateGameValidator(h.gameService, validator.CreateGameValidatorData{
        Username: username,
    }); err != nil {
        return c.JSON(err.StatusCode, map[string]interface{}{"message": err.Message, "details": err.Details})
    }

    gameID := uuid.New().String()
    distributedCards := h.cardsService.DistributeCards()
    kingCards, king := h.playersService.ChooseFirstKing()

    go h.gameService.Matchmaking(c.Request().Context(), username, gameID, distributedCards, kingCards, king)

    err := h.redis.Receive(c.Request().Context(), h.redis.B().Subscribe().Channel("game_creation").Build(), func(msg rueidis.PubSubMessage) {
        message := strings.Split(msg.Message, "|")
        players := strings.Split(message[0], ",")
        if my_slice.Has(players, username) {
            gameID = message[1]
            _, err := h.playersService.PlayersRepo.AddPlayerToGame(username, gameID)
            if err != nil {
                fmt.Println(err)
            }
            unsubscribeErr := h.redis.Do(c.Request().Context(), h.redis.B().Unsubscribe().Channel("game_creation").Build()).Error()
            if unsubscribeErr != nil {
                log.Println("Error while unsubscribing:", unsubscribeErr)
            }
        }
    })

    if err != nil {
        log.Printf("Error in subscribing to %v channel: %v", "game_creation", err)

        if errors.Is(err, context.Canceled) {
            h.gameService.RemovePlayerFromWaitingList(username)
        }
        if errors.Is(c.Request().Context().Err(), context.DeadlineExceeded) {
            return c.JSON(http.StatusRequestTimeout, map[string]interface{}{
                "message": trans.Get("No body have found. please try again later."),
                "gameID":  nil,
            })
        }
        return c.JSON(http.StatusInternalServerError, map[string]interface{}{
            "message": trans.Get("Something went wrong, Please try again later."),
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

    if err := validator.GetGameInformationValidator(h.gameService, validator.GetGameInformationValidatorData{
        Username: username,
        GameID:   gameID,
    }); err != nil {
        return c.JSON(err.StatusCode, map[string]interface{}{"message": err.Message, "details": err.Details})
    }

    gameInformation, err := h.gameService.GameRepo.GetGameInformation(c.Request().Context(), gameID)

    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": trans.Get("Something went wrong, Please try again later.")})
    }

    players := strings.Split(gameInformation["players"].(string), ",")

    uIndex := my_slice.GetIndex(username, players)

    return c.JSON(http.StatusOK, transformer.GameInformationTransformer(h.playersService, h.pointsService, h.cardsService, transformer.GameInformationTransformerData{
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

    gameInformation, err := h.gameService.GameRepo.GetGameInformation(c.Request().Context(), gameID)

    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": trans.Get("Something went wrong, Please try again later.")})
    }

    players := strings.Split(gameInformation["players"].(string), ",")

    uIndex := my_slice.GetIndex(username, players)

    if err := validator.ChooseTrumpValidator(h.gameService, validator.ChooseTrumpValidatorData{
        GameInformation: gameInformation,
        UIndex:          uIndex,
        Trump:           requestBody.Trump,
        GameID:          gameID,
        Username:        username,
    }); err != nil {
        return c.JSON(err.StatusCode, map[string]interface{}{"message": err.Message, "details": err.Details})
    }

    lastMoveTimestamp := strconv.FormatInt(time.Now().Unix(), 10)

    err = h.cardsService.SetTrump(c.Request().Context(), gameID, requestBody.Trump, strconv.Itoa(uIndex), lastMoveTimestamp)

    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]interface{}{
            "message": trans.Get("Something went wrong, Please try again later."),
        })
    }

    return c.JSON(http.StatusOK, map[string]interface{}{
        "trump":        requestBody.Trump,
        "cards":        h.cardsService.GetPlayerCards(gameInformation["cards"].(string), uIndex)[1:],
        "timeRemained": h.playersService.GetTimeRemained(gameInformation["last_move_timestamp"].(string)),
        "turn":         h.playersService.GetTurn(gameInformation["turn"].(string), uIndex),
    })
}

func (h *HokmHandler) GetCards(c echo.Context) error {
    username := c.Get("username").(string)
    gameID := c.Param("gameID")

    if err := validator.GetCardsValidator(h.gameService, validator.GetCardsValidatorData{
        Username: username,
        GameID:   gameID,
    }); err != nil {
        return c.JSON(err.StatusCode, map[string]interface{}{"message": err.Message, "details": err.Details})
    }

    gameInformation, err := h.gameService.GameRepo.GetGameInformation(c.Request().Context(), gameID)

    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": trans.Get("Something went wrong, Please try again later.")})
    }

    trump := gameInformation["trump"].(string)

    if gameInformation["trump"].(string) == "" {
        err := h.redis.Receive(c.Request().Context(), h.redis.B().Subscribe().Channel("choosing_trump").Build(), func(msg rueidis.PubSubMessage) {
            messages := strings.Split(msg.Message, ",")
            messageId := my_slice.HasLike(messages, func(s string) bool {
                return strings.Contains(s, gameID+"|")
            })
            if messageId != -1 {
                data := strings.Split(messages[messageId], "|")
                trump = data[1]
                unsubscribeErr := h.redis.Do(c.Request().Context(), h.redis.B().Unsubscribe().Channel("choosing_trump").Build()).Error()
                if unsubscribeErr != nil {
                    log.Println("Error while unsubscribing:", unsubscribeErr)
                }
            }
        })

        if err != nil {
            log.Printf("Error in subscribing to %v channel: %v", "choosing_trump", err)
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
        "cards": h.cardsService.GetPlayerCards(gameInformation["cards"].(string), uIndex),
        "turn":  h.playersService.GetTurn(gameInformation["turn"].(string), uIndex),
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

    gameInformation, err := h.gameService.GameRepo.GetGameInformation(c.Request().Context(), gameID)

    players := strings.Split(gameInformation["players"].(string), ",")

    uIndex := my_slice.GetIndex(username, players)

    gameState := h.initializeGameState(gameInformation)

    leadSuit := h.determineLeadSuit(requestBody.Card, gameState["leadSuit"].(string))

    if err := validator.PlaceCardValidator(h.gameService, h.cardsService, validator.PlaceCardValidatorData{
        GameInformation: gameInformation,
        Username:        username,
        GameID:          gameID,
        UIndex:          uIndex,
        LeadSuit:        leadSuit,
        Card:            requestBody.Card,
    }); err != nil {
        return c.JSON(err.StatusCode, map[string]interface{}{"message": err.Message, "details": err.Details})
    }

    centerCards := h.cardsService.UpdateCenterCards(gameInformation["center_cards"].(string), requestBody.Card, uIndex)

    cardsWinner := h.pointsService.FindCardsWinner(centerCards, gameState["trump"].(string), leadSuit)

    if cardsWinner != "" {
        h.updateWinnersAndPoints(&gameState, cardsWinner)
        if gameState["roundWinner"].(string) != "" {
            h.startNewRound(&gameState)
        }
    }

    gameState["lastMoveTimestamp"] = strconv.FormatInt(time.Now().Unix(), 10)
    gameState["cards"] = h.cardsService.UpdateUserCards(gameInformation["cards"].(string), requestBody.Card, uIndex)
    gameState["turn"] = h.playersService.GetNewTurn(gameInformation["turn"].(string), gameState["gameWinner"].(string))

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

    if err = h.cardsService.CardsRepo.PlaceCard(c.Request().Context(), params); err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]interface{}{
            "message": trans.Get("Something went wrong, Please try again later."),
        })
    }

    return c.JSON(http.StatusOK, transformer.PlaceCardTransformer(h.playersService, h.pointsService, h.cardsService, transformer.PlaceCardTransformerData{
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
        return h.cardsService.GetCardSuit(card)
    }
    return currentLeadSuit
}

func (h *HokmHandler) updateWinnersAndPoints(gameState *map[string]interface{}, cardsWinner string) {
    points, roundWinner, gameWinner := h.pointsService.UpdatePoints((*gameState)["points"].(string), cardsWinner)
    (*gameState)["points"] = points
    (*gameState)["roundWinner"] = roundWinner
    (*gameState)["gameWinner"] = gameWinner

    if roundWinner == "" {
        (*gameState)["turn"] = cardsWinner
    } else {
        (*gameState)["turn"] = ""
    }

    (*gameState)["king"] = h.playersService.GiveKing(roundWinner, (*gameState)["king"].(string))
    (*gameState)["wasKingChanged"] = (*gameState)["king"] == (*gameState)["king"].(string)
    (*gameState)["leadSuit"] = ""
    (*gameState)["centerCards"] = ",,,"
}

func (h *HokmHandler) startNewRound(gameState *map[string]interface{}) {
    (*gameState)["cards"] = h.cardsService.DistributeCards()
    (*gameState)["isItNewRound"] = true
    if (*gameState)["wasKingChanged"].(bool) {
        (*gameState)["trump"] = ""
    }
}

func (h *HokmHandler) GetUpdate(c echo.Context) error {
    username := c.Get("username").(string)
    gameID := c.Param("gameID")

    if err := validator.GetUpdateValidator(h.gameService, validator.GetUpdateValidatorData{
        Username: username,
        GameID:   gameID,
    }); err != nil {
        return c.JSON(err.StatusCode, map[string]interface{}{"message": err.Message, "details": err.Details})
    }

    gameInformation, err := h.gameService.GameRepo.GetGameInformation(c.Request().Context(), gameID)

    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": trans.Get("Something went wrong, Please try again later.")})
    }

    players := strings.Split(gameInformation["players"].(string), ",")

    uIndex := my_slice.GetIndex(username, players)

    var player int
    var card string

    err = h.redis.Receive(c.Request().Context(), h.redis.B().Subscribe().Channel("placing_card").Build(), func(msg rueidis.PubSubMessage) {
        messages := strings.Split(msg.Message, "|")
        if messages[0] == gameID {
            gameInformation, _ = h.gameService.GameRepo.GetGameInformation(c.Request().Context(), messages[0])
            player, _ = strconv.Atoi(messages[1])
            card = messages[2]
            unsubscribeErr := h.redis.Do(c.Request().Context(), h.redis.B().Unsubscribe().Channel("placing_card").Build()).Error()
            if unsubscribeErr != nil {
                log.Println("Error while unsubscribing:", unsubscribeErr)
            }
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

    return c.JSON(http.StatusOK, transformer.GetUpdateTransformer(h.pointsService, h.playersService, transformer.GetUpdateTransformerData{
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

    _, _ = h.playersService.PlayersRepo.SavePlayer(user, chatInstance)

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
