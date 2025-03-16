import { ParticipantsData } from '$lib/participants.svelte.js';
import { error } from '@sveltejs/kit';

export type Data = {
	participants: ParticipantsData;
};

export async function load({ params, fetch }): Promise<Data> {
	const eventID = params.id;
	if (eventID == '') {
		error(400, 'No event id provided');
	}

	return {
		participants: await ParticipantsData.init(eventID, fetch)
	};
}
