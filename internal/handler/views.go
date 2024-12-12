package handler

import (
    "github.com/labstack/echo/v4"
)

func GetSplashPage(c echo.Context) error {
    return c.File("templates/splash.html")
}

func GetMenuPage(c echo.Context) error {
    return c.File("templates/menu.html")
}

func GetGamePage(c echo.Context) error {
    return c.File("templates/game.html")
}
