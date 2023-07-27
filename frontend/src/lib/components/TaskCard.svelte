<!-- Copyright (C) 2023  Semyon Hoyrish -->
<!-- For more details see the "LICENSE" file -->

<script lang="ts">
	import { browser } from '$app/environment';
	import { goto } from '$app/navigation';
	import type Task from '$lib/types/Task';

	export let task: Task;
</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<!-- svelte-ignore a11y-no-static-element-interactions -->
<div
	style="background-color: {task.solved ? '#1b4b39' : '#1b204b'}"
	class="task-card"
	on:click={() => {
		if (task.id === undefined) {
			browser ?? alert('Error!');
		} else {
			goto(location.href + '/' + task.id.toString());
		}
	}}
>
	<h3 class="title">{task.title}</h3>
	<span class="category">{task.category}</span>
	<span class="points">{task.points}</span>
</div>

<style>
	.task-card {
		color: #ffffff;
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
		margin-bottom: 10px;
		padding: 10px;
		flex-basis: 28%;
		cursor: pointer;
		transition: box-shadow 0.2s;
	}

	.task-card:hover {
		box-shadow: 10px 10px 8px 0px rgba(34, 60, 80, 0.2);
	}
</style>
