<script lang="ts">
	import type { Snippet } from 'svelte';
	import '../app.css';
	let { children }: { children: Snippet } = $props();
</script>

<nav class="pt-4 bg-blue-300 text-blue-950">
	<div class="flex flex-row max-w-prose mx-auto gap-8 text-xl">
		<a href="#/">Eventovátko</a>
	</div>
</nav>

<div class="max-w-prose mx-auto mt-4">
	{@render children()}
</div>

<style lang="postcss">
	:global(body) {
		@apply bg-gray-50;
		@apply text-gray-900;
	}
</style>
