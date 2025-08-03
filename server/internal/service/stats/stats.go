package service

import (
	"context"
	"fmt"
	"log"

	"chat-application/internal/repo/stats"

	"github.com/google/uuid"
)

type StatsService struct {
	statsRepository *stats.StatsRepository
}

type CheckinResult struct {
	StreakCount int `json:"streak_count"`
	IsNewCheckin bool `json:"is_new_checkin"`
	NewAchievements []Achievement `json:"new_achievements,omitempty"`
}

type Achievement struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	Icon string `json:"icon"`
	ThresholdType string `json:"threshold_type"`
	ThresholdValue int `json:"threshold_value"`
	EarnedAt string `json:"earned_at"`
}

type UserProfile struct {
	UserID string `db:"user_id"`
	DailyStreak int `db:"daily_streak"`
	TotalCheckins int `db:"total_checkins"`
	TotalMessages int `db:"total_messages"`
	TotalUpvotes int `db:"total_upvotes"`
	CanReceiveUpvotes bool `db:"can_receive_upvotes"`
	Achievements []Achievement `db:"achievements"`
}

func NewStatsService(statsRepository *stats.StatsRepository) *StatsService {
	return &StatsService{
		statsRepository: statsRepository,
	}
}

func (s *StatsService) ProcessDailyCheckin(ctx context.Context, userID uuid.UUID) (*CheckinResult, error) {
	
	streakCount, IsNewCheckin, err := s.statsRepository.ProcessDailyCheckin(ctx, userID)
	if err != nil {
		log.Printf("ProcessDailyCheckin - Error processing checkin for user %s: %v", userID, err)
		return nil, fmt.Errorf("failed to process daily checkin: %w", err)
	}

	if IsNewCheckin {
		log.Printf("ProcessDailyCheckin - New checkin for user %s, streak count: %d", userID, streakCount)
	}

	var newAchievements []Achievement
	if IsNewCheckin {
		achievements, err := s.statsRepository.CheckAwardsAndAchievements(ctx, userID)
		if err != nil {
			log.Printf("ProcessDailyCheckin - Error checking achievements for user %s: %v", userID, err)
			return nil, fmt.Errorf("failed to check achievements: %w", err)
		} else {
			newAchievements = s.
		}
	
	}

	return &CheckinResult{
		StreakCount: streakCount,
		IsNewCheckin: IsNewCheckin,
		// NewAchievements: newAchievements,
	}, nil
}
