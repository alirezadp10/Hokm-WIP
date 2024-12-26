package repository

import (
    "github.com/redis/rueidis"
    "gorm.io/gorm"
)

type PointsRepository struct {
    sqlite *gorm.DB
    redis  rueidis.Client
}
