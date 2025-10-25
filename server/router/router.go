package router

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	coreHandler "chat-application/internal/api/handler/core"
	statsHandler "chat-application/internal/api/handler/stats"
	userHandler "chat-application/internal/api/handler/user"
	authMiddleware "chat-application/internal/middleware"
)

func SetupRoutes(userHandler *userHandler.UserHandler, coreHandler *coreHandler.CoreHandler, statsHandler *statsHandler.StatsHandler) http.Handler {
	r := chi.NewRouter()

	r.Use(authMiddleware.SecurityHeaders)
	r.Use(authMiddleware.RequestSizeLimit(1 << 20))
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(authMiddleware.Timeout(30 * time.Second))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:5174", "https://yappin.chat"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Route("/api", func(api chi.Router) {
		api.Use(authMiddleware.GetRateLimiter(100))
		api.Use(authMiddleware.ContentTypeJSON)

		api.Route("/users", func(u chi.Router) {
			u.Group(func(r chi.Router) {
				r.Use(authMiddleware.GetRateLimiter(10))
				r.Post("/sign-up", userHandler.CreateUser)
				r.Post("/login", userHandler.Login)
			})

			r.Post("/logout", userHandler.Logout)

			u.Group(func(r chi.Router) {
				r.Use(authMiddleware.JWTAuth)
				r.Get("/me", userHandler.GetCurrentUser)
				r.Put("/username", userHandler.UpdateUsername)
			})
		})

		api.Route("/stats", func(s chi.Router) {
			s.Group(func(r chi.Router) {
				r.Use(authMiddleware.JWTAuth)
				r.Post("/checkin", statsHandler.CheckIn)
				r.Post("/upvote", statsHandler.GivenUpvote)
			})

			s.Group(func(r chi.Router) {
				r.Use(authMiddleware.JWTAuth)
				r.Get("/profile/{userID}", statsHandler.GetUserProfile)
			})
		})

		api.Route("/websoc", func(u chi.Router) {
			u.Group(func(r chi.Router) {
				r.Use(authMiddleware.OptionalJWTAuth)
				// Apply a stricter per-route rate limiter to prevent room-creation spam
				r.Use(authMiddleware.GetRateLimiter(10)) // 10 requests per minute for creating rooms
				r.Post("/create-room", coreHandler.CreateRoom)
			})

			u.Get("/join-room/{roomId}", coreHandler.JoinRoom)
			u.Get("/get-rooms", coreHandler.GetRooms)
			u.Get("/get-clients", coreHandler.GetClients)
		})
	})

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	return r
}
