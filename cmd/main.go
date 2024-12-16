package main

import (
	"log"
	"net/http"
	"Learning-Mode-AI-Snapshot-Service/pkg/router"
	"Learning-Mode-AI-Snapshot-Service/pkg/services"
)

func main() {
	// Initialize Redis connection
	services.InitRedis()

	// Initialize storage folders
	services.InitFolders()

	// Setup routes
	router.SetupRoutes()

	// Start the video processing service
	log.Println("Starting Video Processing Service on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
