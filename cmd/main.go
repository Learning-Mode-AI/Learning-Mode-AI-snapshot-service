package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"

	"github.com/go-redis/redis/v8"
)

// SnapshotRequest struct to parse incoming requests
type SnapshotRequest struct {
	VideoID    string   `json:"video_id"`  // Updated to use Video ID
	VideoURL   string   `json:"video_url"` // URL still needed for downloading
	Timestamps []string `json:"timestamps"`
}

var redisClient *redis.Client
var ctx = context.Background()

// Initialize Redis connection
func initRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis address
		Password: "",               // No password set
		DB:       0,                // Use default DB
	})

	// Test the connection
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Connected to Redis")
}

// Function to capture snapshots using yt-dlp and ffmpeg
func captureSnapshot(videoURL, timestamp string) (string, error) {
	// Define output video file (temporary video file for processing)
	videoFile := "downloaded_video.mp4"

	// Download the video from YouTube using yt-dlp
	downloadCmd := exec.Command("yt-dlp", "-f", "best", "-o", videoFile, videoURL)
	err := downloadCmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to download video: %v", err)
	}

	// Create output filename based on timestamp
	filename := fmt.Sprintf("snapshot_%s.png", strings.ReplaceAll(timestamp, ":", "-"))

	// Use ffmpeg to capture the snapshot at the given timestamp from the downloaded video
	ffmpegCmd := exec.Command("ffmpeg", "-y", "-i", videoFile, "-ss", timestamp, "-vframes", "1", filename)
	err = ffmpegCmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to capture snapshot: %v", err)
	}

	// Optionally, remove the video file to save space
	_ = exec.Command("rm", videoFile).Run()

	return filename, nil
}

// Function to process snapshots and store them in Redis
func processSnapshots(w http.ResponseWriter, r *http.Request) {
	var req SnapshotRequest

	// Parse the request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	for _, timestamp := range req.Timestamps {
		// Capture a snapshot for each timestamp
		imagePath, err := captureSnapshot(req.VideoURL, timestamp)
		if err != nil {
			log.Printf("Failed to capture snapshot for %s: %v", timestamp, err)
			continue
		}

		// Store the snapshot in Redis under the video ID
		err = storeSnapshotInRedis(req.VideoID, timestamp, imagePath)
		if err != nil {
			log.Printf("Failed to store snapshot in Redis for %s: %v", timestamp, err)
		}
	}

	// Minimal response: Success message
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Snapshots stored successfully"})
}

// Function to store snapshot data in Redis under the video ID key
func storeSnapshotInRedis(videoID, timestamp, imagePath string) error {
	// Check if the video ID exists in Redis
	exists, err := redisClient.Exists(ctx, videoID).Result()
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
	err = redisClient.RPush(ctx, snapshotKey, snapshotJSON).Err()
	if err != nil {
		return fmt.Errorf("failed to store snapshot in Redis: %v", err)
	}

	log.Printf("Stored snapshot in Redis for video ID: %s, timestamp: %s", videoID, timestamp)
	return nil
}

func main() {
	// Initialize Redis connection
	initRedis()

	// Handle snapshot processing requests
	http.HandleFunc("/process-snapshots", processSnapshots)

	// Start the video processing service
	log.Println("Starting Video Processing Service on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
