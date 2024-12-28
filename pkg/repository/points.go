package repository

import (
    "github.com/redis/rueidis"
    "gorm.io/gorm"
)

type PointsRepository struct {
    Sqlite *gorm.DB
    Redis  rueidis.Client
}

func NewPointsRepository(sqliteClient *gorm.DB, redisClient rueidis.Client) *PointsRepository {
    return &PointsRepository{
        Sqlite: sqliteClient,
        Redis:  redisClient,
    }
}
