package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

var (
	RedisClient *redis.Client
	Ctx         = context.Background() // Exported context to be used across the package
)

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "redis:6379", // Redis address
		Password: "",               // No password set
		DB:       0,                // Use default DB
	})

	// Test the connection
	_, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Connected to Redis")
}

// StoreSnapshotInRedis saves the snapshot data to Redis under the video ID key
func StoreSnapshotInRedis(videoID, timestamp, imagePath string) error {
	ctx := context.Background()

	// Check if the video ID exists in Redis
	exists, err := RedisClient.Exists(ctx, videoID).Result()
	if err != nil || exists == 0 {
		return fmt.Errorf("failed to find video info in Redis for video ID: %s", videoID)
	}

	// Store snapshots in Redis as a List, so they are structured properly
	snapshotKey := fmt.Sprintf("%s:snapshots", videoID)

	snapshotEntry := map[string]string{
		"timestamp": timestamp,
		"imagePath": imagePath,
	}

	// Serialize snapshot entry to JSON
	snapshotJSON, err := json.Marshal(snapshotEntry)
	if err != nil {
		return fmt.Errorf("failed to serialize snapshot entry: %v", err)
	}

	// Push the snapshot JSON into the Redis list
	err = RedisClient.RPush(ctx, snapshotKey, snapshotJSON).Err()
	if err != nil {
		return fmt.Errorf("failed to store snapshot in Redis: %v", err)
	}

	log.Printf("Stored snapshot in Redis for video ID: %s, timestamp: %s", videoID, timestamp)
	return nil
}
