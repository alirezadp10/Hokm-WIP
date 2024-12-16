package cmd

import (
    "fmt"
    "github.com/alirezadp10/hokm/internal/database"
    "github.com/alirezadp10/hokm/internal/handler"
    "github.com/alirezadp10/hokm/internal/telegram"
    "github.com/labstack/echo/v4"
    "github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
    Use:   "serve",
    Short: "run http server",
    Run:   serve,
}

func init() {
    rootCmd.AddCommand(serveCmd)
}

func serve(cmd *cobra.Command, args []string) {
    db := database.GetNewConnection()
    go telegram.Start(db)

    e := echo.New()
    //e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
    //    return func(c echo.Context) error {
    //        c.Response().Header().Set("ngrok-skip-browser-warning", "true")
    //        return next(c)
    //    }
    //})
    e.Static("/assets", "assets")
    e.GET("/", handler.GetSplashPage)
    e.GET("/menu", handler.GetMenuPage)
    e.GET("/game", handler.GetGamePage)
    e.GET("/game/start", handler.GetGameId)
    e.GET("/game/:gameId", handler.GetGameData)
    e.GET("/game/:gameId/cards", handler.GetYourCards)
    e.POST("/game/:gameId/place", handler.PlaceCard)
    e.GET("/game/:gameId/refresh", handler.GetUpdate)
    fmt.Println("Server is running at 9090")
    e.Logger.Fatal(e.Start("0.0.0.0:9090"))
}
