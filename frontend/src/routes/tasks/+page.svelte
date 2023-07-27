<!-- Copyright (C) 2023  Semyon Hoyrish -->
<!-- For more details see the "LICENSE" file -->

<script lang="ts">
	import { browser } from '$app/environment';
	import TaskCard from '$lib/components/TaskCard.svelte';
	import type Task from '$lib/types/Task';
	import { TaskGetResponseToTaskArray } from '$lib/utils';

	export let data;

	let tasks: Array<Task> = new Array();
	if (data.success && data.data != 'null') {
		tasks = TaskGetResponseToTaskArray(JSON.parse(data.data));
	} else if (data.data != 'null') {
		if (browser) alert(data.data);
	}
</script>

<div class="wrapper">
	{#if !data.success}
		<h3>ERROR</h3>
	{:else}
		{#each tasks as t}
			<TaskCard task={t} />
		{/each}
	{/if}
</div>

<style>
	.wrapper {
		display: flex;
		flex-wrap: wrap;
		justify-content: space-around;
	}
</style>
