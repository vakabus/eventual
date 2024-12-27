<script lang="ts">
	import { base } from '$app/paths';
	import '../app.css';
	let { children } = $props();
</script>

<style lang="postcss">
	:global(body) {
		@apply bg-gray-50;
		@apply text-gray-900;
	}
</style>

<nav class="pt-4 bg-blue-300 text-blue-950">
	<div class="flex flex-row max-w-prose mx-auto gap-8 text-xl">
		<a href="{base}/">Eventovátko</a>
		<a href="https://cs.wikipedia.org/wiki/Kedluben">Kedlubny</a>
	</div>
</nav>

<div class="max-w-prose mx-auto mt-4">
	{@render children()}
</div>
