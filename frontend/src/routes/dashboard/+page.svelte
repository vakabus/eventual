<script lang="ts">
	import type { EventResponse, Profile } from '$lib/types';

	async function fetchEvents(): Promise<EventResponse> {
		return (await fetch('/api/event')).json();
	}

	async function fetchProfile(): Promise<Profile> {
		return (await fetch('/api/profile')).json();
	}

	function logout() {
		// redirect to our logout endpoint, non-SvelteKit, so has to be done via JS
		window.location.pathname = '/auth/logout';
	}
</script>

{#await fetchProfile()}
{:then profile} 
<section class="float-end max-w-64 bg-slate-400 border-solid border-2 rounded-2xl py-2 px-4">
	<h3 class="text-xl mb-2">{profile.name}</h3>
	<img class="max-w-36 mb-2" alt="profilová fotka" src={profile.pictureURL} />
	<button
		onclick={logout}
		type="button"
		class="text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 me-2 mb-2 dark:bg-blue-600 dark:hover:bg-blue-700 focus:outline-none dark:focus:ring-blue-800"
		>Odhlásit</button>
</section>
{/await}

<div class="flex flex-row">
	<h1 class="text-2xl mb-4 me-8">Tvoje akce</h1>
	<button
		type="button"
		class="text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 me-2 mb-2 dark:bg-blue-600 dark:hover:bg-blue-700 focus:outline-none dark:focus:ring-blue-800"
		>Nová akce</button
	>
</div>

{#await fetchEvents()}
	<p>Načítám akce...</p>
{:then events}
	{#if events.events != null && events.events.length > 0}
		<ul>
			{#each events.events as event}
				<li>{event.name}</li>
			{/each}
		</ul>
	{:else}
		<p>Nemáš žádné akce.</p>
	{/if}
{:catch error}
	<p>Chyba při načítání akcí: {error.message}</p>
{/await}
