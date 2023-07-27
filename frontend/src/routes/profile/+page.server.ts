/* Copyright (C) 2023  Semyon Hoyrish */
/* For more details see the "LICENSE" file */

export async function load({ fetch, cookies }) {
	const team = await (
		await fetch('http://localhost:8080/team/get', {
			method: 'POST',
			body: JSON.stringify({
				token: (await cookies.get('TOKEN')) || '',
				teamName: ''
			})
		})
	).json();

	const team_members = await (
		await fetch('http://localhost:8080/team/get_members', {
			method: 'POST',
			body: JSON.stringify({
				token: (await cookies.get('TOKEN')) || '',
				teamName: ''
			})
		})
	).json();

	const team_user_invites = await (
		await fetch('http://localhost:8080/team/get_invites', {
			method: 'POST',
			body: JSON.stringify({
				token: (await cookies.get('TOKEN')) || '',
				teamName: ''
			})
		})
	).json();

	const user_team_invites = await (
		await fetch('http://localhost:8080/user/get_invites', {
			method: 'POST',
			body: JSON.stringify({
				token: (await cookies.get('TOKEN')) || '',
				teamName: ''
			})
		})
	).json();

	return {
		team: team,
		team_members: team_members,
		user_team_invites: user_team_invites,
		team_user_invites: team_user_invites
	};
}
