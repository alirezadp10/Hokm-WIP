package transformer

import (
	"github.com/alirezadp10/hokm/internal/util/my_slice"
	"github.com/alirezadp10/hokm/pkg/service"
)

type GameInformationTransformerData struct {
	GameInformation map[string]string
	Players         []string
	UIndex          int
}

func GameInformationTransformer(playersService *service.PlayersService, cardsService *service.CardsService, data GameInformationTransformerData) map[string]interface{} {
	result := map[string]interface{}{
		"beginnerDirection":    playersService.GetDirection(my_slice.GetIndex(data.Players[0], data.Players), data.UIndex),
		"players":              playersService.GetPlayersWithDirections(data.Players, data.UIndex),
		"points":               cardsService.GetPoints(data.GameInformation["points"], data.UIndex),
		"centerCards":          playersService.GetPlayersCenterCards(data.GameInformation["center_cards"], data.UIndex),
		"turn":                 playersService.GetTurn(data.GameInformation["turn"], data.UIndex),
		"king":                 playersService.GetKing(data.GameInformation["king"], data.UIndex),
		"kingCards":            cardsService.GetKingCards(data.GameInformation["king_cards"]),
		"timeRemained":         playersService.GetTimeRemained(data.GameInformation["last_move_timestamp"]),
		"hasKingCardsFinished": data.GameInformation["has_king_cards_finished"],
		"trump":                data.GameInformation["trump"],
	}

	if result["hasKingCardsFinished"] == "true" {
		result["playerCards"] = cardsService.GetPlayerCards(data.GameInformation["cards"], data.UIndex)
	} else if result["king"] == "down" {
		result["trumpCards"] = cardsService.GetPlayerCards(data.GameInformation["cards"], data.UIndex)[0]
	}

	return result
}
