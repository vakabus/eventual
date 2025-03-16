<script lang="ts">
	import { Carta, MarkdownEditor } from 'carta-md';
	import 'carta-md/default.css';
	import './carta.scss';
	import type { Template } from '$lib/types';

	let { template = $bindable<Template>() } = $props();

	const carta = new Carta({
		sanitizer: false
	});
</script>

<h2
	contenteditable="true"
	onkeydown={(e) => e.key === 'Enter' && e.preventDefault()}
	bind:innerHTML={template.name}
></h2>

<MarkdownEditor bind:value={template.body} mode="tabs" {carta} theme="default" />

<style>
	/* Or in global stylesheet */
	/* Set your monospace font (Required to have the editor working correctly!) */
	:global(.carta-font-code) {
		font-family: '...', monospace;
		font-size: 1.1rem;
	}
</style>
