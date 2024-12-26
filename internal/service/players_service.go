package service

import (
    "fmt"
    "github.com/alirezadp10/hokm/internal/repository"
    "github.com/redis/rueidis"
    "gorm.io/gorm"
    "math/rand"
    "strconv"
    "strings"
    "time"
)

type PlayersService struct {
    sqlite      *gorm.DB
    redis       rueidis.Client
    PlayersRepo repository.PlayersRepository
}

func NewPlayersService(repo repository.PlayersRepository, sqlite *gorm.DB, redis rueidis.Client) *PlayersService {
    return &PlayersService{
        PlayersRepo: repo,
        sqlite:      sqlite,
        redis:       redis,
    }
}

func (s *PlayersService) NextPlayerIndex(index int) int {
    return (index + 1) % 4
}

func (s *PlayersService) GetNewTurn(turnStr string) string {
    turn, _ := strconv.Atoi(turnStr)
    return strconv.Itoa(s.NextPlayerIndex(turn))
}

func (s *PlayersService) GetPlayersWithDirections(players []string, uIndex int) map[string]interface{} {
    return map[string]interface{}{
        "left":  map[string]string{"username": players[(4+(uIndex-1))%4]},
        "down":  map[string]string{"username": players[(4+(uIndex+0))%4]},
        "right": map[string]string{"username": players[(4+(uIndex+1))%4]},
        "up":    map[string]string{"username": players[(4+(uIndex+2))%4]},
    }
}

func (s *PlayersService) GetDirection(pIndex, uIndex int) string {
    directions := []string{"down", "left", "up", "right"}
    diff := (4 + (uIndex - pIndex)) % 4 // Calculate relative direction
    return directions[diff]
}

func (s *PlayersService) GetTurn(turnIndex string, uIndex int) string {
    if turnIndex == "" { // Handle empty turn index
        return ""
    }
    turnI, _ := strconv.Atoi(turnIndex) // Convert turn index to integer
    return s.GetDirection(turnI, uIndex)
}

func (s *PlayersService) GetTimeRemained(lastMoveTimestampStr string) time.Duration {
    lastMoveTimestampInt, err := strconv.ParseInt(lastMoveTimestampStr, 10, 64) // Convert string to integer
    if err != nil {
        fmt.Println("Error parsing timestamp:", err)
        return 15 * time.Second // Default time if parsing fails
    }

    lastMoveTimestamp := time.Unix(lastMoveTimestampInt, 0) // Convert to time.Time
    timeElapsed := time.Since(lastMoveTimestamp)            // Time elapsed since the move

    timeRemaining := 15*time.Second - timeElapsed // Calculate remaining time

    if timeRemaining < 0 { // Ensure non-negative duration
        timeRemaining = 0
    }

    return timeRemaining.Round(time.Second) // Round to the nearest second
}

func (s *PlayersService) ChooseFirstKing() (string, string) {
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

func (s *PlayersService) GetKing(kingIndex string, uIndex int) string {
    kingI, _ := strconv.Atoi(kingIndex) // Convert king index to integer
    return s.GetDirection(kingI, uIndex)
}

func (s *PlayersService) GiveKing(roundWinnerStr, prevKingStr string) string {
    roundWinner, _ := strconv.Atoi(roundWinnerStr)
    prevKing, _ := strconv.Atoi(prevKingStr)
    if roundWinner%2 == prevKing%2 {
        return prevKingStr
    }
    return strconv.Itoa(s.NextPlayerIndex(roundWinner))
}

func (s *PlayersService) GetPlayersCenterCards(centerCards string, uIndex int) map[string]interface{} {
    result := make(map[string]interface{})
    for key, val := range strings.Split(centerCards, ",") {
        result[s.GetDirection(key, uIndex)] = val
    }
    return result
}
