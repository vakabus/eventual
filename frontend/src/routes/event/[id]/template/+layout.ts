import type { Template } from '$lib/types.js';
import { error } from '@sveltejs/kit';

export type TemplateContext = {
	templates: Template[];
};

export async function load({ fetch, params }): Promise<TemplateContext> {
	const eventID = params.id;
	if (eventID == '' || eventID == null) {
		error(400, 'No event id provided');
	}

	const response = await fetch(`/api/event/${eventID}/template`);
	if (!response.ok) {
		error(response.status, 'Failed to fetch a list of templates');
	}

	return {
		templates: await response.json()
	};
}
