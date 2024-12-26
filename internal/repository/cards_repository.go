package repository

import (
    "github.com/redis/rueidis"
    "gorm.io/gorm"
)

type CardsRepository struct {
    sqlite *gorm.DB
    redis  rueidis.Client
}
