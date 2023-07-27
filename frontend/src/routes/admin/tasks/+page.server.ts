/* Copyright (C) 2023  Semyon Hoyrish */
/* For more details see the "LICENSE" file */

export async function load({ fetch, cookies }) {
	return await (
		await fetch('http://localhost:8080/tasks/get_all', {
			method: 'POST',
			body: JSON.stringify({
				token: (await cookies.get('TOKEN')) || '',
				taskId: -1
			})
		})
	).json();
}
