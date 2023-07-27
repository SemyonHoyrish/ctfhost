<!-- Copyright (C) 2023  Semyon Hoyrish -->
<!-- For more details see the "LICENSE" file -->

<script lang="ts">
	import { browser } from '$app/environment';
	import type Settings from '$lib/types/Settings';

	export let data;

	let settings: Settings;

	if (!data.success) {
		browser ?? alert('ERROR! ' + data.data);
	} else {
		settings = JSON.parse(data.data);
	}

	const doSave = async () => {
		const c = confirm('Are you sure?');
		if (!c) return;
		const r = await (
			await fetch('/admin/settings/_set', {
				method: 'POST',
				body: JSON.stringify({ settings: settings })
			})
		).json();
		if (!r.success) {
			alert(r.data);
		}
		location.reload();
	};
</script>

{#if !data.success}
	<h1>ERROR</h1>
{:else}
	<div class="wrapper">
		<label>
			User Can Register
			<input type="checkbox" bind:checked={settings.UserCanRegister} />
		</label>
		<label>
			User Can Manage Team
			<input type="checkbox" bind:checked={settings.UserCanManageTeam} />
		</label>
		<label>
			User Can Manage Account
			<input type="checkbox" bind:checked={settings.UserCanManageAccount} />
		</label>
		<label>
			User Can View Tasks
			<input type="checkbox" bind:checked={settings.UserCanViewTasks} />
		</label>
		<label>
			Max Members Per Team (-1 for unlimited)
			<input type="number" bind:value={settings.MaxMembersPerTeam} />
		</label>
		<label class="main-page-content">
			Main page content:
			<textarea name="" id="" cols="30" rows="20" bind:value={settings.MainPageContent} />
		</label>
		<button class="save" on:click={doSave}>save</button>
	</div>
{/if}

<style>
	.wrapper {
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-content: center;
	}
	.main-page-content {
		display: flex;
		flex-direction: column;
	}
</style>
