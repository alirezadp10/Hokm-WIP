package handlers

import (
    "context"
    "github.com/alirezadp10/hokm/pkg/database"
    "github.com/alirezadp10/hokm/pkg/repository"
    "github.com/alirezadp10/hokm/pkg/service"
    "github.com/labstack/echo/v4"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestPlaceCard(t *testing.T) {
    database.GetNewSqliteConnection()

    e := echo.New()

    gameService := service.NewGameService(&gameRepo{})

    cardsService := service.NewCardsService(&repository.CardsRepository{})

    pointsService := service.NewPointsService(&repository.PointsRepository{}, *cardsService)

    playersService := service.NewPlayersService(&repository.PlayersRepository{})

    NewHokmHandler(gameService, cardsService, pointsService, playersService, nil)

    req := httptest.NewRequest(http.MethodPost, "/game/12345/place", nil)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)
    c.Set("username", "foobar")

    // Call the handler
    //if assert.NoError(t, h.PlaceCard(c)) {
    //    assert.Equal(t, http.StatusOK, rec.Code)
    //    assert.JSONEq(t, `{"message": "hi"}`, rec.Body.String())
    //}
}

var _ repository.GameRepositoryContract = &gameRepo{}

type gameRepo struct {
}

func (r *gameRepo) HasGameFinished(gameID string) (bool, error) {
    //TODO implement me
    panic("implement me")
}

func (r *gameRepo) DoesPlayerBelongToGame(username, gameID string) (bool, error) {
    //TODO implement me
    panic("implement me")
}

func (r *gameRepo) GetGameInformation(ctx context.Context, gameID string) (map[string]interface{}, error) {
    return map[string]interface{}{
        "foo": "bar",
    }, nil
}

func (r *gameRepo) DoesPlayerHaveAnActiveGame(username string) (*string, bool) {
    //TODO implement me
    panic("implement me")
}

func (r *gameRepo) Matchmaking(ctx context.Context, cards []string, username, gameID, lastMoveTimestamps, king, kingCards string) {
    //TODO implement me
    panic("implement me")
}

func (r *gameRepo) RemovePlayerFromWaitingList(ctx context.Context, key, username string) {
    //TODO implement me
    panic("implement me")
}
