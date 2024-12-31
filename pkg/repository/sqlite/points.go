package sqliteRepo

import (
    "gorm.io/gorm"
)

type PointsRepository struct {
    sqlite *gorm.DB
}

func NewPointsRepository(sqliteClient *gorm.DB) *PointsRepository {
    return &PointsRepository{
        sqlite: sqliteClient,
    }
}
