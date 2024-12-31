package validator

import (
    "github.com/alirezadp10/hokm/internal/util/errors"
    "github.com/alirezadp10/hokm/internal/util/trans"
    "github.com/alirezadp10/hokm/pkg/service"
    "net/http"
)

type CreateGameValidatorData struct {
    Username string
}

func CreateGameValidator(playersService service.PlayersService, data CreateGameValidatorData) *errors.ValidationError {
    if gid, ok := playersService.PlayersRepo.DoesPlayerHaveAnyActiveGame(data.Username); ok {
        return &errors.ValidationError{
            Message:    trans.Get("You have already an active game."),
            StatusCode: http.StatusForbidden,
            Details:    *gid,
        }
    }

    return nil
}
