import type { Event } from './types';

class Events {
	events: Event[] = $state([]);

	async updateEvents() {
		this.events = await fetch('/api/event').then((res) => res.json());
		return this.events;
	}

	nameFromId(id: string): string {
		return this.events.find((e) => e.id == id)?.name ?? 'BROKEN PLEASE REPORT';
	}
}

export const events = new Events();
