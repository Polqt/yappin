<script lang="ts">
	import { onMount } from 'svelte';
	import { API_BASE_URL } from '$lib/constants/api';
	import type { LeaderboardEntry } from '$lib/types/leaderboard';
	import { getMedal } from '$lib/utils/leaderboard';
	import Header from '$lib/components/layout/Header.svelte';

	// Component state
	let leaderboard: LeaderboardEntry[] = [];
	let loading = true;
	let error = '';

	// Fetch leaderboard when component loads
	onMount(async () => {
		try {
			const response = await fetch(`${API_BASE_URL}/api/stats/leaderboard?limit=100`, {
				credentials: 'include' // Send cookies
			});

			if (!response.ok) throw new Error('Failed to fetch');

			leaderboard = await response.json();
		} catch (err) {
			error = 'Failed to load leaderboard';
			console.error(err);
		} finally {
			loading = false;
		}
	});
</script>

<div class="min-h-screen bg-neutral-950">
	<Header />

	<div class="mx-auto max-w-6xl p-6 sm:p-8">
		<h1 class="mb-8 text-center text-2xl font-light text-white">🏆 Leaderboard</h1>

		{#if loading}
			<div class="py-12 text-center">
				<div
					class="mx-auto h-12 w-12 animate-spin rounded-full border-2 border-white/20 border-t-white"
				></div>
				<p class="mt-4 text-sm text-neutral-400">Loading...</p>
			</div>
		{:else if error}
			<div class="rounded-lg border border-red-500/20 bg-red-500/10 p-4">
				<p class="text-sm text-red-200">{error}</p>
			</div>
		{:else}
			<div class="overflow-hidden rounded-xl border border-white/10 bg-white/5 backdrop-blur-sm">
				<table class="w-full">
					<thead class="border-b border-white/10 bg-white/5">
						<tr>
							<th class="px-6 py-3 text-left text-xs font-medium text-neutral-400 uppercase"
								>Rank</th
							>
							<th class="px-6 py-3 text-left text-xs font-medium text-neutral-400 uppercase"
								>User</th
							>
							<th class="px-6 py-3 text-right text-xs font-medium text-neutral-400 uppercase"
								>Messages</th
							>
							<th class="px-6 py-3 text-right text-xs font-medium text-neutral-400 uppercase"
								>Upvotes</th
							>
							<th class="px-6 py-3 text-right text-xs font-medium text-neutral-400 uppercase"
								>Streak</th
							>
						</tr>
					</thead>
					<tbody class="divide-y divide-white/10">
						{#each leaderboard as user}
							<tr class="transition-colors hover:bg-white/10 {user.rank <= 3 ? 'bg-white/5' : ''}">
								<td class="px-6 py-4 whitespace-nowrap">
									<span class="text-2xl">{getMedal(user.rank)}</span>
									<span class="ml-2 font-medium text-white">#{user.rank}</span>
								</td>
								<td class="px-6 py-4 font-medium whitespace-nowrap text-white">
									{user.username}
								</td>
								<td class="px-6 py-4 text-right whitespace-nowrap text-neutral-400">
									{user.total_messages}
								</td>
								<td class="px-6 py-4 text-right whitespace-nowrap text-neutral-400">
									{user.total_upvotes} 👍
								</td>
								<td class="px-6 py-4 text-right whitespace-nowrap text-neutral-400">
									{user.daily_streak} 🔥
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		{/if}
	</div>
</div>
