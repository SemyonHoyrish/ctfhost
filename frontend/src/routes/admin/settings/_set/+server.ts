/* Copyright (C) 2023  Semyon Hoyrish */
/* For more details see the "LICENSE" file */

import { json } from '@sveltejs/kit';

export async function POST({ request, cookies }) {
	const d = await request.json();
	const r = await (
		await fetch('http://localhost:8080/settings/set', {
			method: 'POST',
			body: JSON.stringify({
				token: (await cookies.get('TOKEN')) || '',
				settings: d.settings
			})
		})
	).json();

	return json(r, { status: 200 });
}
