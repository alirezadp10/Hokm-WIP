package transformer

import (
    "github.com/alirezadp10/hokm/pkg/service"
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

func PlaceCardTransformer(playersService service.PlayersService, pointsService service.PointsService, cardsService service.CardsService, data PlaceCardTransformerData) map[string]interface{} {
    result := map[string]interface{}{
        "players":           playersService.GetPlayersWithDirections(data.Players, data.UIndex),
        "points":            pointsService.GetPoints(data.Points, data.UIndex),
        "centerCards":       playersService.GetPlayersCenterCards(data.CenterCards, data.UIndex),
        "turn":              playersService.GetTurn(data.Turn, data.UIndex),
        "king":              playersService.GetKing(data.King, data.UIndex),
        "timeRemained":      playersService.GetTimeRemained(data.LastMoveTimestamp),
        "playerCards":       cardsService.GetPlayerCards(data.GameInformation["cards"].(string), data.UIndex),
        "wasKingChanged":    data.WasKingChanged,
        "trump":             data.GameInformation["trump"],
        "whoHasWonTheCards": "",
        "whoHasWonTheRound": "",
        "whoHasWonTheGame":  "",
    }

    if data.CardsWinner != "" {
        cardsWinner, _ := strconv.Atoi(data.CardsWinner)
        result["whoHasWonTheCards"] = playersService.GetDirection(cardsWinner, data.UIndex)
    }

    if data.RoundWinner != "" {
        roundWinner, _ := strconv.Atoi(data.RoundWinner)
        result["whoHasWonTheRound"] = playersService.GetDirection(roundWinner, data.UIndex)
    }

    if data.GameWinner != "" {
        gameWinner, _ := strconv.Atoi(data.GameWinner)
        result["whoHasWonTheGame"] = playersService.GetDirection(gameWinner, data.UIndex)
    }

    return result
}
