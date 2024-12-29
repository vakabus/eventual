<script lang="ts">
	import { goto } from "$app/navigation"
	import { base } from "$app/paths"
	import { onMount } from "svelte";

	onMount(() => {
		fetch('/api/profile').then((response) => {
			if (response.ok) {
				goto(`${base}/dashboard/`)
			} else {
				goto(`${base}/login/`)
			}
		})
	})
</script>

<div class="text-center mx-auto">
	<p>Načítám...</p>	
</div>
