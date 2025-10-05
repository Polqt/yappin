<script lang="ts">
  import { onMount } from 'svelte';
  import { auth } from '$stores/auth';
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';
  import '../app.css';

  let user;
  auth.subscribe(state => user = state.user);

  onMount(async () => {
    await auth.init();
    
    if (!user && !['login', 'signup'].includes($page.route.id?.split('/')[1] || '')) {
      goto('/login');
    }

    if ($auth.isAuthenticated) {
      goto('/dashboard');
    }
  });
</script>

<slot />
