package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"

	"chat-application/db"
	"chat-application/db/migrations"

	userHandler "chat-application/internal/api/handler/user"
	roomRepository "chat-application/internal/repo/room"
	userRepo "chat-application/internal/repo/user"
	userService "chat-application/internal/service/user"
	websoc "chat-application/internal/websocket"

	statsHandler "chat-application/internal/api/handler/stats"
	statsRepo "chat-application/internal/repo/stats"
	statsService "chat-application/internal/service/stats"

	coreHandler "chat-application/internal/api/handler/core"

	pinnedRooms "chat-application/internal/service/pinnedrooms"

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
	webService := websoc.NewCore(dbConn)

	userHandler := userHandler.NewUserHandler(userService)
	coreHandler := coreHandler.NewCoreHandler(webService)
	statsHandler := statsHandler.NewStatsHandler(statsService)

	pinnedRoomService := pinnedRooms.NewPinnedRoomsService(dbConn, webService)
	if err := pinnedRoomService.RefreshPinnedRooms(context.Background()); err != nil {
		log.Printf("Failed to refresh pinned rooms: %v", err)
	}

	go startRoomCleanup(dbConn, webService)

	router := router.SetupRoutes(userHandler, coreHandler, statsHandler)
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func startRoomCleanup(db *sql.DB, websocketCore *websoc.Core) {
	roomRepository := roomRepository.NewRoomRepository(db)
	pinnedRoomsService := pinnedRooms.NewPinnedRoomsService(db, websocketCore)
	ticker := time.NewTicker(5 *time.Minute)
	defer ticker.Stop()

	cleanupRooms(roomRepository, pinnedRoomsService)

	for range ticker.C {
		cleanupRooms(roomRepository, pinnedRoomsService)
	}
}

func cleanupRooms(roomRepository *roomRepository.RoomRepository, pinnedRoomsService *pinnedRooms.PinnedRoomsService) {
	ctx := context.Background()

	deletedCount, err := roomRepository.DeleteExpiredRooms(ctx)
	if err != nil {
		log.Printf("Failed to delete expired rooms: %v", err)
		return
	}

	if deletedCount > 0 {
		log.Printf("Deleted %d expired rooms", deletedCount)
	}

	if err := pinnedRoomsService.RefreshPinnedRooms(ctx); err != nil {
		log.Printf("Failed to refresh pinned rooms: %v", err)
	}
}

