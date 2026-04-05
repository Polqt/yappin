<script lang="ts">
	import { auth } from '$stores/auth';
	import { goto } from '$app/navigation';
	import { ROUTES } from '$lib/constants/api';
	import Button from '$lib/components/common/Button.svelte';
	import Input from '$lib/components/common/Input.svelte';
	import PageBackground from '$lib/components/layout/PageBackground.svelte';
	import { MessageCircle } from 'lucide-svelte';
	import { getErrorMessage } from '$lib/utils/error';

	let email = '';
	let password = '';
	let loading = false;
	let error = '';

	async function handleSubmit() {
		if (!email || !password) {
			error = 'Please fill in all fields';
			return;
		}

		try {
			loading = true;
			error = '';
			await auth.login({ email, password });
			goto(ROUTES.dashboard);
		} catch (err) {
			error = getErrorMessage(err);
		} finally {
			loading = false;
		}
	}
</script>

<PageBackground centered>
	<div class="w-full max-w-md">
		<!-- Logo -->
		<a href={ROUTES.home} class="mb-12 flex items-center justify-center gap-2">
			<MessageCircle class="h-8 w-8 text-white" strokeWidth={1.5} />
			<span class="text-xl font-light text-white">Yappin</span>
		</a>

		<!-- Form Card -->
		<div class="rounded-2xl border border-white/10 bg-white/5 p-8 backdrop-blur-xl">
			<div class="mb-8">
				<h2 class="text-2xl font-light text-white">Welcome back</h2>
				<p class="mt-1 text-sm text-neutral-400">Sign in to continue</p>
			</div>

			<form class="space-y-5" on:submit|preventDefault={handleSubmit}>
				{#if error}
					<div class="rounded-lg border border-red-500/20 bg-red-500/10 px-4 py-3">
						<p class="text-sm text-red-200">{error}</p>
					</div>
				{/if}

				<Input
					type="email"
					bind:value={email}
					label="Email"
					placeholder="you@example.com"
					required
				/>

				<Input
					type="password"
					bind:value={password}
					label="Password"
					placeholder="Enter your password"
					required
				/>

				<Button type="submit" variant="primary" {loading} fullWidth>
					{loading ? 'Signing in...' : 'Sign In'}
				</Button>
			</form>

			<div class="mt-6 text-center">
				<p class="text-sm text-neutral-400">
					New to Yappin?
					<a href={ROUTES.signup} class="font-medium text-white transition hover:text-neutral-300">
						Create account
					</a>
				</p>
			</div>
		</div>
	</div>
</PageBackground>
