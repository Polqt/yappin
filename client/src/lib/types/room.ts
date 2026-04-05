export interface Room {
	id: string;
	name: string;
	is_pinned: boolean;
	created_at: string;
	expires_at: string;
	topic_title?: string;
	topic_description?: string;
	topic_url?: string;
	topic_source?: string;
	creator_username?: string;
	participants: number;
}

export interface RoomPermission {
	role: string;
	can_manage_room: boolean;
	can_manage_channels: boolean;
	can_moderate: boolean;
	can_post: boolean;
	is_muted: boolean;
	is_banned: boolean;
}

export interface RoomMember {
	user_id: string;
	username: string;
	role: string;
	created_at: string;
}

export interface RoomChannel {
	id: string;
	category_id?: string;
	name: string;
	description?: string;
	kind: string;
	position: number;
	is_private: boolean;
}

export interface RoomCategory {
	id: string;
	name: string;
	position: number;
	channels: RoomChannel[];
}

export interface MessageReaction {
	id: string;
	message_id: string;
	user_id: string;
	emoji: string;
	created_at: string;
}

export interface Message {
	id: string;
	content: string;
	room_id: string;
	channel_id: string;
	parent_message_id?: string;
	username: string;
	user_id?: string;
	system: boolean;
	created_at?: string;
	metadata?: Record<string, unknown>;
	reactions?: MessageReaction[];
}

export interface NotificationItem {
	id: string;
	kind: string;
	title: string;
	body: string;
	is_read: boolean;
	created_at: string;
	payload?: Record<string, unknown>;
	message_id?: string;
	room_id?: string;
}

export interface MessageSearchResult {
	id: string;
	room_id: string;
	channel_id: string;
	parent_message_id?: string;
	username: string;
	content: string;
	highlighted: string;
	created_at: string;
}

export interface RoomDetail {
	room: Room;
	categories: RoomCategory[];
	members: RoomMember[];
	current_user?: RoomPermission;
	messages: Message[];
	notifications: NotificationItem[];
	default_channel_id?: string;
	notification_count: number;
	online_member_count: number;
	threaded_reply_count: number;
}

export interface CreateRoomRequest {
	name: string;
	expires_at?: string | null;
}

export interface CreateCategoryRequest {
	name: string;
	position?: number;
}

export interface CreateChannelRequest {
	category_id?: string;
	name: string;
	description?: string;
	kind?: string;
	position?: number;
	is_private?: boolean;
}

export interface UpdateMemberRoleRequest {
	role?: string;
	can_manage_room?: boolean;
	can_manage_channels?: boolean;
	can_moderate?: boolean;
	can_post?: boolean;
	ban?: boolean;
}

export interface PresenceUser {
	user_id?: string;
	username: string;
}

export interface TypingEvent {
	room_id: string;
	channel_id?: string;
	user_id?: string;
	username: string;
	is_typing: boolean;
}

export interface WebSocketEvent {
	type: 'history' | 'message.created' | 'typing' | 'presence' | 'notification';
	message?: Message;
	messages?: Message[];
	typing?: TypingEvent;
	presence?: {
		room_id: string;
		online_users: PresenceUser[];
	};
	notification?: NotificationItem;
}
