export function validateEmail(email: string): boolean {
	const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
	return emailRegex.test(email);
}

export function validatePassword(password: string): { valid: boolean; error?: string } {
	if (password.length < 8) {
		return { valid: false, error: 'Password must be at least 8 characters' };
	}
	if (!/[A-Z]/.test(password)) {
		return { valid: false, error: 'Password must contain at least one uppercase letter' };
	}
	if (!/[a-z]/.test(password)) {
		return { valid: false, error: 'Password must contain at least one lowercase letter' };
	}
	if (!/\d/.test(password)) {
		return { valid: false, error: 'Password must contain at least one number' };
	}
	if (!/[!@#$%^&*(),.?":{}|<>]/.test(password)) {
		return { valid: false, error: 'Password must contain at least one special character' };
	}
	return { valid: true };
}

export function validateUsername(username: string): { valid: boolean; error?: string } {
	if (username.length < 3) {
		return { valid: false, error: 'Username must be at least 3 characters' };
	}
	if (username.length > 20) {
		return { valid: false, error: 'Username must be less than 20 characters' };
	}
	if (!/^[a-zA-Z0-9_]+$/.test(username)) {
		return { valid: false, error: 'Username can only contain letters, numbers, and underscores' };
	}
	return { valid: true };
}
