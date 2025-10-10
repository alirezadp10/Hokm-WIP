package database

import (
	_ "embed"
	"fmt"
	"log"
	"os"

	"github.com/redis/rueidis"
)

func GetNewRedisConnection() *rueidis.Client {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")

	address := fmt.Sprintf("%s:%s", host, port)

	client, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress: []string{address},
	})

	if err != nil {
		log.Fatal("couldn't connect to redis")
	}
	return &client
}
