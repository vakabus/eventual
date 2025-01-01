<script lang="ts">
	import { goto } from '$app/navigation';
	import type { Event } from '$lib/types';

	let name = $state('');
	let description = $state('');
	let loading = $state(false);

	async function createEvent() {
		const response = await fetch('/api/event', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({
				name: name,
				description: description
			})
		});
		const data = (await response.json()) as Event;

		await goto(`#/event/${data.id}`);
	}
</script>

{#if !loading}
	<div>
		<input
			type="text"
			bind:value={name}
			placeholder="Název akce"
			class="block mb-2 w-full p-2 border border-gray-300 rounded-lg"
		/>
		<input
			type="text"
			bind:value={description}
			placeholder="Popis akce"
			class="block mb-2 w-full p-2 border border-gray-300 rounded-lg"
		/>
		<button
			type="button"
			onclick={createEvent}
			class="text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 me-2 mb-2 dark:bg-blue-600 dark:hover:bg-blue-700 focus:outline-none dark:focus:ring-blue-800"
			>Založit novou akci</button
		>
	</div>
{:else}
	<p>Vytvářím novou akci...</p>
{/if}
