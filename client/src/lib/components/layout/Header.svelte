<script lang="ts">
	import { auth } from '$stores/auth';
	import { goto } from '$app/navigation';
	import { MessageCircle, User, Trophy, LogOut } from 'lucide-svelte';

	async function handleLogout() {
		await auth.logout();
		goto('/login');
	}
</script>

<header class="bg-white shadow-md">
	<div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
		<div class="flex h-16 items-center justify-between">
			<!-- Logo -->
			<div class="flex items-center">
				<a href="/dashboard" class="flex items-center gap-2">
					<MessageCircle class="h-8 w-8 text-blue-600" />
					<span class="text-2xl font-bold text-gray-900">Yappin</span>
				</a>
			</div>

			<!-- Navigation -->
			{#if $auth.user}
				<nav class="flex items-center gap-6">
					<a
						href="/dashboard"
						class="flex items-center gap-1 text-gray-700 transition hover:text-blue-600"
					>
						<MessageCircle class="h-5 w-5" />
						<span class="hidden sm:inline">Rooms</span>
					</a>
					<a
						href="/dashboard/leaderboard"
						class="flex items-center gap-1 text-gray-700 transition hover:text-blue-600"
					>
						<Trophy class="h-5 w-5" />
						<span class="hidden sm:inline">Leaderboard</span>
					</a>
					<a
						href="/profile"
						class="flex items-center gap-1 text-gray-700 transition hover:text-blue-600"
					>
						<User class="h-5 w-5" />
						<span class="hidden sm:inline">Profile</span>
					</a>
					<button
						on:click={handleLogout}
						class="flex items-center gap-1 text-gray-700 transition hover:text-red-600"
					>
						<LogOut class="h-5 w-5" />
						<span class="hidden sm:inline">Logout</span>
					</button>

					<div class="ml-4 flex items-center gap-2 rounded-full bg-blue-100 px-3 py-1">
						<div
							class="flex h-8 w-8 items-center justify-center rounded-full bg-blue-600 text-sm font-bold text-white"
						>
							{$auth.user.username.charAt(0).toUpperCase()}
						</div>
						<span class="hidden font-medium text-gray-900 sm:inline">{$auth.user.username}</span>
					</div>
				</nav>
			{/if}
		</div>
	</div>
</header>
