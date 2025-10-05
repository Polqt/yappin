<script lang="ts">
  import { Eye, EyeOff } from 'lucide-svelte';

  export let type: string = 'text';
  export let value: string = '';
  export let label: string = '';
  export let placeholder: string = '';
  export let required: boolean = false;
  export let disabled: boolean = false;
  export let error: string = '';

  let showPassword = false;
</script>

<div class="mb-4">
  {#if label}
    <label class="block text-sm font-medium text-gray-700 mb-1">
      {label}
    </label>
  {/if}
  <div class="relative">
    <input
      {type}
      bind:value
      {placeholder}
      {required}
      {disabled}
      class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 {error ? 'border-red-500' : ''}"
    />
    {#if type === 'password'}
      <button
        type="button"
        class="absolute inset-y-0 right-0 pr-3 flex items-center"
        on:click={() => showPassword = !showPassword}
      >
        {#if showPassword}
          <EyeOff class="h-5 w-5 text-gray-400" />
        {:else}
          <Eye class="h-5 w-5 text-gray-400" />
        {/if}
      </button>
    {/if}
  </div>
  {#if error}
    <p class="mt-1 text-sm text-red-600">{error}</p>
  {/if}
</div>
