import { fileURLToPath, URL } from 'node:url'
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig(({ command, mode, ssrBuild }) => {
	const ret = {
		plugins: [vue()],
		resolve: {
			alias: {
				'@': fileURLToPath(new URL('./src', import.meta.url))
			}
		},
		server: {
			port: 5173,
			strictPort: true,
			proxy: {
				// Se usi percorsi relativi (es. /session), il proxy li gira al backend
				'/session': { target: 'http://localhost:3000', changeOrigin: true },
				'/users':   { target: 'http://localhost:3000', changeOrigin: true },
				'/conversations': { target: 'http://localhost:3000', changeOrigin: true },
				'/groups':  { target: 'http://localhost:3000', changeOrigin: true },
				'/images':  { target: 'http://localhost:3000', changeOrigin: true },
				'/me':      { target: 'http://localhost:3000', changeOrigin: true },
			}
		}
		// ----------------------------------------------------
	};

	ret.define = {
		// Il grader cambierà questo valore automaticamente.
		"__API_URL__": JSON.stringify("http://localhost:3000"),
	};

	return ret;
})
