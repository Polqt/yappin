import type { LoginCredentials, SignupCredentials, User } from '$types/auth'; // Changed imports

const BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

export class AuthService {
	async login(credentials: LoginCredentials): Promise<User> {
		const response = await fetch(`${BASE_URL}/api/users/login`, {
			method: 'POST',
			credentials: 'include',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify(credentials)
		});

		if (!response.ok) {
			const errorText = await response.text();
			throw new Error(errorText || 'Login failed');
		}
		return response.json();
	}

	async signup(credentials: SignupCredentials): Promise<User> {
		const response = await fetch(`${BASE_URL}/api/users/sign-up`, {
			method: 'POST',
		credentials: 'include',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify(credentials)
		});

		if (!response.ok) {
			const errorText = await response.text();
			throw new Error(errorText || 'Signup failed');
		}
		return response.json();
	}

	async logout(): Promise<void> {
		await fetch(`${BASE_URL}/api/users/logout`, {
			method: 'POST',
			credentials: 'include'
		});
	}

	async getCurrentUser(): Promise<User | null> {
		try {
			const response = await fetch(`${BASE_URL}/api/users/me`, {
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
