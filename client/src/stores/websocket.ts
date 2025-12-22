import type { Message } from '$lib/types/room';
import { writable } from 'svelte/store';

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
    const wsUrl = `${import.meta.env.VITE_WS_URL}/join-room/${roomId}?username=${username}`;

    socket = new WebSocket(wsUrl);

    socket.onopen = () => {
      console.log('WebSocket connected');
      update(state => ({
        ...state,
        connected: true,
        error: null
      }))
    }

    socket.onmessage = (event) => {
      try {
        const message = JSON.parse(event.data);
        update(state => ({
          ...state,
          messages: [...state.messages, message]
        }));
      } catch (error) {
        console.error('Error parsing WebSocket message:', error);
      }
    };

    socket.onclose = () => {
      console.log('WebSocket disconnected');
      update(state => ({ ...state, connected: false }));
    };

    socket.onerror = (error) => {
      console.error('WebSocket error:', error);
      update(state => ({
        ...state,
        connected: false,
        error: 'WebSocket error occurred'
      }))
    }
  };

  const disconnect = () => {
    if (socket) {
      socket.close();
      socket = null;
    }

    set({ connected: false, messages: [], error: null });
  }

  const sendMessage = (content: string) => {
    if (socket?.readyState === WebSocket.OPEN) {
      socket.send(content);
    } else {
      console.error('WebSocket is not open. Unable to send message.');
    }
  }

  return {
    subscribe,
    connect,
    disconnect,
    sendMessage
  };
}

export const websocket = createWebSocketStore();
