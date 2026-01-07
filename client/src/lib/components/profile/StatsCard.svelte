<script lang="ts">
	export let label: string;
	export let value: number | string;
	export let icon: string = '📊';
</script>

<div class="rounded-lg bg-white p-6 shadow-md transition hover:shadow-lg">
	<div class="flex items-center justify-between">
		<div>
			<p class="text-sm font-medium text-gray-600">{label}</p>
			<p class="mt-2 text-3xl font-bold text-gray-900">{value}</p>
		</div>
		<div class="text-4xl">{icon}</div>
	</div>
</div>
