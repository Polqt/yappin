<script lang="ts">
  import { signup, setToken } from '$lib/auth';
  import { goto } from '$app/navigation';

  let username = '';
  let email = '';
  let password = '';
  let error = '';

  async function handleSignup() {
    try {
      const data = await signup(username, email, password);
      setToken(data.token);
      goto('/'); // Redirect to home after signup
    } catch (err: any) {
      error = err.message;
    }
  }
</script>

<main class="min-h-screen bg-gray-100 flex items-center justify-center">
  <form on:submit|preventDefault={handleSignup} class="bg-white p-8 rounded shadow-md w-96">
    <h2 class="text-2xl font-bold mb-4">Sign Up</h2>
    {#if error}<p class="text-red-500 mb-4">{error}</p>{/if}
    <input bind:value={username} type="text" placeholder="Username" class="w-full p-2 mb-4 border rounded" required />
    <input bind:value={email} type="email" placeholder="Email" class="w-full p-2 mb-4 border rounded" required />
    <input bind:value={password} type="password" placeholder="Password" class="w-full p-2 mb-4 border rounded" required />
    <button type="submit" class="w-full bg-green-500 text-white p-2 rounded">Sign Up</button>
    <p class="mt-4 text-center"><a href="/login" class="text-blue-500">Already have an account? Login</a></p>
  </form>
</main>
