<script lang="ts">
	import type { MessageReaction } from '$lib/types/room';

	export let reactions: MessageReaction[] = [];
	export let onAddReaction: (emoji: string) => void;

	const emojiOptions = ['👍', '❤️', '😂', '😮', '😢', '👏', '🎉'];

	$: reactionCounts = reactions.reduce(
		(acc, reaction) => {
			acc[reaction.emoji] = (acc[reaction.emoji] || 0) + 1;
			return acc;
		},
		{} as Record<string, number>
	);

	let showPicker = false;
</script>

<div class="mt-2 flex flex-wrap items-center gap-2">
	{#each Object.entries(reactionCounts) as [emoji, count]}
		<button
			on:click={() => onAddReaction(emoji)}
			class="rounded-full border border-white/10 bg-white/5 px-2 py-1 text-xs text-neutral-100 transition hover:bg-white/10"
		>
			{emoji} {count}
		</button>
	{/each}

	<div class="relative">
		<button
			on:click={() => (showPicker = !showPicker)}
			class="rounded-full border border-white/10 bg-white/5 px-2 py-1 text-xs text-neutral-300 transition hover:bg-white/10"
			title="Add reaction"
		>
			+
		</button>

		{#if showPicker}
			<div class="absolute bottom-full z-10 mb-2 flex gap-1 rounded-lg border border-white/10 bg-neutral-900 p-2 shadow-lg">
				{#each emojiOptions as emoji}
					<button
						on:click={() => {
							onAddReaction(emoji);
							showPicker = false;
						}}
						class="text-xl transition-transform hover:scale-110"
					>
						{emoji}
					</button>
				{/each}
			</div>
		{/if}
	</div>
</div>
