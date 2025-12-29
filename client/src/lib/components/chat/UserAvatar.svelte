<script lang="ts">
	import { stringToColor } from "$lib/utils/color";

    export let username: string;
    export let size: 'sm' | 'md' | 'lg' = 'md';

    const initial = username.charAt(0).toUpperCase();
    const bgColor = stringToColor(username)

    const sizeClass = {
        sm: 'w-8 h-8 text-sm',
        md: 'w-10 h-10 text-base',
        lg: 'w-12 h-12 text-lg'
    }
</script>

<div 
    class="flex items-center justify-center rounded-full font-bold text-white {sizeClass[size]}"
    style="background-color: {bgColor}"
    title={username}
>
    {initial}
</div>