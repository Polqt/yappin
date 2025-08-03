package topics

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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

type redditTokenRespose struct {
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

const redditUA = "desktop:gochat:v1.0 (by /u/Polqt)"

func (s *TopicsService) GetRedditToken(ctx context.Context) error {
	clientID := os.Getenv("REDDIT_CLIENT_ID")
	clientSecret := os.Getenv("REDDIT_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		return fmt.Errorf("REDDIT_CLIENT_ID and REDDIT_CLIENT_SECRET must be set")
	}

	data := strings.NewReader("grant_type=client_credentials")
	req, err := http.NewRequestWithContext(ctx, "POST", "https://www.reddit.com/api/v1/access_token", data)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(clientID, clientSecret)
	req.Header.Set("User-Agent", redditUA)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	
	response, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to get reddit token: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get reddit token, status code: %d", response.StatusCode)
	}

	var tokenResponse redditTokenRespose
	if err := json.NewDecoder(response.Body).Decode(&tokenResponse); err != nil {
		return fmt.Errorf("failed to decode reddit token response: %w", err)
	}

	s.redditToken = tokenResponse.AccessToken
	s.tokenExpiry = time.Now().Add(time.Duration(tokenResponse.ExpiresIn-60) * time.Second)

	return nil
} 

// func (s *TopicsService) GetRedditJSON(ctx context.Context, url string, out any) error {

// }

func (s *TopicsService) FetchAllTopics(ctx context.Context) ([]Topic, error) {
	topics := make([]Topic, 0, 5)

	
	return topics, nil
}

