<script lang="ts">
	import { goto } from '$app/navigation';
	import { login as apiLogin } from '$lib/api';
	import { base } from '$app/paths';
	import { type AuthResponse, type ErrorResponse } from '$lib/types';
	import { presentNorEmpty } from '$lib/utils';

	let state: 'input' | 'processing' = $state('input');
	let username = $state('');
	let loginResponse: AuthResponse | ErrorResponse = $state(undefined);

	async function login(event: Event) {
		// Prevent the form from submitting normally
		event.preventDefault();

		state = 'processing';
		loginResponse = await apiLogin(username);
		if (presentNorEmpty(loginResponse, 'token')) {
			localStorage.setItem('token', loginResponse.token);
			await goto(`${base}/dashboard`);
		} else {
			state = 'input';
		}
	}
</script>

<main class="container mx-auto px-6 py-12">
	{#if state == 'input'}
		<form action="#" onsubmit={login}>
			<input type="text" bind:value={username} placeholder="Username" />
			<input type="submit" value="Login" />
			{#if loginResponse != null && presentNorEmpty(loginResponse, 'errorMessage')}
				<div>{loginResponse.errorMessage}</div>
			{/if}
		</form>
	{:else if state == 'processing'}
		<div>Processing authentication...</div>
	{/if}
</main>
