package handler

import (
    "encoding/json"
    "fmt"
    "github.com/alirezadp10/hokm/internal/database/sqlite"
    "github.com/alirezadp10/hokm/internal/request"
    "github.com/alirezadp10/hokm/internal/utils/crypto"
    "github.com/labstack/echo/v4"
    "net/http"
    "strconv"
)

func (h *Handler) GetSplashPage(c echo.Context) error {
    return c.Render(200, "splash.html", nil)
}

func (h *Handler) GetMenuPage(c echo.Context) error {
    var user request.User
    err := json.Unmarshal([]byte(c.QueryParam("user")), &user)
    if err != nil {
        fmt.Println("Error unmarshalling user JSON:", err)
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "unauthorized"})
    }

    encryptedUsername, err := crypto.Encrypt(user.Username)
    if err != nil {
        fmt.Println("Error in encryption:", err)
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "unauthorized"})
    }

    chatInstance, _ := strconv.ParseInt(c.QueryParam("chat_instance"), 10, 64)

    _, _ = sqlite.SavePlayer(h.sqlite, user, chatInstance)

    return c.Render(200, "menu.html", map[string]interface{}{
        "userReferenceKey": encryptedUsername,
    })
}

func (h *Handler) GetGamePage(c echo.Context) error {
    return c.Render(200, "game.html", map[string]interface{}{
        "userReferenceKey": c.QueryParam("user_id"),
        "gameId":           c.QueryParam("game_id"),
    })
}
