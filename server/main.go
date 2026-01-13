package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"chat-application/db"
	"chat-application/db/migrations"

	userHandler "chat-application/internal/api/handler/user"
	"chat-application/internal/constants"
	roomRepository "chat-application/internal/repo/room"
	userRepo "chat-application/internal/repo/user"
	userService "chat-application/internal/service/user"
	websoc "chat-application/internal/websocket"

	statsHandler "chat-application/internal/api/handler/stats"
	statsRepo "chat-application/internal/repo/stats"
	statsService "chat-application/internal/service/stats"

	coreHandler "chat-application/internal/api/handler/core"
	"chat-application/internal/middleware"
	pinnedRooms "chat-application/internal/service/pinnedrooms"
	"chat-application/router"
)

func main() {
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

	// Start WebSocket core to process messages
	log.Println("Starting WebSocket core...")
	go webService.Start()

	go startRoomCleanup(dbConn, webService)

	rateLimiter := middleware.NewRateLimiter(constants.DefaultRateLimit, constants.RateLimitWindow)
	routerWithLimiter := rateLimiter.Middleware(router.SetupRoutes(userHandler, coreHandler, statsHandler))

	// Create server with graceful shutdown support
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      routerWithLimiter,
		ReadTimeout:  constants.HTTPServerTimeout,
		WriteTimeout: constants.HTTPServerTimeout,
		IdleTimeout:  constants.HTTPServerTimeout * 2,
	}

	// Start server in goroutine
	go func() {
		log.Println("Server starting on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Give outstanding requests time to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}

func startRoomCleanup(db *sql.DB, websocketCore *websoc.Core) {
	roomRepository := roomRepository.NewRoomRepository(db)
	pinnedRoomsService := pinnedRooms.NewPinnedRoomsService(db, websocketCore)
	ticker := time.NewTicker(constants.RoomCleanupInterval)
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
