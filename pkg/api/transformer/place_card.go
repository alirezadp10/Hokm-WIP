package transformer

import (
    "encoding/json"
    "strconv"

    "github.com/alirezadp10/hokm/pkg/service"
)

type PlaceCardTransformerData struct {
    Players           []string
    UIndex            int
    Points            string
    Cards             map[int][]string
    CenterCards       string
    Turn              string
    King              string
    Trump             string
    LastMoveTimestamp string
    WasTheKingChanged string
    CardsWinner       string
    RoundWinner       string
    GameWinner        string
    LeadSuit          string
}

func PlaceCardTransformer(playersService *service.PlayersService, cardsService *service.CardsService, data PlaceCardTransformerData) map[string]interface{} {
    cards, _ := json.Marshal(data.Cards)
    result := map[string]interface{}{
        "players":           playersService.GetPlayersWithDirections(data.Players, data.UIndex),
        "points":            cardsService.GetPoints(data.Points, data.UIndex),
        "centerCards":       playersService.GetPlayersCenterCards(data.CenterCards, data.UIndex),
        "turn":              playersService.GetTurn(data.Turn, data.UIndex),
        "king":              playersService.GetKing(data.King, data.UIndex),
        "timeRemained":      playersService.GetTimeRemained(data.LastMoveTimestamp),
        "playerCards":       cardsService.GetPlayerCards(string(cards), data.UIndex),
        "WasTheKingChanged": data.WasTheKingChanged,
        "trump":             data.Trump,
        "leadSuit":          data.LeadSuit,
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
