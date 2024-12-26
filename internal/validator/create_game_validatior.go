package validator

import (
    "github.com/alirezadp10/hokm/internal/service/game_service"
    "github.com/alirezadp10/hokm/internal/util/errors"
    "github.com/alirezadp10/hokm/internal/util/trans"
    "net/http"
)

type CreateGameValidatorData struct {
    Username string
}

func CreateGameValidator(gameService game_service.GameService, data CreateGameValidatorData) *errors.ValidationError {
    if gid, ok := gameService.GameRepo.DoesPlayerHaveAnActiveGame(data.Username); ok {
        return &errors.ValidationError{
            Message:    trans.Get("You have already an active game."),
            StatusCode: http.StatusForbidden,
            Details:    *gid,
        }
    }

    return nil
}
