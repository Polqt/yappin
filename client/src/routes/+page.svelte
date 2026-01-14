<script lang="ts">
	import { auth } from '$stores/auth';
	import { ROUTES } from '$lib/constants/api';
	import PageBackground from '$lib/components/layout/PageBackground.svelte';
	import { MessageCircle, ArrowRight } from 'lucide-svelte';
</script>

<PageBackground centered>
	<!-- Logo and Brand -->
	<div class="mb-12 text-center">
		<div
			class="mb-6 inline-flex items-center justify-center rounded-2xl bg-white/5 p-6 backdrop-blur-sm"
		>
			<MessageCircle class="h-16 w-16 text-white" strokeWidth={1.5} />
		</div>
		<h1 class="mb-3 text-7xl font-light tracking-tight text-white">Yappin</h1>
		<p class="text-lg font-light text-neutral-400">Real-time conversations, simplified</p>
	</div>

	<!-- CTA Buttons -->
	<div class="flex flex-col gap-4 sm:flex-row">
		{#if $auth.loading}
			<div
				class="flex items-center gap-3 rounded-xl border border-white/10 bg-white/5 px-8 py-4 backdrop-blur-sm"
			>
				<div
					class="h-5 w-5 animate-spin rounded-full border-2 border-white/20 border-t-white"
				></div>
				<span class="text-sm font-medium text-white">Loading...</span>
			</div>
		{:else if $auth.user}
			<a
				href={ROUTES.dashboard}
				class="group flex items-center gap-3 rounded-xl border border-white/10 bg-white px-8 py-4 transition-all hover:bg-white/95"
			>
				<span class="text-sm font-medium text-neutral-950">Go to Dashboard</span>
				<ArrowRight
					class="h-4 w-4 text-neutral-950 transition-transform group-hover:translate-x-1"
				/>
			</a>
		{:else}
			<a
				href={ROUTES.login}
				class="group flex items-center gap-3 rounded-xl border border-white/10 bg-white px-8 py-4 transition-all hover:bg-white/95"
			>
				<span class="text-sm font-medium text-neutral-950">Sign In</span>
				<ArrowRight
					class="h-4 w-4 text-neutral-950 transition-transform group-hover:translate-x-1"
				/>
			</a>
			<a
				href={ROUTES.signup}
				class="flex items-center gap-3 rounded-xl border border-white/20 bg-white/5 px-8 py-4 backdrop-blur-sm transition-all hover:border-white/30 hover:bg-white/10"
			>
				<span class="text-sm font-medium text-white">Create Account</span>
			</a>
		{/if}
	</div>
</PageBackground>
