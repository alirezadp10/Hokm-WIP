package api

import (
    "github.com/alirezadp10/hokm/pkg/api/handlers"
    "github.com/alirezadp10/hokm/pkg/api/middleware"
    "github.com/labstack/echo/v4"
    "html/template"
    "io"
)

type Template struct {
    templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
    return t.templates.ExecuteTemplate(w, name, data)
}

func SetupRouter(e *echo.Echo, h *handlers.HokmHandler, authMiddleware *middleware.AuthMiddleware) {
    t := &Template{templates: template.Must(template.ParseGlob("templates/*.html"))}
    e.Renderer = t

    e.Static("/assets", "assets")
    e.File("/favicon.ico", "assets/favicon.ico")
    e.GET("/", h.GetSplashPage)
    e.GET("/menu", h.GetMenuPage)
    e.GET("/game", h.GetGamePage)
    e.POST("/game/start", h.CreateGame, authMiddleware.Handle)
    e.GET("/game/:gameID", h.GetGameInformation, authMiddleware.Handle)
    e.POST("/game/:gameID/choose-trump", h.ChooseTrump, authMiddleware.Handle)
    e.GET("/game/:gameID/cards", h.GetCards, authMiddleware.Handle)
    e.POST("/game/:gameID/place", h.PlaceCard, authMiddleware.Handle)
    e.GET("/game/:gameID/refresh", h.GetUpdate, authMiddleware.Handle)
}
