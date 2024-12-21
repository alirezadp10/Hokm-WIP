package hokm

import (
    "context"
    "encoding/json"
    "fmt"
    "github.com/alirezadp10/hokm/internal/database/redis"
    "github.com/alirezadp10/hokm/internal/utils/my_slice"
    "github.com/redis/rueidis"
    "math/rand"
    "strconv"
    "strings"
    "time"
)

// Cards Global cards list representing all cards in a deck
var Cards = []string{
    "01C", "02C", "03C", "04C", "05C", "06C", "07C", "08C", "09C", "10C", "JC", "QC", "KC", // Clubs
    "01D", "02D", "03D", "04D", "05D", "06D", "07D", "08D", "09D", "10D", "JD", "QD", "KD", // Diamonds
    "01H", "02H", "03H", "04H", "05H", "06H", "07H", "08H", "09H", "10H", "JH", "QH", "KH", // Hearts
    "01S", "02S", "03S", "04S", "05S", "06S", "07S", "08S", "09S", "10S", "JS", "QS", "KS", // Spades
}

// Map to determine card rank
var rank = map[string]int{
    "02": 2, "03": 3, "04": 4, "05": 5, "06": 6, "07": 7, "08": 8, "09": 9, "10": 10, "J": 11, "Q": 12, "K": 13, "01": 14,
}

// chooseFirstKing selects the first player with a king card and returns the cards sequence and the player's index
func chooseFirstKing() (string, string) {
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

// GetKingCards splits a comma-separated king cards string into a slice
func GetKingCards(cards string) []string {
    return strings.Split(cards, ",")
}

// GetKing returns the direction of the king player relative to the current user
func GetKing(kingIndex string, uIndex int) string {
    kingI, _ := strconv.Atoi(kingIndex) // Convert king index to integer
    return GetDirection(kingI, uIndex)
}

// GetTurn returns the direction of the player whose turn it is
func GetTurn(turnIndex string, uIndex int) string {
    if turnIndex == "" { // Handle empty turn index
        return ""
    }
    turnI, _ := strconv.Atoi(turnIndex) // Convert turn index to integer
    return GetDirection(turnI, uIndex)
}

// GetTimeRemained calculates the time remaining for the player's turn
func GetTimeRemained(lastMoveTimestampStr string) time.Duration {
    //fmt.Println("lastMoveTimestampStr")
    //fmt.Println(lastMoveTimestampStr)

    lastMoveTimestampInt, err := strconv.ParseInt(lastMoveTimestampStr, 10, 64) // Convert string to integer
    if err != nil {
        fmt.Println("Error parsing timestamp:", err)
        return 15 * time.Second // Default time if parsing fails
    }

    //fmt.Println("lastMoveTimestampInt")
    //fmt.Println(lastMoveTimestampInt)

    lastMoveTimestamp := time.Unix(lastMoveTimestampInt, 0) // Convert to time.Time
    timeElapsed := time.Since(lastMoveTimestamp)            // Time elapsed since the move

    //fmt.Println("timeElapsed")
    //fmt.Println(timeElapsed)

    timeRemaining := 15*time.Second - timeElapsed // Calculate remaining time

    //fmt.Println("timeRemaining")
    //fmt.Println(timeRemaining)

    if timeRemaining < 0 { // Ensure non-negative duration
        timeRemaining = 0
    }

    x := timeRemaining.Round(time.Second) // Round to the nearest second

    //fmt.Println(x)

    return x
}

// GetPlayersWithDirections maps usernames to relative directions
func GetPlayersWithDirections(players []string, uIndex int) map[string]interface{} {
    return map[string]interface{}{
        "left":  map[string]string{"username": players[(4+(uIndex-1))%4]},
        "down":  map[string]string{"username": players[(4+(uIndex+0))%4]},
        "right": map[string]string{"username": players[(4+(uIndex+1))%4]},
        "up":    map[string]string{"username": players[(4+(uIndex+2))%4]},
    }
}

// GetPoints processes the points string and organizes them by player direction
func GetPoints(pointsString string, uIndex int) map[string]map[string]string {
    points := make(map[string]string)
    err := json.Unmarshal([]byte(pointsString), &points)
    if err != nil {
        fmt.Println("Error unmarshalling:", err)
    }

    var downTotalPoints, rightTotalPoints, downRoundPoints, rightRoundPoints string
    if uIndex%2 == 0 { // Determine point order based on user index
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
        "total":        {"down": downTotalPoints, "right": rightTotalPoints},
        "currentRound": {"down": downRoundPoints, "right": rightRoundPoints},
    }
}

// GetCenterCards processes the center cards string and maps them to directions
func GetCenterCards(centerCards string, uIndex int) map[string]interface{} {
    result := make(map[string]interface{})
    for key, val := range strings.Split(centerCards, ",") {
        result[GetDirection(key, uIndex)] = val
    }
    return result
}

// GetDirection determines the relative direction of a player
func GetDirection(pIndex, uIndex int) string {
    directions := []string{"down", "left", "up", "right"}
    diff := (4 + (uIndex - pIndex)) % 4 // Calculate relative direction
    return directions[diff]
}

// Matchmaking assigns players to a game and initializes game data in Redis
func Matchmaking(ctx context.Context, client rueidis.Client, userId, gameId string) {
    time.Sleep(1 * time.Second)                                   // Simulate delay for matchmaking
    distributedCards := distributeCards()                         // Distribute cards among players
    lastMoveTimestamp := strconv.FormatInt(time.Now().Unix(), 10) // Record timestamp
    kingCards, king := chooseFirstKing()                          // Determine king cards and player
    redis.Matchmaking(ctx, client, distributedCards, userId, gameId, lastMoveTimestamp, king, kingCards)
}

// distributeCards shuffles and deals cards to 4 players
func distributeCards() []string {
    localCards := make([]string, len(Cards))
    copy(localCards, Cards)
    rand.Shuffle(len(localCards), func(i, j int) { localCards[i], localCards[j] = localCards[j], localCards[i] })

    hands := make([][]string, 4) // Initialize hands for 4 players
    for i, card := range localCards {
        player := i % 4 // Determine player index
        hands[player] = append(hands[player], card)
    }

    result := make([]string, 4)
    for i := 0; i < 4; i++ {
        result[i] = `["` + strings.Join(hands[i], `","`) + `"]` // Format hands as JSON strings
    }

    return result
}

// GetPlayerCards parses and chunks the cards for a specific player
func GetPlayerCards(cards string, uIndex int) [][]string {
    gameCards := make(map[int][]string)
    err := json.Unmarshal([]byte(cards), &gameCards)
    if err != nil {
        fmt.Println("Error unmarshalling:", err)
    }

    return chunkCards(gameCards, uIndex)
}

// chunkCards divides a player's cards into groups of 5
func chunkCards(gameCards map[int][]string, uIndex int) [][]string {
    var result [][]string
    var chunk []string
    for i, card := range gameCards[uIndex] {
        chunk = append(chunk, card)
        if (i+1)%5 == 0 { // Create a new chunk every 5 cards
            result = append(result, chunk)
            chunk = []string{}
        }
    }
    result = append(result, chunk) // Add remaining cards
    return result
}

func UpdateCenterCards(cards string, newCard string, uIndex int) string {
    centerCardsList := strings.Split(cards, ",")
    centerCardsList[uIndex] = newCard
    return strings.Join(centerCardsList, ",")
}

func FindCardsWinner(cards, trump, leadSuit string) string {
    centerCards := strings.Split(cards, ",")

    for _, card := range centerCards {
        if card == "" {
            return ""
        }
    }

    highestRank := -1
    winner := 0

    for i, card := range centerCards {
        cardSuit := GetCardSuit(card)        // Extract suit of the card
        cardRank := rank[card[:len(card)-1]] // Get the rank of the card (number or face)

        // Check if the card is a trump
        if cardSuit == trump {
            cardRank += 100 // Increase trump card rank to always beat non-trump cards
        } else if cardSuit != leadSuit {
            continue // Skip non-lead-suit and non-trump cards
        }

        // Update the winner if this card has a higher rank
        if cardRank > highestRank {
            highestRank = cardRank
            winner = i
        }
    }

    return strconv.Itoa(winner)
}

func GetCardSuit(card string) string {
    if len(card) == 0 {
        return ""
    }
    return string(card[len(card)-1])
}

func UpdatePoints(pointsString, cardsWinnerString string) (string, string, string) {
    points := make(map[string]string)
    err := json.Unmarshal([]byte(pointsString), &points)
    if err != nil {
        fmt.Println("Error unmarshalling:", err)
    }

    cardsWinner, _ := strconv.Atoi(cardsWinnerString)

    oldRoundPoints, _ := strconv.Atoi(strings.Split(points["round"], ",")[cardsWinner%2])
    roundPoints := strings.Split(points["round"], ",")
    roundPoints[cardsWinner%2] = strconv.Itoa(oldRoundPoints + 1)
    points["round"] = strings.Join(roundPoints, ",")

    roundWinner := ""
    oldTotalPoints, _ := strconv.Atoi(strings.Split(points["total"], ",")[cardsWinner%2])
    if oldRoundPoints+1 == 7 {
        points["round"] = "0,0"
        totalPoints := strings.Split(points["total"], ",")
        totalPoints[cardsWinner%2] = strconv.Itoa(oldTotalPoints + 1)
        points["total"] = strings.Join(totalPoints, ",")
        roundWinner = strconv.Itoa(cardsWinner % 2)
    }

    gameWinner := ""
    if oldTotalPoints+1 == 7 {
        totalPoints := strings.Split(points["total"], ",")
        totalPoints[cardsWinner%2] = strconv.Itoa(oldTotalPoints + 1)
        points["total"] = strings.Join(totalPoints, ",")
        gameWinner = strconv.Itoa(cardsWinner % 2)
    }

    pointsStr, _ := json.Marshal(points)

    return string(pointsStr), roundWinner, gameWinner
}

func GiveKing(roundWinnerStr, prevKingStr string) string {
    roundWinner, _ := strconv.Atoi(roundWinnerStr)
    prevKing, _ := strconv.Atoi(prevKingStr)
    if roundWinner%2 == prevKing%2 {
        return prevKingStr
    }
    return strconv.Itoa(NextPlayerIndex(roundWinner))
}

func NextPlayerIndex(index int) int {
    return (index + 1) % 4
}

func GetNewTurn(turnStr string) string {
    turn, _ := strconv.Atoi(turnStr)
    return strconv.Itoa(NextPlayerIndex(turn))
}

func UpdateUserCards(cards string, selectedCard string, uIndex int) []string {
    userCards := make(map[int][]string)
    err := json.Unmarshal([]byte(cards), &userCards)
    if err != nil {
        fmt.Println("Error unmarshalling:", err)
    }
    userCards[uIndex] = my_slice.Remove(selectedCard, userCards[uIndex])
    result := make([]string, 4)
    for i := 0; i < 4; i++ {
        result[i] = `["` + strings.Join(userCards[i], `","`) + `"]` // Format hands as JSON strings
    }
    return result
}
