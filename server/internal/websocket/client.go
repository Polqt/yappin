package websocket

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn     *websocket.Conn
	Message  chan *Event
	ID       string `json:"id"`
	RoomID   string `json:"room_id"`
	Username string `json:"username"`
	UserID   string `json:"user_id,omitempty"`
}

type Message struct {
	ID              string         `json:"id,omitempty"`
	Content         string         `json:"content"`
	RoomID          string         `json:"room_id"`
	ChannelID       string         `json:"channel_id,omitempty"`
	ParentMessageID string         `json:"parent_message_id,omitempty"`
	Username        string         `json:"username"`
	UserID          string         `json:"user_id,omitempty"`
	System          bool           `json:"system"`
	CreatedAt       string         `json:"created_at,omitempty"`
	Metadata        map[string]any `json:"metadata,omitempty"`
}

type TypingEvent struct {
	RoomID    string `json:"room_id"`
	ChannelID string `json:"channel_id,omitempty"`
	UserID    string `json:"user_id,omitempty"`
	Username  string `json:"username"`
	IsTyping  bool   `json:"is_typing"`
}

type PresenceEvent struct {
	RoomID      string         `json:"room_id"`
	OnlineUsers []PresenceUser `json:"online_users"`
}

type PresenceUser struct {
	UserID   string `json:"user_id,omitempty"`
	Username string `json:"username"`
}

type NotificationEvent struct {
	ID        string         `json:"id,omitempty"`
	Kind      string         `json:"kind"`
	Title     string         `json:"title"`
	Body      string         `json:"body"`
	RoomID    string         `json:"room_id,omitempty"`
	MessageID string         `json:"message_id,omitempty"`
	Payload   map[string]any `json:"payload,omitempty"`
}

type Event struct {
	Type         string             `json:"type"`
	Message      *Message           `json:"message,omitempty"`
	Messages     []*Message         `json:"messages,omitempty"`
	Typing       *TypingEvent       `json:"typing,omitempty"`
	Presence     *PresenceEvent     `json:"presence,omitempty"`
	Notification *NotificationEvent `json:"notification,omitempty"`
}

type inboundEvent struct {
	Type            string `json:"type"`
	Content         string `json:"content"`
	ChannelID       string `json:"channel_id"`
	ParentMessageID string `json:"parent_message_id"`
	IsTyping        bool   `json:"is_typing"`
}

func (c *Client) ReadMessage(core *Core) {
	defer func() {
		core.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, payload, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
			}
			break
		}

		event := parseInboundEvent(c, payload)
		log.Printf("Received websocket event %s from %s in room %s", event.Type, c.Username, c.RoomID)
		core.Broadcast <- event
	}
}

func parseInboundEvent(client *Client, payload []byte) *Event {
	var inbound inboundEvent
	if err := json.Unmarshal(payload, &inbound); err != nil || inbound.Type == "" {
		return &Event{
			Type: "message.created",
			Message: &Message{
				Content:   strings.TrimSpace(string(payload)),
				RoomID:    client.RoomID,
				Username:  client.Username,
				UserID:    client.UserID,
				CreatedAt: time.Now().UTC().Format(time.RFC3339),
			},
		}
	}

	switch inbound.Type {
	case "typing":
		return &Event{
			Type: "typing",
			Typing: &TypingEvent{
				RoomID:    client.RoomID,
				ChannelID: inbound.ChannelID,
				UserID:    client.UserID,
				Username:  client.Username,
				IsTyping:  inbound.IsTyping,
			},
		}
	default:
		return &Event{
			Type: "message.created",
			Message: &Message{
				Content:         strings.TrimSpace(inbound.Content),
				RoomID:          client.RoomID,
				ChannelID:       inbound.ChannelID,
				ParentMessageID: inbound.ParentMessageID,
				Username:        client.Username,
				UserID:          client.UserID,
				CreatedAt:       time.Now().UTC().Format(time.RFC3339),
			},
		}
	}
}

func (c *Client) WriteMessage() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		event, ok := <-c.Message
		if !ok {
			return
		}
		if err := c.Conn.WriteJSON(event); err != nil {
			log.Printf("error writing websocket event: %v", err)
			return
		}
	}
}
