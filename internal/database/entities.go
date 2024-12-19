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
    Games     []Game `gorm:"foreignkey:PlayerId"`
    UpdatedAt time.Time
    JoinedAt  time.Time
}

type Game struct {
    Id         int64 `gorm:"primaryKey"`
    GameId     string
    PlayerId   int64
    CreatedAt  time.Time
    FinishedAt *time.Time
}
