package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"chat-application/db"
	"chat-application/db/migrations"

	userHandler "chat-application/internal/api/handler/user"
	userRepo "chat-application/internal/repo/user"
	userService "chat-application/internal/service/user"

	statsRepo "chat-application/internal/repo/stats"
	statsService "chat-application/internal/service/stats"
	statsHandler "chat-application/internal/api/handler/stats"

	"chat-application/router"
)

func main(){
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}

	dbConn, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	defer dbConn.Close()
	
	if err := dbConn.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	} 
	log.Println("Database connection established successfully")

	// Run Migrations
	if err := migrations.RunMigrations(dbConn); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	userRepo := userRepo.NewUserRepository(dbConn)
	statsRepository := statsRepo.NewStatsRepository(dbConn)
	
	userService := userService.NewUserService(userRepo)
	statsService := statsService.NewStatsService(statsRepository)

	userHandler := userHandler.NewUserHandler(userService)
	statsHandler := statsHandler.NewStatsHandler(statsService)

	// router := router.SetupRoutes(userHandler)
	router := router.SetupRoutes(userHandler, statsHandler)
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}