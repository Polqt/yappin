<script lang="ts">
	import type { Achievement } from '$lib/types/user';

	export let achievement: Achievement;
</script>

<div
	class="flex flex-col items-center rounded-xl border border-white/10 bg-white/5 p-4 text-center backdrop-blur-sm transition hover:border-white/20 hover:bg-white/10"
>
	<div class="mb-2 text-4xl">{achievement.icon}</div>
	<h4 class="font-medium text-white">{achievement.name}</h4>
	<p class="mt-1 text-xs text-neutral-400">{achievement.description}</p>
	{#if achievement.achieved_at}
		<p class="mt-2 text-xs text-neutral-500">
			{new Date(achievement.achieved_at).toLocaleDateString()}
		</p>
	{/if}
</div>
