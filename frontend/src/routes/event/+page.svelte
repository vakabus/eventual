<script lang="ts">
	let { data } = $props();
	let events = $derived(data.events);
</script>

<div class="flex flex-row">
	<h1 class="text-2xl my-4 me-8">Tvoje akce</h1>
	<a href="#/event/new" class="btn btn-primary mb-4">Nová akce</a>
</div>

{#if events != null && events.length > 0}
	<div class="d-flex flex-row">
		{#each events as event}
			<div class="card border-primary m-2" style="max-width: 20rem;">
				<a class="stretched-link" href="#/event/{event.id}">
					<div class="card-header">Akce</div>
				</a>
				<div class="card-body">
					<h4 class="card-title">{event.name}</h4>
					<p class="card-text">{event.description}</p>
				</div>
			</div>
		{/each}
	</div>
{:else}
	<p>Nemáš žádné akce.</p>
{/if}
