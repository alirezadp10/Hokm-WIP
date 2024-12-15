package main

import (
    "context"
    "fmt"
    "github.com/redis/rueidis"
    "log"
)

func main() {
    client, err := rueidis.NewClient(rueidis.ClientOption{InitAddress: []string{"127.0.0.1:6379"}})
    if err != nil {
        log.Fatalf("could not connect to Redis: %v", err)
    }
    defer client.Close()

    err = client.Receive(context.Background(), client.B().Subscribe().Channel("waiting").Build(), func(msg rueidis.PubSubMessage) {
        fmt.Println("Received message:", msg.Message)
    })
}
