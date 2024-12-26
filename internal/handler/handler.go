package handler

import (
    "github.com/alirezadp10/hokm/internal/service"
    "github.com/redis/rueidis"
    "gorm.io/gorm"
)

type Handler struct {
    sqlite         *gorm.DB
    redis          rueidis.Client
    GameService    service.GameService
    CardsService   service.CardsService
    PlayersService service.PlayersService
    PointsService  service.PointsService
}

func NewHandler(sqlite *gorm.DB,
        redis rueidis.Client,
        gameService service.GameService,
        cardsService service.CardsService,
        pointsService service.PointsService,
        playersService service.PlayersService,
) *Handler {
    return &Handler{
        sqlite:         sqlite,
        redis:          redis,
        GameService:    gameService,
        CardsService:   cardsService,
        PointsService:  pointsService,
        PlayersService: playersService,
    }
}
