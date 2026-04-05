package model

import "time"

type RoomPermissionRes struct {
	Role              string `json:"role"`
	CanManageRoom     bool   `json:"can_manage_room"`
	CanManageChannels bool   `json:"can_manage_channels"`
	CanModerate       bool   `json:"can_moderate"`
	CanPost           bool   `json:"can_post"`
	IsMuted           bool   `json:"is_muted"`
	IsBanned          bool   `json:"is_banned"`
}

type RoomMemberRes struct {
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

type RoomCategoryRes struct {
	ID       string           `json:"id"`
	Name     string           `json:"name"`
	Position int              `json:"position"`
	Channels []RoomChannelRes `json:"channels"`
}

type RoomChannelRes struct {
	ID          string `json:"id"`
	CategoryID  string `json:"category_id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Kind        string `json:"kind"`
	Position    int    `json:"position"`
	IsPrivate   bool   `json:"is_private"`
}

type NotificationRes struct {
	ID        string         `json:"id"`
	Kind      string         `json:"kind"`
	Title     string         `json:"title"`
	Body      string         `json:"body"`
	IsRead    bool           `json:"is_read"`
	CreatedAt time.Time      `json:"created_at"`
	Payload   map[string]any `json:"payload,omitempty"`
	MessageID string         `json:"message_id,omitempty"`
	RoomID    string         `json:"room_id,omitempty"`
}

type MessageSearchRes struct {
	ID              string    `json:"id"`
	RoomID          string    `json:"room_id"`
	ChannelID       string    `json:"channel_id"`
	ParentMessageID string    `json:"parent_message_id,omitempty"`
	Username        string    `json:"username"`
	Content         string    `json:"content"`
	Highlighted     string    `json:"highlighted"`
	CreatedAt       time.Time `json:"created_at"`
}

type MessageRes struct {
	ID              string            `json:"id"`
	Content         string            `json:"content"`
	RoomID          string            `json:"room_id"`
	ChannelID       string            `json:"channel_id"`
	ParentMessageID string            `json:"parent_message_id,omitempty"`
	Username        string            `json:"username"`
	UserID          string            `json:"user_id,omitempty"`
	System          bool              `json:"system"`
	CreatedAt       time.Time         `json:"created_at"`
	Metadata        map[string]any    `json:"metadata,omitempty"`
	Reactions       []MessageReaction `json:"reactions,omitempty"`
}

type RoomDetailRes struct {
	Room               RoomRes            `json:"room"`
	Categories         []RoomCategoryRes  `json:"categories"`
	Members            []RoomMemberRes    `json:"members"`
	CurrentUser        *RoomPermissionRes `json:"current_user,omitempty"`
	Messages           []MessageRes       `json:"messages"`
	Notifications      []NotificationRes  `json:"notifications"`
	DefaultChannelID   string             `json:"default_channel_id,omitempty"`
	NotificationCount  int                `json:"notification_count"`
	OnlineMemberCount  int                `json:"online_member_count"`
	ThreadedReplyCount int                `json:"threaded_reply_count"`
}

type CreateCategoryReq struct {
	Name     string `json:"name"`
	Position int    `json:"position"`
}

type CreateChannelReq struct {
	CategoryID  string `json:"category_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Kind        string `json:"kind"`
	Position    int    `json:"position"`
	IsPrivate   bool   `json:"is_private"`
}

type UpdateMemberRoleReq struct {
	Role              string `json:"role"`
	CanManageRoom     *bool  `json:"can_manage_room,omitempty"`
	CanManageChannels *bool  `json:"can_manage_channels,omitempty"`
	CanModerate       *bool  `json:"can_moderate,omitempty"`
	CanPost           *bool  `json:"can_post,omitempty"`
	Ban               bool   `json:"ban"`
}
