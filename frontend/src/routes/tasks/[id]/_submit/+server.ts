/* Copyright (C) 2023  Semyon Hoyrish */
/* For more details see the "LICENSE" file */

import { json } from '@sveltejs/kit';

export async function POST({ params, request, fetch, cookies }) {
	const d = await request.json();
	const resp = await (
		await fetch('http://localhost:8080/team/submit_flag', {
			method: 'POST',
			body: JSON.stringify({
				token: (await cookies.get('TOKEN')) || '',
				taskId: parseInt(params.id),
				...d
			})
		})
	).json();

	return json(resp, { status: 200 });
}
