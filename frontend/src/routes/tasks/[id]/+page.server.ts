/* Copyright (C) 2023  Semyon Hoyrish */
/* For more details see the "LICENSE" file */

export async function load({ params, fetch, cookies }) {
	if (params.id === undefined) {
		return {};
	}

	return await (
		await fetch('http://localhost:8080/tasks/get', {
			method: 'POST',
			body: JSON.stringify({
				token: (await cookies.get('TOKEN')) || '',
				taskId: parseInt(params.id)
			})
		})
	).json();
}
