import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: {
    host: true,
    strictPort: true,
    port: 8081,
    watch: {
      usePolling: true,
      interval: 100, // Decrease polling interval for faster updates
      followSymlinks: false // Disable symlink following for better performance
    },
    hmr: {
      clientPort: 8081,
      host: 'localhost',
      protocol: 'ws',
      timeout: 10000,
      overlay: true,
    }
  },
  resolve: {
    extensions: ['.mjs', '.js', '.jsx', '.json', '.scss']
  },
  optimizeDeps: {
    esbuildOptions: {
      loader: {
        '.js': 'jsx'
      }
    }
  },
  css: {
    preprocessorOptions: {
      scss: {
        silenceDeprecations: ['mixed-decls', 'color-functions', 'global-builtin', 'import', 'legacy-js-api']
      },
    }
  },
  assetsInclude: ['**/*.woff2']
})

