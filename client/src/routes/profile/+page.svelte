<script lang="ts">
	import { goto } from '$app/navigation';
	import { auth } from '$stores/auth';
	import { onMount } from 'svelte';
	import { API_BASE_URL } from '$lib/constants/api';
	import type { UserProfile } from '$lib/types/user';
	import StatsCard from '$lib/components/profile/StatsCard.svelte';
	import ActivityGraph from '$lib/components/profile/ActivityGraph.svelte';
	import AchievementBadge from '$lib/components/profile/AchievementBadge.svelte';
	import Header from '$lib/components/layout/Header.svelte';

	let userProfile: UserProfile | null = null;
	let loading = true;
	let error = '';

	onMount(async () => {
		if (!$auth.user) {
			goto('/login');
			return;
		}

		try {
			loading = true;
			const response = await fetch(`${API_BASE_URL}/api/stats/profile/${$auth.user.id}`, {
				credentials: 'include'
			});

			if (!response.ok) {
				throw new Error('Failed to load profile');
			}

			userProfile = await response.json();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load profile';
			console.error('Failed to load profile:', err);
		} finally {
			loading = false;
		}
	});
</script>

<div class="min-h-screen bg-neutral-950">
	<Header />

	<div class="py-6 sm:py-8">
		<div class="mx-auto max-w-6xl px-4">
			<h1 class="mb-8 text-2xl font-light text-white">Profile</h1>

			{#if loading}
				<div class="py-12 text-center">
					<div
						class="mx-auto h-12 w-12 animate-spin rounded-full border-2 border-white/20 border-t-white"
					></div>
					<p class="mt-4 text-sm text-neutral-400">Loading profile...</p>
				</div>
			{:else if error}
				<div class="rounded-lg border border-red-500/20 bg-red-500/10 p-4">
					<p class="text-sm text-red-200">{error}</p>
				</div>
			{:else if userProfile}
				<!-- User Info Card -->
				<div class="mb-8 rounded-xl border border-white/10 bg-white/5 p-6 backdrop-blur-sm">
					<div class="flex items-center gap-4">
						<div
							class="flex h-20 w-20 items-center justify-center rounded-full bg-white text-3xl font-bold text-neutral-950"
						>
							{$auth.user?.username.charAt(0).toUpperCase() || 'U'}
						</div>
						<div>
							<h2 class="text-2xl font-medium text-white">{$auth.user?.username}</h2>
							<p class="text-sm text-neutral-400">{$auth.user?.email || 'No email provided'}</p>
						</div>
					</div>
				</div>

				<!-- Stats Grid -->
				<div class="mb-8 grid grid-cols-1 gap-4 sm:gap-6 md:grid-cols-2 lg:grid-cols-4">
					<StatsCard label="Messages Sent" value={userProfile.total_messages || 0} icon="💬" />
					<StatsCard label="Upvotes Received" value={userProfile.total_upvotes || 0} icon="👍" />
					<StatsCard label="Daily Streak" value={userProfile.daily_streak || 0} icon="🔥" />
					<StatsCard label="Level" value={userProfile.level || 1} icon="⭐" />
				</div>

				<!-- Achievements -->
				{#if userProfile.achievements && userProfile.achievements.length > 0}
					<div class="mb-8 rounded-xl border border-white/10 bg-white/5 p-6 backdrop-blur-sm">
						<h3 class="mb-4 text-lg font-medium text-white">Achievements</h3>
						<div class="grid grid-cols-2 gap-4 md:grid-cols-3 lg:grid-cols-4">
							{#each userProfile.achievements as achievement}
								<AchievementBadge {achievement} />
							{/each}
						</div>
					</div>
				{/if}

				<!-- Activity Graph -->
				<div class="rounded-xl border border-white/10 bg-white/5 p-6 backdrop-blur-sm">
					<h3 class="mb-4 text-lg font-medium text-white">Activity</h3>
					<ActivityGraph data={userProfile.activity || []} />
				</div>
			{/if}
		</div>
	</div>
</div>
