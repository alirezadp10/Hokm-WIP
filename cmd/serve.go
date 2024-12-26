package cmd

import (
    "context"
    "fmt"
    "github.com/alirezadp10/hokm/internal/database/redis"
    "github.com/alirezadp10/hokm/internal/database/sqlite"
    "github.com/alirezadp10/hokm/internal/handler"
    "github.com/alirezadp10/hokm/internal/middleware"
    "github.com/alirezadp10/hokm/internal/repository"
    "github.com/alirezadp10/hokm/internal/service"
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
    go telegram.Start(ctx)

    gameRepository := repository.GameRepository{}
    gameService := service.NewGameService(gameRepository, sqliteClient, redisClient)

    cardsRepository := repository.CardsRepository{}
    cardsService := service.NewCardsService(cardsRepository, sqliteClient, redisClient)

    pointsRepository := repository.PointsRepository{}
    pointsService := service.NewPointsService(pointsRepository, sqliteClient, redisClient)

    playersRepository := repository.PlayersRepository{}
    playersService := service.NewPlayersService(playersRepository, sqliteClient, redisClient)

    h := handler.NewHandler(sqliteClient, redisClient, *gameService, *cardsService, *pointsService, *playersService)

    e := echo.New()
    t := &Template{templates: template.Must(template.ParseGlob("templates/*.html"))}
    e.Renderer = t

    e.Static("/assets", "assets")
    e.File("/favicon.ico", "assets/favicon.ico")
    e.GET("/", h.GetSplashPage)
    e.GET("/menu", h.GetMenuPage)
    e.GET("/game", h.GetGamePage)
    e.POST("/game/start", h.CreateGame, middleware.AuthMiddleware(sqliteClient))
    e.GET("/game/:gameID", h.GetGameInformation, middleware.AuthMiddleware(sqliteClient))
    e.POST("/game/:gameID/choose-trump", h.ChooseTrump, middleware.AuthMiddleware(sqliteClient))
    e.GET("/game/:gameID/cards", h.GetCards, middleware.AuthMiddleware(sqliteClient))
    e.POST("/game/:gameID/place", h.PlaceCard, middleware.AuthMiddleware(sqliteClient))
    e.GET("/game/:gameID/refresh", h.GetUpdate, middleware.AuthMiddleware(sqliteClient))

    fmt.Println("Server is running at 9090")
    e.Logger.Fatal(e.Start("0.0.0.0:9090"))
}
