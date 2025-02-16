<script lang="ts">
	import { Carta, MarkdownEditor } from 'carta-md';
	import DOMPurify from 'isomorphic-dompurify';
	import 'carta-md/default.css';
	import type { Template } from '$lib/types';

	let { template = $bindable<Template>() } = $props();

	const carta = new Carta({
		sanitizer: (dirty) => DOMPurify.sanitize(dirty, { USE_PROFILES: { html: true } })
	});
</script>

<h2
	contenteditable="true"
	class="text-2xl font-bold mb-4"
	onkeydown={(e) => e.key === 'Enter' && e.preventDefault()}
	bind:innerHTML={template.name}
></h2>

<MarkdownEditor {carta} bind:value={template.body} />

<style>
	/* Or in global stylesheet */
	/* Set your monospace font (Required to have the editor working correctly!) */
	:global(.carta-font-code) {
		font-family: '...', monospace;
		font-size: 1.1rem;
	}
</style>
