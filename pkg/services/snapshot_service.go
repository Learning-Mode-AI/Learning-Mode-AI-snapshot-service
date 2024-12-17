package services

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

const (
	VideoFolder    = "storage/videos"
	SnapshotFolder = "storage/snapshots"
)

func InitFolders() {
    folders := []string{VideoFolder} // Only initialize the video folder
    for _, folder := range folders {
        if err := os.MkdirAll(folder, os.ModePerm); err != nil {
            log.Fatalf("Failed to create folder %s: %v", folder, err)
        }
    }
    log.Println("Storage folders initialized")
}

// CaptureSnapshots processes multiple timestamps after downloading the video once
func CaptureSnapshots(videoURL string, timestamps []string, videoID string) ([]string, error) {
	// Define per-video folder for snapshots
	videoSnapshotFolder := fmt.Sprintf("%s/%s", SnapshotFolder, videoID)

	// Create the video-specific folder if it doesn't exist
	if err := os.MkdirAll(videoSnapshotFolder, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create folder %s: %v", videoSnapshotFolder, err)
	}

	// Define output video file
	videoFile := fmt.Sprintf("%s/%s.mp4", VideoFolder, videoID)

	// Check if the video file already exists to avoid redundant downloads
	if _, err := os.Stat(videoFile); os.IsNotExist(err) {
		// Download the video from YouTube using yt-dlp
		log.Println("Downloading video...")
		downloadCmd := exec.Command("yt-dlp", "-f", "best", "-o", videoFile, videoURL)
		if err := downloadCmd.Run(); err != nil {
			return nil, fmt.Errorf("failed to download video: %v", err)
		}
		log.Println("Video downloaded successfully.")
	}

	// Loop through timestamps to generate snapshots
	var snapshotPaths []string
	for _, timestamp := range timestamps {
		filename := fmt.Sprintf("snapshot_%s.png", strings.ReplaceAll(timestamp, ":", "-"))
		snapshotPath := fmt.Sprintf("%s/%s", videoSnapshotFolder, filename)

		// Use ffmpeg to capture the snapshot
		log.Printf("Processing snapshot for timestamp %s...", timestamp)
		ffmpegCmd := exec.Command("ffmpeg", "-y", "-i", videoFile, "-ss", timestamp, "-vframes", "1", snapshotPath)
		if err := ffmpegCmd.Run(); err != nil {
			log.Printf("Failed to capture snapshot at %s: %v", timestamp, err)
			continue
		}

		snapshotPaths = append(snapshotPaths, snapshotPath)
	}

	// Optionally remove the video file to save space
	//_ = os.Remove(videoFile)

	return snapshotPaths, nil
}
