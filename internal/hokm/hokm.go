package hokm

import (
    "context"
    "fmt"
    "github.com/alirezadp10/hokm/internal/database/redis"
    "github.com/google/uuid"
    "github.com/redis/rueidis"
    "math/rand"
    "time"
)

// Global cards list
var cards = []string{
    "01C", "02C", "03C", "04C", "05C", "06C", "07C", "08C", "09C", "10C", "JC", "QC", "KC",
    "01D", "02D", "03D", "04D", "05D", "06D", "07D", "08D", "09D", "10D", "JD", "QD", "KD",
    "01H", "02H", "03H", "04H", "05H", "06H", "07H", "08H", "09H", "10H", "JH", "QH", "KH",
    "01S", "02S", "03S", "04S", "05S", "06S", "07S", "08S", "09S", "10S", "JS", "QS", "KS",
}

// SetKingCards Function to randomly select king cards
func SetKingCards() []string {
    // Create a local copy of the cards
    localCards := make([]string, len(cards))
    copy(localCards, cards)

    // Shuffle the local copy
    rand.Seed(time.Now().UnixNano())
    rand.Shuffle(len(localCards), func(i, j int) { localCards[i], localCards[j] = localCards[j], localCards[i] })

    // Result to hold the selected cards
    var result []string

    // Assign cards to players
    for {
        // Take the first card and remove it from the local deck
        card := localCards[0]
        localCards = localCards[1:]

        // Add the card to the result
        result = append(result, card)

        // If the card has "01", stop
        if card[:2] == "01" {
            return result
        }
    }
}

func GetTimeRemained(lastMoveTimestampStr string) time.Duration {
    lastMoveTimestamp, err := time.Parse(time.RFC3339, lastMoveTimestampStr)
    if err != nil {
        fmt.Println("Error parsing timestamp:", err)
        return 15
    }

    return 15*time.Second - time.Since(lastMoveTimestamp)
}

// GetKingsCards get kings cards
func GetKingsCards(cards []string, uIndex int) []interface{} {
    var result []interface{}
    for key, card := range cards {
        result = append(result, map[string]interface{}{
            "direction": GetDirection(key, uIndex),
            "card":      card,
        })
    }
    return result
}

// GetPlayersWithDirections Get players directions
func GetPlayersWithDirections(players []string, uIndex int) map[string]interface{} {
    return map[string]interface{}{
        "left":  map[string]string{"username": players[(4+(uIndex-1))%4]},
        "down":  map[string]string{"username": players[(4+(uIndex+0))%4]},
        "right": map[string]string{"username": players[(4+(uIndex+1))%4]},
        "up":    map[string]string{"username": players[(4+(uIndex+2))%4]},
    }
}

// GetPoints Get Game Points
func GetPoints(points map[string]interface{}, uIndex int) map[string]interface{} {
    var downTotalPoints, rightTotalPoints, downRoundPoints, rightRoundPoints int

    if uIndex%2 == 0 {
        downTotalPoints = points["total"].(map[string]int)["0"]
        downRoundPoints = points["round"].(map[string]int)["0"]
        rightTotalPoints = points["total"].(map[string]int)["1"]
        rightRoundPoints = points["round"].(map[string]int)["1"]
    } else {
        downTotalPoints = points["total"].(map[string]int)["1"]
        downRoundPoints = points["round"].(map[string]int)["1"]
        rightTotalPoints = points["total"].(map[string]int)["0"]
        rightRoundPoints = points["round"].(map[string]int)["0"]
    }

    return map[string]interface{}{
        "total": map[string]int{
            "down": downTotalPoints, "right": rightTotalPoints,
        },
        "currentRound": map[string]int{
            "down": downRoundPoints, "right": rightRoundPoints,
        },
    }
}

// GetCenterCards Get center cards
func GetCenterCards(centerCards map[int]string, uIndex int) map[string]string {
    var result map[string]string
    for key, val := range centerCards {
        result[GetDirection(key, uIndex)] = val
    }
    return result
}

// GetDirection Get directions
func GetDirection(pIndex, uIndex int) string {
    directions := []string{"down", "left", "up", "right"}
    diff := (4 + (uIndex - pIndex)) % 4
    return directions[diff]
}

// Matchmaking Find an open game for a player
func Matchmaking(ctx context.Context, client rueidis.Client, userId string) {
    gameId := uuid.New().String()
    distributedCards := DistributeCards()
    redis.Matchmaking(ctx, client, userId, gameId, distributedCards)
}

func DistributeCards() [][]string {
    localCards := make([]string, len(cards))
    copy(localCards, cards)

    rand.Seed(time.Now().UnixNano())
    rand.Shuffle(len(localCards), func(i, j int) { localCards[i], localCards[j] = localCards[j], localCards[i] })

    hands := make([][]string, 4)
    for i := range hands {
        hands[i] = []string{}
    }

    for i, card := range cards {
        player := i % 4
        hands[player] = append(hands[player], card)
    }

    return hands
}

func GetPlayerCards(cards map[int][]string, uIndex int) [][]string {
    var result [][]string
    var chunk []string
    for i, card := range cards[uIndex] {
        chunk = append(chunk, card)
        if (i+1)%5 == 0 {
            result = append(result, chunk)
            chunk = []string{}
        }
    }
    return result
}
