<script lang="ts">
	import 'bootstrap';
	import type { Snippet } from 'svelte';
	import '../app.scss';
	import NavBar from '$lib/components/NavBar.svelte';
	let { children }: { children: Snippet } = $props();
</script>

<NavBar />

<div class="container">
	{@render children()}
</div>
