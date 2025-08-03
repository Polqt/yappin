package router

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	statsHandler "chat-application/internal/api/handler/stats"
	userHandler "chat-application/internal/api/handler/user"
	authMiddleware "chat-application/middleware"
)

func SetupRoutes(userHandler *userHandler.UserHandler, statsHandler *statsHandler.StatsHandler) http.Handler {
	r := chi.NewRouter()
	
	// Security middleware
	r.Use(authMiddleware.SecurityHeaders)
	r.Use(authMiddleware.RequestSizeLimit(1 << 20)) // 1MB limit
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(authMiddleware.Timeout(30*time.Second))
	
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000", "https://yappin.chat"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders: []string{"Link"},
		AllowCredentials: true,
		MaxAge: 300,
	}))

	// API routes with rate limiting
	r.Route("/api", func(api chi.Router) {
		api.Use(authMiddleware.RateLimiter(100)) // 100 requests per minute
		api.Use(authMiddleware.ContentTypeJSON)
		
		api.Route("/users", func(u chi.Router) {
			// Public endpoints with stricter rate limiting
			u.Group(func(r chi.Router) {
				r.Use(authMiddleware.RateLimiter(10)) // 10 requests per minute for auth endpoints
				r.Post("/sign-up", userHandler.CreateUser)
				r.Post("/login", userHandler.Login)
			})
			
			r.Post("/logout", userHandler.Logout)

			u.Group(func(r chi.Router)  {
				r.Use(authMiddleware.JWTAuth)
				r.Put("/username", userHandler.UpdateUsername)
			})
		})

		api.Route("/stats", func (s chi.Router)  {
			s.Group(func(r chi.Router) {
				r.Use(authMiddleware.JWTAuth)
				r.Post("/checkin", statsHandler.CheckIn)
				r.Post("/upvote", statsHandler.GiveUpvote)
			})

			s.Group(func(r chi.Router)  {
				r.Use(authMiddleware.JWTAuth)
				r.Get("/profile/{userID}", statsHandler.GetUserProfile)
			})
		})
	})

	r.Route("/ws", func(u chi.Router) {
		u.Group(func(r chi.Router) {})
	})
	
	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	return r
}