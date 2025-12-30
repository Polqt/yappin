<script lang="ts">
	export let users: string[] = [];

	$: typingText =
		users.length === 1
			? `${users[0]} is typing...`
			: users.length === 2
				? `${users[0]} and ${users[1]} are typing...`
				: users.length > 2
					? `${users.slice(0, 2).join(', ')} and ${users.length - 2} others are typing...`
					: '';
</script>

{#if users.length > 0}
	<div class="flex items-center gap-2 px-4 py-2 text-sm text-gray-500">
		<div class="flex gap-1">
			<span class="animate-bounce">•</span>
			<span class="animate-bounce" style="animation-delay: 0.2s">•</span>
			<span class="animate-bounce" style="animation-delay: 0.4s">•</span>
		</div>
		<span>{typingText}</span>
	</div>
{/if}
