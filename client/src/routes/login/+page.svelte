<script lang="ts">
  import { login, setToken } from '$lib/auth';
  import { goto } from '$app/navigation';

  let identifier = '';
  let password = '';
  let error = '';

  async function handleLogin() {
    try {
      const data = await login(identifier, password);
      setToken(data.token);
      goto('/'); // Redirect to home after login
    } catch (err: any) {
      error = err.message;
    }
  }
</script>

<main class="min-h-screen bg-gray-100 flex items-center justify-center">
  <form on:submit|preventDefault={handleLogin} class="bg-white p-8 rounded shadow-md w-96">
    <h2 class="text-2xl font-bold mb-4">Login</h2>
    {#if error}<p class="text-red-500 mb-4">{error}</p>{/if}
    <input bind:value={identifier} type="email" placeholder="Email" class="w-full p-2 mb-4 border rounded" required />
    <input bind:value={password} type="password" placeholder="Password" class="w-full p-2 mb-4 border rounded" required />
    <button type="submit" class="w-full bg-blue-500 text-white p-2 rounded">Login</button>
    <p class="mt-4 text-center"><a href="/signup" class="text-green-500">Don't have an account? Sign Up</a></p>
  </form>
</main>
