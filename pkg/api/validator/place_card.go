package validator

import (
	"net/http"
	"strconv"

	"github.com/alirezadp10/hokm/internal/util/errors"
	"github.com/alirezadp10/hokm/internal/util/my_slice"
	"github.com/alirezadp10/hokm/internal/util/trans"
	"github.com/alirezadp10/hokm/pkg/service"
	"github.com/labstack/gommon/log"
)

type PlaceCardValidatorData struct {
	GameInformation map[string]interface{}
	Card            string
	Username        string
	GameID          string
	UIndex          int
	LeadSuit        string
}

func PlaceCardValidator(playersService *service.PlayersService, cardsService *service.CardsService, data PlaceCardValidatorData) *errors.ValidationError {
	isSelectedCardForUser := false

	doesPlayerHaveLeadSuitCard := false

	for _, cards := range cardsService.GetPlayerCards(data.GameInformation["cards"].(string), data.UIndex) {
		for _, card := range cards {
			if card == data.Card {
				isSelectedCardForUser = true
			}
			if cardsService.GetCardSuit(card) == data.LeadSuit {
				doesPlayerHaveLeadSuitCard = true
			}
		}
	}

	if !my_slice.Has(service.Cards, data.Card) {
		return &errors.ValidationError{
			StatusCode: http.StatusBadRequest,
			Message:    trans.Get("Invalid Card."),
		}
	}

	ok, err := playersService.PlayersRepo.DoesPlayerBelongToGame(data.Username, data.GameID)

	if err != nil {
		log.Fatal(err)
		return &errors.ValidationError{
			StatusCode: http.StatusInternalServerError,
			Message:    trans.Get("Something went wrong, Please try again later."),
		}
	}

	if ok == false {
		return &errors.ValidationError{
			StatusCode: http.StatusForbidden,
			Message:    trans.Get("It's not your game."),
		}
	}

	if data.GameInformation["turn"].(string) != strconv.Itoa(data.UIndex) {
		return &errors.ValidationError{
			StatusCode: http.StatusForbidden,
			Message:    trans.Get("It's not your turn."),
		}
	}

	if data.LeadSuit != "" && doesPlayerHaveLeadSuitCard && cardsService.GetCardSuit(data.Card) != data.LeadSuit {
		return &errors.ValidationError{
			StatusCode: http.StatusForbidden,
			Message:    trans.Get("You're not allowed to select this card."),
		}
	}

	if !isSelectedCardForUser {
		return &errors.ValidationError{
			StatusCode: http.StatusForbidden,
			Message:    trans.Get("It's not your card."),
		}
	}

	if data.GameInformation["trump"].(string) == "" {
		return &errors.ValidationError{
			StatusCode: http.StatusForbidden,
			Message:    trans.Get("It's not your turn."),
		}
	}

	return nil
}
