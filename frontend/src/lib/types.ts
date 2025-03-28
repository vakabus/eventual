// Code generated by tygo. DO NOT EDIT.

//////////
// source: types.go

export interface ErrorResponse {
	errorMessage: string;
}
export interface Event {
	id: string;
	name: string;
	description: string;
}
export type Events = Event[];
export interface Profile {
	name: string;
	email: string;
	pictureURL: string;
}
export interface Organizer {
	id: string;
	name: string;
	email: string;
}
export type Organizers = Organizer[];
export interface Participant {
	id: string;
	data: { [key: string]: string };
}
export interface Participants {
	participants: Participant[];
}
export interface Template {
	id: string;
	name: string;
	body: string;
}
export type Templates = Template[];
