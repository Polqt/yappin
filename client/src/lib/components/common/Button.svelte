<script lang="ts">
	export let variant: 'primary' | 'secondary' | 'danger' = 'primary';
	export let type: 'button' | 'submit' = 'button';
	export let disabled = false;
	export let fullWidth = false;

	const baseClasses =
		'px-4 py-2 rounded-md font-medium transition duration-200 focus:outline-none focus:ring-2 focus:ring-offset-2';

	const variantClasses = {
		primary: 'bg-blue-600 text-white hover:bg-blue-700 focus:ring-blue-500',
		secondary: 'bg-gray-200 text-gray-900 hover:bg-gray-300 focus:ring-gray-500',
		danger: 'bg-red-600 text-white hover:bg-red-700 focus:ring-red-500'
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
