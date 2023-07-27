<!-- Copyright (C) 2023  Semyon Hoyrish -->
<!-- For more details see the "LICENSE" file -->

<script lang="ts">
	import type Team from '$lib/types/Team';

	export let data;

	let d: Array<Team> = new Array<Team>();

	if (data.success) {
		d = JSON.parse(data.data);
	}

	let prev_score = -1;
	let current_place = 0;

	const calcPlace = (score: number) => {
		if (prev_score != score) {
			current_place++;
		}
		prev_score = score;
		return current_place;
	};
</script>

<div class="wrapper">
	<h3>Leaderboard</h3>
	{#if !data.success}
		<h3>ERROR</h3>
	{:else if data.data != 'null'}
		<div class="team">
			<span class="place"><b>place</b></span>
			<span class="name"><b>team name</b></span>
			<span class="score"><b>score</b></span>
		</div>
		{#each d as team}
			<div class="team">
				<span class="place">{calcPlace(team.score)}</span>
				<span class="name">{team.name}</span>
				<span class="score">{team.score}</span>
			</div>
		{/each}
	{/if}
</div>

<style>
	.team {
		display: flex;
		justify-content: space-between;
		align-items: center;
		border-bottom: 1px solid grey;
		padding: 10px;
	}
</style>
