import { writable } from 'svelte/store';
import type { Room } from '$lib/types/room';

function createRoomStore() {
	const { subscribe, set, update } = writable<Room[]>([]);

	return {
		subscribe,
		set,
		update,
		add: (room: Room) => update((rooms) => [...rooms, room]),
		remove: (roomId: string) => update((rooms) => rooms.filter((r) => r.id !== roomId)),
		clear: () => set([])
	};
}

export const rooms = createRoomStore();
