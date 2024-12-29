<script lang="ts">
	import '../app.css';
	let { children } = $props();

	function goHome() {
		// redirect to our Google login endpoint, JS because
		// Svelte likes to control links without domain
		window.location.pathname = '/';
	}
</script>

<nav class="pt-4 bg-blue-300 text-blue-950">
	<div class="flex flex-row max-w-prose mx-auto gap-8 text-xl">
		<button type="button" onclick={goHome}>Eventovátko</button>
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
