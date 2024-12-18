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

	// Initialize video metadata in Redis
	exists, err := services.RedisClient.Exists(services.Ctx, req.VideoID).Result()
	if err != nil || exists == 0 {
		videoMetadata := map[string]interface{}{
			"video_id": req.VideoID,
			"created":  time.Now().String(),
		}
		metadataJSON, _ := json.Marshal(videoMetadata)
		services.RedisClient.Set(services.Ctx, req.VideoID, metadataJSON, 0)
		log.Printf("Initialized video ID %s in Redis", req.VideoID)
	}

	// Process snapshots
	snapshotPaths, err := services.CaptureSnapshots(req.VideoURL, req.Timestamps, req.VideoID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Snapshot processing failed: %v", err), http.StatusInternalServerError)
		return
	}

	// Store each snapshot path in Redis
	for i, path := range snapshotPaths {
		err := services.StoreSnapshotInRedis(req.VideoID, req.Timestamps[i], path)
		if err != nil {
			log.Printf("Failed to store snapshot for timestamp %s: %v", req.Timestamps[i], err)
		}
	}

	// Success response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":        "Snapshots stored successfully",
		"snapshot_paths": snapshotPaths,
	})
}