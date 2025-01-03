package repository

import (
	"context"

	"github.com/alirezadp10/hokm/pkg/api/request"
	"github.com/alirezadp10/hokm/pkg/model"
)

type PlaceCardParams struct {
	GameId            string
	Card              string
	CenterCards       string
	LeadSuit          string
	CardsWinner       string
	Points            string
	Turn              string
	King              string
	WasKingChanged    string
	LastMoveTimestamp string
	Trump             string
	IsItNewRound      string
	Cards             map[int][]string
	PlayerIndex       int
}

type CardsRepositoryContract interface {
	SetTrump(ctx context.Context, gameID, trump, uIndex, lastMoveTimestamp string) error
	PlaceCard(ctx context.Context, params PlaceCardParams) error
}

type GameRepositoryContract interface {
	GetGameInformation(ctx context.Context, gameID string) (map[string]interface{}, error)
	Matchmaking(ctx context.Context, cards []string, username, gameID, lastMoveTimestamps, king, kingCards string)
	RemovePlayerFromWaitingList(ctx context.Context, key, username string)
}

type PlayersRepositoryContract interface {
	CheckPlayerExistence(username string) bool
	SavePlayer(user request.User, chatId int64) (*model.Player, error)
	AddPlayerToGame(username, gameID string) (*model.Game, error)
	DoesPlayerBelongToGame(username, gameID string) (bool, error)
	DoesPlayerHaveAnyActiveGame(username string) (*string, bool)
	HasGameFinished(gameID string) (bool, error)
}
