package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/alirezadp10/hokm/internal/util/tests"
	"github.com/alirezadp10/hokm/pkg/mocks"
	"github.com/alirezadp10/hokm/pkg/service"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestNextTurnAfterPlacingCardWithoutCardsWinner(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	e := echo.New()

	c, rec := tests.ApiCall(e, "POST", "/game/:gameID/place", map[string]string{
		"gameID": "12345",
	}, map[string]string{
		"Content-Type": "application/json",
	}, map[string]string{
		"card": "QC",
	})

	c.Set("username", "0")

	mockGameRepo := mocks.NewMockGameRepositoryContract(ctrl)
	mockGameRepo.EXPECT().GetGameInformation(gomock.Any(), gomock.Any()).Return(
		tests.NewGameDataBuilder().BeginingState().SetTrump("C").Build(), nil,
	)

	mockCardRepo := mocks.NewMockCardsRepositoryContract(ctrl)
	mockCardRepo.EXPECT().PlaceCard(gomock.Any(), gomock.Any()).Return(nil)

	mockPlayerRepo := mocks.NewMockPlayersRepositoryContract(ctrl)
	mockPlayerRepo.EXPECT().DoesPlayerBelongToGame(gomock.Any(), gomock.Any()).Return(true, nil)

	gameService := service.NewGameService(nil, nil, mockGameRepo)

	cardsService := service.NewCardsService(nil, nil, mockCardRepo)

	playersService := service.NewPlayersService(nil, nil, mockPlayerRepo)

	h := NewHokmHandler(nil, nil, gameService, cardsService, playersService)

	if assert.NoError(t, h.PlaceCard(c)) {
		resp := make(map[string]interface{})
		json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "right", resp["turn"])
	}
}

func TestNextTurnAfterPlacingCardWithCardsWinner(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	e := echo.New()

	c, rec := tests.ApiCall(e, "POST", "/game/:gameID/place", map[string]string{
		"gameID": "12345",
	}, map[string]string{
		"Content-Type": "application/json",
	}, map[string]string{
		"card": "KS",
	})

	c.Set("username", "3")

	mockGameRepo := mocks.NewMockGameRepositoryContract(ctrl)
	mockGameRepo.EXPECT().GetGameInformation(gomock.Any(), gomock.Any()).Return(
		tests.NewGameDataBuilder().BeginingState().
			SetTrump("C").
			SetCenterCards("JS,04S,06S,").
			SetTurn("3").
			Build(),
		nil,
	)

	mockCardRepo := mocks.NewMockCardsRepositoryContract(ctrl)
	mockCardRepo.EXPECT().PlaceCard(gomock.Any(), gomock.Any()).Return(nil)

	mockPlayerRepo := mocks.NewMockPlayersRepositoryContract(ctrl)
	mockPlayerRepo.EXPECT().DoesPlayerBelongToGame(gomock.Any(), gomock.Any()).Return(true, nil)

	gameService := service.NewGameService(nil, nil, mockGameRepo)

	cardsService := service.NewCardsService(nil, nil, mockCardRepo)

	playersService := service.NewPlayersService(nil, nil, mockPlayerRepo)

	h := NewHokmHandler(nil, nil, gameService, cardsService, playersService)

	if assert.NoError(t, h.PlaceCard(c)) {
		fmt.Print(rec.Body.String())
		resp := make(map[string]interface{})
		json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "down", resp["turn"])
	}
}

func TestPlacingCardAfterGatheringCards(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	e := echo.New()

	c, rec := tests.ApiCall(e, "POST", "/game/:gameID/place", map[string]string{
		"gameID": "12345",
	}, map[string]string{
		"Content-Type": "application/json",
	}, map[string]string{
		"card": "01H",
	})

	c.Set("username", "3")

	mockGameRepo := mocks.NewMockGameRepositoryContract(ctrl)
	mockGameRepo.EXPECT().GetGameInformation(gomock.Any(), gomock.Any()).Return(
		tests.NewGameDataBuilder().BeginingState().
			SetTrump("C").
			SetLeadSuit("S").
			SetWasKingChanged("").
			SetCenterCards("JS,04S,06S,KS").
			SetWhoHasWonTheCards("3").
			SetTurn("3").
			Build(),
		nil,
	)

	mockCardRepo := mocks.NewMockCardsRepositoryContract(ctrl)
	mockCardRepo.EXPECT().PlaceCard(gomock.Any(), gomock.Any()).Return(nil)

	mockPlayerRepo := mocks.NewMockPlayersRepositoryContract(ctrl)
	mockPlayerRepo.EXPECT().DoesPlayerBelongToGame(gomock.Any(), gomock.Any()).Return(true, nil)

	gameService := service.NewGameService(nil, nil, mockGameRepo)

	cardsService := service.NewCardsService(nil, nil, mockCardRepo)

	playersService := service.NewPlayersService(nil, nil, mockPlayerRepo)

	h := NewHokmHandler(nil, nil, gameService, cardsService, playersService)

	if assert.NoError(t, h.PlaceCard(c)) {
		resp := make(map[string]interface{})
		json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "right", resp["turn"])
		assert.Equal(t, "H", resp["leadSuit"])
		assert.Equal(t, map[string]interface{}{
			"down":  "01H",
			"left":  "",
			"right": "",
			"up":    "",
		}, resp["centerCards"])
	}
}
