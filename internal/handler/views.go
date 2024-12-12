package handler

import (
    "github.com/labstack/echo/v4"
)

func GetSplashPage(c echo.Context) error {
    userId := c.QueryParam("userId")
    return c.File("templates/splash.html?userId=" + userId)
}

func GetMenuPage(c echo.Context) error {
    userId := c.QueryParam("userId")
    return c.File("templates/menu.html?userId=" + userId)
}

func GetGamePage(c echo.Context) error {
    userId := c.QueryParam("userId")
    gameId := c.QueryParam("gameId")
    return c.File("templates/game.html?userId=" + userId + "&gameId=" + gameId)
}
