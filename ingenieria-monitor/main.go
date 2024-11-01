package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "time"

    "cloud.google.com/go/pubsub"
    "github.com/go-redis/redis/v8"
)

type StudentResult struct {
    Student    string    `json:"student"`
    Age       uint32    `json:"age"`
    Faculty   string    `json:"faculty"`
    Discipline uint32    `json:"discipline"`
    Result    string    `json:"result"`
    Timestamp time.Time `json:"timestamp"`
}

func main() {
    ctx := context.Background()

    // Initialize Redis client
    rdb := redis.NewClient(&redis.Options{
        Addr: "redis-12311.c266.us-east-1-3.ec2.redns.redis-cloud.com:12311",
        Password: "y08kAqSaJnJt7UgWkPAm4kcIBrwsOYfH",
        DB: 0,
    })

    // Test Redis connection
    pong, err := rdb.Ping(ctx).Result()
    if err != nil {
        log.Fatalf("Failed to connect to Redis: %v", err)
    }
    log.Printf("Redis connected successfully: %s", pong)

    // Initialize Pub/Sub client
    pubsubClient, err := pubsub.NewClient(ctx, "servicio-440415")
    if err != nil {
        log.Fatal(err)
    }
    defer pubsubClient.Close()

    // Get subscription
    sub := pubsubClient.Subscription("ingenieria-monitor-sub")

    log.Println("Starting Ingenieria monitor...")

    // Receive messages
    err = sub.Receive(ctx, func(c context.Context, msg *pubsub.Message) {
        var result StudentResult
        err := json.Unmarshal(msg.Data, &result)
        if err != nil {
            log.Printf("Error unmarshaling message: %v", err)
            msg.Ack()
            return
        }

        // Only process Ingenieria students
        if result.Faculty == "Ingenieria" {
            // Create key based on result (winner/loser)
            key := fmt.Sprintf("ingenieria:%s:discipline:%d", result.Result, result.Discipline)
            err = rdb.Incr(ctx, key).Err()
            if err != nil {
                log.Printf("Error incrementing Redis key: %v", err)
            } else {
                log.Printf("Incremented count for %s", key)
            }
        }

        msg.Ack()
    })

    if err != nil {
        log.Fatal(err)
    }
}