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
		<label for={inputId} class="mb-1 block text-sm font-medium text-gray-700">
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
			class="w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-500 {error
				? 'border-red-500'
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
					<EyeOff class="h-5 w-5 text-gray-400" />
				{:else}
					<Eye class="h-5 w-5 text-gray-400" />
				{/if}
			</button>
		{/if}
	</div>
	{#if error}
		<p class="mt-1 text-sm text-red-600">{error}</p>
	{/if}
</div>
