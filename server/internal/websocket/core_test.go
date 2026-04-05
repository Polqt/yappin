package websocket

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"chat-application/internal/api/model"
	roomRepository "chat-application/internal/repo/room"
	statsRepository "chat-application/internal/repo/stats"

	"github.com/google/uuid"
)

type fakeRoomRepository struct {
	getMessagesFn   func(ctx context.Context, roomID uuid.UUID, limit int, offset int) ([]*roomRepository.Message, error)
	createMessageFn func(ctx context.Context, message *roomRepository.Message) (*roomRepository.Message, error)
}

func (f *fakeRoomRepository) GetDB() *sql.DB { return nil }
func (f *fakeRoomRepository) CreateRoom(ctx context.Context, room *roomRepository.Room) (*roomRepository.Room, error) {
	return room, nil
}
func (f *fakeRoomRepository) GetRoomByID(ctx context.Context, id uuid.UUID) (*roomRepository.Room, error) {
	return nil, nil
}
func (f *fakeRoomRepository) CountActiveRooms(ctx context.Context) (int, error) { return 0, nil }
func (f *fakeRoomRepository) GetAllActiveRooms(ctx context.Context) ([]*roomRepository.Room, error) {
	return nil, nil
}
func (f *fakeRoomRepository) CreateMessage(ctx context.Context, message *roomRepository.Message) (*roomRepository.Message, error) {
	if f.createMessageFn != nil {
		return f.createMessageFn(ctx, message)
	}
	return message, nil
}
func (f *fakeRoomRepository) GetRoomMessages(ctx context.Context, roomID uuid.UUID, limit int, offset int) ([]*roomRepository.Message, error) {
	if f.getMessagesFn != nil {
		return f.getMessagesFn(ctx, roomID, limit, offset)
	}
	return nil, nil
}
func (f *fakeRoomRepository) CountPinnedRooms(ctx context.Context) (int, error)   { return 0, nil }
func (f *fakeRoomRepository) DeleteExpiredRooms(ctx context.Context) (int, error) { return 0, nil }
func (f *fakeRoomRepository) EnsureRoomMembership(ctx context.Context, roomID, userID uuid.UUID) error {
	return nil
}
func (f *fakeRoomRepository) GetRoomMember(ctx context.Context, roomID, userID uuid.UUID) (*roomRepository.RoomMember, error) {
	return nil, nil
}
func (f *fakeRoomRepository) GetRoomMembers(ctx context.Context, roomID uuid.UUID) ([]roomRepository.RoomMember, error) {
	return nil, nil
}
func (f *fakeRoomRepository) UpdateRoomMember(ctx context.Context, member roomRepository.RoomMember) error {
	return nil
}
func (f *fakeRoomRepository) CreateCategory(ctx context.Context, category *roomRepository.RoomCategory) (*roomRepository.RoomCategory, error) {
	return category, nil
}
func (f *fakeRoomRepository) CreateChannel(ctx context.Context, channel *roomRepository.RoomChannel) (*roomRepository.RoomChannel, error) {
	return channel, nil
}
func (f *fakeRoomRepository) GetRoomCategories(ctx context.Context, roomID uuid.UUID) ([]roomRepository.RoomCategory, error) {
	return nil, nil
}
func (f *fakeRoomRepository) GetRoomChannels(ctx context.Context, roomID uuid.UUID) ([]roomRepository.RoomChannel, error) {
	return nil, nil
}
func (f *fakeRoomRepository) GetDefaultChannel(ctx context.Context, roomID uuid.UUID) (*roomRepository.RoomChannel, error) {
	return nil, nil
}
func (f *fakeRoomRepository) GetRoomMessagesByChannel(ctx context.Context, roomID uuid.UUID, channelID *uuid.UUID, limit int, offset int) ([]*roomRepository.Message, error) {
	return f.GetRoomMessages(ctx, roomID, limit, offset)
}
func (f *fakeRoomRepository) SearchMessages(ctx context.Context, roomID uuid.UUID, queryText string, channelID *uuid.UUID, username string, limit int) ([]roomRepository.Message, error) {
	return nil, nil
}
func (f *fakeRoomRepository) CreateNotification(ctx context.Context, notification *roomRepository.Notification) error {
	return nil
}
func (f *fakeRoomRepository) GetNotifications(ctx context.Context, userID uuid.UUID, limit int) ([]roomRepository.Notification, error) {
	return nil, nil
}
func (f *fakeRoomRepository) MarkNotificationRead(ctx context.Context, notificationID, userID uuid.UUID) error {
	return nil
}
func (f *fakeRoomRepository) CreateMentionNotifications(ctx context.Context, roomID uuid.UUID, message *roomRepository.Message) ([]roomRepository.Notification, error) {
	return nil, nil
}
func (f *fakeRoomRepository) AddReaction(ctx context.Context, reaction *model.MessageReaction) error {
	return nil
}
func (f *fakeRoomRepository) GetReactions(ctx context.Context, messageID string) ([]model.MessageReaction, error) {
	return nil, nil
}

type fakeStatsRepository struct {
	incremented []uuid.UUID
}

func (f *fakeStatsRepository) GetOrCreateUserStats(ctx context.Context, userID uuid.UUID) (*statsRepository.UserStats, error) {
	return nil, nil
}
func (f *fakeStatsRepository) GetUserProfile(ctx context.Context, userID uuid.UUID) (*statsRepository.UserStats, error) {
	return nil, nil
}
func (f *fakeStatsRepository) ProcessDailyCheckin(ctx context.Context, userID uuid.UUID) (int, bool, error) {
	return 0, false, nil
}
func (f *fakeStatsRepository) GetUserAchievements(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	return nil, nil
}
func (f *fakeStatsRepository) GetUserAchievementsDetails(ctx context.Context, userID uuid.UUID) ([]statsRepository.Achievement, error) {
	return nil, nil
}
func (f *fakeStatsRepository) GetAllAchievementTypes(ctx context.Context) ([]statsRepository.Achievement, error) {
	return nil, nil
}
func (f *fakeStatsRepository) CheckAwardsAndAchievements(ctx context.Context, userID uuid.UUID) ([]statsRepository.Achievement, error) {
	return nil, nil
}
func (f *fakeStatsRepository) CanUserUpvote(ctx context.Context, fromUserID, toUserID uuid.UUID) (bool, error) {
	return false, nil
}
func (f *fakeStatsRepository) GiveUpvote(ctx context.Context, fromUserID, toUserID uuid.UUID) error {
	return nil
}
func (f *fakeStatsRepository) IncrementMessageCount(ctx context.Context, userID uuid.UUID) error {
	f.incremented = append(f.incremented, userID)
	return nil
}
func (f *fakeStatsRepository) GetLeaderboard(ctx context.Context, limit int) ([]statsRepository.LeaderboardEntry, error) {
	return nil, nil
}

func TestCoreRegisterLoadsRoomHistory(t *testing.T) {
	roomID := uuid.New()
	repo := &fakeRoomRepository{
		getMessagesFn: func(ctx context.Context, gotRoomID uuid.UUID, limit int, offset int) ([]*roomRepository.Message, error) {
			if gotRoomID != roomID {
				t.Fatalf("expected room id %s, got %s", roomID, gotRoomID)
			}
			return []*roomRepository.Message{
				{
					RoomID:    roomID,
					Username:  "alice",
					Content:   "hello from history",
					CreatedAt: time.Now(),
				},
			}, nil
		},
	}

	core := NewCoreWithDependencies(nil, repo, &fakeStatsRepository{})
	core.AddRoom(&Room{
		ID:      roomID.String(),
		Name:    "General",
		Clients: make(map[string]*Client),
	})

	go core.Start()

	client := &Client{
		ID:       "client-1",
		RoomID:   roomID.String(),
		Username: "bob",
		Message:  make(chan *Event, 2),
	}

	core.Register <- client

	select {
	case event := <-client.Message:
		if event.Type != "history" || len(event.Messages) != 1 {
			t.Fatalf("expected history event with one message, got %+v", event)
		}
		if event.Messages[0].Content != "hello from history" {
			t.Fatalf("expected history message, got %q", event.Messages[0].Content)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for history message")
	}
}

func TestCoreBroadcastPersistsAndFanoutsMessage(t *testing.T) {
	roomID := uuid.New()
	userID := uuid.New()
	var persisted *roomRepository.Message

	repo := &fakeRoomRepository{
		createMessageFn: func(ctx context.Context, message *roomRepository.Message) (*roomRepository.Message, error) {
			persisted = message
			return message, nil
		},
	}
	statsRepo := &fakeStatsRepository{}
	core := NewCoreWithDependencies(nil, repo, statsRepo)
	client := &Client{
		ID:       "client-1",
		RoomID:   roomID.String(),
		Username: "listener",
		Message:  make(chan *Event, 2),
	}
	core.AddRoom(&Room{
		ID:      roomID.String(),
		Name:    "General",
		Clients: map[string]*Client{client.ID: client},
	})

	go core.Start()

	core.Broadcast <- &Event{
		Type: "message.created",
		Message: &Message{
			Content:  "hello world",
			RoomID:   roomID.String(),
			Username: "alice",
			UserID:   userID.String(),
		},
	}

	select {
	case event := <-client.Message:
		if event.Type != "message.created" || event.Message == nil {
			t.Fatalf("expected message.created event, got %+v", event)
		}
		if event.Message.Content != "hello world" {
			t.Fatalf("expected broadcast content, got %q", event.Message.Content)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for broadcast")
	}

	if persisted == nil {
		deadline := time.Now().Add(2 * time.Second)
		for persisted == nil && time.Now().Before(deadline) {
			time.Sleep(10 * time.Millisecond)
		}
	}
	if persisted == nil {
		t.Fatal("expected message to be persisted")
	}
	if persisted.Content != "hello world" {
		t.Fatalf("expected persisted content %q, got %q", "hello world", persisted.Content)
	}
	if len(statsRepo.incremented) != 1 || statsRepo.incremented[0] != userID {
		deadline := time.Now().Add(2 * time.Second)
		for len(statsRepo.incremented) != 1 && time.Now().Before(deadline) {
			time.Sleep(10 * time.Millisecond)
		}
	}
	if len(statsRepo.incremented) != 1 || statsRepo.incremented[0] != userID {
		t.Fatalf("expected stats increment for user %s, got %+v", userID, statsRepo.incremented)
	}
}
