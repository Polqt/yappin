<script lang="ts">
	import { onMount } from 'svelte';
	import { rooms } from '$stores/room';
	import { roomService } from '$services/room';
	import { goto } from '$app/navigation';
	import Header from '$lib/components/layout/Header.svelte';

	let loading = true;
	let error = '';

	onMount(async () => {
		loading = true;
		error = '';
		console.log('Dashboard mounted, loading rooms');
		try {
			const fetchedRooms = await roomService.getRooms();
			rooms.set(fetchedRooms);
			console.log('Rooms loaded:', fetchedRooms.length);
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

<div class="min-h-screen bg-neutral-950">
	<Header />

	<main class="p-6 sm:p-8">
		<div class="mx-auto max-w-7xl">
			<div class="mb-8 flex items-center justify-between">
				<h1 class="text-2xl font-light text-white">Rooms</h1>
				<a
					href="/dashboard/create-room"
					class="rounded-lg border border-white/10 bg-white/5 px-4 py-2 text-sm font-medium text-white backdrop-blur-sm transition hover:bg-white/10"
				>
					+ New Room
				</a>
			</div>

			{#if loading}
				<div class="py-12 text-center">
					<div
						class="mx-auto h-12 w-12 animate-spin rounded-full border-2 border-white/20 border-t-white"
					></div>
					<p class="mt-4 text-sm text-neutral-400">Loading rooms...</p>
				</div>
			{:else if error}
				<div class="rounded-lg border border-red-500/20 bg-red-500/10 px-4 py-3">
					<p class="text-sm text-red-200">{error}</p>
				</div>
			{:else if $rooms.length === 0}
				<div class="rounded-xl border border-white/10 bg-white/5 p-12 text-center backdrop-blur-sm">
					<p class="text-lg font-light text-white">No rooms available</p>
					<p class="mt-2 text-sm text-neutral-400">Create one to get started</p>
				</div>
			{:else}
				<div class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
					{#each $rooms as room}
						<div
							class="group rounded-xl border border-white/10 bg-white/5 p-5 backdrop-blur-sm transition hover:border-white/20 hover:bg-white/10"
						>
							<div class="mb-3 flex items-start justify-between">
								<h2 class="text-lg font-medium text-white">{room.name}</h2>
								{#if room.is_pinned}
									<span class="text-base" title="Pinned">📌</span>
								{/if}
							</div>
							{#if room.topic_description}
								<p class="mb-4 line-clamp-2 text-sm text-neutral-400">{room.topic_description}</p>
							{/if}
							<div class="mb-4 space-y-1.5 text-xs text-neutral-500">
								{#if room.creator_username}
									<p>
										by <span class="text-neutral-400">{room.creator_username}</span>
									</p>
								{:else}
									<p>{new Date(room.created_at).toLocaleDateString()}</p>
								{/if}
								<p>
									<span class="text-neutral-400">{room.participants}</span>
									{room.participants === 1 ? 'participant' : 'participants'}
								</p>
								<p>Expires {new Date(room.expires_at).toLocaleDateString()}</p>
							</div>
							<button
								type="button"
								on:click={() => handleJoinRoom(room.id)}
								class="w-full rounded-lg bg-white px-4 py-2 text-sm font-medium text-neutral-950 transition hover:bg-neutral-100"
							>
								Join Room
							</button>
						</div>
					{/each}
				</div>
			{/if}
		</div>
	</main>
</div>
