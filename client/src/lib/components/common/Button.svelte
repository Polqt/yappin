<script lang="ts">
	export let variant: 'primary' | 'secondary' | 'danger' = 'primary';
	export let type: 'button' | 'submit' = 'button';
	export let disabled = false;
	export let fullWidth = false;

	const baseClasses =
		'px-4 py-2 rounded-md font-medium transition duration-200 focus:outline-none focus:ring-2 focus:ring-offset-2';

	const variantClasses = {
		primary: 'bg-white text-neutral-950 hover:bg-neutral-100 focus:ring-white/50 font-medium',
		secondary: 'bg-white/5 text-white border border-white/10 hover:bg-white/10 focus:ring-white/20',
		danger: 'bg-red-500 text-white hover:bg-red-600 focus:ring-red-500/50'
	};

	$: buttonClasses = [
		baseClasses,
		variantClasses[variant],
		disabled ? 'opacity-50 cursor-not-allowed' : '',
		fullWidth ? 'w-full' : ''
	]
		.filter(Boolean)
		.join(' ');
</script>

<button {type} {disabled} class={buttonClasses} on:click>
	<slot />
</button>
