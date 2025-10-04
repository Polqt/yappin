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

  const connect = (roomId: string) => {
    const wsUrl = `${import.meta.env.VITE_WS_URL}/room/${roomId}`;
    socket = new WebSocket(wsUrl);

    socket.onopen = () => {
      update(state => ({ ...state, connected: true, error: null }));
    };

    socket.onmessage = (event) => {
      const message = JSON.parse(event.data);
      update(state => ({
        ...state,
        messages: [...state.messages, message]
      }));
    };

    socket.onclose = () => {
      update(state => ({ ...state, connected: false }));
    };
  };

  return {
    subscribe,
    connect,
    disconnect: () => {
      socket?.close();
      set({ connected: false, messages: [], error: null });
    },
    sendMessage: (content: string) => {
      if (socket?.readyState === WebSocket.OPEN) {
        socket.send(JSON.stringify({ type: 'MESSAGE', content }));
      }
    }
  };
}

export const websocket = createWebSocketStore();
