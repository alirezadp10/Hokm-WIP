package transformer

import (
    "github.com/alirezadp10/hokm/internal/util/my_slice"
    "github.com/alirezadp10/hokm/pkg/api/handlers"
)

type GameInformationTransformerData struct {
    GameInformation map[string]interface{}
    Players         []string
    UIndex          int
}

func GameInformationTransformer(h *handlers.HokmHandler, data GameInformationTransformerData) map[string]interface{} {
    result := map[string]interface{}{
        "beginnerDirection":    h.PlayersService.GetDirection(my_slice.GetIndex(data.Players[0], data.Players), data.UIndex),
        "players":              h.PlayersService.GetPlayersWithDirections(data.Players, data.UIndex),
        "points":               h.PointsService.GetPoints(data.GameInformation["points"].(string), data.UIndex),
        "centerCards":          h.PlayersService.GetPlayersCenterCards(data.GameInformation["center_cards"].(string), data.UIndex),
        "turn":                 h.PlayersService.GetTurn(data.GameInformation["turn"].(string), data.UIndex),
        "king":                 h.PlayersService.GetKing(data.GameInformation["king"].(string), data.UIndex),
        "kingCards":            h.CardsService.GetKingCards(data.GameInformation["king_cards"].(string)),
        "timeRemained":         h.PlayersService.GetTimeRemained(data.GameInformation["last_move_timestamp"].(string)),
        "hasKingCardsFinished": data.GameInformation["has_king_cards_finished"].(string),
        "trump":                data.GameInformation["trump"],
    }

    if result["hasKingCardsFinished"] == "true" {
        result["playerCards"] = h.CardsService.GetPlayerCards(data.GameInformation["cards"].(string), data.UIndex)
    } else if result["king"] == "down" {
        result["trumpCards"] = h.CardsService.GetPlayerCards(data.GameInformation["cards"].(string), data.UIndex)[0]
    }

    return result
}
