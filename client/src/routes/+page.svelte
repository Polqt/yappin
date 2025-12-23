<script lang="ts">
  import { page } from '$app/stores'; 
  import { onMount, onDestroy } from 'svelte'; 
  import { websocket } from '$stores/websocket'; 
  import type { Message } from '$lib/types/room';
	import { auth } from '$stores/auth';
  
  // Get room ID from URL parameter
  let roomId = $page.params.id;
  
  let messageInput = ''; 
  let messages: Message[] = [];
  let messagesContainer: HTMLDivElement;
  let isConnected = false;
  
  // onMount runs when component first appears
  onMount(() => {
    // Connect to WebSocket when page loads
    websocket.connect(roomId);
    
    // Subscribe to WebSocket state changes
    const unsubscribe = websocket.subscribe((state) => {
      messages = state.messages;
      isConnected = state.connected; 
      
      // Auto-scroll to bottom when new message arrives
      if (messagesContainer) {
        setTimeout(() => {
          messagesContainer.scrollTop = messagesContainer.scrollHeight;
        }, 50);
      }
    });
    
    // Return cleanup function
    return unsubscribe;
  });
  
  // onDestroy runs when component is removed
  onDestroy(() => {
    websocket.disconnect(); 
  });
  
  // Send message function
  function sendMessage() {
    if (!messageInput.trim()) return; // Don't send empty messages
    
    websocket.sendMessage(messageInput); // Send to server via WebSocket
    messageInput = '';
  }
  
  // Send message when Enter key is pressed
  function handleKeydown(event: KeyboardEvent) {
    if (event.key === 'Enter' && !event.shiftKey) {
      event.preventDefault(); 
      sendMessage();
    }
  }
</script>

<div class="fixed top-4 right-4 z-50">
  {#if isConnected}
    <div class="bg-green-500 text-white px-4 py-2 rounded-lg shadow-lg">
      ● Connected
    </div>
  {:else}
    <div class="bg-red-500 text-white px-4 py-2 rounded-lg shadow-lg">
      ● Disconnected
    </div>
  {/if}
</div>

<!-- Main chat container -->
<div class="flex flex-col h-screen bg-gray-100">
  <!-- Header -->
  <div class="bg-white shadow-md p-4">
    <div class="max-w-4xl mx-auto flex items-center justify-between">
      <h1 class="text-2xl font-bold text-gray-800">Room: {roomId}</h1>
      <a 
        href="/dashboard" 
        class="text-blue-600 hover:text-blue-800"
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
    <div class="max-w-4xl mx-auto space-y-4">
      {#each messages as message}
        <div 
          class="bg-white rounded-lg shadow p-4"
          class:bg-blue-50={message.user_id === $auth.user?.id}
        >
          <!-- Username -->
          <div class="font-semibold text-sm text-gray-700 mb-1">
            {message.username}
          </div>
          
          <!-- Message content -->
          <div class="text-gray-900">
            {message.content}
          </div>
          
          <!-- Timestamp  -->
          <div class="text-xs text-gray-500 mt-2">
            {new Date(message.created_at).toLocaleTimeString()}
          </div>
        </div>
      {/each}
      
      {#if messages.length === 0}
        <div class="text-center text-gray-500 mt-8">
          No messages yet. Start the conversation!
        </div>
      {/if}
    </div>
  </div>
  
  <!-- Message input area -->
  <div class="bg-white border-t border-gray-200 p-4">
    <div class="max-w-4xl mx-auto flex gap-2">
      <input
        type="text"
        bind:value={messageInput}
        on:keydown={handleKeydown}
        placeholder="Type a message..."
        class="flex-1 px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
        disabled={!isConnected}
      />
      <button
        on:click={sendMessage}
        disabled={!isConnected || !messageInput.trim()}
        class="bg-blue-600 text-white px-6 py-2 rounded-lg hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition"
      >
        Send
      </button>
    </div>
  </div>
</div>