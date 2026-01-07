<script lang="ts">
	import type { Achievement } from '$lib/types/user';

	export let achievement: Achievement;
</script>

<div
	class="flex flex-col items-center rounded-lg bg-gradient-to-br from-yellow-50 to-yellow-100 p-4 text-center shadow transition hover:shadow-md"
>
	<div class="mb-2 text-4xl">{achievement.icon}</div>
	<h4 class="font-semibold text-gray-800">{achievement.name}</h4>
	<p class="mt-1 text-xs text-gray-600">{achievement.description}</p>
	{#if achievement.achieved_at}
		<p class="mt-2 text-xs text-gray-500">
			{new Date(achievement.achieved_at).toLocaleDateString()}
		</p>
	{/if}
</div>
