import type { Profile, Event } from '$lib/types';

export async function load({ fetch }): Promise<{ events: Event[] }> {
	async function fetchEvents(): Promise<Event[]> {
		return (await fetch('/api/event')).json();
	}

	const [events] = await Promise.all([fetchEvents()]);

	return {
		events: events
	};
}
