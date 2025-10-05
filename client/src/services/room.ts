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
  }
};
