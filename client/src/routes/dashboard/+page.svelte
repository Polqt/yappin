<script lang="ts">
  import { onMount } from 'svelte';
  import { rooms } from '$stores/room';
  import { roomService } from '$services/room';
	import { goto } from '$app/navigation';

  onMount(async () => {
    try {
      const fetchedRooms = await roomService.getRooms();
      rooms.set(fetchedRooms);
    } catch (err) {
      console.error('Failed to load rooms:', err);
    }
  });
  async function handleJoinRoom(roomId: string) {
    try{
      await goto(`/room/${roomId}`);
    }catch(err){
      console.error('Failed to join room:', err);
    }
  }

</script>

<main class="min-h-screen bg-gray-100 p-8">
  <div class="max-w-4xl mx-auto">
    <h1 class="text-3xl font-bold text-gray-800 mb-8">Room Dashboard</h1>
    
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      {#each $rooms as room}
        <div class="bg-white rounded-lg shadow-md p-6">
          <h2 class="text-xl font-semibold text-gray-800 mb-2">{room.name}</h2>
          {#if room.description}
            <p class="text-gray-600 mb-4">{room.description}</p>
          {/if}
          <div class="text-sm text-gray-500">
            <p>Participants: {room.participants}</p>
            <p>Created by: {room.createdBy}</p>
          </div>
          <button
            type="button"
            on:click={() => handleJoinRoom(room.id)}
            class="mt-4 bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600 transition-colors"
          >
            Join Room
          </button>
        </div>
      {/each}
    </div>
    
    {#if $rooms.length === 0}
      <p class="text-gray-600">No rooms available.</p>
    {/if}
  </div>
</main>
