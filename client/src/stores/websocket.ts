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

	const connect = (roomId: string, username: string) => {
		const wsUrl = `${WS_BASE_URL}${API_ENDPOINTS.rooms.join(roomId)}?username=${encodeURIComponent(username)}`;

		console.log('Connecting to WebSocket:', wsUrl);

		socket = new WebSocket(wsUrl);

		socket.onopen = () => {
			console.log('WebSocket connected!');
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
			console.log('WebSocket closed:', event.code, event.reason);
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
		if (socket?.readyState === WebSocket.OPEN) {
			socket.send(content);
		} else {
			console.error('WebSocket is not open');
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
