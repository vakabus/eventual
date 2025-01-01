<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';

	onMount(() => {
		fetch('/api/profile').then((response) => {
			if (response.ok) {
				goto(`#/dashboard/`);
			} else {
				goto(`#/login/`);
			}
		});
	});
</script>

<div class="text-center mx-auto">
	<p>Načítám...</p>
</div>
