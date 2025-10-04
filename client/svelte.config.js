import adapter from '@sveltejs/adapter-auto';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

/** @type {import('@sveltejs/kit').Config} */
const config = {
  kit: {
    adapter: adapter(),
    alias: {
      $lib: 'src/lib',
      $components: 'src/lib/components',
      $stores: 'src/stores',
      $services: 'src/services',
      $types: 'src/lib/types',
      $utils: 'src/lib/utils'
    }
  },
  preprocess: vitePreprocess()
};

export default config;