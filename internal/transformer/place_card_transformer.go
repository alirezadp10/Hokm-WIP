package transformer

import (
    "github.com/alirezadp10/hokm/internal/handler"
    "strconv"
)

type PlaceCardTransformerData struct {
    GameInformation   map[string]interface{}
    Players           []string
    UIndex            int
    Points            string
    CenterCards       string
    Turn              string
    King              string
    LastMoveTimestamp string
    WasKingChanged    bool
    CardsWinner       string
    RoundWinner       string
    GameWinner        string
}

func PlaceCardTransformer(h *handler.Handler, data PlaceCardTransformerData) map[string]interface{} {
    result := map[string]interface{}{
        "players":           h.PlayersService.GetPlayersWithDirections(data.Players, data.UIndex),
        "points":            h.PointsService.GetPoints(data.Points, data.UIndex),
        "centerCards":       h.PlayersService.GetPlayersCenterCards(data.CenterCards, data.UIndex),
        "turn":              h.PlayersService.GetTurn(data.Turn, data.UIndex),
        "king":              h.PlayersService.GetKing(data.King, data.UIndex),
        "timeRemained":      h.PlayersService.GetTimeRemained(data.LastMoveTimestamp),
        "playerCards":       h.CardsService.GetPlayerCards(data.GameInformation["cards"].(string), data.UIndex),
        "wasKingChanged":    data.WasKingChanged,
        "trump":             data.GameInformation["trump"],
        "whoHasWonTheCards": "",
        "whoHasWonTheRound": "",
        "whoHasWonTheGame":  "",
    }

    if data.CardsWinner != "" {
        cardsWinner, _ := strconv.Atoi(data.CardsWinner)
        result["whoHasWonTheCards"] = h.PlayersService.GetDirection(cardsWinner, data.UIndex)
    }

    if data.RoundWinner != "" {
        roundWinner, _ := strconv.Atoi(data.RoundWinner)
        result["whoHasWonTheRound"] = h.PlayersService.GetDirection(roundWinner, data.UIndex)
    }

    if data.GameWinner != "" {
        gameWinner, _ := strconv.Atoi(data.GameWinner)
        result["whoHasWonTheGame"] = h.PlayersService.GetDirection(gameWinner, data.UIndex)
    }

    return result
}
