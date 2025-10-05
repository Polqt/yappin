import { writable } from 'svelte/store';
import { authService } from '$services/auth';
import type { AuthState } from '$types/auth';

function createAuthStore() {
	const { subscribe, set, update } = writable<AuthState>({
		user: null,
		loading: true,
		error: null
	});

	return {
		subscribe,
		async init() {
			try {
				const user = await authService.getCurrentUser();
				set({ user, loading: false, error: null });
			} catch {
				set({ user: null, loading: false, error: null });
			}
		},
		async login(credentials: { email: string; password: string }) {
			update((s) => ({ ...s, loading: true, error: null }));
			try {
				const user = await authService.login(credentials);
				set({ user, loading: false, error: null });
			} catch (err) {
				update((s) => ({
					...s,
					loading: false,
					error: err instanceof Error ? err.message : String(err)
				}));
				throw err;
			}
		},
		async signup(credentials: { username: string; email: string; password: string }) {
			update((s) => ({ ...s, loading: true, error: null }));
			try {
				const user = await authService.signup(credentials);
				set({ user, loading: false, error: null });
			} catch (err) {
				update((s) => ({
					...s,
					loading: false,
					error: err instanceof Error ? err.message : String(err)
				}));
				throw err;
			}
		},
		async logout() {
			await authService.logout();
			set({ user: null, loading: false, error: null });
		}
	};
}

export const auth = createAuthStore();
