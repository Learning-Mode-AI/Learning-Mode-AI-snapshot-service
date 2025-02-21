package services

import (
	"Learning-Mode-AI-Snapshot-Service/pkg/config"
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
		Addr:     config.RedisHost, // Redis address
		TLSConfig: &tls.Config{
			// Depending on your certificate setup,
			// you might need to customize this further.
			InsecureSkipVerify: true, // Use caution: this bypasses certificate verification.
		},
	})

	err := RedisClient.Ping(Ctx).Err()
	if err != nil {
		panic(err)
	}
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
