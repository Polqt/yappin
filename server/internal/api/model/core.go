package model

import "time"

type CreateRoomReq struct {
	ID        string     `json:"id,omitempty"`
	Name      string     `json:"name"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

type ClientRes struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type RoomRes struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	IsPinned         bool      `json:"is_pinned"`
	CreatedAt        time.Time `json:"created_at"`
	Expires          time.Time `json:"expires_at"`
	TopicTitle       *string   `json:"topic_title,omitempty"`
	TopicDescription *string   `json:"topic_description,omitempty"`
	TopicURL         *string   `json:"topic_url,omitempty"`
	TopicSource      *string   `json:"topic_source,omitempty"`
	CreatorUsername  *string   `json:"creator_username,omitempty"`
	Participants     int       `json:"participants"`
}

type MessageReaction struct {
	ID        string    `json:"id"`
	MessageID string    `json:"message_id"`
	UserID    string    `json:"user_id"`
	Emoji     string    `json:"emoji"`
	CreatedAt time.Time `json:"created_at"`
}

type RequestAddReaction struct {
	MessageID string `json:"message_id"`
	Emoji     string `json:"emoji"`
}
