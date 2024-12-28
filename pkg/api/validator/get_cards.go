package validator

import (
    "github.com/alirezadp10/hokm/internal/util/errors"
    "github.com/alirezadp10/hokm/internal/util/trans"
    "github.com/alirezadp10/hokm/pkg/service"
    "net/http"
)

type GetCardsValidatorData struct {
    Username string
    GameID   string
}

func GetCardsValidator(gameService service.GameService, data GetCardsValidatorData) *errors.ValidationError {
    ok, err := gameService.GameRepo.DoesPlayerBelongToGame(data.Username, data.GameID)

    if err != nil {
        return &errors.ValidationError{
            StatusCode: http.StatusInternalServerError,
            Message:    trans.Get("Something went wrong. Please try again later."),
        }
    }

    if ok == false {
        return &errors.ValidationError{
            StatusCode: http.StatusForbidden,
            Message:    trans.Get("It's not your game."),
        }
    }

    return nil
}
