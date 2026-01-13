<script lang="ts">
	export let variant: 'primary' | 'secondary' | 'danger' = 'primary';
	export let type: 'button' | 'submit' = 'button';
	export let disabled = false;
	export let fullWidth = false;
	export let loading = false;

	const baseClasses =
		'px-4 py-2 rounded-md font-medium transition duration-200 focus:outline-none focus:ring-2 focus:ring-offset-2 inline-flex items-center justify-center gap-2';

	const variantClasses = {
		primary: 'bg-white text-neutral-950 hover:bg-neutral-100 focus:ring-white/50 font-medium',
		secondary: 'bg-white/5 text-white border border-white/10 hover:bg-white/10 focus:ring-white/20',
		danger: 'bg-red-500 text-white hover:bg-red-600 focus:ring-red-500/50'
	};

	$: isDisabled = disabled || loading;
	$: buttonClasses = [
		baseClasses,
		variantClasses[variant],
		isDisabled ? 'opacity-50 cursor-not-allowed' : '',
		fullWidth ? 'w-full' : ''
	]
		.filter(Boolean)
		.join(' ');
</script>

<button {type} disabled={isDisabled} class={buttonClasses} on:click>
	{#if loading}
		<svg
			class="h-4 w-4 animate-spin"
			xmlns="http://www.w3.org/2000/svg"
			fill="none"
			viewBox="0 0 24 24"
		>
			<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"
			></circle>
			<path
				class="opacity-75"
				fill="currentColor"
				d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
			></path>
		</svg>
	{/if}
	<slot />
</button>
