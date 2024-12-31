package transformer

import (
    "github.com/alirezadp10/hokm/internal/util/my_slice"
    "github.com/alirezadp10/hokm/pkg/service"
)

type GameInformationTransformerData struct {
    GameInformation map[string]interface{}
    Players         []string
    UIndex          int
}

func GameInformationTransformer(playersService service.PlayersService, cardsService service.CardsService, data GameInformationTransformerData) map[string]interface{} {
    result := map[string]interface{}{
        "beginnerDirection":    playersService.GetDirection(my_slice.GetIndex(data.Players[0], data.Players), data.UIndex),
        "players":              playersService.GetPlayersWithDirections(data.Players, data.UIndex),
        "points":               cardsService.GetPoints(data.GameInformation["points"].(string), data.UIndex),
        "centerCards":          playersService.GetPlayersCenterCards(data.GameInformation["center_cards"].(string), data.UIndex),
        "turn":                 playersService.GetTurn(data.GameInformation["turn"].(string), data.UIndex),
        "king":                 playersService.GetKing(data.GameInformation["king"].(string), data.UIndex),
        "kingCards":            cardsService.GetKingCards(data.GameInformation["king_cards"].(string)),
        "timeRemained":         playersService.GetTimeRemained(data.GameInformation["last_move_timestamp"].(string)),
        "hasKingCardsFinished": data.GameInformation["has_king_cards_finished"].(string),
        "trump":                data.GameInformation["trump"],
    }

    if result["hasKingCardsFinished"] == "true" {
        result["playerCards"] = cardsService.GetPlayerCards(data.GameInformation["cards"].(string), data.UIndex)
    } else if result["king"] == "down" {
        result["trumpCards"] = cardsService.GetPlayerCards(data.GameInformation["cards"].(string), data.UIndex)[0]
    }

    return result
}
