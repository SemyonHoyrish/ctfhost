/* Copyright (C) 2023  Semyon Hoyrish */
/* For more details see the "LICENSE" file */

import { json } from '@sveltejs/kit';

// const SESSION_EXPIRES_IN = 60 * 60 * 24; // in seconds

export async function POST({ request, fetch, cookies }) {
	const data = await request.json();

	const action = data.nickname === '' ? 'login' : 'register';

	const j = await (
		await fetch(`http://localhost:8080/user/${action}`, {
			method: 'POST',
			body: JSON.stringify(data)
		})
	).json();

	if (j.success === true) {
		await cookies.set('TOKEN', j.data, {
			path: '/'
			// expirationDate: Date.now() + SESSION_EXPIRES_IN * 1000 // @Fix: idk how
		});
	}

	return json(j, { status: 200 });
}
