<script lang="ts">
	import { type EventData } from './+layout';
	import { type Event } from '$lib/types';
	import type { Snippet } from 'svelte';

	let { data, children }: { data: EventData; children: Snippet } = $props();
	let event: Event = $derived(data.event);
</script>

<header>
	<h1 class="text-2xl">{event.name}</h1>
</header>

<section>
	{@render children()}
</section>
