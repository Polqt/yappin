package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	userHandler "chat-application/internal/api/handler/user"

	statsHandler "chat-application/internal/api/handler/stats"

	authMiddleware "chat-application/middleware"
)

func SetupRoutes(userHandler *userHandler.UserHandler) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000", "https://yappin.chat"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders: []string{"Link"},
		AllowCredentials: true,
		MaxAge: 300,
	}))

	r.Route("/api/users", func(u chi.Router) {
		u.Post("/sign-up", userHandler.CreateUser)
		u.Post("/login", userHandler.Login)
		u.Post("/logout", userHandler.Logout)

		u.Group(func(r chi.Router)  {
			r.Use(authMiddleware.JWTAuth)
			r.Put("/username", userHandler.UpdateUsername)
		})
	})

	r.Route("/api/stats", func (s chi.Router)  {
		s.Group(func(r chi.Router) {
			r.Use(authMiddleware.JWTAuth)
			r.Post("/checkin", statsHandler.CheckIn)
			r.Post("/upvote", statsHandler.GivenUpvote)
		})

		s.Group(func(r chi.Router)  {
			r.Use(authMiddleware.JWTAuth)
			r.Get("/profile/{userID}", statsHandler.GetUserProfile)
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