<script lang="ts">
    import type { MessageReaction } from "$lib/types/room";
    
    export let reactions: MessageReaction[] = [];
    let messageId: string;
    export let onAddReaction: (emoji: string) => void;

    const emojiOptions = ['👍', '❤️', '😂', '😮', '😢', '👏'];

    $: reactionCounts = reactions.reduce((acc, reaction) => {
        acc[reaction.emoji] = (acc[reaction.emoji] || 0) + 1;
        return acc;
    }, {} as Record<string, number>);

    let showPicker = false;

    function handleEmojiClick(emoji: string) {
        onAddReaction(emoji);
        showPicker = false;
    }
</script>

<div class="flex items-center gap-2 mt-2">
    {#each Object.entries(reactionCounts) as [emoji, count]}
        <button
            on:click={() => handleEmojiClick(emoji)}
            aria-label="React with {emoji}"
            class="flex items-center gap-1 px-2 py-1 bg-gray-100 rounded-full hover:bg-gray-200 transition-colors text-sm"
        >
            <span>{emoji}</span>
            <span class="text-gray-600">{count}</span>
        </button>
    {/each}

    <div class="relative">
        <button
            on:click={() => showPicker = !showPicker}
            class="px-2 py-1 text-gray-400 hover:text-gray-600 transition-colors"
            title="Add Reaction"
        >
            +
        </button>

        {#if showPicker}
            <div class="absolute bottom-full mb-2 p-2 bg-white rounded-lg shadow-lg border flex gap-1 z-10">
                {#each emojiOptions as emoji}
                    <button
                        on:click={() => handleEmojiClick(emoji)}
                        class="text-2xl hover:scale-125 transtion-transform"
                    >
                        {emoji}
                    </button>
                {/each}
            </div>
            
        {/if}
    </div>
</div>