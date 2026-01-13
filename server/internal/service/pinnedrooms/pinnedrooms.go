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
	roomRepo      *roomRepository.RoomRepository
	topicsService *topics.TopicsService
	websocketCore *websocket.Core
}

func NewPinnedRoomsService(db *sql.DB, websocketCore *websocket.Core) *PinnedRoomsService {
	return &PinnedRoomsService{
		roomRepo:      roomRepository.NewRoomRepository(db),
		topicsService: topics.NewTopicsService(),
		websocketCore: websocketCore,
	}
}

func getNextMidnightUTC() time.Time {
	now := time.Now().UTC()
	midnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.UTC)
	return midnight
}

func (s *PinnedRoomsService) RefreshPinnedRooms(ctx context.Context) error {
	topics, err := s.topicsService.FetchAllTopics(ctx)
	if err != nil {
		return fmt.Errorf("fetch topics: %w", err)
	}

	expiresAt := getNextMidnightUTC()
	now := time.Now()

	roomNames := []string{"Technology", "Sports", "Entertainment"}

	for i, topic := range topics {
		if i >= len(roomNames) {
			break
		}

		room := &roomRepository.Room{
			Name:             roomNames[i],
			IsPinned:         true,
			TopicTitle:       &topic.Title,
			TopicDescription: &topic.Description,
			TopicURL:         &topic.URL,
			TopicSource:      &topic.Source,
			TopicUpdatedAt:   &now,
			ExpiresAt:        expiresAt,
		}

		createdRoom, err := s.roomRepo.CreateRoom(ctx, room)
		if err != nil {
			log.Printf("Failed to create pinned room: %v", err)
			continue
		}

		s.websocketCore.AddRoom(&websocket.Room{
			ID:               createdRoom.ID.String(),
			Name:             createdRoom.Name,
			Clients:          make(map[string]*websocket.Client),
			IsPinned:         createdRoom.IsPinned,
			TopicTitle:       createdRoom.TopicTitle,
			TopicDescription: createdRoom.TopicDescription,
			TopicURL:         createdRoom.TopicURL,
			TopicSource:      createdRoom.TopicSource,
		})
		log.Printf("Pinned room created: %s", createdRoom.Name)
	}
	return nil
}

func (s *PinnedRoomsService) CheckAndRefreshPinnedRooms(ctx context.Context) error {
	count, err := s.roomRepo.CountPinnedRooms(ctx)
	if err != nil {
		return fmt.Errorf("count pinned rooms: %w", err)
	}

	if count < 3 {
		return s.RefreshPinnedRooms(ctx)
	}
	return nil
}
