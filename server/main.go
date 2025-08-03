package main

import (
	"chat-application/db"
	"chat-application/db/migrations"
	"chat-application/internal/api/handler/user"
	"chat-application/internal/repo/user"
	"chat-application/internal/service/user"
	"log"

	"github.com/joho/godotenv"
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

	userRepo := repository.NewUserRepository(dbConn)
	
	userService := service.NewUserService(userRepo)

	userHandler := handler.NewUserHandler(userService)

	// router := router.SetupRoutes(userHandler)
}