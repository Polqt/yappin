package handler

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"chat-application/internal/api/model"
	"chat-application/internal/middleware"
	roomRepository "chat-application/internal/repo/room"
	statsRepository "chat-application/internal/repo/stats"
	websoc "chat-application/internal/websocket"

	"github.com/google/uuid"
)

type fakeRoomRepository struct {
	createRoomFn       func(ctx context.Context, room *roomRepository.Room) (*roomRepository.Room, error)
	getRoomByIDFn      func(ctx context.Context, id uuid.UUID) (*roomRepository.Room, error)
	countActiveRoomsFn func(ctx context.Context) (int, error)
	getAllActiveFn     func(ctx context.Context) ([]*roomRepository.Room, error)
	createMessageFn    func(ctx context.Context, message *roomRepository.Message) (*roomRepository.Message, error)
	getMessagesFn      func(ctx context.Context, roomID uuid.UUID, limit int, offset int) ([]*roomRepository.Message, error)
}

func (f *fakeRoomRepository) GetDB() *sql.DB { return nil }
func (f *fakeRoomRepository) CreateRoom(ctx context.Context, room *roomRepository.Room) (*roomRepository.Room, error) {
	return f.createRoomFn(ctx, room)
}
func (f *fakeRoomRepository) GetRoomByID(ctx context.Context, id uuid.UUID) (*roomRepository.Room, error) {
	if f.getRoomByIDFn != nil {
		return f.getRoomByIDFn(ctx, id)
	}
	return nil, nil
}
func (f *fakeRoomRepository) CountActiveRooms(ctx context.Context) (int, error) {
	if f.countActiveRoomsFn != nil {
		return f.countActiveRoomsFn(ctx)
	}
	return 0, nil
}
func (f *fakeRoomRepository) GetAllActiveRooms(ctx context.Context) ([]*roomRepository.Room, error) {
	if f.getAllActiveFn != nil {
		return f.getAllActiveFn(ctx)
	}
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

type fakeStatsRepository struct{}

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
	return nil
}
func (f *fakeStatsRepository) GetLeaderboard(ctx context.Context, limit int) ([]statsRepository.LeaderboardEntry, error) {
	return nil, nil
}

func TestCreateRoomSuccess(t *testing.T) {
	roomID := uuid.New()
	var createdBy *uuid.UUID

	repo := &fakeRoomRepository{
		countActiveRoomsFn: func(ctx context.Context) (int, error) {
			return 0, nil
		},
		createRoomFn: func(ctx context.Context, room *roomRepository.Room) (*roomRepository.Room, error) {
			createdBy = room.CreatorID
			room.ID = roomID
			room.CreatedAt = time.Now()
			return room, nil
		},
	}

	core := websoc.NewCoreWithDependencies(nil, repo, &fakeStatsRepository{})
	handler := NewCoreHandlerWithRoomRepository(core, repo)

	body := bytes.NewBufferString(`{"name":"General"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/websoc/create-room", body)
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, uuid.New().String()))
	rec := httptest.NewRecorder()

	handler.CreateRoom(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if createdBy == nil {
		t.Fatal("expected creator id to be set from auth context")
	}

	var response model.CreateRoomReq
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if response.ID != roomID.String() {
		t.Fatalf("expected room id %s, got %s", roomID, response.ID)
	}
}

func TestGetRoomsReturnsParticipantCounts(t *testing.T) {
	roomID := uuid.New()
	repo := &fakeRoomRepository{
		getAllActiveFn: func(ctx context.Context) ([]*roomRepository.Room, error) {
			return []*roomRepository.Room{
				{
					ID:        roomID,
					Name:      "General",
					CreatedAt: time.Now(),
					ExpiresAt: time.Now().Add(time.Hour),
				},
			}, nil
		},
	}

	core := websoc.NewCoreWithDependencies(nil, repo, &fakeStatsRepository{})
	core.AddRoom(&websoc.Room{
		ID:      roomID.String(),
		Name:    "General",
		Clients: map[string]*websoc.Client{"one": {ID: "one"}, "two": {ID: "two"}},
	})
	handler := NewCoreHandlerWithRoomRepository(core, repo)

	req := httptest.NewRequest(http.MethodGet, "/api/websoc/get-rooms", nil)
	rec := httptest.NewRecorder()

	handler.GetRooms(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var rooms []model.RoomRes
	if err := json.NewDecoder(rec.Body).Decode(&rooms); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(rooms) != 1 {
		t.Fatalf("expected one room, got %d", len(rooms))
	}
	if rooms[0].Participants != 2 {
		t.Fatalf("expected 2 participants, got %d", rooms[0].Participants)
	}
}
