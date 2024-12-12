package handler

import (
    "github.com/labstack/echo/v4"
    "net/http"
    "time"
)

func GetGameData(c echo.Context) error {
    //response := map[string]interface{}{
    //    "players": map[string]interface{}{
    //        "up": map[string]interface{}{
    //            "name": "ali",
    //        },
    //        "right": map[string]interface{}{
    //            "name": "maryam",
    //        },
    //        "left": map[string]interface{}{
    //            "name": "sara",
    //        },
    //        "down": map[string]interface{}{
    //            "name": "mammad",
    //        },
    //    },
    //    "points": map[string]interface{}{
    //        "total":        map[string]interface{}{"right": 3, "down": 2},
    //        "currentRound": map[string]interface{}{"right": 3, "down": 2},
    //    },
    //    "firstKingDeterminationCards": []interface{}{},
    //    "centerCards": map[string]interface{}{
    //        "up":    "2S",
    //        "right": "7S",
    //    },
    //    "currentTurn":  "down",
    //    "timeRemained": 14,
    //    "yourCards":    []interface{}{"3H", "3H", "3S", "3S", "4C"},
    //    "judge":        "down",
    //    "trump":        "heart",
    //}
    response := map[string]interface{}{
        "players": map[string]interface{}{
            "up": map[string]interface{}{
                "name": "ali",
            },
            "right": map[string]interface{}{
                "name": "maryam",
            },
            "left": map[string]interface{}{
                "name": "sara",
            },
            "down": map[string]interface{}{
                "name": "mammad",
            },
        },
        "points": map[string]interface{}{
            "total":        map[string]interface{}{"right": 3, "down": 2},
            "currentRound": map[string]interface{}{"right": 3, "down": 2},
        },
        "firstKingDeterminationCards": []interface{}{
            map[string]interface{}{"direction": "up", "card": "2S"},
            map[string]interface{}{"direction": "right", "card": "3C"},
            map[string]interface{}{"direction": "down", "card": "3H"},
            map[string]interface{}{"direction": "left", "card": "3S"},
            map[string]interface{}{"direction": "up", "card": "AC"},
        },
        "centerCards":  map[string]interface{}{"up": "2H", "left": "3H", "right": "4C"},
        "currentTurn":  "right",
        "timeRemained": 14,
        "yourCards":    []interface{}{"3H", "3H", "3S", "3S", "4C"},
        "judge":        "up",
        "trump":        "heart",
    }
    return c.JSON(http.StatusOK, response)
}

func GetYourCards(c echo.Context) error {
    time.Sleep(2 * time.Second)
    response := map[string]interface{}{
        "trump": "heart",
        "cards": []interface{}{
            []interface{}{
                "3C",
                "3H",
                "3S",
                "8S",
                "9D",
            },
            []interface{}{
                "AC",
                "AH",
                "2S",
                "6S",
                "2D",
            },
            []interface{}{
                "JS",
                "KH",
                "QD",
            },
        },
    }
    return c.JSON(http.StatusOK, response)
}

func PlaceCard(c echo.Context) error {
    response := map[string]interface{}{
        "points": map[string]interface{}{
            "total":        map[string]interface{}{"right": 4, "down": 2},
            "currentRound": map[string]interface{}{"right": 0, "down": 3},
        },
        "currentTurn":       "down",
        "timeRemained":      14,
        "judge":             "right",
        "whoHasWonTheCards": "up",
        "whoHasWonTheRound": nil,
        "whoHasWonTheGame":  nil,
        "wasKingChanged":    false,
        //"trumpDeterminationCards": []interface{}{"3H", "3H", "3S", "3S", "4C"},
        "trumpDeterminationCards": nil,
    }
    return c.JSON(http.StatusOK, response)
}

func GetUpdate(c echo.Context) error {
    time.Sleep(2 * time.Second)
    response := map[string]interface{}{
        "lastMove": map[string]interface{}{
            "from": "right",
            "card": "3C",
        },
        "centerCards": map[string]interface{}{"up": "2H", "left": "3H"},
        "points": map[string]interface{}{
            "total":        map[string]interface{}{"right": 4, "down": 2},
            "currentRound": map[string]interface{}{"right": 0, "down": 3},
        },
        "currentTurn":       "down",
        "timeRemained":      14,
        "judge":             "up",
        "whoHasWonTheCards": nil,
        "whoHasWonTheRound": nil,
        "whoHasWonTheGame":  nil,
        "wasKingChanged":    false,
    }
    return c.JSON(http.StatusOK, response)
}
