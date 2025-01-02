package tests

import (
	"encoding/json"
)

type GameDataBuilder struct {
	data map[string]interface{}
}

func NewGameDataBuilder() *GameDataBuilder {
	cards := getSampleCards()
	points := getSamplePoints()

	return &GameDataBuilder{data: map[string]interface{}{
		"players":                 "0,1,2,3",
		"cards":                   string(cards),
		"points":                  string(points),
		"center_cards":            ",,,",
		"trump":                   "C",
		"king":                    "0",
		"lead_suit":               "",
		"is_it_new_round":         "false",
		"turn":                    "1",
		"king_cards":              "01C",
		"was_king_changed":        "false",
		"has_king_cards_finished": "true",
		"who_has_won_the_cards":   "",
		"who_has_won_the_round":   "",
		"who_has_won_the_game":    "",
		"last_move_timestamp":     "1234567890",
	}}
}

func (b *GameDataBuilder) SetPlayers(players string) *GameDataBuilder {
	b.data["players"] = players
	return b
}

func (b *GameDataBuilder) SetCards(cards string) *GameDataBuilder {
	b.data["cards"] = cards
	return b
}

func (b *GameDataBuilder) SetPoints(points string) *GameDataBuilder {
	b.data["points"] = points
	return b
}

func (b *GameDataBuilder) SetCenterCards(centerCards string) *GameDataBuilder {
	b.data["center_cards"] = centerCards
	return b
}

func (b *GameDataBuilder) SetTrump(trump string) *GameDataBuilder {
	b.data["trump"] = trump
	return b
}

func (b *GameDataBuilder) SetKing(king string) *GameDataBuilder {
	b.data["king"] = king
	return b
}

func (b *GameDataBuilder) SetLeadSuit(leadSuit string) *GameDataBuilder {
	b.data["lead_suit"] = leadSuit
	return b
}

func (b *GameDataBuilder) SetIsItNewRound(isItNewRound string) *GameDataBuilder {
	b.data["is_it_new_round"] = isItNewRound
	return b
}

func (b *GameDataBuilder) SetTurn(turn string) *GameDataBuilder {
	b.data["turn"] = turn
	return b
}

func (b *GameDataBuilder) SetKingCards(kingCards string) *GameDataBuilder {
	b.data["king_cards"] = kingCards
	return b
}

func (b *GameDataBuilder) SetWasKingChanged(wasKingChanged string) *GameDataBuilder {
	b.data["was_king_changed"] = wasKingChanged
	return b
}

func (b *GameDataBuilder) SetHasKingCardsFinished(hasKingCardsFinished string) *GameDataBuilder {
	b.data["has_king_cards_finished"] = hasKingCardsFinished
	return b
}

func (b *GameDataBuilder) SetWhoHasWonTheCards(whoHasWonTheCards string) *GameDataBuilder {
	b.data["who_has_won_the_cards"] = whoHasWonTheCards
	return b
}

func (b *GameDataBuilder) SetWhoHasWonTheRound(whoHasWonTheRound string) *GameDataBuilder {
	b.data["who_has_won_the_round"] = whoHasWonTheRound
	return b
}

func (b *GameDataBuilder) SetWhoHasWonTheGame(whoHasWonTheGame string) *GameDataBuilder {
	b.data["who_has_won_the_game"] = whoHasWonTheGame
	return b
}

func (b *GameDataBuilder) SetLastMoveTimestamp(lastMoveTimestamp string) *GameDataBuilder {
	b.data["last_move_timestamp"] = lastMoveTimestamp
	return b
}

func (b *GameDataBuilder) Build() map[string]interface{} {
	return b.data
}

func getSamplePoints() []byte {
	points, _ := json.Marshal(map[string]string{
		"round": "0,0",
		"total": "0,0",
	})
	return points
}

func getSampleCards() []byte {
	cards, _ := json.Marshal(map[string][]string{
		"0": {"01S", "02C", "03C", "04C", "05C", "06C", "07C", "08C", "09C", "10C", "JC", "QC", "KC"},
		"1": {"01C", "02D", "03D", "04D", "05D", "06D", "07D", "08D", "09D", "10D", "JD", "QD", "KD"},
		"2": {"01D", "02H", "03H", "04H", "05H", "06H", "07H", "08H", "09H", "10H", "JH", "QH", "KH"},
		"3": {"01H", "02S", "03S", "04S", "05S", "06S", "07S", "08S", "09S", "10S", "JS", "QS", "KS"},
	})
	return cards
}
