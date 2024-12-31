package handlers

import (
    "context"
    "github.com/alirezadp10/hokm/pkg/repository"
    "github.com/alirezadp10/hokm/pkg/service"
    "github.com/labstack/echo/v4"
    "github.com/redis/rueidis"
    "github.com/stretchr/testify/assert"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestPlaceCard(t *testing.T) {
    e := echo.New()

    gameService := service.NewGameService(nil, nil, &gameMockRepo{})

    cardsService := service.NewCardsService(nil, nil, cardsRepository)

    pointsService := service.NewPointsService(nil, nil, pointsRepository, *cardsService)

    playersService := service.NewPlayersService(nil, nil, playersRepository)

    h := NewHokmHandler(nil, nil, gameService, cardsService, pointsService, playersService)

    req := httptest.NewRequest(http.MethodPost, "/game/12345/place", nil)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)
    c.Set("username", "foobar")

    if assert.NoError(t, h.PlaceCard(c)) {
        assert.Equal(t, http.StatusOK, rec.Code)
        assert.JSONEq(t, `{"message": "hi"}`, rec.Body.String())
    }
}

var _ repository.GameRepositoryContract = &gameMockRepo{}

type gameMockRepo struct {
}

func (g gameMockRepo) HasGameFinished(gameID string) (bool, error) {
    //TODO implement me
    panic("implement me")
}

func (g gameMockRepo) DoesPlayerBelongToGame(username, gameID string) (bool, error) {
    //TODO implement me
    panic("implement me")
}

func (g gameMockRepo) GetGameInformation(ctx context.Context, gameID string) (map[string]interface{}, error) {
    //TODO implement me
    panic("implement me")
}

func (g gameMockRepo) DoesPlayerHaveAnyActiveGame(username string) (*string, bool) {
    //TODO implement me
    panic("implement me")
}

func (g gameMockRepo) Matchmaking(ctx context.Context, cards []string, username, gameID, lastMoveTimestamps, king, kingCards string) {
    //TODO implement me
    panic("implement me")
}

func (g gameMockRepo) RemovePlayerFromWaitingList(ctx context.Context, key, username string) {
    //TODO implement me
    panic("implement me")
}

func (g gameMockRepo) GetGameInf(ctx context.Context, channel string, message func(rueidis.PubSubMessage)) error {
    //TODO implement me
    panic("implement me")
}
