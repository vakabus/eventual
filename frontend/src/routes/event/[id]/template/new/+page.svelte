<script lang="ts">
	import { goto } from '$app/navigation';
	import { error } from '@sveltejs/kit';
	import TemplateEditor from '../TemplateEditor.svelte';

	let { data } = $props();
	let template = $state({
		name: '(nová šablona)',
		body: `---
To: {{ .jméno }} <{{ .email }}>
Reply-To: kurz@psl.cz
Subject: Úvodní informace ke kurzu
---

Ahoj {{ .oslovení }},

nemohl(|a) bys prosím zkontrolovat, že slovesa mají správný rod? Jsi totiž moc šikovn(ý|á) a zvládneš to nejlépe.

Díky a měj se pěkně,
Tvoji organizátoři
`});

	async function create() {
		var response = await fetch(`/api/event/${data.event.id}/template`, {
			method: 'POST',
			body: JSON.stringify(template)
		});

		if (response.ok) {
			let r = await response.json();
			await goto(`#/event/${data.event.id}/template/${r.id}`, { invalidateAll: true });
		} else {
			error(response.status, { message: 'Failed to create template' });
		}
	}
</script>

<TemplateEditor bind:template />
<button
	class="bt-4 bg-blue-500 hove:bg-blue-700 text-white font-bold py-2 px-4 rounded"
	onclick={create}>Create</button
>
