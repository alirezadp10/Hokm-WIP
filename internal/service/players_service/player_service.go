package players_service

import (
    "fmt"
    "strconv"
    "time"
)

func NextPlayerIndex(index int) int {
    return (index + 1) % 4
}

func GetNewTurn(turnStr string) string {
    turn, _ := strconv.Atoi(turnStr)
    return strconv.Itoa(NextPlayerIndex(turn))
}

func GetPlayersWithDirections(players []string, uIndex int) map[string]interface{} {
    return map[string]interface{}{
        "left":  map[string]string{"username": players[(4+(uIndex-1))%4]},
        "down":  map[string]string{"username": players[(4+(uIndex+0))%4]},
        "right": map[string]string{"username": players[(4+(uIndex+1))%4]},
        "up":    map[string]string{"username": players[(4+(uIndex+2))%4]},
    }
}

func GetDirection(pIndex, uIndex int) string {
    directions := []string{"down", "left", "up", "right"}
    diff := (4 + (uIndex - pIndex)) % 4 // Calculate relative direction
    return directions[diff]
}

func GetTurn(turnIndex string, uIndex int) string {
    if turnIndex == "" { // Handle empty turn index
        return ""
    }
    turnI, _ := strconv.Atoi(turnIndex) // Convert turn index to integer
    return GetDirection(turnI, uIndex)
}

func GetTimeRemained(lastMoveTimestampStr string) time.Duration {
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
