<!-- Copyright (C) 2023  Semyon Hoyrish -->
<!-- For more details see the "LICENSE" file -->

<script lang="ts">
	export let data;

	let content = '';

	if (data.success) {
		content = data.data;
	}
</script>

{#if !data.success}
	<h1>Welcome to the CTF!</h1>
{:else}
	{@html content}
{/if}
