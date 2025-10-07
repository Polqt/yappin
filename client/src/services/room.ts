import type { Room } from '$lib/types/room';

const BASE_URL = 'http://localhost:8080';

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

	async createRoom(name: string): Promise<Room> {
		const response = await fetch(`${BASE_URL}/api/websoc/create-room`, {
			method: 'POST',
			credentials: 'include',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({ name })
		});

		if (!response.ok) {
			const errorData = await response.json();
			console.log(errorData);
			throw new Error(errorData.error || 'Failed to create room');
		}
		return response.json();
	}
};
