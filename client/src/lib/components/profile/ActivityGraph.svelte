<script lang="ts">
	import type { ActivityData } from '$lib/types/user';

	export let data: ActivityData[] = [];

	// Calculate max value for scaling
	const maxValue = Math.max(...data.map((d) => d.count), 1);
</script>

<div class="space-y-4">
	{#if data.length === 0}
		<div class="py-8 text-center text-gray-500">
			<p>No activity data available</p>
		</div>
	{:else}
		<div class="space-y-2">
			{#each data as item}
				<div class="flex items-center gap-4">
					<span class="w-24 text-sm text-gray-600">
						{new Date(item.date).toLocaleDateString()}
					</span>
					<div class="flex-1">
						<div class="h-6 rounded bg-gray-200">
							<div
								class="h-full rounded bg-blue-500 transition-all"
								style="width: {(item.count / maxValue) * 100}%"
							></div>
						</div>
					</div>
					<span class="w-12 text-right text-sm font-semibold text-gray-700">
						{item.count}
					</span>
				</div>
			{/each}
		</div>
	{/if}
</div>
