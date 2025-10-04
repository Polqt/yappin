import type { Message, Room } from "$lib/types/room";
import { writable } from "svelte/store";

export const activeRoom = writable<Room | null>(null);
export const messages = writable<Message[]>([]);