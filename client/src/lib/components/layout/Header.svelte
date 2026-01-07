<script lang="ts">
	import { auth } from '$stores/auth';
	import { goto } from '$app/navigation';
	import { MessageCircle, User, Trophy, LogOut } from 'lucide-svelte';

	async function handleLogout() {
		await auth.logout();
		goto('/');
	}
</script>

<header class="border-b border-white/10 bg-neutral-950/80 backdrop-blur-xl">
	<div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
		<div class="flex h-16 items-center justify-between">
			<!-- Logo -->
			<div class="flex items-center">
				<a href="/" class="flex items-center gap-2">
					<MessageCircle class="h-6 w-6 text-white" strokeWidth={1.5} />
					<span class="text-lg font-light text-white">Yappin</span>
				</a>
			</div>

			<!-- Navigation -->
			{#if $auth.user}
				<nav class="flex items-center gap-1">
					<a
						href="/dashboard"
						class="flex items-center gap-2 rounded-lg px-3 py-2 text-sm text-neutral-400 transition hover:bg-white/5 hover:text-white"
					>
						<MessageCircle class="h-4 w-4" strokeWidth={1.5} />
						<span class="hidden sm:inline">Rooms</span>
					</a>
					<a
						href="/dashboard/leaderboard"
						class="flex items-center gap-2 rounded-lg px-3 py-2 text-sm text-neutral-400 transition hover:bg-white/5 hover:text-white"
					>
						<Trophy class="h-4 w-4" strokeWidth={1.5} />
						<span class="hidden sm:inline">Leaderboard</span>
					</a>
					<a
						href="/profile"
						class="flex items-center gap-2 rounded-lg px-3 py-2 text-sm text-neutral-400 transition hover:bg-white/5 hover:text-white"
					>
						<User class="h-4 w-4" strokeWidth={1.5} />
						<span class="hidden sm:inline">Profile</span>
					</a>

					<div class="mx-2 h-6 w-px bg-white/10"></div>

					<button
						on:click={handleLogout}
						class="flex items-center gap-2 rounded-lg px-3 py-2 text-sm text-neutral-400 transition hover:bg-white/5 hover:text-red-400"
					>
						<LogOut class="h-4 w-4" strokeWidth={1.5} />
						<span class="hidden sm:inline">Logout</span>
					</button>

					<div
						class="ml-2 flex items-center gap-2 rounded-lg border border-white/10 bg-white/5 px-3 py-1.5"
					>
						<div
							class="flex h-7 w-7 items-center justify-center rounded-full bg-white text-xs font-medium text-neutral-950"
						>
							{$auth.user.username.charAt(0).toUpperCase()}
						</div>
						<span class="hidden text-sm font-medium text-white sm:inline"
							>{$auth.user.username}</span
						>
					</div>
				</nav>
			{/if}
		</div>
	</div>
</header>
