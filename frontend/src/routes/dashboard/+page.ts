import type { EventResponse, Profile } from '$lib/types';

export async function load({ fetch }): Promise<{ events: EventResponse; profile: Profile }> {
	async function fetchEvents(): Promise<EventResponse> {
		return (await fetch('/api/event')).json();
	}

	async function fetchProfile(): Promise<Profile> {
		return (await fetch('/api/profile')).json();
	}

	const [events, profile] = await Promise.all([fetchEvents(), fetchProfile()]);

	return {
		events: events,
		profile: profile
	};
}
