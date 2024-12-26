import type { AuthRequest, AuthResponse, ErrorResponse } from '$lib/types.ts';

export async function login(username: string): Promise<AuthResponse | ErrorResponse> {
	const request: AuthRequest = {
		username: username
	};
	const response = await fetch('/api/v1/auth', {
		body: JSON.stringify(request),
		method: 'POST',
		headers: {
			'Content-Type': 'application/json'
		}
	});

	return await response.json();
}
