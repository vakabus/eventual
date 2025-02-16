<script lang="ts">
	let { data } = $props();
	let event_id = $derived(data.event.id);
	let templates = $derived(data.templates);
</script>

<div class="flex flex-col gap-4">
	<a
		href="#/event/{event_id}/template/new"
		class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded w-fit"
	>
		Nová šablona
	</a>

	<ul class="flex flex-col gap-2">
		{#each templates as template}
			<li>
				<a
					href="#/event/{event_id}/template/{template.id}"
					class="block p-4 border rounded hover:bg-slate-100"
				>
					{template.name}
				</a>
			</li>
		{/each}
	</ul>
</div>
