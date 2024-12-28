package main

import (
    "context"
    "fmt"
    "github.com/alirezadp10/hokm/pkg/api"
    "github.com/alirezadp10/hokm/pkg/api/handlers"
    "github.com/alirezadp10/hokm/pkg/api/middleware"
    "github.com/alirezadp10/hokm/pkg/database"
    "github.com/alirezadp10/hokm/pkg/repository"
    "github.com/alirezadp10/hokm/pkg/service"
    "github.com/joho/godotenv"
    "github.com/labstack/echo/v4"
)

func main() {
    _ = godotenv.Load()

    sqliteClient := database.GetNewSqliteConnection()
    redisClient := database.GetNewRedisConnection()

    gameRepository := repository.NewGameRepository(sqliteClient, redisClient)
    gameService := service.NewGameService(gameRepository)

    cardsRepository := repository.NewCardsRepository(sqliteClient, redisClient)
    cardsService := service.NewCardsService(cardsRepository)

    pointsRepository := repository.NewPointsRepository(sqliteClient, redisClient)
    pointsService := service.NewPointsService(pointsRepository, *cardsService)

    playersRepository := repository.NewPlayersRepository(sqliteClient, redisClient)
    playersService := service.NewPlayersService(playersRepository)

    redisService := service.NewRedisService(redisClient, context.Background())

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    go service.StartTelegram(ctx)

    hokmHandler := handlers.NewHokmHandler(gameService, cardsService, pointsService, playersService, redisService)

    authMiddleware := middleware.NewAuthMiddleware(*playersRepository)

    e := echo.New()

    api.SetupRouter(e, hokmHandler, authMiddleware)

    fmt.Println("Server is running at 9090")
    e.Start("0.0.0.0:9090")
}
