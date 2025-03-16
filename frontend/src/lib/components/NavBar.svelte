<script lang="ts">
	import { profile } from '$lib/profile.svelte';

	const forceRunes = $state(null);

	function logout() {
		// redirect to our logout endpoint, non-SvelteKit, so has to be done via JS
		window.location.hash = '';
		window.location.pathname = '/auth/logout';
	}
</script>

<nav class="navbar navbar-expand-lg bg-primary" data-bs-theme="dark">
	<div class="container">
		<a class="navbar-brand" href="#/">Eventovátko</a>

		{#if $profile.status == 'logged-in'}
			<ul class="navbar-nav">
				<li class="nav-item dropdown">
					<a
						class="nav-link dropdown-toggle active"
						data-bs-toggle="dropdown"
						href="#/"
						role="button"
						aria-haspopup="true"
						aria-expanded="false"
					>
						<img
							alt="profilová fotka"
							src={$profile.pictureURL}
							class="rounded-circle me-2"
							style="width: 3em; height: 3em;"
						/>
						<span>{$profile.name}</span>
					</a>
					<div class="dropdown-menu">
						<button type="button" class="dropdown-item" onclick={logout}>Odhlásit se</button>
					</div>
				</li>
			</ul>
		{/if}
	</div>
</nav>
