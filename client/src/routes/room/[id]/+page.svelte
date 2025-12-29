<script lang="ts">
    import { page } from '$app/stores';
    import { onMount, onDestroy } from 'svelte';
    import { websocket } from '$stores/websocket';
    import { auth } from '$stores/auth';
    import MessageList from '$lib/components/chat/MessageList.svelte';
    import MessageInput from '$lib/components/chat/MessageInput.svelte';
    import type { Message } from '$lib/types/room';

    // Get room ID from URL
    let roomId = $page.params.id || '';
    
    // State variables
    let messageInput = '';
    let messages: Message[] = [];
    let messagesContainer: HTMLDivElement;
    let isConnected = false;

    // Connect to WebSocket when component mounts
    onMount(() => {
        websocket.connect(roomId, $auth.user?.username || 'Guest');
        
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
        websocket.disconnect();
    });

    // Send message handler
    function sendMessage() {
        if (!messageInput.trim()) return;
        websocket.sendMessage(messageInput);
        messageInput = '';
    }
</script>

<!-- Connection status indicator -->
<div class="fixed right-4 top-4 z-50">
    {#if isConnected}
        <div class="rounded-lg bg-green-500 px-4 py-2 text-white shadow-lg">
            ● Connected
        </div>
    {:else}
        <div class="rounded-lg bg-red-500 px-4 py-2 text-white shadow-lg">
            ● Disconnected
        </div>
    {/if}
</div>

<!-- Main chat container -->
<div class="flex h-screen flex-col bg-gray-100">
    <!-- Header -->
    <div class="bg-white p-4 shadow-md">
        <div class="mx-auto flex max-w-4xl items-center justify-between">
            <h1 class="text-2xl font-bold text-gray-800">
                Room: {roomId.slice(0, 8)}...
            </h1>
            <a 
                href="/dashboard" 
                class="text-blue-600 hover:text-blue-800 transition"
            >
                ← Back to Dashboard
            </a>
        </div>
    </div>

    <!-- Messages area -->
    <div 
        bind:this={messagesContainer} 
        class="flex-1 overflow-y-auto p-4"
    >
        <MessageList {messages} />
    </div>

    <!-- Message input -->
    <MessageInput 
        bind:value={messageInput} 
        {isConnected} 
        onSendMessage={sendMessage} 
    />
</div>