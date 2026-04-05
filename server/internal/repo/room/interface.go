package repository

import (
	"context"
	"database/sql"

	"chat-application/internal/api/model"

	"github.com/google/uuid"
)

type RoomRepositoryInterface interface {
	// CreateRoom creates a new chat room in the database.
	CreateRoom(ctx context.Context, room *Room) (*Room, error)

	// GetDB returns the underlying database connection when direct queries are needed.
	GetDB() *sql.DB

	// GetRoomByID retrieves a room by its unique identifier.
	// Returns nil, nil if the room is not found.
	GetRoomByID(ctx context.Context, id uuid.UUID) (*Room, error)

	// CountActiveRooms returns the count of non-expired rooms.
	CountActiveRooms(ctx context.Context) (int, error)

	// GetAllActiveRooms retrieves all non-expired rooms.
	GetAllActiveRooms(ctx context.Context) ([]*Room, error)

	// CreateMessage creates a new message in a room.
	CreateMessage(ctx context.Context, message *Message) (*Message, error)

	// GetRoomMessages retrieves messages for a room with pagination.
	GetRoomMessages(ctx context.Context, roomID uuid.UUID, limit int, offset int) ([]*Message, error)

	// CountPinnedRooms returns the count of pinned rooms.
	CountPinnedRooms(ctx context.Context) (int, error)

	// DeleteExpiredRooms removes all expired rooms from the database.
	DeleteExpiredRooms(ctx context.Context) (int, error)

	// AddReaction adds a reaction to a message.
	AddReaction(ctx context.Context, reaction *model.MessageReaction) error

	// GetReactions retrieves all reactions for a message.
	GetReactions(ctx context.Context, messageID string) ([]model.MessageReaction, error)

	EnsureRoomMembership(ctx context.Context, roomID, userID uuid.UUID) error
	GetRoomMember(ctx context.Context, roomID, userID uuid.UUID) (*RoomMember, error)
	GetRoomMembers(ctx context.Context, roomID uuid.UUID) ([]RoomMember, error)
	UpdateRoomMember(ctx context.Context, member RoomMember) error
	CreateCategory(ctx context.Context, category *RoomCategory) (*RoomCategory, error)
	CreateChannel(ctx context.Context, channel *RoomChannel) (*RoomChannel, error)
	GetRoomCategories(ctx context.Context, roomID uuid.UUID) ([]RoomCategory, error)
	GetRoomChannels(ctx context.Context, roomID uuid.UUID) ([]RoomChannel, error)
	GetDefaultChannel(ctx context.Context, roomID uuid.UUID) (*RoomChannel, error)
	GetRoomMessagesByChannel(ctx context.Context, roomID uuid.UUID, channelID *uuid.UUID, limit int, offset int) ([]*Message, error)
	SearchMessages(ctx context.Context, roomID uuid.UUID, queryText string, channelID *uuid.UUID, username string, limit int) ([]Message, error)
	CreateNotification(ctx context.Context, notification *Notification) error
	GetNotifications(ctx context.Context, userID uuid.UUID, limit int) ([]Notification, error)
	MarkNotificationRead(ctx context.Context, notificationID, userID uuid.UUID) error
	CreateMentionNotifications(ctx context.Context, roomID uuid.UUID, message *Message) ([]Notification, error)
}

// Ensure RoomRepository implements RoomRepositoryInterface
var _ RoomRepositoryInterface = (*RoomRepository)(nil)
