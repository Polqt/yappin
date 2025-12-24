package middleware

import (
	"context"
	"log"
	"net/http"

	"chat-application/util"

	"github.com/golang-jwt/jwt/v5"
)

type ContextKey string

const UserIDKey ContextKey = "userID"

func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("jwt")
		if err != nil {
			util.WriteErrorResponse(w, http.StatusUnauthorized, "missing auth token")
			return
		}

		tokenString := cookie.Value
		if tokenString == "" {
			util.WriteErrorResponse(w, http.StatusUnauthorized, "missing auth token")
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			secretKey := util.GetEnv("JWT_SECRET_KEY", "")
			if secretKey == "" {
				log.Println("JWT_SECRET_KEY is not set in environment variables")
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			util.WriteErrorResponse(w, http.StatusUnauthorized, "invalid auth token")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			util.WriteErrorResponse(w, http.StatusUnauthorized, "invalid token claims")
			return
		}

		userID, ok := claims["id"].(string)
		if !ok {
			util.WriteErrorResponse(w, http.StatusUnauthorized, "invalid user ID in token")
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func OptionalJWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("OptionalJWTAuth middleware triggered for %s", r.URL.Path)

		cookie, err := r.Cookie("jwt")
		if err != nil {
			log.Println("No JWT cookie found, proceeding without authentication")
			next.ServeHTTP(w, r)
			return
		}

		if cookie.Value == "" {
			log.Println("JWT cookie is empty, proceeding without authentication")
			next.ServeHTTP(w, r)
			return
		}

		log.Println("JWT cookie found")

		token, err := jwt.Parse(cookie.Value, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				log.Printf("Unexpected signing method: %v", t.Header["alg"])
				return nil, jwt.ErrSignatureInvalid
			}

			secretKey := util.GetEnv("JWT_SECRET_KEY", "")
			if secretKey == "" {
				log.Println("JWT_SECRET_KEY is not set in environment variables")
				return nil, jwt.ErrSignatureInvalid
			}

			return []byte(secretKey), nil
		})

		if err != nil {
			log.Printf("Error parsing JWT: %v", err)
			// Don't write error response - this is OPTIONAL auth
			next.ServeHTTP(w, r)
			return
		}

		if !token.Valid {
			log.Println("Invalid JWT token")
			// Don't write error response - this is OPTIONAL auth
			next.ServeHTTP(w, r)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			log.Printf("JWT claims: %v", claims)
			if userID, ok := claims["id"].(string); ok {
				log.Printf("Authenticated user ID: %s", userID)
				ctx := context.WithValue(r.Context(), UserIDKey, userID)
				r = r.WithContext(ctx)
			} else {
				log.Printf("User ID not found in JWT claims")
			}
		} else {
			log.Println("Invalid JWT claims")
		}

		next.ServeHTTP(w, r)
	})
}
