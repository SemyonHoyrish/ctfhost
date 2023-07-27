<!-- Copyright (C) 2023  Semyon Hoyrish -->
<!-- For more details see the "LICENSE" file -->

<script lang="ts">
	let nickname = '';
	let email = '';
	let password = '';

	const doRegister = async () => {
		let j = await (
			await fetch('/auth/_request', {
				method: 'POST',
				body: JSON.stringify({
					nickname: nickname,
					email: email,
					password: password
				})
			})
		).json();

		if (j.success === false) {
			alert('ERROR: ' + j.data);
		} else {
			location.href = '/';
		}
	};
</script>

<div class="wrapper">
	<input type="text" class="nickname" placeholder="nickname" bind:value={nickname} />
	<input type="text" class="email" placeholder="email" bind:value={email} />
	<input type="password" class="password" placeholder="password" bind:value={password} />
	<button class="login" on:click={doRegister}>register</button>
</div>
