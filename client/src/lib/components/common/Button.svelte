<script lang="ts">
	export let variant: 'primary' | 'secondary' | 'danger' = 'primary';
	export let type: 'button' | 'submit' = 'button';
	export let disabled = false;
</script>

<button {type} {disabled} class="btn btn-{variant}" on:click>
	<slot />
</button>
