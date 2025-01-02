package transformer

import (
	"strconv"

	"github.com/alirezadp10/hokm/pkg/service"
)

type GetUpdateTransformerData struct {
	GameInformation map[string]interface{}
	UIndex          int
	PlayerIndex     int
	Card            string
}

func GetUpdateTransformer(cardsService *service.CardsService, playersService *service.PlayersService, data GetUpdateTransformerData) map[string]interface{} {
	result := map[string]interface{}{
		"lastMove": map[string]string{
			"from": playersService.GetDirection(data.PlayerIndex, data.UIndex),
			"card": data.Card,
		},
		"points":            cardsService.GetPoints(data.GameInformation["points"].(string), data.UIndex),
		"centerCards":       playersService.GetPlayersCenterCards(data.GameInformation["center_cards"].(string), data.UIndex),
		"turn":              playersService.GetTurn(data.GameInformation["turn"].(string), data.UIndex),
		"king":              playersService.GetKing(data.GameInformation["king"].(string), data.UIndex),
		"timeRemained":      playersService.GetTimeRemained(data.GameInformation["last_move_timestamp"].(string)),
		"wasKingChanged":    data.GameInformation["was_king_changed"].(string),
		"trump":             data.GameInformation["trump"],
		"whoHasWonTheCards": "",
		"whoHasWonTheRound": "",
		"whoHasWonTheGame":  "",
	}

	cardsWinner, _ := data.GameInformation["who_has_won_the_cards"].(string)
	if cardsWinner != "" {
		cardsWinner, _ := strconv.Atoi(cardsWinner)
		result["whoHasWonTheCards"] = playersService.GetDirection(cardsWinner, data.UIndex)
	}

	roundWinner, _ := data.GameInformation["who_has_won_the_round"].(string)
	if roundWinner != "" {
		roundWinner, _ := strconv.Atoi(roundWinner)
		result["whoHasWonTheRound"] = playersService.GetDirection(roundWinner, data.UIndex)
	}

	gameWinner, _ := data.GameInformation["who_has_won_the_game"].(string)
	if gameWinner != "" {
		gameWinner, _ := strconv.Atoi(gameWinner)
		result["whoHasWonTheGame"] = playersService.GetDirection(gameWinner, data.UIndex)
	}

	return result
}
