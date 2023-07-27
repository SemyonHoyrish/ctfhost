<!-- Copyright (C) 2023  Semyon Hoyrish -->
<!-- For more details see the "LICENSE" file -->

<script lang="ts">
	import type User from '$lib/types/User';
	import type Team from '$lib/types/Team';
	import { onMount } from 'svelte';

	export let data;

	let user: User | null = null;
	let team: Team | null = null;
	let team_members: Array<User> = new Array<User>();

	class UserTeamInvite {
		inviteId = 0;
		teamId = 0;
		teamName = '';
	}
	class TeamUserInvite {
		inviteId = 0;
		userId = 0;
		userNickname = '';
	}

	let user_team_invites: Array<UserTeamInvite> = new Array<UserTeamInvite>();
	let team_user_invites: Array<TeamUserInvite> = new Array<TeamUserInvite>();

	let password: string = '';

	if (data.user && data.user.success) {
		user = JSON.parse(data.user.data);
	}
	if (data.team.success) {
		team = JSON.parse(data.team.data);
	}
	if (data.team_members.success) {
		team_members = JSON.parse(data.team_members.data);
	}
	if (data.user_team_invites.success) {
		user_team_invites = JSON.parse(data.user_team_invites.data);
	}
	if (data.team_user_invites.success) {
		team_user_invites = JSON.parse(data.team_user_invites.data);
	}

	let newTeamName: string = team?.name || '';

	const isUserCaptain = user !== null && team !== null && user.id === team.captainId;

	const removeAccount = async () => {
		const c = confirm('Are you sure?');
		if (!c) return;
		const r = await (
			await fetch('/profile/_remove-account', {
				method: 'POST'
			})
		).json();
		if (!r.success) {
			alert(r.data);
		}
		location.reload();
	};
	const removeTeam = async () => {
		const c = confirm('Are you sure?');
		if (!c) return;
		const r = await (
			await fetch('/profile/_remove-team', {
				method: 'POST'
			})
		).json();
		if (!r.success) {
			alert(r.data);
		}
		location.reload();
	};
	const createTeam = async () => {
		const name = prompt('Enter team name');
		if (name === '') {
			alert('Name cannot be empty');
			return;
		}
		const r = await (
			await fetch('/profile/_create-team', {
				method: 'POST',
				body: JSON.stringify({ newTeamName: name })
			})
		).json();
		if (!r.success) {
			alert(r.data);
		}
		location.reload();
	};

	const updateUserInfo = async () => {
		const r = await (
			await fetch('/profile/_edit-user', {
				method: 'POST',
				body: JSON.stringify({ nickname: user?.nickname, email: user?.email, password: password })
			})
		).json();
		if (!r.success) {
			alert(r.data);
		}
		location.reload();
	};

	const changeTeamName = async () => {
		const r = await (
			await fetch('/profile/_edit-team', {
				method: 'POST',
				body: JSON.stringify({ teamName: team?.name, newTeamName: newTeamName })
			})
		).json();
		if (!r.success) {
			alert(r.data);
		}
		location.reload();
	};

	const inviteByEmail = async () => {
		const p = prompt('Enter email:');
		const r = await (
			await fetch('/profile/_team-invite/invite', {
				method: 'POST',
				body: JSON.stringify({
					userEmail: p,
					userNickname: ''
				})
			})
		).json();
		if (!r.success) {
			alert(r.data);
		}
		location.reload();
	};
	const invitebyNickname = async () => {
		const p = prompt('Enter nickname:');
		const r = await (
			await fetch('/profile/_team-invite/invite', {
				method: 'POST',
				body: JSON.stringify({
					userEmail: '',
					userNickname: p
				})
			})
		).json();
		if (!r.success) {
			alert(r.data);
		}
		location.reload();
	};
	const removeInvite = async (id: number) => {
		if (!confirm('Are you sure to remove invite?')) {
			return;
		}
		const r = await (
			await fetch('/profile/_team-invite/remove', {
				method: 'POST',
				body: JSON.stringify({
					inviteId: id
				})
			})
		).json();
		if (!r.success) {
			alert(r.data);
		}
		location.reload();
	};

	const acceptInvite = async (id: number) => {
		if (!confirm('Are you sure to accept invite?')) {
			return;
		}
		const r = await (
			await fetch('/profile/_user-invite/accept', {
				method: 'POST',
				body: JSON.stringify({
					inviteId: id
				})
			})
		).json();
		if (!r.success) {
			alert(r.data);
		}
		location.reload();
	};
	const declineInvite = async (id: number) => {
		if (!confirm('Are you sure to decline invite?')) {
			return;
		}
		const r = await (
			await fetch('/profile/_user-invite/decline', {
				method: 'POST',
				body: JSON.stringify({
					inviteId: id
				})
			})
		).json();
		if (!r.success) {
			alert(r.data);
		}
		location.reload();
	};

	const leaveTeam = async () => {
		if (!confirm('Are you sure to leave this team?')) {
			return;
		}
		const r = await (
			await fetch('/profile/_user-leave', {
				method: 'POST',
				body: JSON.stringify({})
			})
		).json();
		if (!r.success) {
			alert(r.data);
		}
		location.reload();
	};

	const removeMember = async (id: number) => {
		if (!confirm('Are you sure to remove member?')) {
			return;
		}
		const r = await (
			await fetch('/profile/_team-remove_member', {
				method: 'POST',
				body: JSON.stringify({ userId: id })
			})
		).json();
		if (!r.success) {
			alert(r.data);
		}
		location.reload();
	};

	const userLogout = async () => {
		const r = await (
			await fetch('/profile/_user-logout', {
				method: 'POST',
				body: JSON.stringify({})
			})
		).json();
		if (!r.success) {
			alert(r.data);
		}
		location.reload();
	};
	const userLogoutAll = async () => {
		const r = await (
			await fetch('/profile/_user-logout_all', {
				method: 'POST',
				body: JSON.stringify({})
			})
		).json();
		if (!r.success) {
			alert(r.data);
		}
		location.reload();
	};

	let sessions = 1;

	onMount(() => {
		(async () => {
			const r = await (
				await fetch('/profile/_user-get_sessions_count', {
					method: 'POST',
					body: JSON.stringify({})
				})
			).json();
			if (!r.success) {
				sessions = 0;
			} else {
				sessions = parseInt(r.data);
			}
		})();
	});
</script>

<div class="wrapper">
	<div class="block">
		<h3 class="block-title">User profile</h3>
		{#if user === null}
			<p>User not found</p>
		{:else}
			<div class="info">
				<input type="text" class="nickname" bind:value={user.nickname} />
				<input type="text" class="email" bind:value={user.email} />
				<label>
					Change password:
					<input type="password" class="password" bind:value={password} />
				</label>
				<button class="save" on:click={updateUserInfo}>save</button>
				<div class="session">
					<span class="sessions-count">You have {sessions} opened sessions</span>
					<button class="logout" on:click={userLogout}>logout</button>
					<button class="close-session" on:click={userLogoutAll}>close all sessions</button>
				</div>
				<br />
				<button class="remove-account" on:click={removeAccount}>delete account</button>
			</div>
		{/if}
	</div>
	<div class="block">
		<h3 class="block-title">Team profile</h3>
		{#if team === null}
			<p>Team not found</p>
			<button class="create-team" on:click={createTeam}>create a team</button>
		{:else}
			<div class="info">
				{#if isUserCaptain}
					<input type="text" class="name" bind:value={newTeamName} />
				{:else}
					<span class="name" style="font-weight: bold;">{team.name}</span>
				{/if}
				<span class="score"><b>Score:</b> {team.score}</span>
				<span><b>Team members:</b></span>
				<div class="members">
					{#each team_members as member}
						<div class="member">
							<span class="member-nickname">{member.nickname}</span>
							{#if isUserCaptain}
								<button
									class="remove"
									disabled={member.id === user?.id}
									on:click={() => {
										removeMember(member.id);
									}}
								>
									remove
								</button>
							{/if}
						</div>
					{/each}
				</div>
				<div class="invites">
					<h4>Invited:</h4>
					{#if isUserCaptain}
						<span>invite user</span>
						<div style="display:flex;">
							<button on:click={invitebyNickname}>by nickname</button>
							<button on:click={inviteByEmail}>by email</button>
						</div>
					{/if}
					{#each team_user_invites as i}
						<div class="invite">
							<div class="invited" style="display:flex;">
								<span class="user-nickname">{i.userNickname}</span>
								{#if isUserCaptain}
									<button
										class="remove-invite"
										on:click={() => {
											removeInvite(i.inviteId);
										}}
									>
										remove
									</button>
								{/if}
							</div>
						</div>
					{/each}
				</div>

				{#if isUserCaptain}
					<button class="save" on:click={changeTeamName}>save</button>
					<button class="remove-team" on:click={removeTeam}>delete team</button>
				{:else}
					<button class="leave-team" on:click={leaveTeam}>leave team</button>
				{/if}
			</div>
		{/if}
	</div>
	<div class="block">
		<h3 class="block-title">Invites</h3>
		{#each user_team_invites as i}
			<div class="invite">
				<div class="invited" style="display:flex;">
					<span class="user-nickname">{i.teamName}</span>
					<button
						class="accept"
						on:click={() => {
							acceptInvite(i.inviteId);
						}}
					>
						accept
					</button>
					<button
						class="decline"
						on:click={() => {
							declineInvite(i.inviteId);
						}}
					>
						decline
					</button>
				</div>
			</div>
		{/each}
	</div>
</div>

<style>
	* {
		box-sizing: border-box;
	}
	.block .info {
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
	}
	.block .info label {
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
	}
	.block .info * {
		width: 100%;
	}
	.block {
		max-width: 360px;
	}
	.block * {
		margin-bottom: 5px;
	}

	.members .member {
		display: flex;
		align-items: center;
	}
</style>
