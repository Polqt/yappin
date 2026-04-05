<script lang="ts">
	export let value = '';
	export let isConnected = false;
	export let placeholder = 'Type a message...';
	export let threadLabel = '';
	export let onSendMessage: () => void;
	export let onTyping: (isTyping: boolean) => void = () => {};

	let typingTimeout: ReturnType<typeof setTimeout> | null = null;

	function handleKeyDown(event: KeyboardEvent) {
		if (event.key === 'Enter' && !event.shiftKey) {
			event.preventDefault();
			onSendMessage();
			onTyping(false);
			return;
		}
		onTyping(true);
		if (typingTimeout) clearTimeout(typingTimeout);
		typingTimeout = setTimeout(() => onTyping(false), 1200);
	}
</script>

<div class="border-t border-white/10 bg-neutral-950/80 p-4 backdrop-blur-xl">
	<div class="mx-auto max-w-5xl space-y-3">
		{#if threadLabel}
			<div class="rounded-lg border border-white/10 bg-white/5 px-3 py-2 text-xs text-neutral-300">
				Replying in thread: {threadLabel}
			</div>
		{/if}

		<div class="flex gap-2">
			<textarea
				rows="2"
				bind:value
				on:keydown={handleKeyDown}
				placeholder={isConnected ? placeholder : 'Connecting...'}
				class="flex-1 resize-none rounded-xl border border-white/10 bg-white/5 px-4 py-3 text-white placeholder-neutral-500 backdrop-blur-sm transition focus:border-white/20 focus:outline-none focus:ring-2 focus:ring-white/20 disabled:cursor-not-allowed disabled:opacity-50"
				disabled={!isConnected}
			></textarea>

			<button
				on:click={() => {
					onSendMessage();
					onTyping(false);
				}}
				disabled={!isConnected || !value.trim()}
				class="rounded-xl bg-white px-6 py-2 font-medium text-neutral-950 transition hover:bg-neutral-100 disabled:cursor-not-allowed disabled:opacity-50"
			>
				Send
			</button>
		</div>
	</div>
</div>
