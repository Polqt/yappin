<script lang="ts">
	import type { ActivityData } from '$lib/types/user';

	export let data: ActivityData[] = [];

	// Calculate max value for scaling
	const maxValue = Math.max(...data.map((d) => d.count), 1);
</script>

<div class="space-y-4">
	{#if data.length === 0}
		<div class="py-8 text-center">
			<p class="text-sm text-neutral-400">No activity data available</p>
		</div>
	{:else}
		<div class="space-y-2">
			{#each data as item}
				<div class="flex items-center gap-4">
					<span class="w-24 text-sm text-neutral-400">
						{new Date(item.date).toLocaleDateString()}
					</span>
					<div class="flex-1">
						<div class="h-6 rounded bg-white/10">
							<div
								class="h-full rounded bg-white/50 transition-all"
								style="width: {(item.count / maxValue) * 100}%"
							></div>
						</div>
					</div>
					<span class="w-12 text-right text-sm font-medium text-white">
						{item.count}
					</span>
				</div>
			{/each}
		</div>
	{/if}
</div>
