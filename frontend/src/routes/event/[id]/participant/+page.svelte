<script lang="ts">
	import {
		RevoGrid,
		type AfterEditEvent,
		type BeforeSaveDataDetails,
		type ColumnRegular
	} from '@revolist/svelte-datagrid';

	let { data } = $props();
	let participants = $derived(data.participants);
	let newColumnName = $state('');

	let columns: ColumnRegular[] = $state([
		{
			prop: '',
			size: 50,
			cellTemplate(h, { rowIndex }) {
				return h(
					'span',
					{
						onclick: async () => {
							if (confirm('Opravdu chceš smazat účastníka?')) {
								await participants.delete(rowIndex);
							}
						}
					},
					'🧨'
				);
			}
		},
		...data.participants.keys.filter((key) => key !== '__id__').map((key) => ({
			prop: key,
			name: key,
			sortable: true
		}))
	])

	async function afteredit(ev: AfterEditEvent) {
		await participants.notifyUpdate(ev.detail.rowIndex);
	}

	async function add(ev: Event) {
		await participants.addNew();
	}

	function addColumn() {
		if (newColumnName.trim() === '') return;
		if (newColumnName.includes(' ')) {
			alert('Název sloupce nesmí obsahovat mezery');
			newColumnName = ''
			return;
		}
		
		// Check if column with this name already exists
		const exists = columns.some(col => col.name === newColumnName);
		if (exists) {
			alert('Sloupec s tímto názvem již existuje');
			return;
		}
		
		columns = [
			...columns,
			{
				prop: `${newColumnName}`,
				name: newColumnName,
				sortable: true
			}
		];
		
		newColumnName = '';
	}
</script>

<RevoGrid
	source={participants.get()}
	{columns}
	resize={true}
	theme="material"
	on:afteredit={afteredit}
></RevoGrid>

<div style="margin-top: 1rem">
	<button onclick={add}>Přidat účastníka</button>
	
	<div style="margin-top: 0.5rem">
		<input 
			type="text"
			bind:value={newColumnName}
			placeholder="Název nového sloupce"
		>
		<button onclick={addColumn}>Přidat sloupec</button>
	</div>
</div>
