<script lang="ts">
	function googleLogin() {
		// redirect to our Google login endpoint, JS because
		// Svelte likes to control links without domain
		window.location.pathname = '/auth/google';
	}
</script>

<div>
	<h1 class="mb-4 text-4xl font-extrabold">Eventovátko</h1>
	<p class="mb-2">
		<span class="font-bold">Aktuálně je nástroj ve vývoji a obsahuje nejspíš spoustu bugů.</span> Webová
		appka na zjednodušení organizování akcí.
	</p>
	<button
		onclick={googleLogin}
		type="button"
		class="text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 me-2 mb-2 dark:bg-blue-600 dark:hover:bg-blue-700 focus:outline-none dark:focus:ring-blue-800"
		>Přihlásit se přes Google účet</button
	>
</div>
