import { defineConfig } from 'vite'
import tsconfigPaths from "vite-tsconfig-paths";
import viteCompression from 'vite-plugin-compression';

// https://vitejs.dev/config/
export default defineConfig({
  base: "",
  build: {
    outDir: 'build',
  },
  optimizeDeps: {
    entries: "src/index.tsx"
  },
  server: {
    cors: false
  },
  publicDir: 'public',
  assetsInclude: ['**/*.md'],
  plugins: [tsconfigPaths(),
    viteCompression({
    algorithm: 'gzip',
    disable: false,
    deleteOriginFile: true,
    filter: /\.(js|json|css|svg|md)$/i
  })
],
})
