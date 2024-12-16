package database

import (
    "time"
)

type Player struct {
    Id        int64 `gorm:"primaryKey"`
    ChatId    int64
    FirstName string
    LastName  string
    Username  string
    Score     uint
    UpdatedAt time.Time
    JoinedAt  time.Time
}

type Game struct {
    GameId     string `gorm:"primaryKey"`
    PlayerId   int64  `gorm:"foreignKey"`
    CreatedAt  time.Time
    FinishedAt time.Time
}
