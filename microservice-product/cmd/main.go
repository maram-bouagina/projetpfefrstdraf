package main

import (
	"fmt"
	"log"
	"projet/internal/config"
	"projet/internal/db"
	"projet/routes"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	database, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Connected to database")

	sqlDB, _ := database.DB()
	defer sqlDB.Close()

	app := routes.NewRouter(database)
	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	if err := app.Listen(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
