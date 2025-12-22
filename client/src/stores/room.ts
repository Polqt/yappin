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
    content: string; // The actual message content
    room_id: string; // Which room the message belongs to
    username: string; // Who sent the message
    user_id?: string; // No guests allowed
    system: boolean; // True if the message is a system message like "Jepoy joined the room"
    created_at: string;
}

export type MessageType = 'TEXT' | 'IMAGE' | 'SYSTEM';

export interface CreateRoomRequest {
    name: string;
    expires_at?: string | null;
}
