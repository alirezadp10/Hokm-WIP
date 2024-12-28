package transformer

import (
    "github.com/alirezadp10/hokm/pkg/api/handlers"
    "strconv"
)

type GetUpdateTransformerData struct {
    GameInformation map[string]interface{}
    UIndex          int
    PlayerIndex     int
    Card            string
}

func GetUpdateTransformer(h *handlers.HokmHandler, data GetUpdateTransformerData) map[string]interface{} {
    result := map[string]interface{}{
        "lastMove": map[string]string{
            "from": h.PlayersService.GetDirection(data.PlayerIndex, data.UIndex),
            "card": data.Card,
        },
        "points":            h.PointsService.GetPoints(data.GameInformation["points"].(string), data.UIndex),
        "centerCards":       h.PlayersService.GetPlayersCenterCards(data.GameInformation["center_cards"].(string), data.UIndex),
        "turn":              h.PlayersService.GetTurn(data.GameInformation["turn"].(string), data.UIndex),
        "king":              h.PlayersService.GetKing(data.GameInformation["king"].(string), data.UIndex),
        "timeRemained":      h.PlayersService.GetTimeRemained(data.GameInformation["last_move_timestamp"].(string)),
        "wasKingChanged":    data.GameInformation["was_king_changed"].(string),
        "trump":             data.GameInformation["trump"],
        "whoHasWonTheCards": "",
        "whoHasWonTheRound": "",
        "whoHasWonTheGame":  "",
    }

    cardsWinner, _ := data.GameInformation["who_has_won_the_cards"].(string)
    if cardsWinner != "" {
        cardsWinner, _ := strconv.Atoi(cardsWinner)
        result["whoHasWonTheCards"] = h.PlayersService.GetDirection(cardsWinner, data.UIndex)
    }

    roundWinner, _ := data.GameInformation["who_has_won_the_round"].(string)
    if roundWinner != "" {
        roundWinner, _ := strconv.Atoi(roundWinner)
        result["whoHasWonTheRound"] = h.PlayersService.GetDirection(roundWinner, data.UIndex)
    }

    gameWinner, _ := data.GameInformation["who_has_won_the_game"].(string)
    if gameWinner != "" {
        gameWinner, _ := strconv.Atoi(gameWinner)
        result["whoHasWonTheGame"] = h.PlayersService.GetDirection(gameWinner, data.UIndex)
    }

    return result
}
