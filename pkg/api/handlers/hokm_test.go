package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

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

	c, rec := ApiCall(e, "POST", "/game/:gameID/place", map[string]string{
		"gameID": "12345",
	}, map[string]string{
		"Content-Type": "application/json",
	}, map[string]string{
		"card": "QC",
	})

	c.Set("username", "b")

	mockGameRepo := mocks.NewMockGameRepositoryContract(ctrl)
	mockGameRepo.EXPECT().GetGameInformation(gomock.Any(), "12345").Return(getSampleGameData(), nil)

	mockCardRepo := mocks.NewMockCardsRepositoryContract(ctrl)
	mockCardRepo.EXPECT().PlaceCard(gomock.Any(), gomock.Any()).Return(nil)

	mockPlayerRepo := mocks.NewMockPlayersRepositoryContract(ctrl)
	mockPlayerRepo.EXPECT().DoesPlayerBelongToGame("b", "12345").Return(true, nil)

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

func getSampleGameData() map[string]interface{} {
	cards := getSampleCards()

	points := getSamplePoints()

	return map[string]interface{}{
		"who_has_won_the_cards":   "1",
		"last_move_timestamp":     "1735477626",
		"was_king_changed":        "true",
		"center_cards":            "04C,09C,07C,03C",
		"trump":                   "C",
		"cards":                   string(cards),
		"has_king_cards_finished": "true",
		"king":                    "0",
		"players":                 "a,b,c,d",
		"who_has_won_the_game":    "",
		"lead_suit":               "C",
		"is_it_new_round":         "false",
		"turn":                    "1",
		"who_has_won_the_round":   "",
		"king_cards":              "01C",
		"points":                  string(points),
		"current_turn":            "",
		"players_cards":           "",
	}
}

func getSamplePoints() []byte {
	points, _ := json.Marshal(map[string]string{
		"round": "2,1",
		"total": "0,0",
	})
	return points
}

func getSampleCards() []byte {
	cards, _ := json.Marshal(map[string][]string{
		"0": {"JH", "01C", "02H", "02C", "02S", "07H", "02D", "06H", "06D", "09H", "QS"},
		"1": {"05D", "QC", "05H", "KC", "JS", "01D", "08D", "01S", "QH", "03H", "04D"},
		"2": {"KD", "KS", "JD", "10C", "03D", "10D", "03S", "09S", "10H", "05S", "10S", "JC"},
		"3": {"08S", "09D", "06S", "08H", "01H", "07S", "08C", "04H", "KH", "07D", "04S", "QD"},
	})
	return cards
}

func ApiCall(e *echo.Echo, method, url string, routeParams, headers, body map[string]string) (echo.Context, *httptest.ResponseRecorder) {
	marshaledBody, _ := json.Marshal(body)
	req := httptest.NewRequest(method, url, bytes.NewReader(marshaledBody))
	rec := httptest.NewRecorder()
	for headerKey, headerValue := range headers {
		req.Header.Set(headerKey, headerValue)
	}
	c := e.NewContext(req, rec)
	c.SetPath(url)
	for paramKey, paramValue := range routeParams {
		c.SetParamNames(paramKey)
		c.SetParamValues(paramValue)
	}
	return c, rec
}

func convertByteToStringArray(data []byte) []string {
	var result []string
	_ = json.Unmarshal(data, &result)
	return result
}
