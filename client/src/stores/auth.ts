import { writable, derived } from 'svelte/store';
import type { User } from '$lib/types/user';

interface AuthState {
	user: User | null;
	loading: boolean;
	error: string | null;
	isAuthenticated: boolean;
}

function createAuthStore() {
	const { subscribe, set, update } = writable<AuthState>({
		user: null,
		loading: false,
		error: null,
		isAuthenticated: false
	});

	return {
		subscribe,
		login: async (email: string, password: string) => {
			update((state) => ({ ...state, loading: true, error: null }));
			try {
				const response = await fetch('/api/auth/login', {
					method: 'POST',
					credentials: 'include',
					headers: { 'Content-Type': 'application/json' },
					body: JSON.stringify({ email, password })
				});

				if (!response.ok) throw new Error('Login failed');

				const user = await response.json();
				update((state) => ({
					...state,
					user,
					loading: false,
					isAuthenticated: true
				}));
				return user;
			} catch (err) {
				update((state) => ({
					...state,
					loading: false,
					error: 'Invalid credentials',
					isAuthenticated: false
				}));
				throw err;
			}
		},
		logout: async () => {
			try {
				await fetch('/api/auth/logout', {
					method: 'POST',
					credentials: 'include'
				});
			} finally {
				set({
					user: null,
					loading: false,
					error: null,
					isAuthenticated: false
				});
			}
		},
		checkAuth: async () => {
			update((state) => ({ ...state, loading: true }));
			try {
				const response = await fetch('/api/auth/me', {
					credentials: 'include'
				});

				if (!response.ok) throw new Error();

				const user = await response.json();
				update(() => ({
					user,
					loading: false,
					error: null,
					isAuthenticated: true
				}));
			} catch {
				set({
					user: null,
					loading: false,
					error: null,
					isAuthenticated: false
				});
			}
		}
	};
}

export const auth = createAuthStore();
export const isAuthenticated = derived(auth, ($auth) => $auth.isAuthenticated);
