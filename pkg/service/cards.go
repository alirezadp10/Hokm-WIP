package service

import (
    "context"
    "encoding/json"
    "fmt"
    "github.com/alirezadp10/hokm/internal/util/my_slice"
    "github.com/alirezadp10/hokm/pkg/repository"
    "github.com/redis/rueidis"
    "gorm.io/gorm"
    "math/rand"
    "strings"
)

type CardsService struct {
    sqlite    gorm.DB
    redis     rueidis.Client
    CardsRepo repository.CardsRepositoryContract
}

func NewCardsService(sqliteClient *gorm.DB, redisClient *rueidis.Client, repo repository.CardsRepositoryContract) *CardsService {
    return &CardsService{
        sqlite:    *sqliteClient,
        redis:     *redisClient,
        CardsRepo: repo,
    }
}

var Cards = []string{
    "01C", "02C", "03C", "04C", "05C", "06C", "07C", "08C", "09C", "10C", "JC", "QC", "KC", // Clubs
    "01D", "02D", "03D", "04D", "05D", "06D", "07D", "08D", "09D", "10D", "JD", "QD", "KD", // Diamonds
    "01H", "02H", "03H", "04H", "05H", "06H", "07H", "08H", "09H", "10H", "JH", "QH", "KH", // Hearts
    "01S", "02S", "03S", "04S", "05S", "06S", "07S", "08S", "09S", "10S", "JS", "QS", "KS", // Spades
}

func (s *CardsService) chunkCards(gameCards map[int][]string, uIndex int) [][]string {
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

func (s *CardsService) GetCardSuit(card string) string {
    if len(card) == 0 {
        return ""
    }
    return string(card[len(card)-1])
}

func (s *CardsService) DistributeCards() []string {
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

func (s *CardsService) GetPlayerCards(cards string, uIndex int) [][]string {
    gameCards := make(map[int][]string)
    err := json.Unmarshal([]byte(cards), &gameCards)
    if err != nil {
        fmt.Println("Error unmarshalling:", err)
    }

    return s.chunkCards(gameCards, uIndex)
}

func (s *CardsService) UpdateCenterCards(cards string, newCard string, uIndex int) string {
    centerCardsList := strings.Split(cards, ",")
    centerCardsList[uIndex] = newCard
    return strings.Join(centerCardsList, ",")
}

func (s *CardsService) UpdateUserCards(cards string, selectedCard string, uIndex int) []string {
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

func (s *CardsService) GetKingCards(cards string) []string {
    return strings.Split(cards, ",")
}

func (s *CardsService) SetTrump(ctx context.Context, gameID, trump, uIndex, lastMoveTimestamp string) error {
    return s.CardsRepo.SetTrump(ctx, gameID, trump, uIndex, lastMoveTimestamp)
}
