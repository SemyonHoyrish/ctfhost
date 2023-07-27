<!-- Copyright (C) 2023  Semyon Hoyrish -->
<!-- For more details see the "LICENSE" file -->

<script lang="ts">
	import TaskCard from '$lib/components/TaskCard.svelte';
	import type Task from '$lib/types/Task';
	import { TaskGetResponseToTaskArray } from '$lib/utils';

	export let data;

	let tasks: Array<Task> = new Array();
	if (data.success && data.data != 'null') {
		tasks = TaskGetResponseToTaskArray(JSON.parse(data.data));
	}
</script>

<a href="/admin/tasks/new">new task</a>

{#each tasks as t}
	<TaskCard task={t} />
{/each}
