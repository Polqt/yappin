package repository

import (
	"context"

	"github.com/google/uuid"
)

type UserRepositoryInterface interface {
	// GetUserByID retrieves a user by their unique identifier.
	// Returns nil, nil if the user is not found.
	GetUserByID(ctx context.Context, id uuid.UUID) (*User, error)

	// GetUserByEmail retrieves a user by their email address.
	// Returns nil, nil if the user is not found.
	GetUserByEmail(ctx context.Context, email string) (*User, error)

	// CreateUser creates a new user in the database.
	// Returns the created user with generated ID and timestamps.
	CreateUser(ctx context.Context, user *User) (*User, error)

	// UpdateUsername updates a user's username.
	// Returns the updated user or an error if the username is taken.
	UpdateUsername(ctx context.Context, id uuid.UUID, username string) (*User, error)

	// DeleteUser removes a user from the database by their ID.
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

// Ensure UserRepository implements UserRepositoryInterface
var _ UserRepositoryInterface = (*UserRepository)(nil)
