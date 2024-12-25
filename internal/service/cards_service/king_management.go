package cards_service

import (
    "github.com/alirezadp10/hokm/internal/service/players_service"
    "math/rand"
    "strconv"
    "strings"
)

func ChooseFirstKing() (string, string) {
    localCards := make([]string, len(Cards)) // Create a copy of the cards to shuffle
    copy(localCards, Cards)

    rand.Shuffle(len(localCards), func(i, j int) { localCards[i], localCards[j] = localCards[j], localCards[i] }) // Shuffle cards

    var cardsList []string // To track the sequence of dealt cards
    i := 0

    // Deal cards until a king is found
    for {
        card := localCards[0]
        localCards = localCards[1:] // Remove dealt card from the deck
        cardsList = append(cardsList, card)

        if card[:2] == "01" { // Check if the card is "01" (e.g., Ace of any suit)
            return strings.Join(cardsList, ","), strconv.Itoa(i % 4) // Return the sequence and player index
        }

        i++
    }
}

func GetKingCards(cards string) []string {
    return strings.Split(cards, ",")
}

func GetKing(kingIndex string, uIndex int) string {
    kingI, _ := strconv.Atoi(kingIndex) // Convert king index to integer
    return players_service.GetDirection(kingI, uIndex)
}

func GiveKing(roundWinnerStr, prevKingStr string) string {
    roundWinner, _ := strconv.Atoi(roundWinnerStr)
    prevKing, _ := strconv.Atoi(prevKingStr)
    if roundWinner%2 == prevKing%2 {
        return prevKingStr
    }
    return strconv.Itoa(players_service.NextPlayerIndex(roundWinner))
}
