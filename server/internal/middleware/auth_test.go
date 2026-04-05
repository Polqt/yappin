package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func issueTestToken(t *testing.T, secret, userID string) string {
	t.Helper()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userID,
		"exp": time.Now().Add(time.Hour).Unix(),
	})

	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}

	return signed
}

func TestJWTAuthAcceptsValidCookie(t *testing.T) {
	t.Setenv("JWT_SECRET_KEY", "test-secret")

	called := false
	handler := JWTAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		if got := r.Context().Value(UserIDKey); got != "user-123" {
			t.Fatalf("expected user id in context, got %v", got)
		}
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/users/me", nil)
	req.AddCookie(&http.Cookie{Name: "jwt", Value: issueTestToken(t, "test-secret", "user-123")})
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if !called {
		t.Fatal("expected next handler to be called")
	}
	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, rec.Code)
	}
}

func TestJWTAuthRejectsMissingCookie(t *testing.T) {
	t.Setenv("JWT_SECRET_KEY", "test-secret")

	handler := JWTAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next handler should not be called")
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/users/me", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}
}

func TestJWTAuthRejectsInvalidToken(t *testing.T) {
	t.Setenv("JWT_SECRET_KEY", "test-secret")

	handler := JWTAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next handler should not be called")
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/users/me", nil)
	req.AddCookie(&http.Cookie{Name: "jwt", Value: "not-a-token"})
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}
}

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}
