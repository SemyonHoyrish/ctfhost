<!-- Copyright (C) 2023  Semyon Hoyrish -->
<!-- For more details see the "LICENSE" file -->

<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';

	export let data;

	onMount(() => {
		if (
			!data.logged_in &&
			!(
				location.pathname.includes('/auth/') ||
				location.pathname === '/' ||
				location.pathname.includes('/leaderboard')
			)
		) {
			goto('/auth/login');
		}
	});
</script>

<header class="header">
	<h1>CTFhost</h1>
	<div class="nav">
		<a href="/" class="nav-link">home</a>
		<a href="/tasks" class="nav-link">tasks</a>
		<a href="/leaderboard" class="nav-link">leaderboard</a>
		<a href="/profile" class="nav-link">profile</a>
		{#if data.logged_in && data.user.success && JSON.parse(data.user.data).isAdmin}
			<a href="/admin" class="nav-link">admin</a>
		{/if}

		{#if !data.logged_in}
			<a href="/auth/login" class="nav-link">login</a>
		{/if}
	</div>
</header>

<div class="main">
	<slot />
</div>

<style>
	.header {
		display: flex;
		justify-content: space-around;
		align-items: center;
	}

	.header h1 {
		font-size: 32px;
	}

	.nav a {
		font-size: 18px;
		text-decoration: none;
		color: #000000;
		margin-right: 10px;
		cursor: pointer;
	}
	.nav a:last-child {
		margin-right: 0;
	}

	.main {
		max-width: 960px;
		margin-left: 50%;
		transform: translateX(-50%);
		margin-top: 30px;
	}
</style>
