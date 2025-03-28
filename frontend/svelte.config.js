import adapter from '@sveltejs/adapter-static';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	// Consult https://svelte.dev/docs/kit/integrations
	// for more information about preprocessors
	preprocess: vitePreprocess(),
	dynamicCompileOptions({ filename, compileOptions }) {
		// Dynamically set runes mode per Svelte file
		if (filename.startsWith('src/') && !compileOptions.runes) {
			return { runes: true };
		}
	},
	kit: {
		router: {
			type: 'hash'
		},
		adapter: adapter({
			// default options are shown. On some platforms
			// these options are set automatically — see below
			pages: '../backend/static',
			assets: '../backend/static',
			fallback: undefined,
			precompress: false,
			strict: true
		}),
		paths: {
			base: '/app',
			relative: true
		}
	}
};

export default config;
