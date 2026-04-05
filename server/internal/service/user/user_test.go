package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"chat-application/internal/api/model"
	repository "chat-application/internal/repo/user"
	"chat-application/util"

	"github.com/google/uuid"
)

type fakeUserRepository struct {
	getByEmailFn     func(ctx context.Context, email string) (*repository.User, error)
	createUserFn     func(ctx context.Context, user *repository.User) (*repository.User, error)
	getByIDFn        func(ctx context.Context, id uuid.UUID) (*repository.User, error)
	updateUsernameFn func(ctx context.Context, id uuid.UUID, username string) (*repository.User, error)
	deleteUserFn     func(ctx context.Context, id uuid.UUID) error
}

func (f *fakeUserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*repository.User, error) {
	if f.getByIDFn != nil {
		return f.getByIDFn(ctx, id)
	}
	return nil, nil
}

func (f *fakeUserRepository) GetUserByEmail(ctx context.Context, email string) (*repository.User, error) {
	if f.getByEmailFn != nil {
		return f.getByEmailFn(ctx, email)
	}
	return nil, nil
}

func (f *fakeUserRepository) CreateUser(ctx context.Context, user *repository.User) (*repository.User, error) {
	if f.createUserFn != nil {
		return f.createUserFn(ctx, user)
	}
	return nil, nil
}

func (f *fakeUserRepository) UpdateUsername(ctx context.Context, id uuid.UUID, username string) (*repository.User, error) {
	if f.updateUsernameFn != nil {
		return f.updateUsernameFn(ctx, id, username)
	}
	return nil, nil
}

func (f *fakeUserRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	if f.deleteUserFn != nil {
		return f.deleteUserFn(ctx, id)
	}
	return nil
}

func TestUserServiceLoginSuccess(t *testing.T) {
	t.Setenv("JWT_SECRET_KEY", "test-secret")

	hashedPassword, err := util.HashPassword("super-secret-pass")
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	userID := uuid.New()
	repo := &fakeUserRepository{
		getByEmailFn: func(ctx context.Context, email string) (*repository.User, error) {
			return &repository.User{
				ID:           userID,
				Username:     "alice",
				Email:        email,
				PasswordHash: &hashedPassword,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			}, nil
		},
	}

	service := NewUserService(repo)

	result, err := service.Login(context.Background(), model.RequestLoginUser{
		Email:    "alice@example.com",
		Password: "super-secret-pass",
	})
	if err != nil {
		t.Fatalf("expected login to succeed, got error: %v", err)
	}
	if result.ID != userID.String() {
		t.Fatalf("expected user id %s, got %s", userID, result.ID)
	}
	if result.AccessToken == "" {
		t.Fatal("expected access token to be present")
	}
}

func TestUserServiceCreateUserDuplicateError(t *testing.T) {
	t.Setenv("JWT_SECRET_KEY", "test-secret")

	repo := &fakeUserRepository{
		createUserFn: func(ctx context.Context, user *repository.User) (*repository.User, error) {
			return nil, errors.New("duplicate key value violates unique constraint")
		},
	}

	service := NewUserService(repo)

	_, err := service.CreateUser(context.Background(), model.RequestCreateUser{
		Username: "alice",
		Email:    "alice@example.com",
		Password: "Super-secret-pass1",
	})
	if err == nil {
		t.Fatal("expected duplicate user error")
	}
	if err.Error() != "username or email already exists" {
		t.Fatalf("unexpected error: %v", err)
	}
}
