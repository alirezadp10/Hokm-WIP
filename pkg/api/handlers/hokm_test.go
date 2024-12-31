package handlers

import (
    "github.com/alirezadp10/hokm/pkg/mocks"
    "github.com/alirezadp10/hokm/pkg/service"
    "github.com/golang/mock/gomock"
    "github.com/labstack/echo/v4"
    "github.com/stretchr/testify/assert"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestPlaceCard(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    e := echo.New()

    req := httptest.NewRequest(http.MethodPost, "/game/12345/place", nil)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)

    c.Set("username", "foobar")

    gameMockRepo := mocks.NewMockGameRepositoryContract(ctrl)
    gameMockRepo.EXPECT().GetGameInformation(c, "1234").Return(map[string]interface{}{
        "foo": "bar",
    }, nil)

    gameService := service.NewGameService(nil, nil, gameMockRepo)

    cardsService := service.NewCardsService(nil, nil, mocks.NewMockCardsRepositoryContract(ctrl))

    playersService := service.NewPlayersService(nil, nil, mocks.NewMockPlayersRepositoryContract(ctrl))

    h := NewHokmHandler(nil, nil, gameService, cardsService, playersService)

    if assert.NoError(t, h.PlaceCard(c)) {
        assert.Equal(t, http.StatusOK, rec.Code)
        assert.JSONEq(t, `{"message": "hi"}`, rec.Body.String())
    }
}
