package stats

import (
	"context"

	"github.com/google/uuid"
)
type StatsRepositoryInterface interface {
	// GetOrCreateUserStats retrieves or creates user statistics.
	GetOrCreateUserStats(ctx context.Context, userID uuid.UUID) (*UserStats, error)

	// GetUserProfile retrieves user profile statistics.
	GetUserProfile(ctx context.Context, userID uuid.UUID) (*UserStats, error)

	// ProcessDailyCheckin handles daily check-in logic.
	// Returns streak count, whether it's a new check-in, and any error.
	ProcessDailyCheckin(ctx context.Context, userID uuid.UUID) (int, bool, error)

	// GetUserAchievements retrieves achievement IDs for a user.
	GetUserAchievements(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)

	// GetUserAchievementsDetails retrieves detailed achievement info for a user.
	GetUserAchievementsDetails(ctx context.Context, userID uuid.UUID) ([]Achievement, error)

	// GetAllAchievementTypes retrieves all available achievement types.
	GetAllAchievementTypes(ctx context.Context) ([]Achievement, error)

	// CheckAwardsAndAchievements checks and awards new achievements.
	CheckAwardsAndAchievements(ctx context.Context, userID uuid.UUID) ([]Achievement, error)

	// CanUserUpvote checks if a user can give an upvote to another user.
	CanUserUpvote(ctx context.Context, fromUserID, toUserID uuid.UUID) (bool, error)

	// GiveUpvote records an upvote from one user to another.
	GiveUpvote(ctx context.Context, fromUserID, toUserID uuid.UUID) error

	// IncrementMessageCount increments the message count for a user.
	IncrementMessageCount(ctx context.Context, userID uuid.UUID) error

	// GetLeaderboard retrieves the top users by score.
	GetLeaderboard(ctx context.Context, limit int) ([]LeaderboardEntry, error)
}

// Ensure StatsRepository implements StatsRepositoryInterface
var _ StatsRepositoryInterface = (*StatsRepository)(nil)
