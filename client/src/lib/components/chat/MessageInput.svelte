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

<div class="border-t border-gray-200 bg-white p-4 shadow-lg">
	<div class="mx-auto flex max-w-4xl gap-2">
		<input
			type="text"
			bind:value
			on:keydown={handleKeyDown}
			placeholder={isConnected ? 'Type a message...' : 'Connecting...'}
			class="flex-1 rounded-lg border border-gray-300 px-4 py-2 transition focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:bg-gray-100"
			disabled={!isConnected}
		/>

		<button
			on:click={handleClick}
			disabled={!isConnected || !value.trim()}
			class="rounded-lg bg-blue-600 px-6 py-2 font-semibold text-white transition hover:bg-blue-700 disabled:cursor-not-allowed disabled:bg-gray-400"
		>
			Send
		</button>
	</div>
</div>
