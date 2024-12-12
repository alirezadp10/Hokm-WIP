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
    db := database.Connection()

    _ = db.AutoMigrate(&database.Player{})

    go telegram.Start(db)

    e := echo.New()

    e.Static("/assets", "assets")

    e.GET("/", handler.GetSplashPage)

    e.GET("/menu", handler.GetMenuPage)

    e.GET("/game", handler.GetGamePage)

    e.GET("/game/:gameId", handler.GetGameData)

    e.GET("/game/:gameId/cards", handler.GetYourCards)

    e.POST("/game/:gameId/place", handler.PlaceCard)

    e.GET("/game/:gameId/refresh", handler.GetUpdate)

    fmt.Println("Server is running at 7070")

    err := e.Start("localhost:7070")

    if err != nil {
        fmt.Println("Error starting server:", err)
    }
}
