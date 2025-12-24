<script lang="ts">
	import { onMount } from 'svelte';
	import { rooms } from '$stores/room';
	import { roomService } from '$services/room';
	import { goto } from '$app/navigation';

	let loading = true;
	let error = '';

	onMount(async () => {
		try {
			const fetchedRooms = await roomService.getRooms();
			rooms.set(fetchedRooms);
		} catch (err) {
			error = 'Failed to load rooms';
			console.error('Failed to load rooms:', err);
		} finally {
			loading = false;
		}
	});

	async function handleJoinRoom(roomId: string) {
		try {
			await goto(`/room/${roomId}`);
		} catch (err) {
			console.error('Failed to join room:', err);
		}
	}
</script>

<main class="min-h-screen bg-gray-100 p-8">
	<div class="mx-auto max-w-4xl">
		<div class="mb-8 flex items-center justify-between">
			<h1 class="text-3xl font-bold text-gray-800">Room Dashboard</h1>
			<a
				href="/dashboard/create-room"
				class="rounded bg-green-500 px-4 py-2 text-white transition-colors hover:bg-green-600"
			>
				Create Room
			</a>
		</div>

		{#if loading}
			<div class="py-12 text-center">
				<p class="text-gray-600">Loading rooms...</p>
			</div>
		{:else if error}
			<div class="rounded-md border border-red-200 bg-red-50 px-4 py-3 text-red-700">
				<p class="text-sm">{error}</p>
			</div>
		{:else if $rooms.length === 0}
			<div class="py-12 text-center">
				<p class="text-gray-600">No rooms available. Create one to get started!</p>
			</div>
		{:else}
			<div class="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
				{#each $rooms as room}
					<div class="rounded-lg bg-white p-6 shadow-md">
						<div class="mb-2 flex items-center justify-between">
							<h2 class="text-xl font-semibold text-gray-800">{room.name}</h2>
							{#if room.is_pinned}
								<span class="text-lg" title="Pinned Room">📌</span>
							{/if}
						</div>
						{#if room.topic_description}
							<p class="mb-4 text-gray-600">{room.topic_description}</p>
						{/if}
						<div class="mb-4 space-y-1 text-sm text-gray-500">
							{#if room.creator_username}
								<p>
									Created by: <span class="font-medium text-gray-700">{room.creator_username}</span>
								</p>
							{:else}
								<p>Created: {new Date(room.created_at).toLocaleDateString()}</p>
							{/if}
							<p>
								Participants: <span class="font-medium text-gray-700">{room.participants}</span>
							</p>
							<p class="text-xs">Expires: {new Date(room.expires_at).toLocaleDateString()}</p>
						</div>
						<button
							type="button"
							on:click={() => handleJoinRoom(room.id)}
							class="w-full rounded bg-blue-500 px-4 py-2 text-white transition-colors hover:bg-blue-600"
						>
							Join Room
						</button>
					</div>
				{/each}
			</div>
		{/if}
	</div>
</main>
