package redisRepo

import (
    "github.com/redis/rueidis"
)

type PointsRepository struct {
    redis *rueidis.Client
}

func NewPointsRepository(redisClient *rueidis.Client) *PointsRepository {
    return &PointsRepository{
        redis: redisClient,
    }
}
