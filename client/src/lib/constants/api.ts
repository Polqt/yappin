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
		detail: (roomId: string) => `/api/websoc/rooms/${roomId}`,
		search: (roomId: string) => `/api/websoc/rooms/${roomId}/search`,
		categories: (roomId: string) => `/api/websoc/rooms/${roomId}/categories`,
		channels: (roomId: string) => `/api/websoc/rooms/${roomId}/channels`,
		members: (roomId: string, userId: string) => `/api/websoc/rooms/${roomId}/members/${userId}`,
		join: (roomId: string) => `/api/websoc/join-room/${roomId}`,
		notifications: '/api/websoc/notifications',
		markNotificationRead: (notificationId: string) => `/api/websoc/notifications/${notificationId}/read`,
		reactions: '/api/websoc/reactions'
	},
	stats: {
		profile: (userId: string) => `/api/stats/profile/${userId}`,
		leaderboard: '/api/stats/leaderboard',
		checkin: '/api/stats/checkin',
		upvote: '/api/stats/upvote'
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
