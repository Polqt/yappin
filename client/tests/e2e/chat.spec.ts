import { expect, test } from '@playwright/test';

const fakeUser = { id: 'user-1', username: 'alice', email: 'alice@example.com' };
const fakeRooms = [
	{
		id: 'room-1',
		name: 'General Hub',
		is_pinned: false,
		created_at: new Date().toISOString(),
		expires_at: new Date(Date.now() + 86_400_000).toISOString(),
		participants: 2
	}
];

const fakeRoomDetail = {
	room: fakeRooms[0],
	categories: [
		{
			id: 'cat-1',
			name: 'General',
			position: 0,
			channels: [
				{
					id: 'chan-1',
					name: 'lobby',
					description: 'Default channel',
					kind: 'text',
					position: 0,
					is_private: false
				}
			]
		}
	],
	members: [{ user_id: 'user-1', username: 'alice', role: 'owner', created_at: new Date().toISOString() }],
	current_user: {
		role: 'owner',
		can_manage_room: true,
		can_manage_channels: true,
		can_moderate: true,
		can_post: true,
		is_muted: false,
		is_banned: false
	},
	messages: [
		{
			id: 'msg-1',
			content: 'Welcome to the lobby',
			room_id: 'room-1',
			channel_id: 'chan-1',
			username: 'system',
			system: true,
			created_at: new Date().toISOString()
		}
	],
	notifications: [],
	default_channel_id: 'chan-1',
	notification_count: 0,
	online_member_count: 1,
	threaded_reply_count: 0
};

test.beforeEach(async ({ page }) => {
	await page.addInitScript(() => {
		class MockWebSocket {
			static instances: MockWebSocket[] = [];
			url: string;
			readyState = 1;
			onopen: ((event: Event) => void) | null = null;
			onmessage: ((event: MessageEvent) => void) | null = null;
			onclose: ((event: CloseEvent) => void) | null = null;
			onerror: ((event: Event) => void) | null = null;

			constructor(url: string) {
				this.url = url;
				MockWebSocket.instances.push(this);
				setTimeout(() => {
					this.onopen?.(new Event('open'));
					this.onmessage?.(
						new MessageEvent('message', {
							data: JSON.stringify({
								type: 'history',
								messages: [
									{
										id: 'msg-1',
										content: 'Welcome to the lobby',
										room_id: 'room-1',
										channel_id: 'chan-1',
										username: 'system',
										system: true,
										created_at: new Date().toISOString()
									}
								]
							})
						})
					);
				}, 10);
			}

			send(raw: string) {
				const payload = JSON.parse(raw);
				if (payload.type === 'message.send') {
					setTimeout(() => emitRealtimeMessage(payload.content, payload.channel_id), 20);
				}
			}

			close() {
				this.readyState = 3;
				this.onclose?.(new CloseEvent('close', { wasClean: true, code: 1000 }));
			}
		}

		function emitRealtimeMessage(content: string, channelId: string) {
			const socket = MockWebSocket.instances[0];
			socket?.onmessage?.(
				new MessageEvent('message', {
					data: JSON.stringify({
						type: 'message.created',
						message: {
							id: `msg-${Date.now()}`,
							content,
							room_id: 'room-1',
							channel_id: channelId,
							username: 'alice',
							user_id: 'user-1',
							system: false,
							created_at: new Date().toISOString()
						}
					})
				})
			);
		}

		// @ts-expect-error test shim
		window.WebSocket = MockWebSocket;
		// @ts-expect-error test shim
		window.__mockRealtime = { emitRealtimeMessage };
	});

	await page.route('**/api/users/me', async (route) => {
		await route.fulfill({ json: fakeUser });
	});

	await page.route('**/api/websoc/get-rooms', async (route) => {
		await route.fulfill({ json: fakeRooms });
	});

	await page.route('**/api/websoc/create-room', async (route) => {
		await route.fulfill({
			json: {
				id: 'room-2',
				name: 'Created Room'
			}
		});
	});

	await page.route('**/api/websoc/rooms/room-1', async (route) => {
		await route.fulfill({ json: fakeRoomDetail });
	});

	await page.route('**/api/websoc/rooms/room-1/search**', async (route) => {
		await route.fulfill({
			json: [
				{
					id: 'msg-1',
					room_id: 'room-1',
					channel_id: 'chan-1',
					username: 'system',
					content: 'Welcome to the lobby',
					highlighted: 'Welcome to the <mark>lobby</mark>',
					created_at: new Date().toISOString()
				}
			]
		});
	});

	await page.route('**/api/websoc/notifications', async (route) => {
		await route.fulfill({ json: [] });
	});

	await page.route('**/api/websoc/reactions', async (route) => {
		await route.fulfill({ json: { ok: true } });
	});
});

test('login page loads', async ({ page }) => {
	await page.goto('/login');
	await expect(page.getByRole('button', { name: /sign in/i })).toBeVisible();
});

test('room creation page renders', async ({ page }) => {
	await page.goto('/dashboard/create-room');
	await expect(page.getByText(/create a new room/i)).toBeVisible();
});

test('live chat messaging renders in the room workspace', async ({ page }) => {
	await page.goto('/room/room-1');
	await expect(page.getByText(/workspace/i)).toBeVisible();
	await expect(page.getByText('Welcome to the lobby')).toBeVisible();
	const composer = page.getByPlaceholder(/message #lobby/i);
	await expect(composer).toBeEnabled();
	await composer.fill('hello realtime world');
	await page.getByRole('button', { name: /send/i }).click();
	await page.evaluate(() => {
		// @ts-expect-error test shim
		window.__mockRealtime.emitRealtimeMessage('hello realtime world', 'chan-1');
	});
	await expect(page.getByText('hello realtime world')).toBeVisible();
});
