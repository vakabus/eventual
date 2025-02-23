<script lang="ts">
	import { type EventData } from './+layout';
	import { type Event } from '$lib/types';
	import { Carta, Markdown } from 'carta-md';

	let { data }: { data: EventData } = $props();
	let event: Event = $derived(data.event);

	let carta = new Carta({
		sanitizer: false,
	});
</script>

<!-- Note: we use tailwind/typography styling for the markdown -->
<div class="prose">
	<Markdown bind:value={event.description} {carta} />
</div>
