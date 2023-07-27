/* Copyright (C) 2023  Semyon Hoyrish */
/* For more details see the "LICENSE" file */

export async function load({ cookies }) {
	if ((await cookies.get('TOKEN')) !== undefined) {
		const user = await (
			await fetch('http://localhost:8080/user/get', {
				method: 'POST',
				body: JSON.stringify({
					token: (await cookies.get('TOKEN')) || '',
					nickname: ''
				})
			})
		).json();
		if (user.success === false) {
			return { logged_in: false };
		}

		return { logged_in: true, user: user };
	}
	return { logged_in: false };
}
