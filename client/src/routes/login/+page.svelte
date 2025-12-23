<script lang="ts">
  import { auth } from '$stores/auth';
  import { goto } from '$app/navigation';
  import Button from '$lib/components/common/Button.svelte';
  import Input from '$lib/components/common/Input.svelte';
  import { MessageCircle } from 'lucide-svelte';

  let email = '';
  let password = '';
  let loading = false;
  let error = '';

  async function handleSubmit() {
    if (!email || !password) {
      error = 'Please fill in all fields';
      return;
    }
    try {
      loading = true;
      error = '';
      await auth.login({ email, password });
      goto('/dashboard');
    } catch (err) {
      error = 'Invalid email or password';
    } finally {
      loading = false;
    }
  }
</script>

<div class="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
  <div class="max-w-md w-full bg-white rounded-lg shadow-xl p-8">
    <div class="text-center mb-8">
      <MessageCircle class="mx-auto h-12 w-12 text-blue-600" />
      <h2 class="mt-4 text-3xl font-bold text-gray-900">
        Welcome Back
      </h2>
      <p class="mt-2 text-sm text-gray-600">
        Sign in to your Yappin account
      </p>
    </div>
    
    <form class="space-y-6" on:submit|preventDefault={handleSubmit}>
      {#if error}
        <div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-md">
          <p class="text-sm">{error}</p>
        </div>
      {/if}

      <Input
        type="email"
        bind:value={email}
        label="Email Address"
        placeholder="Enter your email"
        required
        {error}
      />

      <Input
        type="password"
        bind:value={password}
        label="Password"
        placeholder="Enter your password"
        required
      />

      <Button
        type="submit"
        variant="primary"
        disabled={loading}
      >
        {loading ? 'Signing in...' : 'Sign In'}
      </Button>
    </form>

    <div class="mt-6 text-center">
      <p class="text-sm text-gray-600">
        Don't have an account?
        <a href="/signup" class="font-medium text-blue-600 hover:text-blue-500 transition duration-200">
          Sign up
        </a>
      </p>
    </div>
  </div>
</div>
