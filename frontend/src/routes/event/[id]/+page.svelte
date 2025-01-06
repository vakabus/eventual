<script lang="ts">
	import { type EventData } from './+layout';
	import { type Event } from '$lib/types';
	import { marked } from 'marked';

	let { data }: { data: EventData } = $props();
	let event: Event = $derived(data.event);
	let renderedHTML: string = $derived(marked.parse(event.description));

	let content: HTMLParagraphElement;

	$effect(() => {
		content.innerHTML = renderedHTML;
	});
</script>

<!-- Note: we use tailwind/typography styling for the markdown -->
<p bind:this={content} class="prose"></p>
