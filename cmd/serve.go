package cmd

import (
    "context"
    "fmt"
    "github.com/alirezadp10/hokm/internal/database/redis"
    "github.com/alirezadp10/hokm/internal/database/sqlite"
    "github.com/alirezadp10/hokm/internal/handler"
    "github.com/alirezadp10/hokm/internal/middleware"
    "github.com/alirezadp10/hokm/internal/telegram"
    "github.com/labstack/echo/v4"
    "github.com/spf13/cobra"
    "html/template"
    "io"
)

var serveCmd = &cobra.Command{
    Use:   "serve",
    Short: "run http server",
    Run:   serve,
}

func init() {
    rootCmd.AddCommand(serveCmd)
}

type Template struct {
    templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
    return t.templates.ExecuteTemplate(w, name, data)
}

func serve(cmd *cobra.Command, args []string) {
    sqliteClient := sqlite.GetNewConnection()
    redisClient := redis.GetNewConnection()

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    go telegram.Start(ctx, sqliteClient)

    h := handler.NewHandler(sqliteClient, redisClient)

    e := echo.New()

    t := &Template{templates: template.Must(template.ParseGlob("templates/*.html"))}
    e.Renderer = t

    e.Static("/assets", "assets")
    e.File("/favicon.ico", "assets/favicon.ico")
    e.GET("/", h.GetSplashPage)
    e.GET("/menu", h.GetMenuPage)
    e.GET("/game", h.GetGamePage)
    e.POST("/game/start", h.GetGameId, middleware.AuthMiddleware(sqliteClient))
    e.GET("/game/:gameId", h.GetGameData, middleware.AuthMiddleware(sqliteClient))
    e.POST("/game/:gameId/choose-trump", h.ChooseTrump, middleware.AuthMiddleware(sqliteClient))
    e.GET("/game/:gameId/cards", h.GetYourCards, middleware.AuthMiddleware(sqliteClient))
    e.POST("/game/:gameId/place", h.PlaceCard, middleware.AuthMiddleware(sqliteClient))
    e.GET("/game/:gameId/refresh", h.GetUpdate, middleware.AuthMiddleware(sqliteClient))

    fmt.Println("Server is running at 9090")
    e.Logger.Fatal(e.Start("0.0.0.0:9090"))
}
