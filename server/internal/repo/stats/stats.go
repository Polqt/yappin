package stats

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
)

type UserStats struct {
	UserID uuid.UUID `db:"user_id"`
	DailyStreak int `db:"daily_streak"`
	TotalCheckins int `db:"total_checkins"`
	TotalMessages int `db:"total_messages"`
	TotalUpvotes int `db:"total_upvotes"`
	LastCheckinDate *time.Time `db:"last_checkin_date"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type DailyCheckin struct {
	ID uuid.UUID `db:"id"`
	UserID uuid.UUID `db:"user_id"`
	CheckinDate time.Time `db:"checkin_date"`
	StreakCount int `db:"streak_count"`
	CreatedAt time.Time `db:"created_at"`
}

type Upvote struct {
	ID uuid.UUID `db:"id"`
	FromUserID uuid.UUID `db:"from_user_id"`
	ToUserID uuid.UUID `db:"to_user_id"`
	CreateAt time.Time `db:"created_at"`
}

type Achievement struct {
	ID uuid.UUID `db:"id"`
	Name string `db:"name"`
	Description string `db:"description"`
	Icon string `db:"icon"`
	ThresholdType string `db:"threshold_type"`
	ThresholdValue int `db:"threshold_value"`

}

type StatsRepository struct {
	db *sql.DB
} 

func NewStatsRepository(db *sql.DB) *StatsRepository {
	return &StatsRepository{
		db: db,
	}
}

func (r *StatsRepository) GetOrCreateUserStats(ctx context.Context, userID uuid.UUID) (*UserStats, error) {
	stats := &UserStats{}

	query := `
		SELECT user_id, daily_streak, total_checkins, total_messages,
			total_upvotes, last_checkin_date, created_at, updated_at
		FROM user_stats
		WHERE user_id = $1
	`

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&stats.UserID, &stats.DailyStreak, &stats.TotalCheckins,
		&stats.TotalMessages, &stats.TotalUpvotes, &stats.LastCheckinDate,
		&stats.CreatedAt, &stats.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		insertQuery := `
			INSERT INTO user_stats (user_id)
			VALUES ($1)
			RETURNING user_id, daily_streak, total_checkins, total_messages,
				total_upvotes, last_checkin_date, created_at, updated_at
		`

		err = r.db.QueryRowContext(ctx, insertQuery, userID).Scan(
			&stats.UserID, &stats.DailyStreak, &stats.TotalCheckins,
			&stats.TotalMessages, &stats.TotalUpvotes, &stats.LastCheckinDate,
			&stats.CreatedAt, &stats.UpdatedAt,
		)
	}
	return stats, err
}

func (r *StatsRepository) GetUserProfile(ctx context.Context, userID uuid.UUID) (*UserStats, error) {
	return r.GetOrCreateUserStats(ctx, userID)
}

func (r *StatsRepository) ProcessDailyCheckin(ctx context.Context, userID uuid.UUID) (int, bool, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, false, err
	}
	defer tx.Rollback()

	today := time.Now().UTC().Truncate(24 * time.Hour)

	stats, err := r.GetOrCreateUserStats(ctx, userID)
	if err != nil {
		return 0, false, err
	}

	if stats.LastCheckinDate != nil {
		lastCheckIn := stats.LastCheckinDate.UTC().Truncate(24 * time.Hour)
		if lastCheckIn.Equal(today) {
			return stats.DailyStreak, false, nil // Already checked in today
		}
	}

	newStreak := 1
	if stats.LastCheckinDate != nil {
		yesterday := today.Add(-24 * time.Hour)
		lastCheckin := stats.LastCheckinDate.UTC().Truncate(24 * time.Hour)
		if lastCheckin.Equal(yesterday) {
			newStreak = stats.DailyStreak + 1
		}
	}

	updateQuery := `
		INSERT INTO user_stats (user_id, checkin_date, streak_count)
		VALUES ($1, $2, $3)
	`
	_, err = tx.ExecContext(ctx, updateQuery, userID, today, newStreak)
	if err != nil {
		return 0, false, err
	}

	err = tx.Commit()
	return newStreak, true, err
}

func (r *StatsRepository) GetUserAchievements(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	query := `
		SELECT achievement_type_id
		FROM user_achievements
		WHERE user_id = $1
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var achievementID []uuid.UUID
	
	for rows.Next() {
		var achID uuid.UUID
		if err := rows.Scan(&achID); err != nil {
			return nil, err
		}
		achievementID = append(achievementID, achID)
	}
	return achievementID, rows.Err()
}

func (r *StatsRepository) GetAllAchievementTypes(ctx context.Context) ([]Achievement, error) {
	query := `
		SELECT id, name, description, icon, threshold_type, threshold_value
		FROM achievement_types
		ORDER BY threshold_value DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		log.Printf("GetAllAchievementTypes - Error querying achievement types: %v", err)
		return nil, err
	}
	defer rows.Close()

	var achievements []Achievement
	for rows.Next() {
		var ach Achievement
		err := rows.Scan(&ach.ID, &ach.Name, &ach.Description, &ach.Icon, &ach.ThresholdType, &ach.ThresholdValue)
		
		if err != nil {
			log.Printf("GetAllAchievementTypes - Error scanning achievement type: %v", err)
			return nil, err
		}
		achievements = append(achievements, ach)
	}

	return achievements, rows.Err()
}

func (r *StatsRepository) CheckAwardsAndAchievements(ctx context.Context, userID uuid.UUID) ([]Achievement, error) {
	stats, err := r.GetOrCreateUserStats(ctx, userID)
	if err != nil {
		log.Printf("CheckAwardsAndAchievements - Error getting user stats for user %s: %v", userID, err)
		return nil, err
	}

	achievementTypes, err := r.GetAllAchievementTypes(ctx)
	if err != nil {
		log.Printf("CheckAwardsAndAchievements - Error getting achievement types for user %s: %v", userID, err)
		return nil, err
	}

	earnedAchievements, err := r.GetUserAchievements(ctx, userID)
	if err != nil {
		log.Printf("CheckAwardsAndAchievements - Error getting user achievements for user %s: %v", userID, err)
		return nil, err
	}

	newAchievements := []Achievement{}

	for _, achType := range achievementTypes {
		if r.hasAchievement(earnedAchievements, achType.ID) {
			continue // User already has this achievement
		}

		var currentValue int
		switch achType.ThresholdType {
		case "daily_streak":
			currentValue = stats.DailyStreak
		case "messages":
			currentValue = stats.TotalMessages
		case "upvotes":
			currentValue = stats.TotalUpvotes
		default:
			continue
		}
	
		if currentValue >= achType.ThresholdValue {
			if err := r.awardAchievement(ctx, userID, achType.ID); err != nil {
				log.Printf("CheckAwardsAndAchievements - Error awarding achievement %s to user %s: %v", achType.ID, userID, err)
				continue
			}

			achievement := Achievement {
				ID: achType.ID,
				Name: achType.Name,
				Description: achType.Description,
				Icon: achType.Icon,
				ThresholdType: achType.ThresholdType,
				ThresholdValue: achType.ThresholdValue,
			}

			newAchievements = append(newAchievements, achievement)
		}
	}

	return newAchievements, nil
}

func (r *StatsRepository) hasAchievement(earnedAchievements []uuid.UUID, achievementID uuid.UUID) bool {
	for _, achID := range earnedAchievements {
		if achID == achievementID {
			return true
		}
	}
	return false
}

func (r *StatsRepository) awardAchievement(ctx context.Context, userID, achievementID uuid.UUID) error {
	query := `
		INSERT INTO user_achievements (user_id, achievement_type_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, achievement_type_id) DO NOTHING
	`

	_, err := r.db.ExecContext(ctx, query, userID, achievementID)
	return err
}