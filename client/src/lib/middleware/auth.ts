import { writable } from 'svelte/store';
import { browser } from '$app/environment';

export interface RouteGuard {
	requiresAuth: boolean;
	redirectTo?: string;
}

export const routeGuards = writable<Map<string, RouteGuard>>(new Map());

// Add route guard
export function addRouteGuard(path: string, guard: RouteGuard) {
	routeGuards.update((guards) => {
		guards.set(path, guard);
		return guards;
	});
}

// Check if route requires authentication
export function checkRouteAuth(pathname: string, isAuthenticated: boolean): boolean {
	// Default protected routes
	const protectedPaths = ['/dashboard', '/profile', '/room'];
	const publicPaths = ['/login', '/signup'];

	const isProtectedRoute = protectedPaths.some((path) => pathname.startsWith(path));
	const isPublicRoute = publicPaths.some((path) => pathname.startsWith(path));

	// Return whether access is allowed (component will handle redirects)
	if (isProtectedRoute && !isAuthenticated) {
		return false;
	}

	if (isPublicRoute && isAuthenticated) {
		return false;
	}

	return true;
}

// Initialize default route guards
if (browser) {
	addRouteGuard('/dashboard', { requiresAuth: true, redirectTo: '/login' });
	addRouteGuard('/profile', { requiresAuth: true, redirectTo: '/login' });
	addRouteGuard('/room', { requiresAuth: false });
	addRouteGuard('/login', { requiresAuth: false, redirectTo: '/dashboard' });
	addRouteGuard('/signup', { requiresAuth: false, redirectTo: '/dashboard' });
}
