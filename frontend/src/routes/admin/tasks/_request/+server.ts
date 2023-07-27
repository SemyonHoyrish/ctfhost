/* Copyright (C) 2023  Semyon Hoyrish */
/* For more details see the "LICENSE" file */

import { json } from '@sveltejs/kit';

export async function POST({ request, fetch, cookies }) {
	const data = await request.json();

	const resp = await (
		await fetch(`http://localhost:8080/tasks/${data.action}`, {
			method: 'POST',
			body: JSON.stringify({
				token: (await cookies.get('TOKEN')) || '',
				task: data.task
			})
		})
	).json();

	return json(resp, { status: 200 });
}
