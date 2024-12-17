package handler

import (
    "github.com/redis/rueidis"
    "gorm.io/gorm"
)

type Handler struct {
    sqliteConnection *gorm.DB
    redisConnection  rueidis.Client
}

func NewHandler(sqlite *gorm.DB, redis rueidis.Client) *Handler {
    return &Handler{sqliteConnection: sqlite, redisConnection: redis}
}
