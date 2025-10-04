import type { User } from '$lib/types/user';
import axios from 'axios';

const BASE_URL = import.meta.env.VITE_API_URL;

export class AuthService {
	private static instance: AuthService;

	private constructor() {}

	static getInstance(): AuthService {
		if (!AuthService.instance) {
			AuthService.instance = new AuthService();
		}
		return AuthService.instance;
	}

	async login(email: string, password: string): Promise<User> {
		const { data } = await axios.post(
			`${BASE_URL}/auth/login`,
			{
				email,
				password
			},
			{ withCredentials: true }
		);
		return data;
	}

	async logout(): Promise<void> {
		await axios.post(`${BASE_URL}/auth/logout`, {}, { withCredentials: true });
	}

	async getCurrentUser(): Promise<User | null> {
		try {
			const { data } = await axios.get(`${BASE_URL}/auth/me`, { withCredentials: true });
			return data;
		} catch {
			return null;
		}
	}
}

export const authService = AuthService.getInstance();
