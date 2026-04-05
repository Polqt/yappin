package websocket

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"sync"
	"time"

	roomRepository "chat-application/internal/repo/room"
	statsRepository "chat-application/internal/repo/stats"

	"github.com/google/uuid"
)

type Room struct {
	ID               string             `json:"id"`
	Name             string             `json:"name"`
	Clients          map[string]*Client `json:"clients"`
	History          []*Message
	IsPinned         bool    `json:"is_pinned"`
	TopicTitle       *string `json:"topic_title,omitempty"`
	TopicDescription *string `json:"topic_description,omitempty"`
	TopicURL         *string `json:"topic_url,omitempty"`
	TopicSource      *string `json:"topic_source,omitempty"`
	mu               sync.RWMutex
}

func (r *Room) AddMessage(msg *Message) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(r.History) >= 100 {
		r.History = r.History[1:]
	}
	r.History = append(r.History, msg)
}

type Core struct {
	Rooms           map[string]*Room
	roomsMu         sync.RWMutex
	Register        chan *Client
	Unregister      chan *Client
	Broadcast       chan *Event
	RoomRepository  roomRepository.RoomRepositoryInterface
	StatsRepository statsRepository.StatsRepositoryInterface
	db              *sql.DB
}

func NewCore(db *sql.DB) *Core {
	return NewCoreWithDependencies(
		db,
		roomRepository.NewRoomRepository(db),
		statsRepository.NewStatsRepository(db),
	)
}

func NewCoreWithDependencies(
	db *sql.DB,
	roomRepo roomRepository.RoomRepositoryInterface,
	statsRepo statsRepository.StatsRepositoryInterface,
) *Core {
	return &Core{
		Rooms:           make(map[string]*Room),
		Register:        make(chan *Client),
		Unregister:      make(chan *Client),
		Broadcast:       make(chan *Event, 16),
		RoomRepository:  roomRepo,
		StatsRepository: statsRepo,
		db:              db,
	}
}

func (c *Core) GetDB() *sql.DB {
	return c.db
}

func (c *Core) GetRoom(roomID string) (*Room, bool) {
	c.roomsMu.RLock()
	defer c.roomsMu.RUnlock()
	room, ok := c.Rooms[roomID]
	return room, ok
}

func (c *Core) AddRoom(room *Room) {
	c.roomsMu.Lock()
	defer c.roomsMu.Unlock()
	c.Rooms[room.ID] = room
}

func (c *Core) DeleteRoom(roomID string) {
	c.roomsMu.Lock()
	defer c.roomsMu.Unlock()
	delete(c.Rooms, roomID)
}

func (c *Core) Start() {
	for {
		select {
		case client := <-c.Register:
			c.registerClient(client)
		case client := <-c.Unregister:
			c.unregisterClient(client)
		case event := <-c.Broadcast:
			c.handleEvent(event)
		}
	}
}

func (c *Core) registerClient(client *Client) {
	room, ok := c.GetRoom(client.RoomID)
	if !ok {
		return
	}

	room.mu.Lock()
	if _, exists := room.Clients[client.ID]; !exists {
		room.Clients[client.ID] = client
	}
	room.mu.Unlock()

	go func() {
		roomUUID, err := uuid.Parse(client.RoomID)
		if err != nil {
			log.Printf("error parsing room ID: %v", err)
			return
		}

		defaultChannel, err := c.RoomRepository.GetDefaultChannel(context.Background(), roomUUID)
		if err != nil {
			log.Printf("error fetching default channel: %v", err)
			return
		}

		var channelID *uuid.UUID
		if defaultChannel != nil {
			channelID = &defaultChannel.ID
		}

		messages, err := c.RoomRepository.GetRoomMessagesByChannel(context.Background(), roomUUID, channelID, 100, 0)
		if err != nil {
			log.Printf("error fetching room messages: %v", err)
			return
		}

		history := make([]*Message, 0, len(messages))
		for _, msg := range messages {
			history = append(history, mapRepositoryMessage(msg))
		}

		client.Message <- &Event{
			Type:     "history",
			Messages: history,
		}
		client.Message <- &Event{
			Type:     "presence",
			Presence: c.buildPresenceSnapshot(client.RoomID),
		}
		c.emitPresence(client.RoomID)
	}()
}

func (c *Core) unregisterClient(client *Client) {
	room, ok := c.GetRoom(client.RoomID)
	if !ok {
		return
	}

	room.mu.Lock()
	if _, exists := room.Clients[client.ID]; exists {
		delete(room.Clients, client.ID)
		close(client.Message)
	}
	room.mu.Unlock()
	c.emitPresence(client.RoomID)
}

func (c *Core) handleEvent(event *Event) {
	if event == nil {
		return
	}

	switch event.Type {
	case "typing":
		if event.Typing != nil {
			c.fanout(event.Typing.RoomID, event, "")
		}
	case "notification":
		if event.Notification != nil {
			c.fanout(event.Notification.RoomID, event, "")
		}
	case "message.created":
		c.handleMessageCreated(event)
	}
}

func (c *Core) handleMessageCreated(event *Event) {
	if event.Message == nil {
		return
	}

	room, ok := c.GetRoom(event.Message.RoomID)
	if !ok {
		return
	}

	message := event.Message
	if message.CreatedAt == "" {
		message.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	}

	room.AddMessage(message)
	c.fanout(message.RoomID, &Event{Type: "message.created", Message: message}, "")

	go func(msg *Message) {
		roomUUID, err := uuid.Parse(msg.RoomID)
		if err != nil {
			log.Printf("error parsing room ID: %v", err)
			return
		}

		var channelID *uuid.UUID
		if msg.ChannelID != "" {
			if parsed, err := uuid.Parse(msg.ChannelID); err == nil {
				channelID = &parsed
			}
		}

		var parentMessageID *uuid.UUID
		if msg.ParentMessageID != "" {
			if parsed, err := uuid.Parse(msg.ParentMessageID); err == nil {
				parentMessageID = &parsed
			}
		}

		var userID *uuid.UUID
		if msg.UserID != "" {
			if parsedUserID, err := uuid.Parse(msg.UserID); err == nil {
				userID = &parsedUserID
			}
		}

		metadataBytes := []byte(`{}`)
		if len(msg.Metadata) > 0 {
			if encoded, err := json.Marshal(msg.Metadata); err == nil {
				metadataBytes = encoded
			}
		}

		dbMessage := &roomRepository.Message{
			RoomID:          roomUUID,
			ChannelID:       channelID,
			ParentMessageID: parentMessageID,
			UserID:          userID,
			Username:        msg.Username,
			Content:         msg.Content,
			IsSystem:        msg.System,
			Metadata:        metadataBytes,
		}

		createdMessage, err := c.RoomRepository.CreateMessage(context.Background(), dbMessage)
		if err != nil {
			log.Printf("error creating message in database: %v", err)
			return
		}

		msg.ID = createdMessage.ID.String()

		if userID != nil {
			if err := c.StatsRepository.IncrementMessageCount(context.Background(), *userID); err != nil {
				log.Printf("error incrementing message count: %v", err)
			} else {
				go func() {
					_, err := c.StatsRepository.CheckAwardsAndAchievements(context.Background(), *userID)
					if err != nil {
						log.Printf("error checking awards and achievements: %v", err)
					}
				}()
			}
		}

		notifications, err := c.RoomRepository.CreateMentionNotifications(context.Background(), roomUUID, createdMessage)
		if err != nil {
			log.Printf("error creating mention notifications: %v", err)
			return
		}
		for _, notification := range notifications {
			payload := map[string]any{}
			if len(notification.Payload) > 0 {
				_ = json.Unmarshal(notification.Payload, &payload)
			}
			c.Broadcast <- &Event{
				Type: "notification",
				Notification: &NotificationEvent{
					ID:        notification.ID.String(),
					Kind:      notification.Kind,
					Title:     notification.Title,
					Body:      notification.Body,
					RoomID:    roomUUID.String(),
					MessageID: createdMessage.ID.String(),
					Payload:   payload,
				},
			}
		}
	}(message)
}

func (c *Core) fanout(roomID string, event *Event, excludeClientID string) {
	room, ok := c.GetRoom(roomID)
	if !ok {
		return
	}

	room.mu.RLock()
	defer room.mu.RUnlock()
	for _, client := range room.Clients {
		if excludeClientID != "" && client.ID == excludeClientID {
			continue
		}
		select {
		case client.Message <- event:
		default:
			log.Printf("dropping websocket event %s for client %s due to full channel", event.Type, client.ID)
		}
	}
}

func (c *Core) emitPresence(roomID string) {
	snapshot := c.buildPresenceSnapshot(roomID)
	if snapshot == nil {
		return
	}
	c.fanout(roomID, &Event{
		Type:     "presence",
		Presence: snapshot,
	}, "")
}

func (c *Core) buildPresenceSnapshot(roomID string) *PresenceEvent {
	room, ok := c.GetRoom(roomID)
	if !ok {
		return nil
	}

	room.mu.RLock()
	defer room.mu.RUnlock()

	users := make([]PresenceUser, 0, len(room.Clients))
	for _, client := range room.Clients {
		users = append(users, PresenceUser{
			UserID:   client.UserID,
			Username: client.Username,
		})
	}

	return &PresenceEvent{
		RoomID:      roomID,
		OnlineUsers: users,
	}
}

func mapRepositoryMessage(msg *roomRepository.Message) *Message {
	message := &Message{
		ID:        msg.ID.String(),
		Content:   msg.Content,
		RoomID:    msg.RoomID.String(),
		Username:  msg.Username,
		System:    msg.IsSystem,
		CreatedAt: msg.CreatedAt.UTC().Format(time.RFC3339),
	}
	if msg.UserID != nil {
		message.UserID = msg.UserID.String()
	}
	if msg.ChannelID != nil {
		message.ChannelID = msg.ChannelID.String()
	}
	if msg.ParentMessageID != nil {
		message.ParentMessageID = msg.ParentMessageID.String()
	}
	if len(msg.Metadata) > 0 {
		var metadata map[string]any
		if err := json.Unmarshal(msg.Metadata, &metadata); err == nil {
			message.Metadata = metadata
		}
	}
	return message
}
