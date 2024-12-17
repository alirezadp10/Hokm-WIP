package handler

import (
    "context"
    "github.com/redis/rueidis"
    "gorm.io/gorm"
)

type Handler struct {
    context          context.Context
    sqliteConnection *gorm.DB
    redisConnection  rueidis.Client
}

func NewHandler(sqlite *gorm.DB, redis rueidis.Client, context context.Context) *Handler {
    return &Handler{sqliteConnection: sqlite, redisConnection: redis, context: context}
}
