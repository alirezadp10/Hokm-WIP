package handler

import (
    "github.com/redis/rueidis"
    "gorm.io/gorm"
)

type Handler struct {
    sqlite *gorm.DB
    redis  rueidis.Client
}

func NewHandler(sqlite *gorm.DB, redis rueidis.Client) *Handler {
    return &Handler{sqlite: sqlite, redis: redis}
}
