<script lang="ts">
	import { page } from '$app/stores';
	import { onDestroy, onMount } from 'svelte';
	import Header from '$lib/components/layout/Header.svelte';
	import MessageList from '$lib/components/chat/MessageList.svelte';
	import MessageInput from '$lib/components/chat/MessageInput.svelte';
	import TypingIndicator from '$lib/components/chat/TypingIndicator.svelte';
	import { auth } from '$stores/auth';
	import { websocket } from '$stores/websocket';
	import { roomService } from '$services/room';
	import type {
		CreateCategoryRequest,
		CreateChannelRequest,
		Message,
		MessageSearchResult,
		NotificationItem,
		RoomChannel,
		RoomDetail,
		RoomMember,
		UpdateMemberRoleRequest
	} from '$lib/types/room';

	let roomId = $page.params.id ?? '';
	let roomDetail: RoomDetail | null = null;
	let selectedChannelId = '';
	let messageInput = '';
	let searchQuery = '';
	let searchResults: MessageSearchResult[] = [];
	let loading = true;
	let error = '';
	let activeThreadParent: Message | null = null;
	let showNotifications = false;
	let creatingCategory = false;
	let creatingChannel = false;
	let categoryRequest: CreateCategoryRequest = { name: '', position: 0 };
	let channelRequest: CreateChannelRequest = {
		name: '',
		description: '',
		position: 0,
		kind: 'text'
	};

	$: liveMessages = $websocket.messages;
	$: mergedMessages = mergeMessages(roomDetail?.messages ?? [], liveMessages);
	$: selectedChannel =
		roomDetail?.categories.flatMap((category) => category.channels).find((channel) => channel.id === selectedChannelId) ??
		null;
	$: threadMessages = activeThreadParent
		? mergedMessages.filter((message) => message.parent_message_id === activeThreadParent?.id)
		: [];
	$: notifications = dedupeNotifications([
		...(roomDetail?.notifications ?? []),
		...$websocket.notifications
	]);
	$: activeTypingUsers = $websocket.typingUsers.filter((username) => username !== $auth.user?.username);

	onMount(async () => {
		await loadRoom();
	});

	onDestroy(() => {
		websocket.disconnect();
	});

	async function loadRoom() {
		try {
			loading = true;
			error = '';
			roomDetail = await roomService.getRoomDetail(roomId);
			selectedChannelId =
				roomDetail.default_channel_id ??
				roomDetail.categories[0]?.channels[0]?.id ??
				'';
			activeThreadParent = null;

			const username = $auth.user?.username || 'Guest';
			websocket.connect(roomId, username, $auth.user?.id);
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load room';
		} finally {
			loading = false;
		}
	}

	async function handleCreateCategory() {
		if (!categoryRequest.name.trim()) return;
		await roomService.createCategory(roomId, categoryRequest);
		categoryRequest = { name: '', position: 0 };
		creatingCategory = false;
		await loadRoom();
	}

	async function handleCreateChannel() {
		if (!channelRequest.name.trim()) return;
		await roomService.createChannel(roomId, {
			...channelRequest,
			category_id: channelRequest.category_id || roomDetail?.categories[0]?.id
		});
		channelRequest = { name: '', description: '', position: 0, kind: 'text' };
		creatingChannel = false;
		await loadRoom();
	}

	async function handleSearch() {
		if (!searchQuery.trim()) {
			searchResults = [];
			return;
		}
		searchResults = await roomService.searchMessages(roomId, searchQuery, selectedChannelId);
	}

	async function handleReact(messageId: string, emoji: string) {
		await roomService.addReaction(messageId, emoji);
		await loadRoom();
	}

	function handleReply(message: Message) {
		activeThreadParent = message;
	}

	function handleSendMessage() {
		if (!messageInput.trim() || !selectedChannelId) return;
		websocket.sendMessage({
			content: messageInput,
			channelId: selectedChannelId,
			parentMessageId: activeThreadParent?.id
		});
		messageInput = '';
	}

	function handleTyping(isTyping: boolean) {
		if (!selectedChannelId) return;
		websocket.sendTyping(selectedChannelId, isTyping);
	}

	async function markNotificationRead(notification: NotificationItem) {
		await roomService.markNotificationRead(notification.id);
		roomDetail = roomDetail
			? {
					...roomDetail,
					notifications: roomDetail.notifications.map((item) =>
						item.id === notification.id ? { ...item, is_read: true } : item
					)
				}
			: roomDetail;
	}

	async function updateMember(member: RoomMember, role: string, ban = false) {
		const payload: UpdateMemberRoleRequest = { role, ban };
		await roomService.updateMemberRole(roomId, member.user_id, payload);
		await loadRoom();
	}

	function selectChannel(channel: RoomChannel) {
		selectedChannelId = channel.id;
		activeThreadParent = null;
	}

	function mergeMessages(initialMessages: Message[], realtimeMessages: Message[]) {
		const map = new Map<string, Message>();
		for (const message of [...initialMessages, ...realtimeMessages]) {
			map.set(message.id || `${message.username}-${message.created_at}`, message);
		}
		return Array.from(map.values()).sort((left, right) =>
			(left.created_at || '').localeCompare(right.created_at || '')
		);
	}

	function dedupeNotifications(items: NotificationItem[]) {
		const map = new Map(items.map((item) => [item.id, item]));
		return Array.from(map.values()).sort((left, right) =>
			right.created_at.localeCompare(left.created_at)
		);
	}
</script>

<div class="min-h-screen bg-[radial-gradient(circle_at_top,_rgba(34,197,94,0.18),_transparent_32%),linear-gradient(180deg,_#0b1220_0%,_#050816_100%)]">
	<Header />

	{#if loading}
		<div class="flex min-h-[calc(100vh-4rem)] items-center justify-center">
			<div class="rounded-2xl border border-white/10 bg-white/5 px-6 py-4 text-sm text-neutral-300">
				Loading collaboration workspace...
			</div>
		</div>
	{:else if error || !roomDetail}
		<div class="mx-auto max-w-3xl p-6">
			<div class="rounded-2xl border border-red-500/20 bg-red-500/10 p-6 text-red-100">
				{error || 'Room not found'}
			</div>
		</div>
	{:else}
		<div class="mx-auto grid min-h-[calc(100vh-4rem)] max-w-[1600px] grid-cols-12 gap-4 p-4">
			<aside class="col-span-12 rounded-3xl border border-white/10 bg-white/[0.04] p-4 backdrop-blur-xl lg:col-span-3">
				<div class="mb-4">
					<p class="text-xs uppercase tracking-[0.24em] text-emerald-300">Workspace</p>
					<h1 class="mt-2 text-2xl font-semibold text-white">{roomDetail.room.name}</h1>
					<p class="mt-2 text-sm text-neutral-400">
						{roomDetail.room.topic_description || 'Channels, threads, moderation, and live presence in one space.'}
					</p>
				</div>

				<div class="mb-4 rounded-2xl border border-white/10 bg-white/5 p-3 text-sm text-neutral-300">
					<div class="flex items-center justify-between">
						<span>Your role</span>
						<span class="rounded-full bg-emerald-400/15 px-2 py-1 text-xs text-emerald-200">
							{roomDetail.current_user?.role || 'guest'}
						</span>
					</div>
					<div class="mt-2 text-xs text-neutral-500">
						{roomDetail.online_member_count} online • {roomDetail.threaded_reply_count} threaded replies
					</div>
				</div>

				<div class="space-y-4">
					{#each roomDetail.categories as category}
						<div class="rounded-2xl border border-white/10 bg-black/10 p-3">
							<div class="mb-2 flex items-center justify-between">
								<h2 class="text-xs uppercase tracking-[0.2em] text-neutral-400">{category.name}</h2>
							</div>
							<div class="space-y-2">
								{#each category.channels as channel}
									<button
										type="button"
										on:click={() => selectChannel(channel)}
										class="flex w-full items-start justify-between rounded-xl px-3 py-2 text-left transition {selectedChannelId === channel.id ? 'bg-emerald-400/15 text-white' : 'bg-white/5 text-neutral-300 hover:bg-white/10'}"
									>
										<div>
											<div class="font-medium"># {channel.name}</div>
											{#if channel.description}
												<div class="text-xs text-neutral-500">{channel.description}</div>
											{/if}
										</div>
										{#if channel.is_private}
											<span class="text-[10px] uppercase tracking-[0.2em] text-amber-300">Private</span>
										{/if}
									</button>
								{/each}
							</div>
						</div>
					{/each}
				</div>

				{#if roomDetail.current_user?.can_manage_channels}
					<div class="mt-4 space-y-3 rounded-2xl border border-white/10 bg-white/5 p-3">
						<div class="flex gap-2">
							<button class="rounded-xl bg-white/10 px-3 py-2 text-xs text-white" on:click={() => (creatingCategory = !creatingCategory)}>
								New Category
							</button>
							<button class="rounded-xl bg-white/10 px-3 py-2 text-xs text-white" on:click={() => (creatingChannel = !creatingChannel)}>
								New Channel
							</button>
						</div>

						{#if creatingCategory}
							<div class="space-y-2">
								<input bind:value={categoryRequest.name} placeholder="Category name" class="w-full rounded-xl border border-white/10 bg-black/20 px-3 py-2 text-sm text-white" />
								<button class="rounded-xl bg-emerald-400 px-3 py-2 text-xs font-medium text-neutral-950" on:click={handleCreateCategory}>
									Create Category
								</button>
							</div>
						{/if}

						{#if creatingChannel}
							<div class="space-y-2">
								<select bind:value={channelRequest.category_id} class="w-full rounded-xl border border-white/10 bg-black/20 px-3 py-2 text-sm text-white">
									<option value="">Select category</option>
									{#each roomDetail.categories as category}
										<option value={category.id}>{category.name}</option>
									{/each}
								</select>
								<input bind:value={channelRequest.name} placeholder="Channel name" class="w-full rounded-xl border border-white/10 bg-black/20 px-3 py-2 text-sm text-white" />
								<input bind:value={channelRequest.description} placeholder="Description" class="w-full rounded-xl border border-white/10 bg-black/20 px-3 py-2 text-sm text-white" />
								<button class="rounded-xl bg-cyan-300 px-3 py-2 text-xs font-medium text-neutral-950" on:click={handleCreateChannel}>
									Create Channel
								</button>
							</div>
						{/if}
					</div>
				{/if}
			</aside>

			<section class="col-span-12 flex min-h-[80vh] flex-col rounded-3xl border border-white/10 bg-white/[0.04] backdrop-blur-xl lg:col-span-6">
				<div class="border-b border-white/10 px-5 py-4">
					<div class="flex flex-col gap-4 xl:flex-row xl:items-center xl:justify-between">
						<div>
							<p class="text-xs uppercase tracking-[0.24em] text-cyan-300">Channel</p>
							<h2 class="mt-1 text-xl font-semibold text-white">
								{selectedChannel ? `# ${selectedChannel.name}` : 'Select a channel'}
							</h2>
							<p class="mt-1 text-sm text-neutral-400">
								{selectedChannel?.description || 'Thread replies, mentions, embeds, and realtime collaboration.'}
							</p>
						</div>

						<div class="flex flex-col gap-2 sm:flex-row">
							<input
								bind:value={searchQuery}
								on:keydown={(event) => event.key === 'Enter' && handleSearch()}
								placeholder="Search messages"
								class="rounded-xl border border-white/10 bg-black/20 px-4 py-2 text-sm text-white"
							/>
							<button class="rounded-xl border border-white/10 bg-white/10 px-4 py-2 text-sm text-white" on:click={handleSearch}>
								Search
							</button>
							<button class="rounded-xl border border-white/10 bg-white/10 px-4 py-2 text-sm text-white" on:click={() => (showNotifications = !showNotifications)}>
								Notifications ({notifications.filter((item) => !item.is_read).length})
							</button>
						</div>
					</div>
				</div>

				<div class="flex-1 overflow-y-auto px-5 py-4">
					{#if searchResults.length > 0}
						<div class="mb-4 rounded-2xl border border-cyan-400/20 bg-cyan-400/10 p-4">
							<h3 class="mb-2 text-sm font-medium text-cyan-100">Search Results</h3>
							<div class="space-y-2">
								{#each searchResults as result}
									<button class="block w-full rounded-xl bg-black/20 px-3 py-2 text-left text-sm text-white" on:click={() => (selectedChannelId = result.channel_id)}>
										<div class="font-medium">{result.username}</div>
										<div class="text-neutral-300">{@html result.highlighted}</div>
									</button>
								{/each}
							</div>
						</div>
					{/if}

					<MessageList
						messages={mergedMessages}
						currentUserId={$auth.user?.id}
						activeChannelId={selectedChannelId}
						searchQuery={searchQuery}
						onReply={handleReply}
						onReact={handleReact}
					/>
				</div>

				<TypingIndicator users={activeTypingUsers} />

				<MessageInput
					bind:value={messageInput}
					isConnected={$websocket.connected && !!selectedChannelId}
					placeholder={selectedChannel ? `Message #${selectedChannel.name}` : 'Select a channel to chat'}
					threadLabel={activeThreadParent ? `${activeThreadParent.username}: ${activeThreadParent.content.slice(0, 40)}` : ''}
					onSendMessage={handleSendMessage}
					onTyping={handleTyping}
				/>
			</section>

			<aside class="col-span-12 space-y-4 lg:col-span-3">
				<div class="rounded-3xl border border-white/10 bg-white/[0.04] p-4 backdrop-blur-xl">
					<h3 class="mb-3 text-sm font-semibold uppercase tracking-[0.2em] text-neutral-400">Online Presence</h3>
					<div class="space-y-2">
						{#each $websocket.onlineUsers as user}
							<div class="flex items-center justify-between rounded-xl bg-white/5 px-3 py-2 text-sm text-white">
								<span>{user.username}</span>
								<span class="h-2.5 w-2.5 rounded-full bg-emerald-400"></span>
							</div>
						{/each}
					</div>
				</div>

				<div class="rounded-3xl border border-white/10 bg-white/[0.04] p-4 backdrop-blur-xl">
					<h3 class="mb-3 text-sm font-semibold uppercase tracking-[0.2em] text-neutral-400">Thread</h3>
					{#if activeThreadParent}
						<div class="mb-3 rounded-2xl border border-white/10 bg-white/5 p-3 text-sm text-white">
							<div class="font-medium">{activeThreadParent.username}</div>
							<div class="mt-1 text-neutral-300">{activeThreadParent.content}</div>
						</div>
						<MessageList
							messages={threadMessages.map((message) => ({ ...message, channel_id: selectedChannelId }))}
							currentUserId={$auth.user?.id}
							activeChannelId={selectedChannelId}
							onReply={handleReply}
							onReact={handleReact}
						/>
					{:else}
						<p class="text-sm text-neutral-400">Select “Reply in thread” on any message to open its thread.</p>
					{/if}
				</div>

				<div class="rounded-3xl border border-white/10 bg-white/[0.04] p-4 backdrop-blur-xl">
					<h3 class="mb-3 text-sm font-semibold uppercase tracking-[0.2em] text-neutral-400">Moderation</h3>
					<div class="space-y-2">
						{#each roomDetail.members as member}
							<div class="rounded-2xl border border-white/10 bg-white/5 p-3">
								<div class="flex items-center justify-between gap-2">
									<div>
										<div class="text-sm font-medium text-white">{member.username}</div>
										<div class="text-xs uppercase tracking-[0.16em] text-neutral-500">{member.role}</div>
									</div>
									{#if roomDetail.current_user?.can_moderate}
										<div class="flex gap-2">
											<button class="rounded-lg bg-white/10 px-2 py-1 text-[11px] text-white" on:click={() => updateMember(member, 'mod')}>
												Make Mod
											</button>
											<button class="rounded-lg bg-red-500/20 px-2 py-1 text-[11px] text-red-100" on:click={() => updateMember(member, member.role, true)}>
												Ban
											</button>
										</div>
									{/if}
								</div>
							</div>
						{/each}
					</div>
				</div>

				{#if showNotifications}
					<div class="rounded-3xl border border-white/10 bg-white/[0.04] p-4 backdrop-blur-xl">
						<h3 class="mb-3 text-sm font-semibold uppercase tracking-[0.2em] text-neutral-400">Notification Center</h3>
						<div class="space-y-2">
							{#each notifications as notification}
								<button
									on:click={() => markNotificationRead(notification)}
									class="block w-full rounded-2xl border px-3 py-3 text-left transition {notification.is_read ? 'border-white/10 bg-white/5 text-neutral-400' : 'border-emerald-400/20 bg-emerald-400/10 text-white'}"
								>
									<div class="text-sm font-medium">{notification.title}</div>
									<div class="mt-1 text-xs">{notification.body}</div>
								</button>
							{/each}
						</div>
					</div>
				{/if}
			</aside>
		</div>
	{/if}
</div>
