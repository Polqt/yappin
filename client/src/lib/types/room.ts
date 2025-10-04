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
  id: string;
  content: string;
  userId: string;
  roomId: string;
  type: MessageType;
  createdAt: string;
}

export type MessageType = 'TEXT' | 'IMAGE' | 'SYSTEM';