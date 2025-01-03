package model

import "time"

type Game struct {
	Id         int64 `gorm:"primaryKey"`
	GameId     string
	PlayerId   int64
	CreatedAt  time.Time
	FinishedAt *time.Time
}
