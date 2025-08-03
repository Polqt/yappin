package service

import (
	"net/http"
	"strings"
	"time"
)

type TopicsService struct {
	client *http.Client
	redditToken string
	tokenExpiry time.Time
}

type Topic struct {
	Title string
	Description string
	URL string
	Source string
}

type redditToken struct {
	AccessToken string `json:"access_token"`
	TokenType string `json:"token_type"`
	ExpiresIn int `json:"expires_in"`
}

func NewTopicsService() *TopicsService {
	return &TopicsService{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func cleanText(text string) string  {
	decoded := strings.ReplaceAll(text, "\\n", "\n")
	return strings.TrimSpace(decoded)
}

