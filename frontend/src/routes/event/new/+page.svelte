<script lang="ts">
	import { goto } from '$app/navigation';
	import EventEditor from '$lib/EventEditor.svelte';
	import type { Event } from '$lib/types';

	let event: Event = $state({ id: '', description: '', name: '' });
	let loading = $state(false);
	let editor: EventEditor;

	async function createEvent() {
		const newEvent = await editor.createEvent();
		await goto(`#/event/${newEvent.id}`);
	}
</script>

{#if !loading}
	<div>
		<EventEditor {event} bind:this={editor} />
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
