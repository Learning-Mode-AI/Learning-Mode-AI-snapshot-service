package router

import (
	"Learning-Mode-AI-Snapshot-Service/pkg/handlers"
	"net/http"
)

func SetupRoutes() {
	http.HandleFunc("/process-snapshots", handlers.ProcessSnapshots)
}
