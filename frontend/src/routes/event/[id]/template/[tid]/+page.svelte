<script lang="ts">
	import TemplateEditor from '../TemplateEditor.svelte';
	import type { Template } from '$lib/types';
	import { page } from '$app/state';
	import { goto, invalidate } from '$app/navigation';
	import { error } from '@sveltejs/kit';
	import { debounce } from 'lodash';

	let { data } = $props();

	// we need to use a state variable here, because we need to invalidate the template when it changes
	let template: Template = $state(
		data.templates.find((t) => t.id == page.params.tid) ||
			error(404, { message: 'Template not found' })
	);

	const debouncedChange = debounce(async () => {
		const response = await fetch(`/api/event/${data.event.id}/template/${page.params.tid}`, {
			method: 'POST',
			body: JSON.stringify(template)
		});
		if (response.ok) {
			invalidate(`/api/event/${data.event.id}/template`);
		} else {
			error(response.status, { message: 'Failed to update template' });
		}
	}, 500);

	$effect(() => {
		// effect triggers
		const triggers = [template.name, template.body];

		setTimeout(async () => {
			await debouncedChange();
		}, 0);
	});

	async function remove() {
		const response = await fetch(`/api/event/${data.event.id}/template/${page.params.tid}`, {
			method: 'DELETE'
		});
		if (response.ok) {
			await goto(`#/event/${data.event.id}/template`, { invalidateAll: true });
		} else {
			error(response.status, { message: 'Failed to delete template' });
		}
	}
</script>

<TemplateEditor bind:template />
<button
	class="mt-4 bg-red-500 hover:bg-red-700 text-white font-bold py-2 px-4 rounded"
	onclick={remove}>Odstranit šablonu</button
>
