<script lang="ts">
	/**
	 * Reusable alert component for displaying error, success, warning, and info messages.
	 */
	export let variant: 'error' | 'success' | 'warning' | 'info' = 'error';
	export let message: string = '';

	const variantClasses = {
		error: 'border-red-500/20 bg-red-500/10 text-red-200',
		success: 'border-green-500/20 bg-green-500/10 text-green-200',
		warning: 'border-yellow-500/20 bg-yellow-500/10 text-yellow-200',
		info: 'border-blue-500/20 bg-blue-500/10 text-blue-200'
	};
</script>

{#if message}
	<div class="rounded-lg border px-4 py-3 {variantClasses[variant]}">
		<p class="text-sm">{message}</p>
	</div>
{/if}
