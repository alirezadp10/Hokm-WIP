package database

import (
    "time"
)

type Player struct {
    Id        int64 `gorm:"primaryKey"`
    ChatId    string
    FirstName string
    LastName  string
    Username  string
    Score     uint
    Avatar    string
    UpdatedAt time.Time
    JoinedAt  time.Time
}
