package hokm

import (
    "testing"
)

func TestSetKingCards(t *testing.T) {
    result := SetKingCards()

    if len(result) == 0 {
        t.Error("Result is empty; expected at least one card to be assigned")
    }

    for _, entry := range result {
        cardMap, ok := entry.(map[string]interface{})
        if !ok {
            t.Errorf("Invalid format: expected map[string]interface{}, got %T", entry)
        }

        if _, ok = cardMap["player"]; !ok {
            t.Error("Missing 'direction' key in result entry")
        }

        if _, ok = cardMap["card"]; !ok {
            t.Error("Missing 'card' key in result entry")
        }
    }

    lastCard := result[len(result)-1].(map[string]interface{})["card"].(string)
    if lastCard[:2] != "01" {
        t.Errorf("Expected the last card to contain '01', but got %s", lastCard)
    }
}

func TestGetDirections(t *testing.T) {
    result := GetPlayersWithDirections([]string{"1", "2", "3", "4"}, 0)
    if result["left"] != 4 || result["down"] != 1 || result["right"] != 2 || result["up"] != 3 {
        t.Error("Result is wrong")
    }

    result = GetPlayersWithDirections([]string{"1", "2", "3", "4"}, 1)
    if result["left"] != 1 || result["down"] != 2 || result["right"] != 3 || result["up"] != 4 {
        t.Error("Result is wrong")
    }

    result = GetPlayersWithDirections([]string{"1", "2", "3", "4"}, 2)
    if result["left"] != 2 || result["down"] != 3 || result["right"] != 4 || result["up"] != 1 {
        t.Error("Result is wrong")
    }

    result = GetPlayersWithDirections([]string{"1", "2", "3", "4"}, 3)
    if result["left"] != 3 || result["down"] != 4 || result["right"] != 1 || result["up"] != 2 {
        t.Error("Result is wrong")
    }
}
