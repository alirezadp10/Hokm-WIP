package transformer

import (
	"strconv"

	"github.com/alirezadp10/hokm/pkg/service"
)

type GetUpdateTransformerData struct {
	GameInformation map[string]string
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
		"points":            cardsService.GetPoints(data.GameInformation["points"], data.UIndex),
		"centerCards":       playersService.GetPlayersCenterCards(data.GameInformation["center_cards"], data.UIndex),
		"turn":              playersService.GetTurn(data.GameInformation["turn"], data.UIndex),
		"king":              playersService.GetKing(data.GameInformation["king"], data.UIndex),
		"timeRemained":      playersService.GetTimeRemained(data.GameInformation["last_move_timestamp"]),
		"wasKingChanged":    data.GameInformation["was_the_king_changed"],
		"trump":             data.GameInformation["trump"],
		"whoHasWonTheCards": "",
		"whoHasWonTheRound": "",
		"whoHasWonTheGame":  "",
	}

	cardsWinner := data.GameInformation["who_has_won_the_cards"]
	if cardsWinner != "" {
		cardsWinner, _ := strconv.Atoi(cardsWinner)
		result["whoHasWonTheCards"] = playersService.GetDirection(cardsWinner, data.UIndex)
	}

	roundWinner := data.GameInformation["who_has_won_the_round"]
	if roundWinner != "" {
		roundWinner, _ := strconv.Atoi(roundWinner)
		result["whoHasWonTheRound"] = playersService.GetDirection(roundWinner, data.UIndex)
	}

	gameWinner := data.GameInformation["who_has_won_the_game"]
	if gameWinner != "" {
		gameWinner, _ := strconv.Atoi(gameWinner)
		result["whoHasWonTheGame"] = playersService.GetDirection(gameWinner, data.UIndex)
	}

	return result
}
