import type { CreateRoomRequest, Room } from '$lib/types/room';
import { API_BASE_URL, API_ENDPOINTS } from '$lib/constants/api';
import { handleApiError } from '$lib/utils/error';

export const roomService = {
	async getRooms(): Promise<Room[]> {
		const response = await fetch(`${API_BASE_URL}${API_ENDPOINTS.rooms.list}`, {
			credentials: 'include'
		});

		if (!response.ok) {
			await handleApiError(response);
		}
		return response.json();
	},

	async createRoom(request: CreateRoomRequest): Promise<Room> {
		const body: { name: string; expires_at?: string } = {
			name: request.name
		};

		if (request.expires_at) {
			body.expires_at = new Date(request.expires_at).toISOString();
		}

		const response = await fetch(`${API_BASE_URL}${API_ENDPOINTS.rooms.create}`, {
			method: 'POST',
			credentials: 'include',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(body)
		});

		if (!response.ok) {
			await handleApiError(response);
		}
		return response.json();
	}
};
