import type { Event } from '$lib/types';
import { events } from '$lib/events.svelte';

type Data = {
	events: Event[];
};

export async function load(): Promise<Data> {
	return {
		events: await events.updateEvents()
	};
}
