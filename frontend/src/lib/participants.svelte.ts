import { error } from '@sveltejs/kit';
import type { Participant, Participants } from './types';

type OurParticipant = {
	__id__: string;
	[key: string]: string;
};

export class Ucastnik {
	id: string = $state('');
	eventID: string = $state('');
	data: Record<string, string> = $state({});

	constructor(eventID: string, participant: Participant) {
		this.id = participant.id;
		this.eventID = eventID;
		this.data = participant.data;
	}

	get keys(): string[] {
		return Object.keys(this.data);
	}

	async delete() {
		if (!this.id) {
			error(500, 'Participant ID is not set when deleting');
		}

		const response = await fetch(`/api/event/${this.eventID}/participant/${this.id}`, {
			method: 'DELETE'
		});

		if (!response.ok) {
			console.error(response);
			error(500, 'Nepodařilo se odstranit účastníka.');
		}
	}

	async save() {
		const payload: Participant = {
			id: this.id,
			data: this.data
		};

		let url;
		if (this.id == '') {
			url = `/api/event/${this.eventID}/participant`;
		} else {
			url = `/api/event/${this.eventID}/participant/${this.id}`;
		}

		const response = await fetch(url, {
			method: 'POST',
			body: JSON.stringify(payload)
		});

		if (!response.ok) {
			console.error(response);
			error(500, 'Nepodařilo se uložit účastníka.');
		} else {
			const data = await response.json();
			if (this.id == '') {
				this.id = data.id;
			} else if (this.id != data.id) {
				error(500, 'ID of the participant changed when saving');
			}
		}
	}
}

export class ParticipantsData {
	participants: Ucastnik[] = $state([]);
	keys: string[] = $state([]);
	eventID: string = $state('');
	fetch: typeof window.fetch;

	private constructor(eventID: string, data: Participants, fetch: typeof window.fetch) {
		this.participants = data.participants.map((p) => new Ucastnik(eventID, p));
		this.keys = [...new Set(this.participants.flatMap((p) => p.keys))];
		this.eventID = eventID;
		this.fetch = fetch;
	}

	get(): Record<string, string>[] {
		return this.participants.map((p) => p.data);
	}

	static async init(
		eventID: string,
		fetch: typeof window.fetch = window.fetch
	): Promise<ParticipantsData> {
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

		await toDelete.delete();
	}

	async addNew() {
		// add locally to trigger UI update
		const newIndex = this.participants.length;
		const newParticipant = new Ucastnik(this.eventID, {
			id: '',
			data: {}
		});
		this.participants = [...this.participants, newParticipant];

		await newParticipant.save();
	}

	async notifyUpdate(index: number) {
		if (this.participants[index].id == '') {
			// it's a row that was never actually saved on the server, so we can just ignore it
			alert('This row was never saved on the server, this should not happen.');
			return;
		}

		await this.participants[index].save();
	}
}
