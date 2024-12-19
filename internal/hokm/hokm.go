package hokm

import (
    "context"
    "encoding/json"
    "fmt"
    "github.com/alirezadp10/hokm/internal/database/redis"
    "github.com/redis/rueidis"
    "math/rand"
    "strconv"
    "strings"
    "time"
)

// Global cards list
var cards = []string{
    "01C", "02C", "03C", "04C", "05C", "06C", "07C", "08C", "09C", "10C", "JC", "QC", "KC",
    "01D", "02D", "03D", "04D", "05D", "06D", "07D", "08D", "09D", "10D", "JD", "QD", "KD",
    "01H", "02H", "03H", "04H", "05H", "06H", "07H", "08H", "09H", "10H", "JH", "QH", "KH",
    "01S", "02S", "03S", "04S", "05S", "06S", "07S", "08S", "09S", "10S", "JS", "QS", "KS",
}

// chooseFirstJudge Function to randomly select king cards
func chooseFirstJudge() (string, string) {
    // Create a local copy of the cards
    localCards := make([]string, len(cards))
    copy(localCards, cards)

    // Shuffle the local copy
    rand.Shuffle(len(localCards), func(i, j int) { localCards[i], localCards[j] = localCards[j], localCards[i] })

    // Result to hold the selected cards
    var cardsList []string

    i := 0
    // Assign cards to players
    for {
        // Take the first card and remove it from the local deck
        card := localCards[0]
        localCards = localCards[1:]

        // Add the card to the cards
        cardsList = append(cardsList, card)

        // If the card has "01", stop
        if card[:2] == "01" {
            return `["` + strings.Join(cardsList, `","`) + `"]`, strconv.Itoa(i % 4)
        }

        i++
    }
}

func GetJudgeCards(judgeCardsString string) []string {
    var judgeCards []string
    err := json.Unmarshal([]byte(judgeCardsString), &judgeCards)
    if err != nil {
        fmt.Println("Error unmarshalling:", err)
    }
    return judgeCards
}

func GetJudge(judgeIndex string, uIndex int) string {
    judgeI, _ := strconv.Atoi(judgeIndex)
    return GetDirection(judgeI, uIndex)
}

func GetTurn(turnIndex string, uIndex int) string {
    if turnIndex == "" {
        return ""
    }
    turnI, _ := strconv.Atoi(turnIndex)
    return GetDirection(turnI, uIndex)
}

func GetTimeRemained(lastMoveTimestampStr string) time.Duration {
    // Convert the string to an integer (Unix timestamp)
    lastMoveTimestampInt, err := strconv.ParseInt(lastMoveTimestampStr, 10, 64)
    if err != nil {
        fmt.Println("Error parsing timestamp:", err)
        return 15 * time.Second // Return default duration if parsing fails
    }

    // Convert Unix timestamp to time.Time
    lastMoveTimestamp := time.Unix(lastMoveTimestampInt, 0)

    // Calculate time remaining
    timeElapsed := time.Since(lastMoveTimestamp)
    timeRemaining := 15*time.Second - timeElapsed

    // Ensure non-negative duration
    if timeRemaining < 0 {
        timeRemaining = 0
    }

    timeRemaining = timeRemaining.Round(time.Second)

    return timeRemaining
}

// GetjudgeCards get kings cards
func GetjudgeCards(cards string, uIndex int) []interface{} {
    var result []interface{}
    if cards == "" {
        return result
    }

    judgeCards := make(map[int]string)
    err := json.Unmarshal([]byte(cards), &judgeCards)
    if err != nil {
        fmt.Println("Error unmarshalling:", err)
    }

    for key, card := range judgeCards {
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
func GetPoints(pointsString string, uIndex int) map[string]map[string]string {
    points := make(map[string]string)
    err := json.Unmarshal([]byte(pointsString), &points)
    if err != nil {
        fmt.Println("Error unmarshalling:", err)
    }

    var downTotalPoints, rightTotalPoints, downRoundPoints, rightRoundPoints string

    if uIndex%2 == 0 {
        downTotalPoints = strings.Split(points["total"], ",")[0]
        downRoundPoints = strings.Split(points["round"], ",")[0]
        rightTotalPoints = strings.Split(points["total"], ",")[1]
        rightRoundPoints = strings.Split(points["round"], ",")[1]
    } else {
        downTotalPoints = strings.Split(points["total"], ",")[1]
        downRoundPoints = strings.Split(points["round"], ",")[1]
        rightTotalPoints = strings.Split(points["total"], ",")[0]
        rightRoundPoints = strings.Split(points["round"], ",")[0]
    }

    return map[string]map[string]string{
        "total": {
            "down": downTotalPoints, "right": rightTotalPoints,
        },
        "currentRound": {
            "down": downRoundPoints, "right": rightRoundPoints,
        },
    }
}

// GetCenterCards Get center cards
func GetCenterCards(centerCards string, uIndex int) map[string]interface{} {
    result := make(map[string]interface{})
    for key, val := range strings.Split(centerCards, ",") {
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
func Matchmaking(ctx context.Context, client rueidis.Client, userId, gameId string) {
    time.Sleep(1 * time.Second)
    distributedCards := distributeCards()

    lastMoveTimestamp := strconv.FormatInt(time.Now().Unix(), 10)

    judgeCards, judge := chooseFirstJudge()

    redis.Matchmaking(ctx, client, distributedCards, userId, gameId, lastMoveTimestamp, judge, judgeCards)
}

func distributeCards() []string {
    localCards := make([]string, len(cards))
    copy(localCards, cards)

    rand.Shuffle(len(localCards), func(i, j int) { localCards[i], localCards[j] = localCards[j], localCards[i] })

    hands := make([][]string, 4)
    for i := range hands {
        hands[i] = []string{}
    }

    for i, card := range cards {
        player := i % 4
        hands[player] = append(hands[player], card)
    }

    result := make([]string, 4)

    for i := 0; i < 4; i++ {
        result[i] = `["` + strings.Join(hands[i], `","`) + `"]`
    }

    return result
}

func GetPlayerCards(gameCardsString string, uIndex int) [][]string {
    gameCards := make(map[int][]string)
    err := json.Unmarshal([]byte(gameCardsString), &gameCards)
    if err != nil {
        fmt.Println("Error unmarshalling:", err)
    }

    return chunkCards(gameCards, uIndex)
}

func chunkCards(gameCards map[int][]string, uIndex int) [][]string {
    var result [][]string
    var chunk []string
    for i, card := range gameCards[uIndex] {
        chunk = append(chunk, card)
        if (i+1)%5 == 0 {
            result = append(result, chunk)
            chunk = []string{}
        }
    }
    result = append(result, chunk)
    return result
}
