<script lang="ts">
	import { onMount } from 'svelte';
	import { auth } from '$stores/auth';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { checkRouteAuth } from '$lib/middleware/auth';
	import '../app.css';

	const publicRoutes = ['/login', '/signup', '/'];

	onMount(async () => {
		await auth.init();
	});

	// Watch for auth state changes and enforce route guards
	$: {
		if (!$auth.loading) {
			const currentPath = $page.url.pathname;
			const isPublicRoute = publicRoutes.includes(currentPath);
			const hasUser = $auth.user !== null;

			// Use middleware to check route access
			checkRouteAuth(currentPath, hasUser);

			// Fallback route protection (skip root page)
			if (currentPath !== '/') {
				if (!hasUser && !isPublicRoute) {
					goto('/login');
				} else if (hasUser && (currentPath === '/login' || currentPath === '/signup')) {
					goto('/dashboard');
				}
			}
		}
	}
</script>

<slot />
