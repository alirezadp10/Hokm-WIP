package validator

import (
    "github.com/alirezadp10/hokm/internal/service/game_service"
    "github.com/alirezadp10/hokm/internal/util/errors"
    "github.com/alirezadp10/hokm/internal/util/trans"
    "net/http"
)

type GetGameInformationValidatorData struct {
    Username string
    GameID   string
}

func GetGameInformationValidator(gameService game_service.GameService, data GetGameInformationValidatorData) *errors.ValidationError {
    ok, err := gameService.GameRepo.HasGameFinished(data.GameID)

    if err != nil {
        return &errors.ValidationError{
            StatusCode: http.StatusInternalServerError,
            Message:    trans.Get("Something went wrong. Please try again later."),
        }
    }

    if ok == true {
        return &errors.ValidationError{
            StatusCode: http.StatusForbidden,
            Message:    trans.Get("The game has already finished."),
        }
    }

    ok, err = gameService.GameRepo.DoesPlayerBelongToGame(data.Username, data.GameID)

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
