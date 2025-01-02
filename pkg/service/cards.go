package service

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/alirezadp10/hokm/internal/util/my_slice"
	"github.com/alirezadp10/hokm/pkg/repository"
	"github.com/redis/rueidis"
	"gorm.io/gorm"
)

type CardsService struct {
	sqlite    *gorm.DB
	redis     *rueidis.Client
	CardsRepo repository.CardsRepositoryContract
}

func NewCardsService(sqliteClient *gorm.DB, redisClient *rueidis.Client, repo repository.CardsRepositoryContract) *CardsService {
	return &CardsService{
		sqlite:    sqliteClient,
		redis:     redisClient,
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

var rank = map[string]int{
	"02": 2, "03": 3, "04": 4, "05": 5, "06": 6, "07": 7, "08": 8, "09": 9, "10": 10, "J": 11, "Q": 12, "K": 13, "01": 14,
}

func (s *CardsService) GetPoints(pointsString string, uIndex int) map[string]map[string]string {
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
		"total": {"down": downTotalPoints, "right": rightTotalPoints},
		"round": {"down": downRoundPoints, "right": rightRoundPoints},
	}
}

func (s *CardsService) FindCardsWinner(centerCards, trump, leadSuit string) string {
	for _, card := range strings.Split(centerCards, ",") {
		if card == "" {
			return ""
		}
	}

	highestRank := -1
	winner := 0

	for i, card := range strings.Split(centerCards, ",") {
		cardSuit := s.GetCardSuit(card)      // Extract suit of the card
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

func (s *CardsService) UpdatePoints(pointsString, cardsWinnerString string) (string, string, string) {
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
