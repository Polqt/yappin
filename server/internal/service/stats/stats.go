package service

import (
	"context"
	"fmt"
	"log"

	stats "chat-application/internal/repo/stats"

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
			newAchievements = s.convertAchievement(achievements)
		}
	}

	return &CheckinResult{
		StreakCount: streakCount,
		IsNewCheckin: IsNewCheckin,
		NewAchievements: newAchievements,
	}, nil
}

func (s *StatsService) convertAchievement(repositoryAchievements []stats.Achievement) []Achievement {
	var achievements []Achievement
	for _, ach := range repositoryAchievements {
		earnedAt := ""
		if ach.EarnedAt != nil {
			earnedAt = ach.EarnedAt.Format("2006-01-02T15:04:05Z")
		}

		achievements = append(achievements, Achievement{
			ID: ach.ID.String(),
			Name: ach.Name,
			Description: ach.Description,
			Icon: ach.Icon,
			ThresholdType: ach.ThresholdType,
			ThresholdValue: ach.ThresholdValue,
			EarnedAt: earnedAt,
		})
	}
	return achievements
}

func (s *StatsService) GetUserProfile(ctx context.Context, userID, viewerID uuid.UUID) (*UserProfile, error) {
	stats, err := s.statsRepository.GetUserProfile(ctx, userID)
	if err != nil {
		log.Printf("GetUserProfile - Error retrieving profile for user %s: %v", userID, err)
		return nil, fmt.Errorf("failed to retrieve user profile: %w", err)
	}

	achievements, err := s.statsRepository.GetUserAchievementsDetails(ctx, userID)
	if err != nil {
		log.Printf("GetUserProfile - Error retrieving achievements for user %s: %v", userID, err)
		return nil, fmt.Errorf("failed to retrieve user achievements: %w", err)
	}

	canUpvote := false
	if viewerID != userID {
		canUpvote, err = s.statsRepository.CanUserUpvote(ctx, viewerID, userID)
		if err != nil {
			log.Printf("GetUserProfile - Error checking upvote permission for viewer %s on user %s: %v", viewerID, userID, err)
			canUpvote = false
		}
	}

	return &UserProfile{
		UserID: userID.String(),
		DailyStreak: stats.DailyStreak,
		TotalCheckins: stats.TotalCheckins,
		TotalMessages: stats.TotalMessages,
		TotalUpvotes: stats.TotalUpvotes,
		CanReceiveUpvotes: canUpvote,
		Achievements: s.convertAchievement(achievements),
	}, nil
}

func (s *StatsService) GivenUpvote(ctx context.Context, fromUserID, toUserID uuid.UUID) error {
	if fromUserID == toUserID {
		return fmt.Errorf("user cannot upvote themselves")
	}

	canUpvote, err := s.statsRepository.CanUserUpvote(ctx, fromUserID, toUserID)
	if err != nil {
		log.Printf("GivenUpvote - Error checking upvote permission for user %s on user %s: %v", fromUserID, toUserID, err)
		return fmt.Errorf("failed to check upvote permission: %w", err)
	}

	if !canUpvote {
		log.Printf("GivenUpvote - User %s cannot upvote user %s", fromUserID, toUserID)
		return fmt.Errorf("user %s cannot upvote user %s", fromUserID, toUserID)
	}

	log.Printf("GivenUpvote - User %s is giving an upvote to user %s", fromUserID, toUserID)

	err = s.statsRepository.GiveUpvote(ctx, fromUserID, toUserID)
	if err != nil {
		log.Printf("GivenUpvote - Error giving upvote from user %s to user %s: %v", fromUserID, toUserID, err)
		return fmt.Errorf("failed to give upvote: %w", err)
	}

	go func() {
		_, err := s.statsRepository.CheckAwardsAndAchievements(context.Background(), toUserID)
		if err != nil {
			log.Printf("GivenUpvote - Error checking achievements after upvote for user %s: %v", toUserID, err)
		}
	}()

	log.Printf("GivenUpvote - Upvote successfully given from user %s to user %s", fromUserID, toUserID)
	return nil
}