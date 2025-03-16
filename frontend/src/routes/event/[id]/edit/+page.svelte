<script lang="ts">
	import { goto, invalidateAll } from '$app/navigation';
	import EventEditor from '$lib/components/EventEditor.svelte';
	import type { Event } from '$lib/types';
	import type { EventData } from '../+layout';

	let { data }: { data: EventData } = $props();
	let event: Event = $derived(data.event);
	let updating = $state(false);
	let editor: EventEditor = $state();

	async function updateEvent() {
		updating = true;
		editor.pushUpdate();
		await goto(`#/event/${event.id}`, {
			invalidateAll: true
		});
	}
</script>

{#if !updating}
	<div>
		<EventEditor {event} bind:this={editor} />
		<button type="button" onclick={updateEvent} class="btn btn-primary">Změnit data</button>
	</div>
{:else}
	<p>Upravuji existující akci...</p>
{/if}
