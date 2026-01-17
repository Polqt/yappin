import axios from 'axios';
import { API_BASE_URL, API_ENDPOINTS } from '$lib/constants/api';
import type { LeaderboardEntry } from '$lib/types/leaderboard';

export interface UserProfile {
	id: string;
	username: string;
	email: string;
	stats: UserStats;
	achievements: Achievement[];
	dailyCheckins: DailyCheckin[];
}

export interface UserStats {
	messages_sent: number;
	rooms_joined: number;
	rooms_created: number;
	upvotes_received: number;
	current_streak: number;
	longest_streak: number;
}

export interface Achievement {
	id: string;
	type: string;
	name: string;
	description: string;
	earned_at: string;
}

export interface DailyCheckin {
	id: string;
	checkin_date: string;
	created_at: string;
}

export async function getUserProfile(): Promise<UserProfile> {
	const response = await axios.get(`${API_BASE_URL}${API_ENDPOINTS.stats.profile}`, {
		withCredentials: true
	});
	return response.data;
}

/**
 * Fetches the leaderboard data.
 * @param period - Time period for leaderboard (all-time, weekly, monthly)
 * @param limit - Maximum number of entries to return
 */
export async function getLeaderboard(
	period: 'all-time' | 'weekly' | 'monthly' = 'all-time',
	limit: number = 10
): Promise<LeaderboardEntry[]> {
	const response = await axios.get(`${API_BASE_URL}${API_ENDPOINTS.stats.leaderboard}`, {
		params: { period, limit },
		withCredentials: true
	});
	return response.data;
}

export async function recordCheckin(): Promise<{ success: boolean; streak: number }> {
	const response = await axios.post(
		`${API_BASE_URL}${API_ENDPOINTS.stats.checkin}`,
		{},
		{
			withCredentials: true
		}
	);
	return response.data;
}

export async function getUserStats(userId: string): Promise<UserStats> {
	const response = await axios.get(`${API_BASE_URL}${API_ENDPOINTS.stats.userStats(userId)}`, {
		withCredentials: true
	});
	return response.data;
}
