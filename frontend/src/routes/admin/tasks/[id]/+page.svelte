<!-- Copyright (C) 2023  Semyon Hoyrish -->
<!-- For more details see the "LICENSE" file -->

<script lang="ts">
	import { browser } from '$app/environment';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import Task from '$lib/types/Task';
	import { TaskGetResponseToTask, TaskToBackendTask } from '$lib/utils';

	export let data;

	const id_p: string = $page.params.id;

	let task: Task = new Task(0, '', 0, '');

	if (id_p === undefined) {
		browser ?? alert('ERROR!');
	} else if (id_p === 'new') {
		task = new Task(0, '', 0, '');
	} else {
		if (!data.success) {
			browser ?? alert('ERROR: ' + data.data);
		} else {
			task = TaskGetResponseToTask(JSON.parse(data.data));
		}
	}

	const doSave = async () => {
		const c = confirm('Are you sure?');
		if (!c) return;
		const r = await (
			await fetch('/admin/tasks/_request', {
				method: 'POST',
				body: JSON.stringify({
					action: id_p === 'new' ? 'new' : 'edit',
					task: TaskToBackendTask(task)
				})
			})
		).json();
		if (!r.success) {
			alert(r.data);
		}
		if (id_p === 'new') {
			goto('/admin/tasks');
		} else {
			location.reload();
		}
	};
	const doRemove = async () => {
		const c = confirm('Are you sure?');
		if (!c) return;
		const r = await (
			await fetch('/admin/tasks/_request', {
				method: 'POST',
				body: JSON.stringify({
					action: 'remove',
					task: TaskToBackendTask(task)
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
	<div class="wrapper">
		<h3 class="title" contenteditable bind:textContent={task.title} />
		<!-- @ToDo: Description editor -->
		Description_EDITOR
		<input type="number" class="points" bind:value={task.points} />
		<input type="text" class="category" placeholder="category" bind:value={task.category} />
		<input type="text" class="flag" placeholder="flag" bind:value={task.flag} />
		<label>
			visible:
			<input type="checkbox" bind:checked={task.visible} />
		</label>
		<button class="save" on:click={doSave}>save</button>
		<button class="remove" on:click={doRemove}>remove</button>
	</div>
{/if}

<style>
	.title {
		border: 1px solid grey;
		padding: 2px;
	}
</style>
