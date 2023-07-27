/* Copyright (C) 2023  Semyon Hoyrish */
/* For more details see the "LICENSE" file */

import { json } from '@sveltejs/kit';

export async function POST({ fetch, cookies }) {
	const resp = await (
		await fetch('http://localhost:8080/user/logout', {
			method: 'POST',
			body: JSON.stringify({
				token: (await cookies.get('TOKEN')) || '',
				user: {}
			})
		})
	).json();

	return json(resp, { status: 200 });
}
