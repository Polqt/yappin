<script lang="ts">
	import { auth } from '$stores/auth';
	import { goto } from '$app/navigation';
	import Button from '$lib/components/common/Button.svelte';
	import Input from '$lib/components/common/Input.svelte';
	import { MessageCircle } from 'lucide-svelte';
	import { getErrorMessage } from '$lib/utils/error';

	let username = '';
	let email = '';
	let password = '';
	let loading = false;
	let error = '';

	async function handleSubmit() {
		if (!username || !email || !password) {
			error = 'Please fill in all fields';
			return;
		}
		if (password.length < 6) {
			error = 'Password must be at least 6 characters';
			return;
		}

		try {
			loading = true;
			error = '';
			await auth.signup({ username, email, password });
			goto('/dashboard');
		} catch (err) {
			error = getErrorMessage(err);
		} finally {
			loading = false;
		}
	}
</script>

<div
	class="flex min-h-screen items-center justify-center bg-gradient-to-br from-green-50 to-emerald-100 px-4 py-12 sm:px-6 lg:px-8"
>
	<div class="w-full max-w-md rounded-lg bg-white p-8 shadow-xl">
		<div class="mb-8 text-center">
			<MessageCircle class="mx-auto h-12 w-12 text-green-600" />
			<h2 class="mt-4 text-3xl font-bold text-gray-900">Join Yappin</h2>
			<p class="mt-2 text-sm text-gray-600">Create your account to start chatting</p>
		</div>

		<form class="space-y-6" on:submit|preventDefault={handleSubmit}>
			{#if error}
				<div class="rounded-md border border-red-200 bg-red-50 px-4 py-3 text-red-700">
					<p class="text-sm">{error}</p>
				</div>
			{/if}

			<Input
				type="text"
				bind:value={username}
				label="Username"
				placeholder="Choose a username"
				required
			/>

			<Input
				type="email"
				bind:value={email}
				label="Email Address"
				placeholder="Enter your email"
				required
			/>

			<Input
				type="password"
				bind:value={password}
				label="Password"
				placeholder="Create a password (min 6 chars)"
				required
			/>

			<Button type="submit" variant="primary" disabled={loading}>
				{loading ? 'Creating account...' : 'Create Account'}
			</Button>
		</form>

		<div class="mt-6 text-center">
			<p class="text-sm text-gray-600">
				Already have an account?
				<a
					href="/login"
					class="font-medium text-green-600 transition duration-200 hover:text-green-500"
				>
					Sign in
				</a>
			</p>
		</div>
	</div>
</div>
