package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

type SnapshotRequest struct {
	VideoURL   string   `json:"video_url"`
	Timestamps []string `json:"timestamps"`
}

type SnapshotResponse struct {
	Timestamp string `json:"timestamp"`
	ImagePath string `json:"image_path"`
}

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

	// Optionally, you could remove the video file to save space
	_ = exec.Command("rm", videoFile).Run()

	return filename, nil
}

func processSnapshots(w http.ResponseWriter, r *http.Request) {
	var req SnapshotRequest

	// Parse the request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Prepare a response structure to store the results
	var snapshots []SnapshotResponse

	for _, timestamp := range req.Timestamps {
		// Capture a snapshot for each timestamp
		imagePath, err := captureSnapshot(req.VideoURL, timestamp)
		if err != nil {
			log.Printf("Failed to capture snapshot for %s: %v", timestamp, err)
			continue
		}

		// Add the snapshot result to the response
		snapshots = append(snapshots, SnapshotResponse{
			Timestamp: timestamp,
			ImagePath: imagePath,
		})
	}

	// Return the snapshots as a JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(snapshots)
}

func main() {
	http.HandleFunc("/process-snapshots", processSnapshots)
	log.Println("Starting Video Processing Service on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
