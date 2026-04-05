import axios from 'axios';
import { API_BASE_URL, API_ENDPOINTS } from '$lib/constants/api';
import type { LeaderboardEntry } from '$lib/types/leaderboard';
import type { UserProfile } from '$lib/types/user';
import { authService } from '$services/auth';

export interface CheckinResult {
	streak_count: number;
	is_new_checkin: boolean;
	new_achievements?: Array<{
		id: string;
		name: string;
		description: string;
		icon: string;
		threshold_type: string;
		threshold_value: number;
		earned_at: string;
	}>;
}

export async function getUserProfile(userId: string): Promise<UserProfile> {
	const response = await axios.get(`${API_BASE_URL}${API_ENDPOINTS.stats.profile(userId)}`, {
		withCredentials: true
	});
	return response.data;
}

export async function getLeaderboard(limit: number = 10): Promise<LeaderboardEntry[]> {
	const response = await axios.get(`${API_BASE_URL}${API_ENDPOINTS.stats.leaderboard}`, {
		params: { limit },
		withCredentials: true
	});
	return response.data;
}

export async function recordCheckin(): Promise<CheckinResult> {
	const response = await axios.post(
		`${API_BASE_URL}${API_ENDPOINTS.stats.checkin}`,
		{},
		{
			withCredentials: true
		}
	);
	return response.data;
}

export async function getCurrentUserProfile(): Promise<UserProfile> {
	const currentUser = await authService.getCurrentUser();
	if (!currentUser?.id) {
		throw new Error('User is not authenticated');
	}

	return getUserProfile(currentUser.id);
}
