package handlers

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/alirezadp10/hokm/internal/util/tests"
	"github.com/alirezadp10/hokm/pkg/mocks"
	"github.com/alirezadp10/hokm/pkg/service"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestPlaceCard(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	e := echo.New()

	c, rec := tests.ApiCall(e, "POST", "/game/:gameID/place", map[string]string{
		"gameID": "12345",
	}, map[string]string{
		"Content-Type": "application/json",
	}, map[string]string{
		"card": "QD",
	})

	c.Set("username", "1")

	mockGameRepo := mocks.NewMockGameRepositoryContract(ctrl)
	mockGameRepo.EXPECT().GetGameInformation(gomock.Any(), "12345").Return(tests.NewGameDataBuilder().Build(), nil)

	mockCardRepo := mocks.NewMockCardsRepositoryContract(ctrl)
	mockCardRepo.EXPECT().PlaceCard(gomock.Any(), gomock.Any()).Return(nil)

	mockPlayerRepo := mocks.NewMockPlayersRepositoryContract(ctrl)
	mockPlayerRepo.EXPECT().DoesPlayerBelongToGame("1", "12345").Return(true, nil)

	gameService := service.NewGameService(nil, nil, mockGameRepo)

	cardsService := service.NewCardsService(nil, nil, mockCardRepo)

	playersService := service.NewPlayersService(nil, nil, mockPlayerRepo)

	h := NewHokmHandler(nil, nil, gameService, cardsService, playersService)

	if assert.NoError(t, h.PlaceCard(c)) {
		resp := make(map[string]interface{})
		json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, resp["turn"], "down")
	}
}
