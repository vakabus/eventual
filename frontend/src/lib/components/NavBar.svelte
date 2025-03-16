<script lang="ts">
	import { profile } from '$lib/profile.svelte';
	import { page } from '$app/state';
	import { events } from '$lib/events.svelte';

	function logout() {
		// redirect to our logout endpoint, non-SvelteKit, so has to be done via JS
		window.location.hash = '';
		window.location.pathname = '/auth/logout';
	}

	function breadcrumbLabel(url: string) {
		const universalUrl = url.slice(2).replace(/[0-9]+/g, 'X');
		return (
			{
				'event/': 'Akce',
				'event/X/': events.nameFromId(url.slice(8).replace(/\/.*/, '')),
				'event/X/participant/': 'Účastníci',
				'event/X/template/': 'Vzory mailů',
				'event/X/template/X/': 'Úprava',
				'event/X/edit/': 'Upravit',
				'event/new/': 'Nová akce'
			}[universalUrl] ?? 'FIXME'
		);
	}

	const breadcrumbs = $derived(
		(function () {
			const hash = page.url.hash;

			if (!hash.includes('event')) {
				return [];
			}

			const parts = hash.split('/').slice(1);
			const breadcrumbs = [];
			for (let i = 0; i < parts.length; i++) {
				if (parts[i] == '') {
					continue;
				}

				const link = `#/${parts.slice(0, i + 1).join('/')}/`;
				breadcrumbs.push({ link: link, label: breadcrumbLabel(link) });
			}
			return breadcrumbs;
		})()
	);
</script>

<nav class="navbar navbar-expand-lg bg-primary" data-bs-theme="dark">
	<div class="container d-flex flex-row align-items-center">
		<a class="navbar-brand" href="#/">Eventovátko</a>
		<nav class="me-auto" aria-label="breadcrumb">
			<ol class="breadcrumb mb-0">
				{#each breadcrumbs as bc}
					<li class="breadcrumb-item"><a href={bc.link}>{bc.label}</a></li>
				{/each}
			</ol>
		</nav>

		{#if $profile.status == 'logged-in'}
			<ul class="navbar-nav">
				<li class="nav-item dropdown d-flex align-items-center">
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
							style="width: fit-content; height: 2em;"
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
