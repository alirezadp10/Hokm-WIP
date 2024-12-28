package service

import (
    "context"
    "github.com/redis/rueidis"
    "log"
)

type RedisService struct {
    client rueidis.Client
    ctx    context.Context
}

func NewRedisService(client rueidis.Client, ctx context.Context) *RedisService {
    return &RedisService{client: client, ctx: ctx}
}

func (s *RedisService) Subscribe(ctx context.Context, channel string, message func(rueidis.PubSubMessage)) error {
    err := s.client.Receive(ctx, s.client.B().Subscribe().Channel(channel).Build(), message)

    if err != nil {
        log.Printf("Error in subscribing to %v channel: %v", channel, err)
        return err
    }

    return nil
}

func (s *RedisService) Unsubscribe(ctx context.Context, channel string) {
    unsubscribeErr := s.client.Do(ctx, s.client.B().Unsubscribe().Channel(channel).Build()).Error()
    if unsubscribeErr != nil {
        log.Println("Error while unsubscribing:", unsubscribeErr)
    }
}
