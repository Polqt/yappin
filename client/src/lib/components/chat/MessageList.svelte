<script lang="ts">
	import type { Message } from '$lib/types/room';
	import { auth } from '$stores/auth';
	import UserAvatar from './UserAvatar.svelte';

	export let messages: Message[] = [];
</script>

<div class="mx-auto max-w-4xl space-y-4">
	{#each messages as message}
		<div
			class="flex gap-3 rounded-lg bg-white p-4 shadow transition hover:shadow-md"
			class:bg-blue-50={message.user_id === $auth.user?.id}
		>
			<!-- Avatar -->
			<UserAvatar username={message.username} size="md" />

			<!-- Message content -->
			<div class="flex-1">
				<div class="mb-1 flex items-center justify-between">
					<span class="text-sm font-semibold text-gray-700">
						{message.username}
					</span>
					{#if message.system}
						<span class="rounded-full bg-gray-200 px-2 py-0.5 text-xs text-gray-600">System</span>
					{/if}
				</div>

				<div class="text-gray-900">
					{message.content}
				</div>

				{#if message.created_at}
					<div class="mt-2 text-xs text-gray-500">
						{new Date(message.created_at).toLocaleTimeString()}
					</div>
				{/if}
			</div>
		</div>
	{/each}

	{#if messages.length === 0}
		<div class="mt-8 text-center text-gray-500">
			<p class="text-lg">No messages yet</p>
			<p class="text-sm">Start the conversation!</p>
		</div>
	{/if}
</div>
