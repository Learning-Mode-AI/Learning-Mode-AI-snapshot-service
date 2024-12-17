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

func CaptureSnapshot(videoURL, timestamp, videoID string) (string, error) {
    // Define per-video folder for snapshots
    videoSnapshotFolder := fmt.Sprintf("%s/%s", SnapshotFolder, videoID)

    // Create the video-specific folder if it doesn't exist
    if err := os.MkdirAll(videoSnapshotFolder, os.ModePerm); err != nil {
        return "", fmt.Errorf("failed to create folder %s: %v", videoSnapshotFolder, err)
    }

    // Define output video file (temporary video file for processing)
    videoFile := fmt.Sprintf("%s/%s.mp4", VideoFolder, videoID)

    // Download the video from YouTube using yt-dlp
    downloadCmd := exec.Command("yt-dlp", "-f", "best", "-o", videoFile, videoURL)
    err := downloadCmd.Run()
    if err != nil {
        return "", fmt.Errorf("failed to download video: %v", err)
    }

    // Create output snapshot filename
    filename := fmt.Sprintf("snapshot_%s.png", strings.ReplaceAll(timestamp, ":", "-"))
    snapshotPath := fmt.Sprintf("%s/%s", videoSnapshotFolder, filename)

    // Use ffmpeg to capture the snapshot at the given timestamp
    ffmpegCmd := exec.Command("ffmpeg", "-y", "-i", videoFile, "-ss", timestamp, "-vframes", "1", snapshotPath)
    err = ffmpegCmd.Run()
    if err != nil {
        return "", fmt.Errorf("failed to capture snapshot: %v", err)
    }

    // Optionally, remove the video file to save space
    _ = exec.Command("rm", videoFile).Run()

    return snapshotPath, nil
}
