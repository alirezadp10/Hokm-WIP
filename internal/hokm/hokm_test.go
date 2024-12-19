package hokm

import (
    "testing"
)

func TestGetDirections(t *testing.T) {
    // Expected results
    expected := []map[string]string{
        {"left": "4", "down": "1", "right": "2", "up": "3"},
        {"left": "1", "down": "2", "right": "3", "up": "4"},
        {"left": "2", "down": "3", "right": "4", "up": "1"},
        {"left": "3", "down": "4", "right": "1", "up": "2"},
    }

    for i := 0; i < 4; i++ {
        // Call the function
        result := GetPlayersWithDirections([]string{"1", "2", "3", "4"}, i)

        // Validate results
        for direction, username := range expected[i] {
            if result[direction].(map[string]string)["username"] != username {
                t.Errorf("For direction %s, expected %s but got %s",
                    direction, username, result[direction].(map[string]string)["username"])
            }
        }
    }
}
