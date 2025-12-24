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

export interface Message {
	content: string;
	room_id: string;
	username: string;
	user_id?: string;
	system: boolean;
	created_at?: string;
	reactions?: MessageReaction[];
}

export interface MessageReaction {
	id: string;
	message_id: string;
	user_id: string;
	emoji: string;
	created_at: string;
}

export type MessageType = 'TEXT' | 'IMAGE' | 'SYSTEM';

export interface CreateRoomRequest {
	name: string;
	expires_at?: string | null;
}
