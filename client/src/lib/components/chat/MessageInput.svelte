<script lang="ts">
	export let value = '';
	export let isConnected = false;
	export let onSendMessage: () => void;

	function handleKeyDown(event: KeyboardEvent) {
		if (event.key === 'Enter' && !event.shiftKey) {
			event.preventDefault();
			console.log('Enter key pressed, calling onSendMessage');
			onSendMessage();
		}
	}

	function handleClick() {
		console.log('Send button clicked');
		console.log('Value:', value);
		console.log('isConnected:', isConnected);
		console.log('onSendMessage:', onSendMessage);
		onSendMessage();
	}
</script>

<div class="border-t border-white/10 bg-neutral-950/80 p-4 backdrop-blur-xl">
	<div class="mx-auto flex max-w-4xl gap-2">
		<input
			type="text"
			bind:value
			on:keydown={handleKeyDown}
			placeholder={isConnected ? 'Type a message...' : 'Connecting...'}
			class="flex-1 rounded-lg border border-white/10 bg-white/5 px-4 py-2 text-white placeholder-neutral-500 backdrop-blur-sm transition focus:border-white/20 focus:ring-2 focus:ring-white/20 focus:outline-none disabled:cursor-not-allowed disabled:opacity-50"
			disabled={!isConnected}
		/>

		<button
			on:click={handleClick}
			disabled={!isConnected || !value.trim()}
			class="rounded-lg bg-white px-6 py-2 font-medium text-neutral-950 transition hover:bg-neutral-100 disabled:cursor-not-allowed disabled:opacity-50"
		>
			Send
		</button>
	</div>
</div>
