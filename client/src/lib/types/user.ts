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

export interface Achievement {
	id: string;
	name: string;
	description: string;
	icon: string;
	achieved_at: string;
}

export interface UserProfile {
	user_id: string;
	username: string;
	total_messages: number;
	total_upvotes: number;
	daily_streak: number;
	level: number;
	achievements?: Achievement[];
	activity?: ActivityData[];
}

export interface ActivityData {
	date: string;
	count: number;
}
