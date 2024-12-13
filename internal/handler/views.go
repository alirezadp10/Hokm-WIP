package handler

import (
    "github.com/labstack/echo/v4"
    "net/http"
)

func GetSplashPage(c echo.Context) error {
    cookie := new(http.Cookie)
    cookie.Name = "userId"
    cookie.Value = c.QueryParam("userId")
    cookie.Path = "/"
    cookie.HttpOnly = true
    cookie.MaxAge = 3600
    c.SetCookie(cookie)
    return c.File("templates/splash.html")
}

func GetMenuPage(c echo.Context) error {
    return c.File("templates/menu.html")
}

func GetGamePage(c echo.Context) error {
    return c.File("templates/game.html")
}
