package main

import (
	"fmt"

	"github.com/alirezadp10/hokm/pkg/api"
	"github.com/alirezadp10/hokm/pkg/api/handlers"
	"github.com/alirezadp10/hokm/pkg/api/middleware"
	"github.com/alirezadp10/hokm/pkg/database"
	redisRepo "github.com/alirezadp10/hokm/pkg/repository/redis"
	sqliteRepo "github.com/alirezadp10/hokm/pkg/repository/sqlite"
	"github.com/alirezadp10/hokm/pkg/service"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	_ = godotenv.Load()

	sqliteClient := database.GetNewSqliteConnection()
	redisClient := database.GetNewRedisConnection()

	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()
	// go service.StartTelegram(ctx)

	gameRepository := redisRepo.NewGameRepository(redisClient)
	gameService := service.NewGameService(sqliteClient, redisClient, gameRepository)

	cardsRepository := redisRepo.NewCardsRepository(redisClient)
	cardsService := service.NewCardsService(sqliteClient, redisClient, cardsRepository)

	playersRepository := sqliteRepo.NewPlayersRepository(sqliteClient)
	playersService := service.NewPlayersService(sqliteClient, redisClient, playersRepository)

	hokmHandler := handlers.NewHokmHandler(sqliteClient, redisClient, gameService, cardsService, playersService)

	authMiddleware := middleware.NewAuthMiddleware(playersRepository)

	e := echo.New()

	api.SetupRouter(e, hokmHandler, authMiddleware)

	fmt.Println("Server is running at 9090")
	e.Start("0.0.0.0:9090")
}
