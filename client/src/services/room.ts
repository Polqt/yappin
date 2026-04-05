import type {
	CreateCategoryRequest,
	CreateChannelRequest,
	CreateRoomRequest,
	MessageSearchResult,
	NotificationItem,
	Room,
	RoomDetail,
	UpdateMemberRoleRequest
} from '$lib/types/room';
import { API_BASE_URL, API_ENDPOINTS } from '$lib/constants/api';
import { handleApiError } from '$lib/utils/error';

async function readJson<T>(response: Response): Promise<T> {
	if (!response.ok) {
		await handleApiError(response);
	}
	return response.json() as Promise<T>;
}

export const roomService = {
	async getRooms(): Promise<Room[]> {
		const response = await fetch(`${API_BASE_URL}${API_ENDPOINTS.rooms.list}`, {
			credentials: 'include'
		});
		return readJson<Room[]>(response);
	},

	async getRoomById(roomId: string): Promise<Room | null> {
		const rooms = await this.getRooms();
		return rooms.find((room) => room.id === roomId) || null;
	},

	async getRoomDetail(roomId: string): Promise<RoomDetail> {
		const response = await fetch(`${API_BASE_URL}${API_ENDPOINTS.rooms.detail(roomId)}`, {
			credentials: 'include'
		});
		return readJson<RoomDetail>(response);
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

		return readJson<Room>(response);
	},

	async createCategory(roomId: string, request: CreateCategoryRequest) {
		const response = await fetch(`${API_BASE_URL}${API_ENDPOINTS.rooms.categories(roomId)}`, {
			method: 'POST',
			credentials: 'include',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(request)
		});

		return readJson(response);
	},

	async createChannel(roomId: string, request: CreateChannelRequest) {
		const response = await fetch(`${API_BASE_URL}${API_ENDPOINTS.rooms.channels(roomId)}`, {
			method: 'POST',
			credentials: 'include',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(request)
		});

		return readJson(response);
	},

	async updateMemberRole(roomId: string, userId: string, request: UpdateMemberRoleRequest) {
		const response = await fetch(`${API_BASE_URL}${API_ENDPOINTS.rooms.members(roomId, userId)}`, {
			method: 'PUT',
			credentials: 'include',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(request)
		});

		return readJson(response);
	},

	async searchMessages(roomId: string, query: string, channelId?: string): Promise<MessageSearchResult[]> {
		const url = new URL(`${API_BASE_URL}${API_ENDPOINTS.rooms.search(roomId)}`);
		url.searchParams.set('query', query);
		if (channelId) {
			url.searchParams.set('channel_id', channelId);
		}

		const response = await fetch(url.toString(), {
			credentials: 'include'
		});
		return readJson<MessageSearchResult[]>(response);
	},

	async getNotifications(): Promise<NotificationItem[]> {
		const response = await fetch(`${API_BASE_URL}${API_ENDPOINTS.rooms.notifications}`, {
			credentials: 'include'
		});
		return readJson<NotificationItem[]>(response);
	},

	async markNotificationRead(notificationId: string) {
		const response = await fetch(
			`${API_BASE_URL}${API_ENDPOINTS.rooms.markNotificationRead(notificationId)}`,
			{
				method: 'PUT',
				credentials: 'include'
			}
		);
		return readJson(response);
	},

	async addReaction(messageId: string, emoji: string) {
		const response = await fetch(`${API_BASE_URL}${API_ENDPOINTS.rooms.reactions}`, {
			method: 'POST',
			credentials: 'include',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({ message_id: messageId, emoji })
		});
		return readJson(response);
	}
};
