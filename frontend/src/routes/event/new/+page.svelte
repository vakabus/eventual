<script lang="ts">
	import { goto } from '$app/navigation';
	import EventEditor from '$lib/components/EventEditor.svelte';
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
		<button type="button" onclick={createEvent} class="btn btn-primary">Založit novou akci</button>
	</div>
{:else}
	<p>Vytvářím novou akci...</p>
{/if}
