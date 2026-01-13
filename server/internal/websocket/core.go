package websocket

import (
	"context"
	"database/sql"
	"log"
	"sync"

	"chat-application/internal/constants"
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

// AddMessage adds a message to the room history with a size limit
func (r *Room) AddMessage(msg *Message) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(r.History) >= constants.MaxRoomHistory {
		r.History = r.History[1:]
	}
	r.History = append(r.History, msg)
}

// GetHistory returns a copy of the room history
func (r *Room) GetHistory() []*Message {
	r.mu.RLock()
	defer r.mu.RUnlock()
	history := make([]*Message, len(r.History))
	copy(history, r.History)
	return history
}

type Core struct {
	Rooms           map[string]*Room
	roomsMu         sync.RWMutex
	Register        chan *Client
	Unregister      chan *Client
	Broadcast       chan *Message
	RoomRepository  *roomRepository.RoomRepository
	StatsRepository *statsRepository.StatsRepository
	db              *sql.DB
}

func NewCore(db *sql.DB) *Core {
	return &Core{
		Rooms:           make(map[string]*Room),
		Register:        make(chan *Client),
		Unregister:      make(chan *Client),
		Broadcast:       make(chan *Message, 5),
		RoomRepository:  roomRepository.NewRoomRepository(db),
		StatsRepository: statsRepository.NewStatsRepository(db),
		db:              db,
	}
}

func (c *Core) GetDB() *sql.DB {
	return c.db
}

// GetRoom safely retrieves a room by ID
func (c *Core) GetRoom(roomID string) (*Room, bool) {
	c.roomsMu.RLock()
	defer c.roomsMu.RUnlock()
	room, ok := c.Rooms[roomID]
	return room, ok
}

// AddRoom safely adds a room
func (c *Core) AddRoom(room *Room) {
	c.roomsMu.Lock()
	defer c.roomsMu.Unlock()
	c.Rooms[room.ID] = room
}

// DeleteRoom safely removes a room
func (c *Core) DeleteRoom(roomID string) {
	c.roomsMu.Lock()
	defer c.roomsMu.Unlock()
	delete(c.Rooms, roomID)
}

// GetAllRooms returns a snapshot of all rooms
func (c *Core) GetAllRooms() map[string]*Room {
	c.roomsMu.RLock()
	defer c.roomsMu.RUnlock()
	rooms := make(map[string]*Room, len(c.Rooms))
	for k, v := range c.Rooms {
		rooms[k] = v
	}
	return rooms
}

func (c *Core) Start() {
	for {
		select {
		case client := <-c.Register:
			room, ok := c.GetRoom(client.RoomID)
			if ok {
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

					messages, err := c.RoomRepository.GetRoomMessages(context.Background(), roomUUID, 100, 0)
					if err != nil {
						log.Printf("error fetching room messages: %v", err)
						return
					}

					for _, msg := range messages {
						userID := ""
						if msg.UserID != nil {
							userID = msg.UserID.String()
						}

						websocketMsg := &Message{
							Content:  msg.Content,
							RoomID:   client.RoomID,
							Username: msg.Username,
							UserID:   userID,
							System:   msg.IsSystem,
						}
						client.Message <- websocketMsg
					}
				}()
			}

		case client := <-c.Unregister:
			room, ok := c.GetRoom(client.RoomID)
			if ok {
				room.mu.Lock()
				if _, exists := room.Clients[client.ID]; exists {
					delete(room.Clients, client.ID)
					close(client.Message)
				}
				room.mu.Unlock()
			}

		case message := <-c.Broadcast:
			room, ok := c.GetRoom(message.RoomID)
			if ok {
				room.AddMessage(message)

				go func(msg *Message) {
					roomUUID, err := uuid.Parse(msg.RoomID)
					if err != nil {
						log.Printf("error parsing room ID: %v", err)
						return
					}

					var userID *uuid.UUID
					if msg.UserID != "" {
						if parsedUserID, err := uuid.Parse(msg.UserID); err == nil {
							userID = &parsedUserID
						}
					}

					dbMessage := &roomRepository.Message{
						RoomID:   roomUUID,
						UserID:   userID,
						Username: msg.Username,
						Content:  msg.Content,
						IsSystem: msg.System,
					}

					if _, err := c.RoomRepository.CreateMessage(context.Background(), dbMessage); err != nil {
						log.Printf("error creating message in database: %v", err)
						return
					}

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
				}(message)

				// Send to all clients in room
				room.mu.RLock()
				for _, client := range room.Clients {
					client.Message <- message
				}
				room.mu.RUnlock()
			}
		}
	}
}
