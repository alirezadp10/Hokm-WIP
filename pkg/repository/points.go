package repository

import (
    "github.com/redis/rueidis"
    "gorm.io/gorm"
)

type PointsRepository struct {
    sqlite *gorm.DB
    redis  *rueidis.Client
}

func NewPointsRepository(sqliteClient *gorm.DB, redisClient *rueidis.Client) *PointsRepository {
    return &PointsRepository{
        sqlite: sqliteClient,
        redis:  redisClient,
    }
}
