import type { CreateRoomRequest, Room } from '$lib/types/room';

const BASE_URL = 'http://localhost:8081';

export const roomService = {
  async getRooms(): Promise<Room[]> {
    const response = await fetch(`${BASE_URL}/api/websoc/get-rooms`, {
      credentials: 'include'
    });

    if (!response.ok) {
      throw new Error('Failed to fetch rooms');
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

		const response = await fetch(`${BASE_URL}/api/websoc/create-room`, {
			method: 'POST',
			credentials: 'include',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(body)
		});

		if (!response.ok) {
			const errorData = await response.json();
			throw new Error(errorData.error || 'Failed to create room');
		}
		return response.json();
	}
};
