<script lang="ts">
	import { onMount } from 'svelte';
	import { auth } from '$stores/auth';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import '../app.css';

	const publicRoutes = ['/login', '/signup'];

	onMount(async () => {
		await auth.init();
	});

	$: {
		if (!$auth.loading) {
			const isPublicRoute = publicRoutes.some((route) => $page.url.pathname.startsWith(route));
			const hasUser = $auth.user !== null;

			if (!hasUser && !isPublicRoute) {
				goto('/login');
			} else if (hasUser && isPublicRoute) {
				goto('/dashboard');
			}
		}
	}
</script>

<slot />
