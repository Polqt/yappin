<script lang="ts">
  import { goto } from '$app/navigation';
  import { roomService } from '$services/room';
  import Button from '$lib/components/common/Button.svelte';
  import Input from '$lib/components/common/Input.svelte';

  let name = '';
  let loading = false;
  let error = '';

  async function handleSubmit() {
    if (!name) {
      error = 'Room name is required';
      return;
    }

    try {
      loading = true;
      error = '';
      await roomService.createRoom(name);
      goto('/dashboard');
    } catch (err) {
      error = err instanceof Error ? err.message : 'Failed to create room. Please try again.';
    } finally {
      loading = false;
    }
  }
</script>

<div class="min-h-screen bg-gray-100 flex items-center justify-center">
  <div class="max-w-md w-full bg-white rounded-lg shadow-xl p-8">
    <div class="text-center mb-8">
      <h2 class="mt-4 text-3xl font-bold text-gray-900">
        Create a New Room
      </h2>
      <p class="mt-2 text-sm text-gray-600">
        Start a new conversation.
      </p>
    </div>

    <form class="space-y-6" on:submit|preventDefault={handleSubmit}>
      {#if error}
        <div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-md">
          <p class="text-sm">{error}</p>
        </div>
      {/if}

      <Input
              type="text"
              bind:value={name}
              label="Room Name"
              placeholder="Enter a name for your room"
              required
      />

      <Button
              type="submit"
              variant="primary"
              disabled={loading}
              class="w-full bg-green-600 hover:bg-green-700 text-white font-medium py-2 px-4 rounded-md transition duration-200"
      >
        {loading ? 'Creating...' : 'Create Room'}
      </Button>
    </form>
    <div class="mt-6 text-center">
      <a href="/dashboard" class="font-medium text-blue-600 hover:text-blue-500 transition duration-200">
        Back to Dashboard
      </a>
    </div>
  </div>
</div>