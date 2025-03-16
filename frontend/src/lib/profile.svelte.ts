import { writable } from 'svelte/store';

export type Profile =
	| {
			status: 'logged-in';
			name: string;
			pictureURL: string;
	  }
	| {
			status: 'logged-out';
	  }
	| {
			status: 'initializing';
	  };

function initializeProfile(set: (value: Profile) => void) {
	fetch('/api/profile').then(async (response) => {
		if (response.ok) {
			const json = await response.json();
			set({
				status: 'logged-in',
				name: json.name,
				pictureURL: json.pictureURL
			});
		} else {
			set({
				status: 'logged-out'
			});
		}
	});
}

export const profile = writable<Profile>({ status: 'initializing' }, initializeProfile);
