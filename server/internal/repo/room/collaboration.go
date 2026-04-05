package repository

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

var mentionRegex = regexp.MustCompile(`@([a-zA-Z0-9_]+)`)

type RoomMember struct {
	RoomID            uuid.UUID
	UserID            uuid.UUID
	Username          string
	Role              string
	CanManageRoom     bool
	CanManageChannels bool
	CanModerate       bool
	CanPost           bool
	MutedUntil        *time.Time
	BannedAt          *time.Time
	CreatedAt         time.Time
}

type RoomCategory struct {
	ID        uuid.UUID
	RoomID    uuid.UUID
	Name      string
	Position  int
	CreatedAt time.Time
}

type RoomChannel struct {
	ID          uuid.UUID
	RoomID      uuid.UUID
	CategoryID  *uuid.UUID
	Name        string
	Description string
	Kind        string
	Position    int
	IsPrivate   bool
	CreatedAt   time.Time
}

type Notification struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	RoomID    *uuid.UUID
	MessageID *uuid.UUID
	Kind      string
	Title     string
	Body      string
	Payload   []byte
	IsRead    bool
	CreatedAt time.Time
}

func (r *RoomRepository) EnsureRoomMembership(ctx context.Context, roomID, userID uuid.UUID) error {
	query := `
		INSERT INTO room_members (room_id, user_id, role, can_post)
		VALUES ($1, $2, 'member', TRUE)
		ON CONFLICT (room_id, user_id) DO NOTHING
	`
	_, err := r.db.ExecContext(ctx, query, roomID, userID)
	return err
}

func (r *RoomRepository) GetRoomMember(ctx context.Context, roomID, userID uuid.UUID) (*RoomMember, error) {
	query := `
		SELECT rm.room_id, rm.user_id, u.username, rm.role, rm.can_manage_room, rm.can_manage_channels,
			rm.can_moderate, rm.can_post, rm.muted_until, rm.banned_at, rm.created_at
		FROM room_members rm
		JOIN users u ON u.id = rm.user_id
		WHERE rm.room_id = $1 AND rm.user_id = $2
	`

	var member RoomMember
	err := r.db.QueryRowContext(ctx, query, roomID, userID).Scan(
		&member.RoomID,
		&member.UserID,
		&member.Username,
		&member.Role,
		&member.CanManageRoom,
		&member.CanManageChannels,
		&member.CanModerate,
		&member.CanPost,
		&member.MutedUntil,
		&member.BannedAt,
		&member.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &member, nil
}

func (r *RoomRepository) GetRoomMembers(ctx context.Context, roomID uuid.UUID) ([]RoomMember, error) {
	query := `
		SELECT rm.room_id, rm.user_id, u.username, rm.role, rm.can_manage_room, rm.can_manage_channels,
			rm.can_moderate, rm.can_post, rm.muted_until, rm.banned_at, rm.created_at
		FROM room_members rm
		JOIN users u ON u.id = rm.user_id
		WHERE rm.room_id = $1 AND rm.banned_at IS NULL
		ORDER BY rm.created_at ASC
	`
	rows, err := r.db.QueryContext(ctx, query, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []RoomMember
	for rows.Next() {
		var member RoomMember
		if err := rows.Scan(
			&member.RoomID,
			&member.UserID,
			&member.Username,
			&member.Role,
			&member.CanManageRoom,
			&member.CanManageChannels,
			&member.CanModerate,
			&member.CanPost,
			&member.MutedUntil,
			&member.BannedAt,
			&member.CreatedAt,
		); err != nil {
			return nil, err
		}
		members = append(members, member)
	}
	return members, rows.Err()
}

func (r *RoomRepository) UpdateRoomMember(ctx context.Context, member RoomMember) error {
	query := `
		UPDATE room_members
		SET role = $3,
			can_manage_room = $4,
			can_manage_channels = $5,
			can_moderate = $6,
			can_post = $7,
			banned_at = $8,
			updated_at = NOW()
		WHERE room_id = $1 AND user_id = $2
	`
	_, err := r.db.ExecContext(ctx, query,
		member.RoomID,
		member.UserID,
		member.Role,
		member.CanManageRoom,
		member.CanManageChannels,
		member.CanModerate,
		member.CanPost,
		member.BannedAt,
	)
	return err
}

func (r *RoomRepository) CreateCategory(ctx context.Context, category *RoomCategory) (*RoomCategory, error) {
	query := `
		INSERT INTO room_categories (room_id, name, position)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`
	err := r.db.QueryRowContext(ctx, query, category.RoomID, category.Name, category.Position).Scan(
		&category.ID,
		&category.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (r *RoomRepository) CreateChannel(ctx context.Context, channel *RoomChannel) (*RoomChannel, error) {
	query := `
		INSERT INTO room_channels (room_id, category_id, name, description, kind, position, is_private)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at
	`
	err := r.db.QueryRowContext(ctx, query,
		channel.RoomID,
		channel.CategoryID,
		channel.Name,
		channel.Description,
		channel.Kind,
		channel.Position,
		channel.IsPrivate,
	).Scan(&channel.ID, &channel.CreatedAt)
	if err != nil {
		return nil, err
	}
	return channel, nil
}

func (r *RoomRepository) GetRoomCategories(ctx context.Context, roomID uuid.UUID) ([]RoomCategory, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, room_id, name, position, created_at
		FROM room_categories
		WHERE room_id = $1
		ORDER BY position ASC, created_at ASC
	`, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []RoomCategory
	for rows.Next() {
		var category RoomCategory
		if err := rows.Scan(&category.ID, &category.RoomID, &category.Name, &category.Position, &category.CreatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, rows.Err()
}

func (r *RoomRepository) GetRoomChannels(ctx context.Context, roomID uuid.UUID) ([]RoomChannel, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, room_id, category_id, name, description, kind, position, is_private, created_at
		FROM room_channels
		WHERE room_id = $1
		ORDER BY position ASC, created_at ASC
	`, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var channels []RoomChannel
	for rows.Next() {
		var channel RoomChannel
		if err := rows.Scan(
			&channel.ID,
			&channel.RoomID,
			&channel.CategoryID,
			&channel.Name,
			&channel.Description,
			&channel.Kind,
			&channel.Position,
			&channel.IsPrivate,
			&channel.CreatedAt,
		); err != nil {
			return nil, err
		}
		channels = append(channels, channel)
	}
	return channels, rows.Err()
}

func (r *RoomRepository) GetDefaultChannel(ctx context.Context, roomID uuid.UUID) (*RoomChannel, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, room_id, category_id, name, description, kind, position, is_private, created_at
		FROM room_channels
		WHERE room_id = $1
		ORDER BY position ASC, created_at ASC
		LIMIT 1
	`, roomID)

	var channel RoomChannel
	err := row.Scan(
		&channel.ID,
		&channel.RoomID,
		&channel.CategoryID,
		&channel.Name,
		&channel.Description,
		&channel.Kind,
		&channel.Position,
		&channel.IsPrivate,
		&channel.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &channel, nil
}

func (r *RoomRepository) GetRoomMessagesByChannel(ctx context.Context, roomID uuid.UUID, channelID *uuid.UUID, limit int, offset int) ([]*Message, error) {
	query := `
		SELECT id, room_id, user_id, username, content, is_system, created_at, channel_id, parent_message_id, metadata
		FROM messages
		WHERE room_id = $1
	`
	args := []any{roomID}
	if channelID != nil {
		query += ` AND channel_id = $2`
		args = append(args, *channelID)
	}
	query += ` ORDER BY created_at DESC LIMIT $` + fmt.Sprintf("%d", len(args)+1) + ` OFFSET $` + fmt.Sprintf("%d", len(args)+2)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(
			&msg.ID,
			&msg.RoomID,
			&msg.UserID,
			&msg.Username,
			&msg.Content,
			&msg.IsSystem,
			&msg.CreatedAt,
			&msg.ChannelID,
			&msg.ParentMessageID,
			&msg.Metadata,
		); err != nil {
			return nil, err
		}
		messages = append(messages, &msg)
	}
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}
	return messages, rows.Err()
}

func (r *RoomRepository) SearchMessages(ctx context.Context, roomID uuid.UUID, queryText string, channelID *uuid.UUID, username string, limit int) ([]Message, error) {
	base := `
		SELECT m.id, m.room_id, m.user_id, m.username, m.content, m.is_system, m.created_at, m.channel_id, m.parent_message_id, m.metadata
		FROM messages m
		WHERE m.room_id = $1 AND m.content ILIKE $2
	`
	args := []any{roomID, "%" + queryText + "%"}
	param := 3
	if channelID != nil {
		base += fmt.Sprintf(" AND m.channel_id = $%d", param)
		args = append(args, *channelID)
		param++
	}
	if username != "" {
		base += fmt.Sprintf(" AND m.username = $%d", param)
		args = append(args, username)
		param++
	}
	base += fmt.Sprintf(" ORDER BY m.created_at DESC LIMIT $%d", param)
	args = append(args, limit)

	rows, err := r.db.QueryContext(ctx, base, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(
			&msg.ID,
			&msg.RoomID,
			&msg.UserID,
			&msg.Username,
			&msg.Content,
			&msg.IsSystem,
			&msg.CreatedAt,
			&msg.ChannelID,
			&msg.ParentMessageID,
			&msg.Metadata,
		); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, rows.Err()
}

func (r *RoomRepository) CreateNotification(ctx context.Context, notification *Notification) error {
	query := `
		INSERT INTO notifications (user_id, room_id, message_id, kind, title, body, payload, is_read)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at
	`
	err := r.db.QueryRowContext(ctx, query,
		notification.UserID,
		notification.RoomID,
		notification.MessageID,
		notification.Kind,
		notification.Title,
		notification.Body,
		notification.Payload,
		notification.IsRead,
	).Scan(&notification.ID, &notification.CreatedAt)
	return err
}

func (r *RoomRepository) GetNotifications(ctx context.Context, userID uuid.UUID, limit int) ([]Notification, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, room_id, message_id, kind, title, body, payload, is_read, created_at
		FROM notifications
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []Notification
	for rows.Next() {
		var notification Notification
		if err := rows.Scan(
			&notification.ID,
			&notification.UserID,
			&notification.RoomID,
			&notification.MessageID,
			&notification.Kind,
			&notification.Title,
			&notification.Body,
			&notification.Payload,
			&notification.IsRead,
			&notification.CreatedAt,
		); err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}
	return notifications, rows.Err()
}

func (r *RoomRepository) MarkNotificationRead(ctx context.Context, notificationID, userID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE notifications
		SET is_read = TRUE
		WHERE id = $1 AND user_id = $2
	`, notificationID, userID)
	return err
}

func (r *RoomRepository) CreateMentionNotifications(ctx context.Context, roomID uuid.UUID, message *Message) ([]Notification, error) {
	if message.UserID == nil {
		return nil, nil
	}

	matches := mentionRegex.FindAllStringSubmatch(message.Content, -1)
	if len(matches) == 0 {
		return nil, nil
	}

	seen := map[string]struct{}{}
	usernames := make([]string, 0, len(matches))
	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		name := strings.ToLower(match[1])
		if _, ok := seen[name]; ok {
			continue
		}
		seen[name] = struct{}{}
		usernames = append(usernames, name)
	}
	if len(usernames) == 0 {
		return nil, nil
	}

	query := `
		SELECT u.id, u.username
		FROM users u
		JOIN room_members rm ON rm.user_id = u.id
		WHERE rm.room_id = $1
			AND LOWER(u.username) = ANY($2)
			AND rm.banned_at IS NULL
	`
	rows, err := r.db.QueryContext(ctx, query, roomID, pq.Array(usernames))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []Notification
	for rows.Next() {
		var mentionedUserID uuid.UUID
		var username string
		if err := rows.Scan(&mentionedUserID, &username); err != nil {
			return nil, err
		}
		if mentionedUserID == *message.UserID {
			continue
		}

		notification := Notification{
			UserID:    mentionedUserID,
			RoomID:    &roomID,
			MessageID: &message.ID,
			Kind:      "mention",
			Title:     "You were mentioned",
			Body:      fmt.Sprintf("%s mentioned you in %s", message.Username, message.Content),
			Payload:   []byte(fmt.Sprintf(`{"room_id":"%s","message_id":"%s","username":"%s"}`, roomID.String(), message.ID.String(), message.Username)),
		}
		if err := r.CreateNotification(ctx, &notification); err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}

	return notifications, nil
}
