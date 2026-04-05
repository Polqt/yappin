<script lang="ts">
	import DOMPurify from 'dompurify';
	import { marked } from 'marked';
	import type { Message } from '$lib/types/room';
	import UserAvatar from './UserAvatar.svelte';
	import MessageReaction from './MessageReaction.svelte';

	export let messages: Message[] = [];
	export let currentUserId: string | undefined;
	export let activeChannelId = '';
	export let onReply: (message: Message) => void = () => {};
	export let onReact: (messageId: string, emoji: string) => void = () => {};
	export let searchQuery = '';

	$: visibleMessages = messages.filter((message) =>
		activeChannelId ? message.channel_id === activeChannelId : true
	);

	function renderContent(message: Message) {
		const raw = message.content;
		const highlighted =
			searchQuery.trim().length > 0
				? raw.replaceAll(searchQuery, `<mark>${searchQuery}</mark>`)
				: raw;
		return DOMPurify.sanitize(marked.parse(highlighted, { breaks: true }) as string);
	}

	function mediaUrl(content: string) {
		const match = content.match(/https?:\/\/\S+/);
		if (!match) return '';
		return match[0];
	}
</script>

<div class="space-y-3">
	{#each visibleMessages as message, index (`${message.id || index}`)}
		<div
			class="rounded-xl border border-white/10 p-4 transition hover:border-white/20 {message.user_id === currentUserId ? 'bg-white/10' : 'bg-white/5'}"
		>
			<div class="flex gap-3">
				<UserAvatar username={message.username} size="md" />

				<div class="min-w-0 flex-1">
					<div class="mb-2 flex flex-wrap items-center gap-2">
						<span class="text-sm font-medium text-white">{message.username}</span>
						{#if message.parent_message_id}
							<span class="rounded-full border border-white/10 bg-white/5 px-2 py-0.5 text-[11px] text-neutral-400">
								Thread reply
							</span>
						{/if}
						{#if message.system}
							<span class="rounded-full border border-white/10 bg-white/5 px-2 py-0.5 text-[11px] text-neutral-400">
								System
							</span>
						{/if}
						{#if message.created_at}
							<span class="text-xs text-neutral-500">
								{new Date(message.created_at).toLocaleTimeString()}
							</span>
						{/if}
					</div>

					<div class="prose prose-invert max-w-none break-words text-sm">
						{@html renderContent(message)}
					</div>

					{#if mediaUrl(message.content)}
						<div class="mt-3 overflow-hidden rounded-xl border border-white/10 bg-black/20">
							{#if /\.(png|jpe?g|gif|webp|svg)$/i.test(mediaUrl(message.content))}
								<img src={mediaUrl(message.content)} alt="Embedded media" class="max-h-80 w-full object-cover" />
							{:else if /\.(mp4|webm|ogg)$/i.test(mediaUrl(message.content))}
								<video src={mediaUrl(message.content)} controls class="max-h-80 w-full">
									<track kind="captions" />
								</video>
							{:else}
								<a
									href={mediaUrl(message.content)}
									target="_blank"
									rel="noreferrer"
									class="block px-4 py-3 text-sm text-cyan-300 hover:text-cyan-200"
								>
									Open embedded link
								</a>
							{/if}
						</div>
					{/if}

					<div class="mt-3 flex flex-wrap items-center gap-2">
						<button
							on:click={() => onReply(message)}
							class="rounded-full border border-white/10 bg-white/5 px-3 py-1 text-xs text-neutral-300 transition hover:bg-white/10"
						>
							Reply in thread
						</button>
					</div>

					<MessageReaction
						reactions={message.reactions ?? []}
						onAddReaction={(emoji) => onReact(message.id, emoji)}
					/>
				</div>
			</div>
		</div>
	{/each}

	{#if visibleMessages.length === 0}
		<div class="rounded-xl border border-dashed border-white/10 bg-white/[0.03] p-10 text-center">
			<p class="text-lg font-light text-white">No messages yet</p>
			<p class="mt-2 text-sm text-neutral-400">Start the conversation in this channel.</p>
		</div>
	{/if}
</div>
