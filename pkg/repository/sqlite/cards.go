package sqliteRepo

import (
    "context"
    _ "embed"
    "github.com/alirezadp10/hokm/pkg/repository"
    "gorm.io/gorm"
)

var _ repository.CardsRepositoryContract = &CardsRepository{}

type CardsRepository struct {
    sqlite gorm.DB
}

func NewCardsRepository(sqliteClient *gorm.DB) *CardsRepository {
    return &CardsRepository{
        sqlite: *sqliteClient,
    }
}

func (c CardsRepository) SetTrump(ctx context.Context, gameID, trump, uIndex, lastMoveTimestamp string) error {
    //TODO implement me
    panic("implement me")
}

func (c CardsRepository) PlaceCard(ctx context.Context, params repository.PlaceCardParams) error {
    //TODO implement me
    panic("implement me")
}
