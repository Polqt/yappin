<script lang="ts">
	import { Eye, EyeOff } from 'lucide-svelte';

	export let type: string = 'text';
	export let value: string | null | undefined = '';
	export let label: string = '';
	export let placeholder: string = '';
	export let required: boolean = false;
	export let disabled: boolean = false;
	export let error: string = '';

	let showPassword = false;
	let inputId = `input-${Math.random().toString(36).substring(2, 9)}`;

	$: inputType = type === 'password' && showPassword ? 'text' : type;
	$: inputValue = value || '';
</script>

<div class="mb-4">
	{#if label}
		<label for={inputId} class="mb-2 block text-sm font-medium text-neutral-300">
			{label}
		</label>
	{/if}
	<div class="relative">
		<input
			id={inputId}
			type={inputType}
			bind:value={inputValue}
			on:input={(e) => {
				const target = e.target as HTMLInputElement;
				value = target.value || null;
			}}
			{placeholder}
			{required}
			{disabled}
			class="w-full rounded-lg border border-white/10 bg-white/5 px-4 py-3 text-sm text-white placeholder-neutral-500 backdrop-blur-sm transition focus:border-white/30 focus:ring-2 focus:ring-white/20 focus:outline-none {error
				? 'border-red-500/50'
				: ''}"
		/>
		{#if type === 'password'}
			<button
				type="button"
				class="absolute inset-y-0 right-0 flex items-center pr-3"
				on:click={() => (showPassword = !showPassword)}
				aria-label={showPassword ? 'Hide password' : 'Show password'}
			>
				{#if showPassword}
					<EyeOff class="h-5 w-5 text-neutral-400" />
				{:else}
					<Eye class="h-5 w-5 text-neutral-400" />
				{/if}
			</button>
		{/if}
	</div>
	{#if error}
		<p class="mt-1 text-sm text-red-400">{error}</p>
	{/if}
</div>
