<script lang="ts">
	import { goto } from "$app/navigation";
	import { auth } from "$stores/auth";
	import { onMount } from "svelte";


    let userStats = {
        totalMessaes: 0,
        roomsJoined: 0,
        checkInStreak: 0,
        upvotesReceived: 0,
        level: 1,
        xp: 0
    }

    let achievements: any[] = [];
    let loading = true;

    onMount(async () => { 
        if (!$auth.user) {
            goto("/login");
            return;
        }

        try {
            const response = await fetch(``);
        } catch (error) {
            
        }
    })
</script>