export interface User {
	id: string;
	username: string;
	email?: string;
	AccessToken?: string; 
}

export interface LoginCredentials {
	email: string;
	password: string;
}

export interface SignupCredentials extends LoginCredentials {
	username: string;
}

export interface AuthResponse {
	user: User;
	token: string;
}

export interface AuthState {
	user: User | null;
	loading: boolean;
	error: string | null;
}
