export interface Room {
	id: string;
	name: string;
	description?: string;
	createdBy: string;
	participants: number;
	createdAt: string;
	updatedAt: string;
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
