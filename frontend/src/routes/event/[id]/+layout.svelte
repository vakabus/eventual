<script lang="ts">
	import { type EventData } from './+layout';
	import { type Event } from '$lib/types';
	import type { Snippet } from 'svelte';
	import { page } from '$app/state';

	let { data, children }: { data: EventData; children: Snippet } = $props();
	let event: Event = $derived(data.event);

	let tabsActive = $derived({
		participant: page.route.id?.includes('/participant') ? 'active' : '',
		template: page.route.id?.includes('/template') ? 'active' : '',
		edit: page.route.id?.includes('/edit') ? 'active' : ''
	});
	let defaultTabActive = $derived(Object.values(tabsActive).every((v) => v === '') ? 'active' : '');
</script>

<h1 class="my-4">{event.name}</h1>

<nav class="mb-4">
	<ul class="nav nav-tabs" role="tablist">
		<li class="nav-item" role="presentation">
			<a
				class="nav-link {defaultTabActive}"
				href="#/event/{data.event.id}"
				aria-selected="false"
				role="tab"
				tabindex="-1">Přehled</a
			>
		</li>
		<li class="nav-item" role="presentation">
			<a
				class="nav-link {tabsActive.edit}"
				href="#/event/{data.event.id}/edit"
				aria-selected="false"
				role="tab"
				tabindex="-1">Editovat</a
			>
		</li>
		<li class="nav-item" role="presentation">
			<a
				class="nav-link {tabsActive.participant}"
				href="#/event/{data.event.id}/participant"
				aria-selected="true"
				role="tab">Účastníci</a
			>
		</li>
		<li class="nav-item" role="presentation">
			<a
				class="nav-link {tabsActive.template}"
				href="#/event/{data.event.id}/template"
				aria-selected="false"
				tabindex="-1"
				role="tab">Vzory mailů</a
			>
		</li>
	</ul>
</nav>

<section>
	{@render children()}
</section>
