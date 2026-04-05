import type { Message, NotificationItem, PresenceUser, WebSocketEvent } from '$lib/types/room';
import { writable } from 'svelte/store';
import { WS_BASE_URL, API_ENDPOINTS } from '$lib/constants/api';

type ConnectionState = 'disconnected' | 'connecting' | 'connected' | 'reconnecting';

interface WebSocketState {
	connectionState: ConnectionState;
	connected: boolean;
	messages: Message[];
	onlineUsers: PresenceUser[];
	typingUsers: string[];
	notifications: NotificationItem[];
	error: string | null;
}

const MAX_RECONNECT_ATTEMPTS = 5;
const RECONNECT_DELAY_MS = 2000;

function createWebSocketStore() {
	const { subscribe, set, update } = writable<WebSocketState>({
		connectionState: 'disconnected',
		connected: false,
		messages: [],
		onlineUsers: [],
		typingUsers: [],
		notifications: [],
		error: null
	});

	let socket: WebSocket | null = null;
	let reconnectAttempts = 0;
	let reconnectTimeout: ReturnType<typeof setTimeout> | null = null;
	let currentRoomId: string | null = null;
	let currentUsername: string | null = null;
	let currentUserId: string | undefined;

	const clearReconnectTimeout = () => {
		if (reconnectTimeout) {
			clearTimeout(reconnectTimeout);
			reconnectTimeout = null;
		}
	};

	const connect = (roomId: string, username: string, userId?: string) => {
		currentRoomId = roomId;
		currentUsername = username;
		currentUserId = userId;
		reconnectAttempts = 0;
		clearReconnectTimeout();

		update((state) => ({
			...state,
			messages: [],
			typingUsers: [],
			error: null
		}));

		performConnect(roomId, username, userId);
	};

	const performConnect = (roomId: string, username: string, userId?: string) => {
		if (socket && socket.readyState === WebSocket.OPEN) {
			socket.close();
		}

		let wsUrl = `${WS_BASE_URL}${API_ENDPOINTS.rooms.join(roomId)}?username=${encodeURIComponent(username)}`;
		if (userId) {
			wsUrl += `&user_id=${encodeURIComponent(userId)}`;
		}

		update((state) => ({
			...state,
			connectionState: reconnectAttempts > 0 ? 'reconnecting' : 'connecting',
			error: null
		}));

		socket = new WebSocket(wsUrl);

		socket.onopen = () => {
			reconnectAttempts = 0;
			update((state) => ({
				...state,
				connectionState: 'connected',
				connected: true,
				error: null
			}));
		};

		socket.onmessage = (event) => {
			try {
				const payload = JSON.parse(event.data) as WebSocketEvent;
				update((state) => applyEvent(state, payload));
			} catch {
				update((state) => ({
					...state,
					error: 'Failed to parse realtime event'
				}));
			}
		};

		socket.onclose = (event) => {
			update((state) => ({
				...state,
				connectionState: 'disconnected',
				connected: false
			}));

			if (
				!event.wasClean &&
				currentRoomId &&
				currentUsername &&
				reconnectAttempts < MAX_RECONNECT_ATTEMPTS
			) {
				reconnectAttempts++;
				const delay = RECONNECT_DELAY_MS * reconnectAttempts;
				reconnectTimeout = setTimeout(() => {
					if (currentRoomId && currentUsername) {
						performConnect(currentRoomId, currentUsername, currentUserId);
					}
				}, delay);
			}
		};

		socket.onerror = () => {
			update((state) => ({
				...state,
				connected: false,
				error: 'WebSocket connection error'
			}));
		};
	};

	const disconnect = () => {
		clearReconnectTimeout();
		currentRoomId = null;
		currentUsername = null;
		currentUserId = undefined;
		reconnectAttempts = MAX_RECONNECT_ATTEMPTS;

		if (socket) {
			socket.close();
			socket = null;
		}
		set({
			connectionState: 'disconnected',
			connected: false,
			messages: [],
			onlineUsers: [],
			typingUsers: [],
			notifications: [],
			error: null
		});
	};

	const sendMessage = (payload: { content: string; channelId?: string; parentMessageId?: string }) => {
		if (socket?.readyState !== WebSocket.OPEN) {
			update((state) => ({ ...state, error: 'Cannot send message: not connected' }));
			return;
		}

		socket.send(
			JSON.stringify({
				type: 'message.send',
				content: payload.content,
				channel_id: payload.channelId,
				parent_message_id: payload.parentMessageId
			})
		);
	};

	const sendTyping = (channelId: string, isTyping: boolean) => {
		if (socket?.readyState !== WebSocket.OPEN) {
			return;
		}
		socket.send(
			JSON.stringify({
				type: 'typing',
				channel_id: channelId,
				is_typing: isTyping
			})
		);
	};

	return {
		subscribe,
		connect,
		disconnect,
		sendMessage,
		sendTyping
	};
}

function applyEvent(state: WebSocketState, event: WebSocketEvent): WebSocketState {
	switch (event.type) {
		case 'history':
			return {
				...state,
				messages: event.messages ?? []
			};
		case 'message.created':
			return event.message
				? {
						...state,
						messages: [...state.messages.filter((item) => item.id !== event.message?.id), event.message],
						typingUsers: state.typingUsers.filter((name) => name !== event.message?.username)
					}
				: state;
		case 'typing':
			if (!event.typing) return state;
			return {
				...state,
				typingUsers: event.typing.is_typing
					? Array.from(new Set([...state.typingUsers, event.typing.username]))
					: state.typingUsers.filter((name) => name !== event.typing?.username)
			};
		case 'presence':
			return {
				...state,
				onlineUsers: event.presence?.online_users ?? []
			};
		case 'notification':
			return event.notification
				? {
						...state,
						notifications: [event.notification, ...state.notifications]
					}
				: state;
		default:
			return state;
	}
}

export const websocket = createWebSocketStore();
