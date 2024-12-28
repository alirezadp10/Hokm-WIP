package service

import (
    "encoding/json"
    "fmt"
    "github.com/alirezadp10/hokm/pkg/repository"
    "strconv"
    "strings"
)

type PointsService struct {
    PointsRepo repository.PointsRepository
}

func NewPointsService(repo *repository.PointsRepository) *PointsService {
    return &PointsService{
        PointsRepo: *repo,
    }
}

var rank = map[string]int{
    "02": 2, "03": 3, "04": 4, "05": 5, "06": 6, "07": 7, "08": 8, "09": 9, "10": 10, "J": 11, "Q": 12, "K": 13, "01": 14,
}

func (s *PointsService) GetPoints(pointsString string, uIndex int) map[string]map[string]string {
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

func (s *PointsService) FindCardsWinner(centerCards, trump, leadSuit string) string {
    for _, card := range strings.Split(centerCards, ",") {
        if card == "" {
            return ""
        }
    }

    highestRank := -1
    winner := 0

    for i, card := range strings.Split(centerCards, ",") {
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

func (s *PointsService) UpdatePoints(pointsString, cardsWinnerString string) (string, string, string) {
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
