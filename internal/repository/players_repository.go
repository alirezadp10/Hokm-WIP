package repository

import (
    "github.com/redis/rueidis"
    "gorm.io/gorm"
)

type PlayersRepository struct {
    sqlite *gorm.DB
    redis  rueidis.Client
}
