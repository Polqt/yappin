<script lang="ts">
  import { auth } from '$lib/stores/auth';
  import { goto } from '$app/navigation';
  import Button from '$lib/components/common/Button.svelte';
  import Input from '$lib/components/common/Input.svelte';

  let email = '';
  let password = '';
  let loading = false;
  let error = '';

  async function handleSubmit() {
    try {
      loading = true;
      await auth.login(email, password);
      goto('/');
    } catch (err) {
      error = 'Invalid credentials';
    } finally {
      loading = false;
    }
  }
</script>

<div class="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
  <div class="max-w-md w-full space-y-8">
    <div>
      <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">
        Sign in to your account
      </h2>
    </div>
    
    <form class="mt-8 space-y-6" on:submit|preventDefault={handleSubmit}>
      {#if error}
        <div class="rounded-md bg-red-50 p-4">
          <p class="text-sm text-red-700">{error}</p>
        </div>
      {/if}

      <div class="rounded-md shadow-sm -space-y-px">
        <Input
          type="email"
          required
          bind:value={email}
          label="Email address"
          placeholder="Email address"
        />

        <Input
          type="password"
          required
          bind:value={password}
          label="Password"
          placeholder="Password"
        />
      </div>

      <Button
        type="submit"
        variant="primary"
        disabled={loading}
        class="w-full"
      >
        {loading ? 'Signing in...' : 'Sign in'}
      </Button>
    </form>
  </div>
</div>