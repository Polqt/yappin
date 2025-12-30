<script lang="ts">
	import type { MessageReaction } from '$lib/types/room';

	export let reactions: MessageReaction[] = [];
	let messageId: string;
	export let onAddReaction: (emoji: string) => void;

	const emojiOptions = ['👍', '❤️', '😂', '😮', '😢', '👏'];

	$: reactionCounts = reactions.reduce(
		(acc, reaction) => {
			acc[reaction.emoji] = (acc[reaction.emoji] || 0) + 1;
			return acc;
		},
		{} as Record<string, number>
	);

	let showPicker = false;

	function handleEmojiClick(emoji: string) {
		onAddReaction(emoji);
		showPicker = false;
	}
</script>

<div class="mt-2 flex items-center gap-2">
	{#each Object.entries(reactionCounts) as [emoji, count]}
		<button
			on:click={() => handleEmojiClick(emoji)}
			aria-label="React with {emoji}"
			class="flex items-center gap-1 rounded-full bg-gray-100 px-2 py-1 text-sm transition-colors hover:bg-gray-200"
		>
			<span>{emoji}</span>
			<span class="text-gray-600">{count}</span>
		</button>
	{/each}

	<div class="relative">
		<button
			on:click={() => (showPicker = !showPicker)}
			class="px-2 py-1 text-gray-400 transition-colors hover:text-gray-600"
			title="Add Reaction"
		>
			+
		</button>

		{#if showPicker}
			<div
				class="absolute bottom-full z-10 mb-2 flex gap-1 rounded-lg border bg-white p-2 shadow-lg"
			>
				{#each emojiOptions as emoji}
					<button
						on:click={() => handleEmojiClick(emoji)}
						class="text-2xl transition-transform hover:scale-125"
					>
						{emoji}
					</button>
				{/each}
			</div>
		{/if}
	</div>
</div>
