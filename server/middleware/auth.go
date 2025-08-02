package middleware

import (
	"chat-application/util"
	"net/http"
	"context"
	"log"

	"github.com/golang-jwt/jwt/v5"
)

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
			return []byte(util.GetEnv("secretKey", "")), nil
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

		ctx := context.WithValue(r.Context(), "userID", userID)
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

		log.Printf("JWT cookie found, value: %s", cookie.Value)

		token, err := jwt.Parse(cookie.Value, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				log.Printf("Unexpected signing method: %v", t.Header["alg"])
				return nil, jwt.ErrSignatureInvalid
			}
			
			secretKey := util.GetEnv("secretKey", "")
			if secretKey == "" {
				log.Println("Secret key is not set in environment variables")
				return nil, jwt.ErrSignatureInvalid
			}

			return []byte(secretKey), nil
		})

		if err != nil {
			log.Printf("Error parsing JWT: %v", err)
			util.WriteErrorResponse(w, http.StatusUnauthorized, "invalid auth token")
			next.ServeHTTP(w, r)
			return
		}

		if !token.Valid {
			log.Println("Invalid JWT token")
			util.WriteErrorResponse(w, http.StatusUnauthorized, "invalid auth token")
			next.ServeHTTP(w, r)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			log.Printf("JWT claims: %v", claims)
			if userID, ok := claims["id"].(string); ok {
				log.Printf("Authenticated user ID: %s", userID)
				ctx := context.WithValue(r.Context(), "userID", userID)
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