package handler

import (
    "github.com/labstack/echo/v4"
)

func (h *Handler) GetSplashPage(c echo.Context) error {
    return c.Render(200, "splash.html", map[string]interface{}{
        "userReferenceKey": c.QueryParam("user_id"),
    })
}

func (h *Handler) GetMenuPage(c echo.Context) error {
    return c.Render(200, "menu.html", map[string]interface{}{
        "userReferenceKey": c.QueryParam("user_id"),
    })
}

func (h *Handler) GetGamePage(c echo.Context) error {
    return c.Render(200, "game.html", map[string]interface{}{
        "userReferenceKey": c.QueryParam("user_id"),
        "gameId":           c.QueryParam("game_id"),
    })
}
