package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	roomRepository "chat-application/internal/repo/room"
	"chat-application/internal/service/topics"
	"chat-application/internal/websocket"
)

type PinnedRoomsService struct {
	roomRepo *roomRepository.Room
	topicsService *topics.TopicsService
	websocketCore *websocket.Core
}

func NewPinnedRoomsService(db *sql.DB, websocketCore *websocket.Core) *PinnedRoomsService {
	return &PinnedRoomsService{
		roomRepo: roomRepository.NewRoomRepository(db),
		topicsService: topics.NewTopicsService(),
		websocketCore: websocketCore,
	}
}

