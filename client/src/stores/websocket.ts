import type { Message } from '$lib/types/room';
import { writable } from 'svelte/store';
import { WS_BASE_URL, API_ENDPOINTS } from '$lib/constants/api';

interface WebSocketState {
	connected: boolean;
	messages: Message[];
	error: string | null;
}

function createWebSocketStore() {
	const { subscribe, set, update } = writable<WebSocketState>({
		connected: false,
		messages: [],
		error: null
	});

	let socket: WebSocket | null = null;

	const connect = (roomId: string, username: string, userId?: string) => {
		// Close existing connection if any
		if (socket && socket.readyState === WebSocket.OPEN) {
			console.log('Closing existing WebSocket connection');
			socket.close();
		}

		let wsUrl = `${WS_BASE_URL}${API_ENDPOINTS.rooms.join(roomId)}?username=${encodeURIComponent(username)}`;
		if (userId) {
			wsUrl += `&client_id=${encodeURIComponent(userId)}`;
		}

		console.log('Connecting to WebSocket:', wsUrl);

		socket = new WebSocket(wsUrl);

		socket.onopen = () => {
			console.log('WebSocket connected successfully');
			update((state) => ({
				...state,
				connected: true,
				error: null
			}));
		};

		socket.onmessage = (event) => {
			console.log('Message received:', event.data);
			try {
				const message = JSON.parse(event.data);
				update((state) => ({
					...state,
					messages: [...state.messages, message]
				}));
			} catch (error) {
				console.error('Error parsing WebSocket message:', error);
			}
		};

		socket.onclose = (event) => {
			console.log('WebSocket closed. Code:', event.code, 'Reason:', event.reason || 'None');
			update((state) => ({ ...state, connected: false }));
		};

		socket.onerror = (error) => {
			console.error('WebSocket error:', error);
			update((state) => ({
				...state,
				connected: false,
				error: 'WebSocket connection error'
			}));
		};
	};

	const disconnect = () => {
		if (socket) {
			socket.close();
			socket = null;
		}
		set({ connected: false, messages: [], error: null });
	};

	const sendMessage = (content: string) => {
		console.log('sendMessage called in websocket store. Content:', content);
		console.log('Socket state:', socket?.readyState);
		console.log('WebSocket.OPEN constant:', WebSocket.OPEN);

		if (socket?.readyState === WebSocket.OPEN) {
			console.log('Sending message:', content);
			socket.send(content);
			console.log('Message sent successfully');
		} else {
			console.error('WebSocket is not open. Current state:', socket?.readyState);
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
