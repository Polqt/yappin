export interface User {
  id: string;
  email: string;
  username: string;
  avatar?: string;
  createdAt: string;
  updatedAt: string;
}

export interface UserStats {
  totalMessages: number;
  roomsJoined: number;
  upvotesReceived: number;
}