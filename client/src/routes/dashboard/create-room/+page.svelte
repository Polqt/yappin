<script lang="ts">
	import { goto } from '$app/navigation';
	import { roomService } from '$services/room';
	import Button from '$lib/components/common/Button.svelte';
	import Input from '$lib/components/common/Input.svelte';
	import type { CreateRoomRequest } from '$lib/types/room';
	import { getErrorMessage } from '$lib/utils/error';

	let request: CreateRoomRequest = {
		name: '',
		expires_at: undefined
	};

	let loading = false;
	let error = '';

	async function handleSubmit() {
		if (!request.name) {
			error = 'Room name is required';
			return;
		}

		try {
			loading = true;
			error = '';
			await roomService.createRoom(request);
			goto('/dashboard');
		} catch (err) {
			error = getErrorMessage(err);
		} finally {
			loading = false;
		}
	}
</script>

<div class="flex min-h-screen items-center justify-center bg-gray-100">
	<div class="w-full max-w-md rounded-lg bg-white p-8 shadow-xl">
		<div class="mb-8 text-center">
			<h2 class="mt-4 text-3xl font-bold text-gray-900">Create a New Room</h2>
			<p class="mt-2 text-sm text-gray-600">Start a new conversation.</p>
		</div>

		<form class="space-y-6" on:submit|preventDefault={handleSubmit}>
			{#if error}
				<div class="rounded-md border border-red-200 bg-red-50 px-4 py-3 text-red-700">
					<p class="text-sm">{error}</p>
				</div>
			{/if}

			<Input
				type="text"
				bind:value={request.name}
				label="Room Name"
				placeholder="Enter a name for your room"
				required
			/>

			<Input type="datetime-local" bind:value={request.expires_at} label="Expires At (Optional)" />

			<Button type="submit" variant="primary" disabled={loading}>
				{loading ? 'Creating...' : 'Create Room'}
			</Button>
		</form>
		<div class="mt-6 text-center">
			<a
				href="/dashboard"
				class="font-medium text-blue-600 transition duration-200 hover:text-blue-500"
			>
				Back to Dashboard
			</a>
		</div>
	</div>
</div>
