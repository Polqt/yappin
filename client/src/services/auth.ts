import type { LoginCredentials, SignupCredentials, User } from '$types/auth';
import { API_BASE_URL, API_ENDPOINTS } from '$lib/constants/api';
import { handleApiError } from '$lib/utils/error';

export class AuthService {
	async login(credentials: LoginCredentials): Promise<User> {
		const response = await fetch(`${API_BASE_URL}${API_ENDPOINTS.auth.login}`, {
			method: 'POST',
			credentials: 'include',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify(credentials)
		});

		if (!response.ok) {
			await handleApiError(response);
		}
		return response.json();
	}

	async signup(credentials: SignupCredentials): Promise<User> {
		const response = await fetch(`${API_BASE_URL}${API_ENDPOINTS.auth.signup}`, {
			method: 'POST',
			credentials: 'include',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify(credentials)
		});

		if (!response.ok) {
			await handleApiError(response);
		}
		return response.json();
	}

	async logout(): Promise<void> {
		await fetch(`${API_BASE_URL}${API_ENDPOINTS.auth.logout}`, {
			method: 'POST',
			credentials: 'include'
		});
	}

	async getCurrentUser(): Promise<User | null> {
		try {
			const response = await fetch(`${API_BASE_URL}${API_ENDPOINTS.auth.me}`, {
				credentials: 'include'
			});

			if (!response.ok) return null;
			return response.json();
		} catch {
			return null;
		}
	}
}

export const authService = new AuthService();
