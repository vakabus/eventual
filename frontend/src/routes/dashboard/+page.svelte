<script lang="ts">
	let { data } = $props();
	let profile = $derived(data.profile);
	let events = $derived(data.events);

	function logout() {
		// redirect to our logout endpoint, non-SvelteKit, so has to be done via JS
		window.location.hash = '';
		window.location.pathname = '/auth/logout';
	}
</script>

<section class="float-end max-w-64 bg-slate-400 border-solid border-2 rounded-2xl py-2 px-4">
	<h3 class="text-xl mb-2">{profile.name}</h3>
	<img class="max-w-36 mb-2" alt="profilová fotka" src={profile.pictureURL} />
	<button
		onclick={logout}
		type="button"
		class="text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 me-2 mb-2 dark:bg-blue-600 dark:hover:bg-blue-700 focus:outline-none dark:focus:ring-blue-800"
		>Odhlásit</button
	>
</section>

<div class="flex flex-row">
	<h1 class="text-2xl mb-4 me-8">Tvoje akce</h1>
	<a
		href="#/event/new"
		class="text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 me-2 mb-2 dark:bg-blue-600 dark:hover:bg-blue-700 focus:outline-none dark:focus:ring-blue-800"
		>Nová akce</a
	>
</div>

{#if events != null && events.length > 0}
	<div class="flex flex-wrap flex-row gap-4">
		{#each events as event}
			<a
				class="block max-w-80 min-w-52 rounded overflow-hidden shadow-lg"
				href="#/event/{event.id}"
			>
				<div class="px-6 py-4 cursor-pointer">
					<div class="font-bold text-xl mb-2">{event.name}</div>
					<p class="text-gray-700 text-base">
						{event.description}
					</p>
				</div>
			</a>
		{/each}
	</div>
{:else}
	<p>Nemáš žádné akce.</p>
{/if}
