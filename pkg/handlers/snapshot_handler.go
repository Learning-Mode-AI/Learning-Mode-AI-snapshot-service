package handlers

import (
	"Learning-Mode-AI-Snapshot-Service/pkg/services"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func ProcessSnapshots(w http.ResponseWriter, r *http.Request) {
	var req struct {
		VideoID    string   `json:"video_id"`
		VideoURL   string   `json:"video_url"`
		Timestamps []string `json:"timestamps"`
	}

	// Parse the request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Ensure the video ID is initialized in Redis
	exists, err := services.RedisClient.Exists(services.Ctx, req.VideoID).Result()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to check Redis for video ID: %v", err), http.StatusInternalServerError)
		return
	}

	if exists == 0 {
		videoMetadata := map[string]interface{}{
			"video_id": req.VideoID,
			"created":  fmt.Sprintf("%v", time.Now()),
		}

		videoMetadataJSON, err := json.Marshal(videoMetadata)
		if err != nil {
			http.Error(w, "Failed to initialize Redis for video ID", http.StatusInternalServerError)
			return
		}

		err = services.RedisClient.Set(services.Ctx, req.VideoID, videoMetadataJSON, 0).Err()
		if err != nil {
			http.Error(w, "Failed to store video metadata in Redis", http.StatusInternalServerError)
			return
		}
		log.Printf("Initialized video ID in Redis: %s", req.VideoID)
	}

	// Process each timestamp
	for _, timestamp := range req.Timestamps {
		// Capture a snapshot for each timestamp
		imagePath, err := services.CaptureSnapshot(req.VideoURL, timestamp, req.VideoID)
		if err != nil {
			log.Printf("Failed to capture snapshot for %s: %v", timestamp, err)
			continue
		}

		// Store the snapshot in Redis under the video ID
		err = services.StoreSnapshotInRedis(req.VideoID, timestamp, imagePath)
		if err != nil {
			log.Printf("Failed to store snapshot in Redis for %s: %v", timestamp, err)
		}
	}

	// Minimal response: Success message
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Snapshots stored successfully"})
}
