package services

import (
	"fmt"
	"os/exec"
	"strings"
	"log"
)

func CaptureSnapshot(videoURL, timestamp string) (string, error) {
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

	log.Printf("Saving snapshot to %s", filename)


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