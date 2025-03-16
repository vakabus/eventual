<script lang="ts">
	import { goto } from '$app/navigation';
	import { profile } from '$lib/profile.svelte';

	profile.subscribe((p) => {
		if (p.status == 'logged-in') {
			goto(`#/event/`);
		} else {
			goto(`#/login/`);
		}
	});
</script>

<div class="text-center mx-auto">
	<p>Načítám...</p>
</div>
