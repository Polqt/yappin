<script lang="ts">
  import { auth } from '$stores/auth';
  import { goto } from '$app/navigation';
  import Button from '$lib/components/common/Button.svelte';
  import Input from '$lib/components/common/Input.svelte';
  import { MessageCircle } from 'lucide-svelte';

  let username = '';
  let email = '';
  let password = '';
  let loading = false;
  let error = '';

  async function handleSubmit() {
    if (!username || !email || !password) {
      error = 'Please fill in all fields';
      return;
    }
    if (password.length < 6) {
      error = 'Password must be at least 6 characters';
      return;
    }
    try {
      loading = true;
      error = '';
      await auth.signup({ username, email, password });
      goto('/');
    } catch (err) {
      error = 'Signup failed. Email may already be in use.';
    } finally {
      loading = false;
    }
  }
</script>

<div class="min-h-screen bg-gradient-to-br from-green-50 to-emerald-100 flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
  <div class="max-w-md w-full bg-white rounded-lg shadow-xl p-8">
    <div class="text-center mb-8">
      <MessageCircle class="mx-auto h-12 w-12 text-green-600" />
      <h2 class="mt-4 text-3xl font-bold text-gray-900">
        Join Yappin
      </h2>
      <p class="mt-2 text-sm text-gray-600">
        Create your account to start chatting
      </p>
    </div>
    
    <form class="space-y-6" on:submit|preventDefault={handleSubmit}>
      {#if error}
        <div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-md">
          <p class="text-sm">{error}</p>
        </div>
      {/if}

      <Input
        type="text"
        bind:value={username}
        label="Username"
        placeholder="Choose a username"
        required
      />

      <Input
        type="email"
        bind:value={email}
        label="Email Address"
        placeholder="Enter your email"
        required
      />

      <Input
        type="password"
        bind:value={password}
        label="Password"
        placeholder="Create a password (min 6 chars)"
        required
      />

      <Button
        type="submit"
        variant="primary"
        disabled={loading}
        class="w-full bg-green-600 hover:bg-green-700 text-white font-medium py-2 px-4 rounded-md transition duration-200"
      >
        {loading ? 'Creating account...' : 'Create Account'}
      </Button>
    </form>

    <div class="mt-6 text-center">
      <p class="text-sm text-gray-600">
        Already have an account?
        <a href="/login" class="font-medium text-green-600 hover:text-green-500 transition duration-200">
          Sign in
        </a>
      </p>
    </div>
  </div>
</div>
