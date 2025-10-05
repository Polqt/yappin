<script lang="ts">
  export let variant: 'primary' | 'secondary' | 'danger' = 'primary';
  export let type: 'button' | 'submit' = 'button';
  export let disabled = false;
</script>

<button
  {type}
  {disabled}
  class="px-4 py-2 rounded-md font-medium transition duration-200 focus:outline-none focus:ring-2 focus:ring-offset-2 {variant === 'primary' ? 'bg-blue-600 text-white hover:bg-blue-700 focus:ring-blue-500' : variant === 'secondary' ? 'bg-gray-200 text-gray-900 hover:bg-gray-300 focus:ring-gray-500' : 'bg-red-600 text-white hover:bg-red-700 focus:ring-red-500'} {disabled ? 'opacity-50 cursor-not-allowed' : ''}"
  on:click
>
  <slot />
</button>
