import type { Event, EventResponse } from '$lib/types';
import { error } from '@sveltejs/kit';

export type EventData = {
	event: Event;
};

export async function load({ params, fetch }): Promise<EventData> {
	const eventID = params.id;
	if (eventID == '') {
		error(400, 'No event id provided');
	}

	const response = await fetch(`/api/event?id=${eventID}`);
	const event = ((await response.json()) as EventResponse).events[0];

	if (!response.ok) {
		error(
			400,
			`Failed to load event id=${window.location.hash}, got status ${response.status} with message ${await response.text()}`
		);
	}

	return {
		event: event
	};
}

// prevent prerendering of this page
// export const ssr = false;
