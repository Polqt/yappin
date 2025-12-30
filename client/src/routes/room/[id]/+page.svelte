<script lang="ts">
	import { page } from '$app/stores';
	import { onMount, onDestroy } from 'svelte';
	import { websocket } from '$stores/websocket';
	import { auth } from '$stores/auth';
	import { roomService } from '$services/room';
	import MessageList from '$lib/components/chat/MessageList.svelte';
	import MessageInput from '$lib/components/chat/MessageInput.svelte';
	import type { Message, Room } from '$lib/types/room';

	// Get room ID from URL
	let roomId = $page.params.id || '';

	// State variables
	let messageInput = '';
	let messages: Message[] = [];
	let messagesContainer: HTMLDivElement;
	let isConnected = false;
	let room: Room | null = null;
	let loadingRoom = true;

	// Connect to WebSocket when component mounts
	onMount(async () => {
		console.log('Room page mounted. RoomId:', roomId, 'Username:', $auth.user?.username);

		// Wait for auth to be ready
		if ($auth.loading) {
			console.log('Waiting for auth to load...');
			await new Promise((resolve) => setTimeout(resolve, 100));
		}

		// Fetch room details first
		try {
			room = await roomService.getRoomById(roomId);
			loadingRoom = false;
			console.log('Room loaded:', room?.name);
		} catch (error) {
			console.error('Failed to load room:', error);
			loadingRoom = false;
		}

		// Ensure we have valid roomId and username before connecting
		if (!roomId || roomId === 'undefined') {
			console.error('Invalid roomId:', roomId);
			return;
		}

		const username = $auth.user?.username || 'Guest';
		const userId = $auth.user?.id;
		console.log('Connecting with roomId:', roomId, 'username:', username, 'userId:', userId);

		// Then connect to WebSocket
		websocket.connect(roomId, username, userId);

		const unsubscribe = websocket.subscribe((state) => {
			messages = state.messages;
			isConnected = state.connected;

			// Auto-scroll to bottom
			if (messagesContainer) {
				setTimeout(() => {
					messagesContainer.scrollTop = messagesContainer.scrollHeight;
				}, 50);
			}
		});

		return unsubscribe;
	});

	// Disconnect when leaving page
	onDestroy(() => {
		console.log('Room page destroyed, disconnecting WebSocket');
		websocket.disconnect();
	});

	// Send message handler
	function sendMessage() {
		console.log('sendMessage called. messageInput:', messageInput);
		if (!messageInput.trim()) {
			console.log('Message is empty, returning');
			return;
		}
		console.log('Calling websocket.sendMessage with:', messageInput);
		websocket.sendMessage(messageInput);
		console.log('Message sent, clearing input');
		messageInput = '';
	}
</script>

<!-- Connection status indicator -->
<div class="fixed right-4 top-4 z-50">
	{#if isConnected}
		<div class="rounded-lg bg-green-500 px-4 py-2 text-white shadow-lg">● Connected</div>
	{:else}
		<div class="rounded-lg bg-red-500 px-4 py-2 text-white shadow-lg">● Disconnected</div>
	{/if}
</div>

<!-- Main chat container -->
<div class="flex h-screen flex-col bg-gray-100">
	<!-- Header -->
	<div class="bg-white p-4 shadow-md">
		<div class="mx-auto flex max-w-4xl items-center justify-between">
			{#if loadingRoom}
				<h1 class="text-2xl font-bold text-gray-800">Loading room...</h1>
			{:else if room}
				<div>
					<h1 class="text-2xl font-bold text-gray-800">{room.name}</h1>
					{#if room.topic_description}
						<p class="text-sm text-gray-600">{room.topic_description}</p>
					{/if}
				</div>
			{:else}
				<h1 class="text-2xl font-bold text-gray-800">Room not found</h1>
			{/if}
			<a href="/dashboard" class="text-blue-600 transition hover:text-blue-800">
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
