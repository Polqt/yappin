const API_BASE = 'http://localhost:8080';

export async function signup(username: string, email: string, password: string) {
	const res = await fetch(`${API_BASE}/api/users/sign-up`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ username, email, password })
	});
	if (!res.ok) throw new Error(await res.text());
	return res.json();
}

export async function login(email: string, password: string) {
	const res = await fetch(`${API_BASE}/api/users/login`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ email, password })
	});
	if (!res.ok) throw new Error(await res.text());
	return res.json();
}

export function setToken(token: string) {
	localStorage.setItem('token', token);
}
