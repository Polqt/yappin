export const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';
export const WS_BASE_URL = import.meta.env.VITE_WS_URL || 'ws://localhost:8080';

export const API_ENDPOINTS = {
	auth: {
		login: '/api/users/login',
		signup: '/api/users/sign-up',
		logout: '/api/users/logout',
		me: '/api/users/me'
	},
	rooms: {
		list: '/api/websoc/get-rooms',
		create: '/api/websoc/create-room',
		join: (roomId: string) => `/api/websoc/join-room/${roomId}`
	},
	stats: {
		profile: '/api/users/profile',
		leaderboard: '/api/stats/leaderboard',
		checkin: '/api/stats/checkin',
		userStats: (userId: string) => `/api/stats/user/${userId}`
	}
} as const;

export const ROUTES = {
	home: '/',
	login: '/login',
	signup: '/signup',
	dashboard: '/dashboard',
	createRoom: '/dashboard/create-room',
	leaderboard: '/dashboard/leaderboard',
	profile: '/profile',
	room: (roomId: string) => `/room/${roomId}`
} as const;
