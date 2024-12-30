package handlers

import (
    "context"
    "github.com/DATA-DOG/go-sqlmock"
    "github.com/alirezadp10/hokm/pkg/database"
    "github.com/alirezadp10/hokm/pkg/repository"
    "github.com/alirezadp10/hokm/pkg/service"
    "github.com/labstack/echo/v4"
    "github.com/redis/rueidis"
    "github.com/redis/rueidis/mock"
    "github.com/stretchr/testify/assert"
    "go.uber.org/mock/gomock"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestPlaceCard(t *testing.T) {
    database.GetNewSqliteConnection()

    sqlDB, _, _ := sqlmock.New()
    sqliteClient, _ := gorm.Open(sqlite.Dialector{Conn: sqlDB}, &gorm.Config{})
    //sMock.ExpectQuery("SELECT * FROM `games` WHERE `id` = ?").WithArgs(1).WillReturnRows(
    //    sqlmock.NewRows([]string{"id"}).AddRow(1),
    //)

    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    redisClient := mock.NewClient(ctrl)
    //redisClient.EXPECT().Do(context.Background(), mock.Match("GET", "key")).Return(mock.Result(mock.RedisString("val")))

    e := echo.New()

    h := setupTestServer(sqliteClient, redisClient)

    // Create a request to test the GET /example endpoint
    req := httptest.NewRequest(http.MethodGet, "/game/12345/place", nil)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)
    c.Set("username", "foobar")

    // Call the handler
    if assert.NoError(t, h.PlaceCard(c)) {
        assert.Equal(t, http.StatusOK, rec.Code)
        assert.JSONEq(t, `{"message": "hi"}`, rec.Body.String())
    }
}

func setupTestServer(sqliteClient *gorm.DB, redisClient rueidis.Client) *HokmHandler {
    gameRepository := repository.NewGameRepository(sqliteClient, redisClient)
    gameService := service.NewGameService(gameRepository)

    cardsRepository := repository.NewCardsRepository(sqliteClient, redisClient)
    cardsService := service.NewCardsService(cardsRepository)

    pointsRepository := repository.NewPointsRepository(sqliteClient, redisClient)
    pointsService := service.NewPointsService(pointsRepository, *cardsService)

    playersRepository := repository.NewPlayersRepository(sqliteClient, redisClient)
    playersService := service.NewPlayersService(playersRepository)

    redisService := service.NewRedisService(redisClient, context.Background())

    hokmHandler := NewHokmHandler(gameService, cardsService, pointsService, playersService, redisService)

    return hokmHandler
}
