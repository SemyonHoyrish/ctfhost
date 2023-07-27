/* Copyright (C) 2023  Semyon Hoyrish */
/* For more details see the "LICENSE" file */

export async function load({ fetch }) {
	return await (
		await fetch('http://localhost:8080/team/get_leaderboard', {
			method: 'POST',
			body: JSON.stringify({})
		})
	).json();
}
