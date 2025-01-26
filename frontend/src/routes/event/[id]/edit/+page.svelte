<script lang="ts">
	import { goto, invalidateAll } from '$app/navigation';
	import EventEditor from '$lib/EventEditor.svelte';
	import type { Event } from '$lib/types';
	import type { EventData } from '../+layout';

    let { data }: { data: EventData } = $props()
    let event: Event = $derived(data.event)
	let updating = $state(false);
	let editor: EventEditor = $state();

	async function updateEvent() {
        updating = true
		editor.pushUpdate()
		await goto(`#/event/${event.id}`, {
			invalidateAll: true
		})
	}
</script>

{#if !updating}
	<div>
		<EventEditor {event} bind:this={editor} />
		<button
			type="button"
			onclick={updateEvent}
			class="text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 me-2 mb-2 dark:bg-blue-600 dark:hover:bg-blue-700 focus:outline-none dark:focus:ring-blue-800"
			>Změnit data</button>
	</div>
{:else}
	<p>Upravuji existující akci...</p>
{/if}
