//eslint-disable-next-line
export function present(obj: any, property: string): boolean {
	return obj[property] !== undefined && obj[property] !== null;
}

//eslint-disable-next-line
export function presentNorEmpty(obj: any, property: string): boolean {
	return present(obj, property) && obj[property] !== '';
}
