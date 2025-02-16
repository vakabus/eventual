<script lang="ts">
	import {
		RevoGrid,
		type AfterEditEvent,
		type BeforeSaveDataDetails,
		type ColumnRegular
	} from '@revolist/svelte-datagrid';

	let { data } = $props();
	let participants = $derived(data.participants);

	const columns: ColumnRegular[] = [
		{
			prop: '',
			size: 50,
			cellTemplate(h, { rowIndex }) {
				return h(
					'span',
					{
						onclick: async () => {
							await participants.delete(rowIndex);
						}
					},
					'🧨'
				);
			}
		},
		{
			prop: 'id',
			name: 'ID',
			readonly: true,
			sortable: true
		},
		{
			prop: 'name',
			name: 'Jméno',
			sortable: true,
			editable: true
		},
		{
			prop: 'email',
			name: 'Email',
			sortable: true,
			editable: true
		}
	];

	async function afteredit(ev: BeforeSaveDataDetails) {
		await participants.notifyUpdate(ev.detail.rowIndex);
	}

	async function add(ev: Event) {
		await participants.addNew();
	}
</script>

<RevoGrid
	source={participants.get()}
	{columns}
	resize={true}
	theme="material"
	on:afteredit={afteredit}
></RevoGrid>
<button onclick={add}>Přidat účastníka</button>
