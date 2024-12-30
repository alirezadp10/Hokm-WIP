package handlers

import (
    "context"
    "github.com/DATA-DOG/go-sqlmock"
    "github.com/alirezadp10/hokm/pkg/repository"
    "github.com/alirezadp10/hokm/pkg/service"
    "github.com/redis/rueidis/mock"
    "go.uber.org/mock/gomock"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "testing"
)

func TestPlaceCard(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    redisClient := mock.NewClient(ctrl)
    redisClient.EXPECT().Do(context.Background(), mock.Match("GET", "key")).Return(mock.Result(mock.RedisString("val")))

    sqlDB, sMock, _ := sqlmock.New()
    sqliteClient, _ := gorm.Open(sqlite.Dialector{Conn: sqlDB}, &gorm.Config{})
    sMock.ExpectQuery("SELECT * FROM `games` WHERE `id` = ?").WithArgs(1).WillReturnRows(
        sqlmock.NewRows([]string{"id"}).AddRow(1),
    )

    gameRepository := repository.NewGameRepository(sqliteClient, redisClient)
    gameService := service.NewGameService(gameRepository)

    cardsRepository := repository.NewCardsRepository(sqliteClient, redisClient)
    cardsService := service.NewCardsService(cardsRepository)

    pointsRepository := repository.NewPointsRepository(sqliteClient, redisClient)
    pointsService := service.NewPointsService(pointsRepository, *cardsService)

    playersRepository := repository.NewPlayersRepository(sqliteClient, redisClient)
    playersService := service.NewPlayersService(playersRepository)

    redisService := service.NewRedisService(redisClient, context.Background())

    NewHokmHandler(gameService, cardsService, pointsService, playersService, redisService)
}
