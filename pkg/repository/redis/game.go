package redisRepo

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/alirezadp10/hokm/pkg/repository"
	"github.com/redis/rueidis"
)

//go:embed lua/matchmaking.lua
var matchmakingScript string

var _ repository.GameRepositoryContract = &GameRepository{}

type GameRepository struct {
	redis rueidis.Client
}

func NewGameRepository(redisClient *rueidis.Client) *GameRepository {
	return &GameRepository{
		redis: *redisClient,
	}
}

func (r *GameRepository) GetGameInformation(ctx context.Context, gameID string) (map[string]string, error) {
	result := make(map[string]string)

	gameFields := []string{
		"who_has_won_the_cards",
		"last_move_timestamp",
		"was_the_king_changed",
		"center_cards",
		"trump",
		"cards",
		"has_king_cards_finished",
		"king",
		"players",
		"who_has_won_the_game",
		"lead_suit",
		"is_it_new_round",
		"turn",
		"who_has_won_the_round",
		"king_cards",
		"points",
	}

	command := r.redis.B().Hmget().Key("game:" + gameID).Field(gameFields...).Build()
	information, err := r.redis.Do(ctx, command).ToArray()
	if err != nil {
		log.Fatalf("could not resolve game information: %v", err)
		return nil, err
	}

	for key, data := range information {
		var value map[string]interface{}
		err = json.Unmarshal([]byte(data.String()), &value)
		if err != nil {
			log.Fatalf("could not resolve game information: %v", err)
			return nil, err
		}
		result[gameFields[key]] = value["Value"].(string)
	}

	return result, nil
}

func (r *GameRepository) Matchmaking(ctx context.Context, username string) {
	command := r.redis.B().Eval().Script(matchmakingScript).Numkeys(2).Key("matchmaking", "matchmaking").Arg(username).Build()
	_, err := r.redis.Do(ctx, command).ToArray()
	if err != nil {
		log.Fatalf("could not execute Lua script: %v", err)
	}
}

func (r *GameRepository) AddPlayerToWaitingList(ctx context.Context, username string) {
	ts := float64(time.Now().Unix())
	command := r.redis.B().Zadd().Key("waitingList").ScoreMember().ScoreMember(ts, username).Build()
	err := r.redis.Do(ctx, command).Error()
	if err != nil {
		log.Printf("Error in adding player to list: %v", err)
	}
}

func (r *GameRepository) RemovePlayerFromWaitingList(ctx context.Context, username string) {
	command := r.redis.B().Zrem().Key("waitingList").Member(username).Build()
	err := r.redis.Do(ctx, command).Error()
	if err != nil {
		log.Printf("Error in removing player from list: %v", err)
	}
}

func (r *GameRepository) CheckAnyExistingGameForPlayer(ctx context.Context, username string) (string, bool, error) {
	command := r.redis.B().Hget().Key("playersGame").Field(username).Build()
	val, err := r.redis.Do(ctx, command).ToString()
	if err != nil {
		if errors.Is(err, rueidis.Nil) {
			return "", false, nil
		}
		return "", false, fmt.Errorf("hget playersGame[%s] failed: %w", username, err)
	}
	return val, true, nil
}
