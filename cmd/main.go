package main

import (
	"chat-websocket/internal/config"
	"chat-websocket/internal/database"
	"chat-websocket/internal/infrastructure/repository"
	"chat-websocket/internal/presentation/http/router"
	"log"
	"net/http"
)

func main() {
	// Load environment variables
	config.LoadEnv()

	// Initialize Database
	db := database.InitDB()
	socketPathRepo := repository.NewSocketPathRepository(db)

	socketPath, err := socketPathRepo.FindAll()

	if err != nil {
		log.Fatalf("Failed to fetch paths: %v", err)
	}
	var paths []string
	for _, path := range socketPath {
		paths = append(paths, path.Path)
	}
	router.SetupRoutes(paths)
	port := config.GetEnv("APP_PORT", "8080")
	log.Println("Server started on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
