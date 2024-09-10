package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

type SnapshotRequest struct {
	VideoURL   string   `json:"video_url"`
	Timestamps []string `json:"timestamps"`
}

var redisClient *redis.Client

func initRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis address
		Password: "",               // no password set
		DB:       0,                // use default DB
	})
}

func captureSnapshot(videoURL, timestamp string) (string, error) {
	// Create output filename based on timestamp
	filename := fmt.Sprintf("snapshot_%s.png", strings.ReplaceAll(timestamp, ":", "-"))

	// Use yt-dlp to download the video at the given timestamp
	cmd := exec.Command("yt-dlp", "--skip-download", "--write-thumbnail", "--postprocessor-args", fmt.Sprintf("-ss %s", timestamp), videoURL)
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to download snapshot: %v", err)
	}

	return filename, nil
}

func storeSnapshotInRedis(ctx context.Context, timestamp, filename string) error {
	// Store the snapshot in Redis with an expiration time (e.g., 1 hour)
	err := redisClient.Set(ctx, timestamp, filename, 1*time.Hour).Err()
	if err != nil {
		return fmt.Errorf("failed to store snapshot in Redis: %v", err)
	}
	return nil
}

func processSnapshots(w http.ResponseWriter, r *http.Request) {
	var req SnapshotRequest

	// Parse the request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	ctx := context.Background()

	for _, timestamp := range req.Timestamps {
		// Capture a snapshot
		filename, err := captureSnapshot(req.VideoURL, timestamp)
		if err != nil {
			log.Printf("Failed to capture snapshot for %s: %v", timestamp, err)
			continue
		}

		// Store the snapshot in Redis
		err = storeSnapshotInRedis(ctx, timestamp, filename)
		if err != nil {
			log.Printf("Failed to store snapshot in Redis for %s: %v", timestamp, err)
			continue
		}
	}

	// Return success
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Snapshots processed successfully.")
}

func main() {
	initRedis()

	http.HandleFunc("/process-snapshots", processSnapshots)
	log.Println("Starting Video Processing Service on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
