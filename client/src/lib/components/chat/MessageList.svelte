<script lang="ts">
	import type { Message } from '$lib/types/room';
	import { auth } from '$stores/auth';
	import UserAvatar from './UserAvatar.svelte';

	export let messages: Message[] = [];
</script>

<div class="mx-auto max-w-4xl space-y-3">
	{#each messages as message, index (`${message.user_id || 'anon'}-${message.created_at || index}`)}
		<div
			class="flex gap-3 rounded-xl border border-white/10 p-4 backdrop-blur-sm transition hover:border-white/20 {message.user_id ===
			$auth.user?.id
				? 'bg-white/10'
				: 'bg-white/5'}"
		>
			<!-- Avatar -->
			<UserAvatar username={message.username} size="md" />

			<!-- Message content -->
			<div class="flex-1">
				<div class="mb-1 flex items-center justify-between">
					<span class="text-sm font-medium text-white">
						{message.username}
					</span>
					{#if message.system}
						<span
							class="rounded-full border border-white/10 bg-white/5 px-2 py-0.5 text-xs text-neutral-400"
							>System</span
						>
					{/if}
				</div>

				<div class="text-neutral-200">
					{message.content}
				</div>

				{#if message.created_at}
					<div class="mt-2 text-xs text-neutral-500">
						{new Date(message.created_at).toLocaleTimeString()}
					</div>
				{/if}
			</div>
		</div>
	{/each}

	{#if messages.length === 0}
		<div class="mt-8 text-center">
			<p class="text-lg font-light text-white">No messages yet</p>
			<p class="text-sm text-neutral-400">Start the conversation!</p>
		</div>
	{/if}
</div>
