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

<div class="flex flex-row">
	<nav class="max-w-fit p-2 me-4 bg-slate-400">
		<a href="#/event/{data.event.id}/edit" class="block">Upravit událost</a>
		<a href="#/event/{data.event.id}/participant" class="block">Účastníci</a>
	</nav>
	<section>
		{@render children()}
	</section>	
</div>
