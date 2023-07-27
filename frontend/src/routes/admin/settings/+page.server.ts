/* Copyright (C) 2023  Semyon Hoyrish */
/* For more details see the "LICENSE" file */

export async function load({ cookies }) {
	return await (
		await fetch('http://localhost:8080/settings/get', {
			method: 'POST',
			body: JSON.stringify({
				token: (await cookies.get('TOKEN')) || '',
				settings: {}
			})
		})
	).json();
}
