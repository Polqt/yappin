package websocket

import (
	"context"
	"database/sql"
	"log"

	roomRepository "chat-application/internal/repo/room"
	statsRepository "chat-application/internal/repo/stats"

	"github.com/google/uuid"
)

type Room struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Clients map[string]*Client `json:"clients"`
	History []*Message 
	IsPinned bool `json:"is_pinned"`
	TopicTitle *string `json:"topic_title,omitempty"`
	TopicDescription *string `json:"topic_description,omitempty"`
	TopicURL *string `json:"topic_url,omitempty"`
	TopicSource *string `json:"topic_source,omitempty"`
}

type Core struct {
	Rooms map[string]*Room
	Register chan *Client
	Unregister chan *Client
	Broadcast chan *Message
	RoomRepository *roomRepository.RoomRepository
	StatsRepository *statsRepository.StatsRepository
	db *sql.DB
}

func NewCore(db *sql.DB) *Core {
	return &Core{
		Rooms: make(map[string]*Room),
		Register: make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast: make(chan *Message, 5),
		RoomRepository: roomRepository.NewRoomRepository(db),
		StatsRepository: statsRepository.NewStatsRepository(db),
		db: db,
	}
}

func (c *Core) GetDB() *sql.DB  {
	return c.db
}

func (c *Core) Start() {
	for {
		select {
		case client := <-c.Register:
			if room, ok := c.Rooms[client.RoomID]; ok {
				if _, ok := room.Clients[client.ID]; !ok {
					room.Clients[client.ID] = client
				}
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
							Content: msg.Content,
							RoomID: client.RoomID,
							Username: msg.Username,
							UserID: userID,
							System: msg.IsSystem,
						}
						client.Message <- websocketMsg
					}
				}()
			}
		
		case client := <-c.Unregister:
			if room, ok := c.Rooms[client.RoomID]; ok {
				if _, ok := room.Clients[client.ID]; ok {
					delete(c.Rooms[client.RoomID].Clients, client.ID)
					close(client.Message)
				}
			}

		case message := <-c.Broadcast:
			if room, ok := c.Rooms[message.RoomID]; ok {
				room.History = append(room.History, message)

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
						RoomID: roomUUID,
						UserID: userID,
						Username: msg.Username,
						Content: msg.Content,
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
				for _, client := range room.Clients {
					client.Message <- message
				}
			}
		}
		
	
	}
}