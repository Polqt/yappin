package repository

import (
	"context"

	"chat-application/internal/api/model"

	"github.com/google/uuid"
)

type RoomRepositoryInterface interface {
	// CreateRoom creates a new chat room in the database.
	CreateRoom(ctx context.Context, room *Room) (*Room, error)

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
}

// Ensure RoomRepository implements RoomRepositoryInterface
var _ RoomRepositoryInterface = (*RoomRepository)(nil)
