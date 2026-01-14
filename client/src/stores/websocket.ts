import type { Message } from '$lib/types/room';
import { writable } from 'svelte/store';
import { WS_BASE_URL, API_ENDPOINTS } from '$lib/constants/api';
import { wsLogger } from '$lib/utils/logger';

type ConnectionState = 'disconnected' | 'connecting' | 'connected' | 'reconnecting';

interface WebSocketState {
	connectionState: ConnectionState;
	connected: boolean;
	messages: Message[];
	error: string | null;
}

// Reconnection configuration
const MAX_RECONNECT_ATTEMPTS = 5;
const RECONNECT_DELAY_MS = 3000;

function createWebSocketStore() {
	const { subscribe, set, update } = writable<WebSocketState>({
		connectionState: 'disconnected',
		connected: false,
		messages: [],
		error: null
	});

	let socket: WebSocket | null = null;
	let reconnectAttempts = 0;
	let reconnectTimeout: ReturnType<typeof setTimeout> | null = null;
	let currentRoomId: string | null = null;
	let currentUsername: string | null = null;
	let currentUserId: string | undefined = undefined;

	const clearReconnectTimeout = () => {
		if (reconnectTimeout) {
			clearTimeout(reconnectTimeout);
			reconnectTimeout = null;
		}
	};

	const connect = (roomId: string, username: string, userId?: string) => {
		// Store connection params for reconnection
		currentRoomId = roomId;
		currentUsername = username;
		currentUserId = userId;
		reconnectAttempts = 0;

		performConnect(roomId, username, userId);
	};

	const performConnect = (roomId: string, username: string, userId?: string) => {
		// Close existing connection if any
		if (socket && socket.readyState === WebSocket.OPEN) {
			wsLogger.log('Closing existing WebSocket connection');
			socket.close();
		}

		let wsUrl = `${WS_BASE_URL}${API_ENDPOINTS.rooms.join(roomId)}?username=${encodeURIComponent(username)}`;
		if (userId) {
			wsUrl += `&client_id=${encodeURIComponent(userId)}`;
		}

		wsLogger.log('Connecting to WebSocket:', wsUrl);

		update((state) => ({
			...state,
			connectionState: reconnectAttempts > 0 ? 'reconnecting' : 'connecting',
			error: null
		}));

		socket = new WebSocket(wsUrl);

		socket.onopen = () => {
			wsLogger.log('WebSocket connected successfully');
			reconnectAttempts = 0;
			update((state) => ({
				...state,
				connectionState: 'connected',
				connected: true,
				error: null
			}));
		};

		socket.onmessage = (event) => {
			wsLogger.log('Message received:', event.data);
			try {
				const message = JSON.parse(event.data);
				update((state) => ({
					...state,
					messages: [...state.messages, message]
				}));
			} catch (error) {
				wsLogger.error('Error parsing WebSocket message:', error);
			}
		};

		socket.onclose = (event) => {
			wsLogger.log('WebSocket closed. Code:', event.code, 'Reason:', event.reason || 'None');

			update((state) => ({
				...state,
				connectionState: 'disconnected',
				connected: false
			}));

			// Attempt reconnection if not a clean close and we have connection params
			if (
				!event.wasClean &&
				currentRoomId &&
				currentUsername &&
				reconnectAttempts < MAX_RECONNECT_ATTEMPTS
			) {
				reconnectAttempts++;
				const delay = RECONNECT_DELAY_MS * reconnectAttempts;
				wsLogger.log(
					`Attempting reconnection ${reconnectAttempts}/${MAX_RECONNECT_ATTEMPTS} in ${delay}ms`
				);

				update((state) => ({
					...state,
					connectionState: 'reconnecting',
					error: `Reconnecting... (${reconnectAttempts}/${MAX_RECONNECT_ATTEMPTS})`
				}));

				reconnectTimeout = setTimeout(() => {
					if (currentRoomId && currentUsername) {
						performConnect(currentRoomId, currentUsername, currentUserId);
					}
				}, delay);
			}
		};

		socket.onerror = (error) => {
			wsLogger.error('WebSocket error:', error);
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
		reconnectAttempts = MAX_RECONNECT_ATTEMPTS; // Prevent auto-reconnect

		if (socket) {
			socket.close();
			socket = null;
		}
		set({ connectionState: 'disconnected', connected: false, messages: [], error: null });
	};

	const sendMessage = (content: string) => {
		wsLogger.log('Sending message:', content);

		if (socket?.readyState === WebSocket.OPEN) {
			socket.send(content);
			wsLogger.log('Message sent successfully');
		} else {
			wsLogger.error('WebSocket is not open. Current state:', socket?.readyState);
			update((state) => ({
				...state,
				error: 'Cannot send message: not connected'
			}));
		}
	};

	return {
		subscribe,
		connect,
		disconnect,
		sendMessage
	};
}

export const websocket = createWebSocketStore();
