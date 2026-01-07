<script lang="ts">
	export let label: string;
	export let value: number | string;
	export let icon: string = '📊';
</script>

<div
	class="rounded-xl border border-white/10 bg-white/5 p-6 backdrop-blur-sm transition hover:bg-white/10"
>
	<div class="flex items-center justify-between">
		<div>
			<p class="text-sm font-medium text-neutral-400">{label}</p>
			<p class="mt-2 text-3xl font-light text-white">{value}</p>
		</div>
		<div class="text-4xl">{icon}</div>
	</div>
</div>
