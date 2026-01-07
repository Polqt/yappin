<script lang="ts">
	import { auth } from '$stores/auth';
	import { MessageCircle, ArrowRight } from 'lucide-svelte';
</script>

<div class="relative min-h-screen overflow-hidden bg-neutral-950">
	<div class="absolute inset-0 bg-gradient-to-br from-neutral-900 via-neutral-950 to-black"></div>

	<div
		class="absolute inset-0 bg-[url('data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iNjAiIGhlaWdodD0iNjAiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+PGRlZnM+PHBhdHRlcm4gaWQ9ImdyaWQiIHdpZHRoPSI2MCIgaGVpZ2h0PSI2MCIgcGF0dGVyblVuaXRzPSJ1c2VyU3BhY2VPblVzZSI+PHBhdGggZD0iTSAxMCAwIEwgMCAwIDAgMTAiIGZpbGw9Im5vbmUiIHN0cm9rZT0icmdiYSgyNTUsMjU1LDI1NSwwLjAzKSIgc3Ryb2tlLXdpZHRoPSIxIi8+PC9wYXR0ZXJuPjwvZGVmcz48cmVjdCB3aWR0aD0iMTAwJSIgaGVpZ2h0PSIxMDAlIiBmaWxsPSJ1cmwoI2dyaWQpIi8+PC9zdmc+')] opacity-40"
	></div>

	<div class="relative z-10 flex min-h-screen flex-col items-center justify-center px-4">
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
					href="/dashboard"
					class="group flex items-center gap-3 rounded-xl border border-white/10 bg-white px-8 py-4 transition-all hover:bg-white/95"
				>
					<span class="text-sm font-medium text-neutral-950">Go to Dashboard</span>
					<ArrowRight
						class="h-4 w-4 text-neutral-950 transition-transform group-hover:translate-x-1"
					/>
				</a>
			{:else}
				<a
					href="/login"
					class="group flex items-center gap-3 rounded-xl border border-white/10 bg-white px-8 py-4 transition-all hover:bg-white/95"
				>
					<span class="text-sm font-medium text-neutral-950">Sign In</span>
					<ArrowRight
						class="h-4 w-4 text-neutral-950 transition-transform group-hover:translate-x-1"
					/>
				</a>
				<a
					href="/signup"
					class="flex items-center gap-3 rounded-xl border border-white/20 bg-white/5 px-8 py-4 backdrop-blur-sm transition-all hover:border-white/30 hover:bg-white/10"
				>
					<span class="text-sm font-medium text-white">Create Account</span>
				</a>
			{/if}
		</div>

		<!-- Features -->
		<div class="mt-24 grid max-w-4xl grid-cols-1 gap-6 sm:grid-cols-3">
			<div class="rounded-xl border border-white/10 bg-white/5 p-6 backdrop-blur-sm">
				<div class="mb-3 text-3xl">⚡</div>
				<h3 class="mb-2 text-sm font-medium text-white">Real-time</h3>
				<p class="text-xs font-light text-neutral-400">
					Instant messaging with WebSocket technology
				</p>
			</div>
			<div class="rounded-xl border border-white/10 bg-white/5 p-6 backdrop-blur-sm">
				<div class="mb-3 text-3xl">🔒</div>
				<h3 class="mb-2 text-sm font-medium text-white">Secure</h3>
				<p class="text-xs font-light text-neutral-400">End-to-end encrypted conversations</p>
			</div>
			<div class="rounded-xl border border-white/10 bg-white/5 p-6 backdrop-blur-sm">
				<div class="mb-3 text-3xl">🎯</div>
				<h3 class="mb-2 text-sm font-medium text-white">Simple</h3>
				<p class="text-xs font-light text-neutral-400">Clean interface, zero distractions</p>
			</div>
		</div>
	</div>
</div>
