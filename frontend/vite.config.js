import { dirname, resolve } from 'node:path'
import { fileURLToPath } from 'node:url'
import { defineConfig } from 'vite'

const __dirname = dirname(fileURLToPath(import.meta.url))

export default defineConfig({
    resolve: {
	alias: {
		'~bootstrap': resolve(__dirname, 'node_modules/bootstrap'),
	}
    },
    build: {
	rollupOptions: {
	    input: {
		main: resolve(__dirname, 'index.html'),
		prefs: resolve(__dirname, 'prefs.html'),
		about: resolve(__dirname, 'about.html')
	    },
	},
    },
})

