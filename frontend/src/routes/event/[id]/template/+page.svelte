<script lang="ts">
	let { data } = $props();
	let event_id = $derived(data.event.id);
	let templates = $derived(data.templates);
</script>

<div class="flex flex-col gap-4">
	<a href="#/event/{event_id}/template/new" class="btn btn-primary mb-4"> Nová šablona </a>

	<ul class="d-flex gap-2">
		{#each templates as template}
			<li>
				<a href="#/event/{event_id}/template/{template.id}">{template.name}</a>
			</li>
		{/each}
	</ul>
</div>
