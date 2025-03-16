<script lang="ts">
	function googleLogin() {
		// redirect to our Google login endpoint, JS because
		// Svelte likes to control links without domain
		window.location.pathname = '/auth/google';
	}
</script>

<p>Webová appka na zjednodušení organizování akcí.</p>

<p><b>Aktuálně je nástroj ve vývoji a obsahuje nejspíš spoustu bugů.</b></p>

<button onclick={googleLogin} type="button" class="btn btn-primary"
	>Přihlásit se přes Google účet</button
>
