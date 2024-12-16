package handlers

import (
	"Learning-Mode-AI-Snapshot-Service/pkg/services"
	"encoding/json"
	"log"
	"net/http"
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

	// Capture a snapshot for each timestamp
	for _, timestamp := range req.Timestamps {
		// Store the snapshot in Redis under the video ID
		imagePath, err := services.CaptureSnapshot(req.VideoURL, timestamp)
		if err != nil {
			log.Printf("Failed to capture snapshot: %v", err)
			continue
		}

		if err := services.StoreSnapshotInRedis(req.VideoID, timestamp, imagePath); err != nil {
			log.Printf("Failed to store snapshot in Redis: %v", err)
		}
	}

	// Minimal response: Success message
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Snapshots processed successfully"})
}