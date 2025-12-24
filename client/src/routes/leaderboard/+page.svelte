<script lang="ts">
    import { onMount } from 'svelte';
    import { API_BASE_URL } from '$lib/constants/api';
	import type { LeaderboardEntry } from '$lib/types/leaderboard';
	import { getMedal } from '$lib/utils/leaderboard';

    
    // Component state
    let leaderboard: LeaderboardEntry[] = [];
    let loading = true;
    let error = '';
    
    // Fetch leaderboard when component loads
    onMount(async () => {
        try {
            const response = await fetch(`${API_BASE_URL}/api/stats/leaderboard?limit=100`, {
                credentials: 'include'  // Send cookies
            });
            
            if (!response.ok) throw new Error('Failed to fetch');
            
            leaderboard = await response.json();
        } catch (err) {
            error = 'Failed to load leaderboard';
            console.error(err);
        } finally {
            loading = false;
        }
    });

</script>

<div class="max-w-4xl mx-auto p-8">
    <h1 class="text-4xl font-bold mb-8 text-center">🏆 Leaderboard</h1>
    
    {#if loading}
        <div class="text-center py-12">
            <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500 mx-auto"></div>
            <p class="mt-4 text-gray-600">Loading...</p>
        </div>
    {:else if error}
        <div class="bg-red-50 border border-red-200 rounded-lg p-4 text-red-700">
            {error}
        </div>
    {:else}
        <div class="bg-white rounded-lg shadow-lg overflow-hidden">
            <table class="w-full">
                <thead class="bg-gray-50 border-b">
                    <tr>
                        <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Rank</th>
                        <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">User</th>
                        <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">Messages</th>
                        <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">Upvotes</th>
                        <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">Streak</th>
                    </tr>
                </thead>
                <tbody class="divide-y divide-gray-200">
                    {#each leaderboard as user, index}
                        <tr class:bg-yellow-50={user.rank <= 3} class="hover:bg-gray-50 transition-colors">
                            <td class="px-6 py-4 whitespace-nowrap">
                                <span class="text-2xl">{getMedal(user.rank)}</span>
                                <span class="ml-2 font-semibold">#{user.rank}</span>
                            </td>
                            <td class="px-6 py-4 whitespace-nowrap font-medium text-gray-900">
                                {user.username}
                            </td>
                            <td class="px-6 py-4 whitespace-nowrap text-right text-gray-600">
                                {user.total_messaes}
                            </td>
                            <td class="px-6 py-4 whitespace-nowrap text-right text-gray-600">
                                {user.total_upvotes} 👍
                            </td>
                            <td class="px-6 py-4 whitespace-nowrap text-right text-gray-600">
                                {user.daily_streak} 🔥
                            </td>
                        </tr>
                    {/each}
                </tbody>
            </table>
        </div>
    {/if}
</div>
