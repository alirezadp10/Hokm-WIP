package validator

import (
	"net/http"
	"strconv"

	"github.com/alirezadp10/hokm/internal/util/errors"
	"github.com/alirezadp10/hokm/internal/util/my_slice"
	"github.com/alirezadp10/hokm/internal/util/trans"
	"github.com/alirezadp10/hokm/pkg/service"
)

type ChooseTrumpValidatorData struct {
	GameInformation map[string]interface{}
	UIndex          int
	Username        string
	GameID          string
	Trump           string
}

func ChooseTrumpValidator(s *service.PlayersService, data ChooseTrumpValidatorData) *errors.ValidationError {
	if !my_slice.Has([]string{"H", "D", "C", "S"}, data.Trump) {
		return &errors.ValidationError{
			StatusCode: http.StatusBadRequest,
			Message:    trans.Get("Invalid trump."),
		}
	}

	ok, err := s.PlayersRepo.DoesPlayerBelongToGame(data.Username, data.GameID)

	if err != nil {
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

	if data.GameInformation["king"].(string) != strconv.Itoa(data.UIndex) {
		return &errors.ValidationError{
			StatusCode: http.StatusForbidden,
			Message:    trans.Get("You're not king in this round."),
		}
	}

	if data.GameInformation["has_king_cards_finished"].(string) == "true" {
		return &errors.ValidationError{
			StatusCode: http.StatusForbidden,
			Message:    trans.Get("You're not allowed to choose a trump at the moment."),
		}
	}

	return nil
}
