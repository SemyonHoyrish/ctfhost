<!-- Copyright (C) 2023  Semyon Hoyrish -->
<!-- For more details see the "LICENSE" file -->

<script lang="ts">
	import { page } from '$app/stores';
	import { TaskGetResponseToTask } from '$lib/utils';
	import type Task from '$lib/types/Task';

	export let data;

	const id_s = $page.params.id;

	let task: Task | null = null;

	if (id_s === undefined) {
		alert('ERROR!');
	} else if (data.success) {
		const id = parseInt(id_s);

		task = TaskGetResponseToTask(JSON.parse(data.data));
	}

	let flag = '';

	const submit = async () => {
		const r = await (
			await fetch(`/tasks/${id_s}/_submit`, {
				method: 'POST',
				body: JSON.stringify({
					flag: flag
				})
			})
		).json();
		if (!r.success) {
			alert(r.data);
		}
		location.reload();
	};
</script>

{#if task === null}
	<h1>ERROR</h1>
{:else}
	{#if task.solved}
		<h2>Congratulations!!!</h2>
	{/if}
	<div class="task">
		<h2 class="title">{task.title}</h2>
		<span class="point">{task.points}</span>
		<span class="category">{task.category}</span>
		<div class="description">{task.getDescriptionHTML()}</div>
		<div class="flag-submit">
			<input
				type="text"
				class="flag"
				placeholder={task.solved ? 'submitted' : 'flag'}
				disabled={task.solved}
				bind:value={flag}
			/>
			<button class="submit" on:click={submit} disabled={task.solved}>submit</button>
		</div>
	</div>
{/if}
