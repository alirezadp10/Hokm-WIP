package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

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
)

type HokmHandler struct {
	sqlite         *gorm.DB
	redis          *rueidis.Client
	gameService    *service.GameService
	cardsService   *service.CardsService
	playersService *service.PlayersService
}

func NewHokmHandler(sqliteClient *gorm.DB, redisClient *rueidis.Client, gameService *service.GameService, cardsService *service.CardsService, playersService *service.PlayersService) *HokmHandler {
	return &HokmHandler{
		sqlite:         sqliteClient,
		redis:          redisClient,
		gameService:    gameService,
		cardsService:   cardsService,
		playersService: playersService,
	}
}

func (h *HokmHandler) CreateGame(c echo.Context) error {
	username := c.Get("username").(string)

	if err := validator.CreateGameValidator(h.playersService, validator.CreateGameValidatorData{
		Username: username,
	}); err != nil {
		return c.JSON(err.StatusCode, map[string]interface{}{"message": err.Message, "details": err.Details})
	}

	gameID := uuid.New().String()
	distributedCards := h.cardsService.DistributeCards()
	kingCards, king := h.playersService.ChooseFirstKing()

	go h.gameService.Matchmaking(c.Request().Context(), username, gameID, distributedCards, kingCards, king)

	err := (*h.redis).Receive(c.Request().Context(), (*h.redis).B().Subscribe().Channel("game_creation").Build(), func(msg rueidis.PubSubMessage) {
		message := strings.Split(msg.Message, "|")
		players := strings.Split(message[0], ",")
		if my_slice.Has(players, username) {
			gameID = message[1]
			_, err := h.playersService.PlayersRepo.AddPlayerToGame(username, gameID)
			if err != nil {
				fmt.Println(err)
			}
			unsubscribeErr := (*h.redis).Do(c.Request().Context(), (*h.redis).B().Unsubscribe().Channel("game_creation").Build()).Error()
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

	if err := validator.GetGameInformationValidator(h.playersService, validator.GetGameInformationValidatorData{
		Username: username,
		GameID:   gameID,
	}); err != nil {
		return c.JSON(err.StatusCode, map[string]interface{}{"message": err.Message, "details": err.Details})
	}

	gameInformation, err := h.gameService.GameRepo.GetGameInformation(c.Request().Context(), gameID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": trans.Get("Something went wrong, Please try again later.")})
	}

	players := strings.Split(gameInformation["players"], ",")

	uIndex := my_slice.GetIndex(username, players)

	return c.JSON(http.StatusOK, transformer.GameInformationTransformer(h.playersService, h.cardsService, transformer.GameInformationTransformerData{
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

	players := strings.Split(gameInformation["players"], ",")

	uIndex := my_slice.GetIndex(username, players)

	if err := validator.ChooseTrumpValidator(h.playersService, validator.ChooseTrumpValidatorData{
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
		"cards":        h.cardsService.GetPlayerCards(gameInformation["cards"], uIndex)[1:],
		"timeRemained": h.playersService.GetTimeRemained(gameInformation["last_move_timestamp"]),
		"turn":         h.playersService.GetTurn(gameInformation["turn"], uIndex),
	})
}

func (h *HokmHandler) GetCards(c echo.Context) error {
	username := c.Get("username").(string)
	gameID := c.Param("gameID")

	if err := validator.GetCardsValidator(h.playersService, validator.GetCardsValidatorData{
		Username: username,
		GameID:   gameID,
	}); err != nil {
		return c.JSON(err.StatusCode, map[string]interface{}{"message": err.Message, "details": err.Details})
	}

	gameInformation, err := h.gameService.GameRepo.GetGameInformation(c.Request().Context(), gameID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": trans.Get("Something went wrong, Please try again later.")})
	}

	trump := gameInformation["trump"]

	if gameInformation["trump"] == "" {
		err := (*h.redis).Receive(c.Request().Context(), (*h.redis).B().Subscribe().Channel("choosing_trump").Build(), func(msg rueidis.PubSubMessage) {
			messages := strings.Split(msg.Message, ",")
			messageId := my_slice.HasLike(messages, func(s string) bool {
				return strings.Contains(s, gameID+"|")
			})
			if messageId != -1 {
				data := strings.Split(messages[messageId], "|")
				trump = data[1]
				unsubscribeErr := (*h.redis).Do(c.Request().Context(), (*h.redis).B().Unsubscribe().Channel("choosing_trump").Build()).Error()
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

	players := strings.Split(gameInformation["players"], ",")

	uIndex := my_slice.GetIndex(username, players)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"cards": h.cardsService.GetPlayerCards(gameInformation["cards"], uIndex),
		"turn":  h.playersService.GetTurn(gameInformation["turn"], uIndex),
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

	gameInformation, _ := h.gameService.GameRepo.GetGameInformation(c.Request().Context(), gameID)

	players := strings.Split(gameInformation["players"], ",")

	uIndex := my_slice.GetIndex(username, players)

	leadSuit := h.getLeadSuit(requestBody.Card, gameInformation["lead_suit"], gameInformation["who_has_won_the_cards"])

	if err := validator.PlaceCardValidator(h.playersService, h.cardsService, validator.PlaceCardValidatorData{
		GameInformation: gameInformation,
		Username:        username,
		GameID:          gameID,
		UIndex:          uIndex,
		LeadSuit:        leadSuit,
		Card:            requestBody.Card,
	}); err != nil {
		return c.JSON(err.StatusCode, map[string]interface{}{"message": err.Message, "details": err.Details})
	}

	if gameInformation["who_has_won_the_cards"] != "" {
		gameInformation["who_has_won_the_cards"] = ""
		gameInformation["center_cards"] = ",,,"
	}

	centerCards := h.cardsService.UpdateCenterCards(gameInformation["center_cards"], requestBody.Card, uIndex)

	cardsWinner := h.cardsService.FindCardsWinner(centerCards, gameInformation["trump"], leadSuit)

	if cardsWinner != "" {
		pervKing := gameInformation["king"]
		points, roundWinner, gameWinner := h.cardsService.UpdatePoints(gameInformation["points"], cardsWinner)
		gameInformation["points"] = points
		gameInformation["who_has_won_the_round"] = roundWinner
		gameInformation["who_has_won_the_game"] = gameWinner

		gameInformation["king"] = h.playersService.GiveKing(roundWinner, gameInformation["king"])
		gameInformation["was_the_king_changed"] = my_bool.ToString(pervKing != gameInformation["king"])
	}

	gameInformation["last_move_timestamp"] = strconv.FormatInt(time.Now().Unix(), 10)
	newCards := h.cardsService.UpdateUserCards(gameInformation["cards"], requestBody.Card, uIndex)

	if gameInformation["who_has_won_the_round"] != "" {
		newCards = h.cardsService.DistributeCards()
		gameInformation["trump"] = ""
	}

	gameInformation["turn"] = h.playersService.GetNewTurn(
		gameInformation["turn"],
		cardsWinner,
		gameInformation["king"],
		gameInformation["was_the_king_changed"],
	)

	params := repository.PlaceCardParams{
		GameId:            gameID,
		Card:              requestBody.Card,
		CenterCards:       centerCards,
		LeadSuit:          leadSuit,
		CardsWinner:       cardsWinner,
		Points:            gameInformation["points"],
		Turn:              gameInformation["turn"],
		King:              gameInformation["king"],
		WasKingChanged:    gameInformation["was_the_king_changed"],
		LastMoveTimestamp: gameInformation["last_move_timestamp"],
		Trump:             gameInformation["trump"],
		Cards:             newCards,
		PlayerIndex:       uIndex,
	}

	if err := h.cardsService.CardsRepo.PlaceCard(c.Request().Context(), params); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": trans.Get("Something went wrong, Please try again later."),
		})
	}

	return c.JSON(http.StatusOK, transformer.PlaceCardTransformer(h.playersService, h.cardsService, transformer.PlaceCardTransformerData{
		Players:           players,
		UIndex:            uIndex,
		CenterCards:       centerCards,
		LeadSuit:          leadSuit,
		CardsWinner:       cardsWinner,
		Cards:             newCards,
		Points:            gameInformation["points"],
		Turn:              gameInformation["turn"],
		King:              gameInformation["king"],
		LastMoveTimestamp: gameInformation["last_move_timestamp"],
		WasKingChanged:    gameInformation["was_the_king_changed"],
		RoundWinner:       gameInformation["who_has_won_the_round"],
		GameWinner:        gameInformation["who_has_won_the_game"],
	}))
}

func (h *HokmHandler) getLeadSuit(card, currentLeadSuit, whoHasWonTheCards string) string {
	if currentLeadSuit == "" || whoHasWonTheCards != "" {
		return h.cardsService.GetCardSuit(card)
	}
	return currentLeadSuit
}

func (h *HokmHandler) GetUpdate(c echo.Context) error {
	username := c.Get("username").(string)
	gameID := c.Param("gameID")

	if err := validator.GetUpdateValidator(h.playersService, validator.GetUpdateValidatorData{
		Username: username,
		GameID:   gameID,
	}); err != nil {
		return c.JSON(err.StatusCode, map[string]interface{}{"message": err.Message, "details": err.Details})
	}

	gameInformation, err := h.gameService.GameRepo.GetGameInformation(c.Request().Context(), gameID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": trans.Get("Something went wrong, Please try again later.")})
	}

	players := strings.Split(gameInformation["players"], ",")

	uIndex := my_slice.GetIndex(username, players)

	var player int
	var card string

	err = (*h.redis).Receive(c.Request().Context(), (*h.redis).B().Subscribe().Channel("placing_card").Build(), func(msg rueidis.PubSubMessage) {
		messages := strings.Split(msg.Message, "|")
		if messages[0] == gameID {
			gameInformation, _ = h.gameService.GameRepo.GetGameInformation(c.Request().Context(), messages[0])
			player, _ = strconv.Atoi(messages[1])
			card = messages[2]
			unsubscribeErr := (*h.redis).Do(c.Request().Context(), (*h.redis).B().Unsubscribe().Channel("placing_card").Build()).Error()
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

	return c.JSON(http.StatusOK, transformer.GetUpdateTransformer(h.cardsService, h.playersService, transformer.GetUpdateTransformerData{
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
