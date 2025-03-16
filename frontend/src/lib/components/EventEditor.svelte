<script lang="ts">
	import type { Event } from '$lib/types';
	let { event }: { event: Event } = $props();

	export async function pushUpdate() {
		const response = await fetch(`/api/event/${event.id}`, {
			method: 'POST',
			body: JSON.stringify(event)
		});
		if (!response.ok) console.error(response);
	}

	export async function createEvent(): Promise<Event> {
		const response = await fetch('/api/event', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({
				name: event.name,
				description: event.description
			})
		});
		if (response.ok) return await response.json();
		else {
			console.error(response);
			throw 'error';
		}
	}
</script>

<div class="form-group">
	<label for="name" class="form-group-label">Název akce</label>
	<input
		id="name"
		type="text"
		bind:value={event.name}
		placeholder="Název akce"
		class="form-control"
	/>
</div>

<div class="form-group">
	<label for="description" class="form-group-label">Popis akce (markdown)</label>
	<textarea
		id="description"
		bind:value={event.description}
		placeholder="Popis akce"
		class="form-control"
	></textarea>
</div>
