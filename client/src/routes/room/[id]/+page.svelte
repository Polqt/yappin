<script lang="ts">
	import { page } from '$app/stores';
	import { onMount, onDestroy } from 'svelte';
	import { websocket } from '$stores/websocket';
	import { auth } from '$stores/auth';
	import { roomService } from '$services/room';
	import MessageList from '$lib/components/chat/MessageList.svelte';
	import MessageInput from '$lib/components/chat/MessageInput.svelte';
	import type { Message, Room } from '$lib/types/room';

	let roomId = $page.params.id || '';
	let messageInput = '';
	let messages: Message[] = [];
	let messagesContainer: HTMLDivElement;
	let isConnected = false;
	let room: Room | null = null;
	let loadingRoom = true;

	// Load room data and connect to WebSocket
	onMount(async () => {
		// Wait for auth to load
		if ($auth.loading) {
			await new Promise((resolve) => setTimeout(resolve, 100));
		}

		// Validate roomId
		if (!roomId || roomId === 'undefined') {
			console.error('Invalid roomId:', roomId);
			return;
		}

		// Load room details
		try {
			room = await roomService.getRoomById(roomId);
		} catch (error) {
			console.error('Failed to load room:', error);
		} finally {
			loadingRoom = false;
		}

		// Connect to WebSocket
		const username = $auth.user?.username || 'Guest';
		const userId = $auth.user?.id;
		websocket.connect(roomId, username, userId);
	});

	// Subscribe to WebSocket state changes
	$: {
		const state = $websocket;
		messages = state.messages;
		isConnected = state.connected;

		// Auto-scroll to bottom when messages change
		if (messagesContainer && messages.length > 0) {
			setTimeout(() => {
				messagesContainer.scrollTop = messagesContainer.scrollHeight;
			}, 50);
		}
	}

	// Cleanup on component destroy
	onDestroy(() => {
		websocket.disconnect();
	});

	// Send message
	function sendMessage() {
		if (!messageInput.trim()) return;
		websocket.sendMessage(messageInput);
		messageInput = '';
	}
</script>

<!-- Connection status indicator -->
<div class="fixed right-4 top-4 z-50">
	{#if isConnected}
		<div
			class="rounded-lg border border-green-500/20 bg-green-500/10 px-4 py-2 text-green-200 backdrop-blur-sm"
		>
			● Connected
		</div>
	{:else}
		<div
			class="rounded-lg border border-red-500/20 bg-red-500/10 px-4 py-2 text-red-200 backdrop-blur-sm"
		>
			● Disconnected
		</div>
	{/if}
</div>

<!-- Main chat container -->
<div class="flex h-screen flex-col bg-neutral-950">
	<!-- Header -->
	<div class="border-b border-white/10 bg-neutral-950/80 p-4 backdrop-blur-xl">
		<div class="mx-auto flex max-w-4xl items-center justify-between">
			{#if loadingRoom}
				<h1 class="text-xl font-light text-white">Loading room...</h1>
			{:else if room}
				<div>
					<h1 class="text-xl font-medium text-white">{room.name}</h1>
					{#if room.topic_description}
						<p class="text-sm text-neutral-400">{room.topic_description}</p>
					{/if}
				</div>
			{:else}
				<h1 class="text-xl font-light text-white">Room not found</h1>
			{/if}
			<a href="/dashboard" class="text-sm text-neutral-400 transition hover:text-white">
				← Back to Dashboard
			</a>
		</div>
	</div>

	<!-- Messages area -->
	<div bind:this={messagesContainer} class="flex-1 overflow-y-auto p-4">
		<MessageList {messages} />
	</div>

	<!-- Message input -->
	<MessageInput bind:value={messageInput} {isConnected} onSendMessage={sendMessage} />
</div>
