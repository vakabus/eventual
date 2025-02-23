import { error } from '@sveltejs/kit';
import type { Participant, Participants } from './types';

type OurParticipant = {
	__id__: string;
	[key: string]: string;
}

function convertBack(p: OurParticipant): Participant {
	return {
		id: p.__id__,
		data: p
	}
}

export class ParticipantsData {
	participants: OurParticipant[] = $state([]);
	keys: string[] = $state([])
	eventID: string
	fetch: typeof window.fetch;

	private constructor(eventID: string, data: Participants, fetch: typeof window.fetch) {
		this.participants = data.participants.map((p) => ({ ...p.data, __id__: p.id }));
		this.keys = [...new Set(data.participants.flatMap((p) => Object.keys(p.data)))];
		this.eventID = eventID;
		this.fetch = fetch;
	}

	get(): OurParticipant[] {
		return this.participants;
	}

	static async init(eventID: string, fetch: typeof window.fetch = window.fetch): Promise<ParticipantsData> {
		const response = await fetch(`/api/event/${eventID}/participant`);
		if (response.ok) {
			const data: Participants = await response.json();
			return new ParticipantsData(eventID, data, fetch);
		} else {
			error(500, 'Failed to get a list of participants');
		}
	}

	async delete(index: number) {
		const toDelete = this.participants[index];

		// delete locally to update the UI
		this.participants = this.participants.filter((_, i) => i != index);

		// delete on the server
		if (toDelete.id == '') {
			// it's a row that was never actually saved on the server, so we can just ignore it
			return;
		}

		const response = await this.fetch(`/api/event/${this.eventID}/participant/${toDelete.__id__}`, {
			method: 'DELETE'
		});
		if (!response.ok) {
			console.error(response);
			error(500, 'Nepodařilo se uložit odstranění účastníka.');
		}
	}

	async addNew() {
		// add locally to trigger UI update
		const newIndex = this.participants.length;
		this.participants = [
			...this.participants,
			{
				__id__: '',
			}
		];

		// add on the server
		const response = await this.fetch(`/api/event/${this.eventID}/participant`, {
			method: 'POST',
			body: JSON.stringify(convertBack(this.participants[newIndex]))
		});
		if (!response.ok) {
			console.log(response);
			error(500, 'Nepodařilo se přidat účastníka.');
		} else {
			const data = await response.json();
			this.participants[newIndex].id = data.id;
		}
	}

	async notifyUpdate(index: number) {
		if (this.participants[index].id == '') {
			// it's a row that was never actually saved on the server, so we can just ignore it
			alert('This row was never saved on the server, this should not happen.');
			return;
		}

		const response = await this.fetch(
			`/api/event/${this.eventID}/participant/${this.participants[index].__id__}`,
			{
				method: 'POST',
				body: JSON.stringify(convertBack(this.participants[index]))
			}
		);

		if (!response.ok) {
			console.error(response);
			error(500, 'failed saving the data');
		}
	}
}
