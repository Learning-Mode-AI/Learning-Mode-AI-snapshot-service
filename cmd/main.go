package main

import (
	"Learning-Mode-AI-Snapshot-Service/pkg/config"
	"Learning-Mode-AI-Snapshot-Service/pkg/router"
	"Learning-Mode-AI-Snapshot-Service/pkg/services"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	config.InitConfig()
	services.InitRedis()
}

func main() {
	// Initialize storage folders
	services.InitFolders()

	// Setup routes
	router.SetupRoutes()

	// Start the video processing service
	log.Println("Starting Video Processing Service on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
