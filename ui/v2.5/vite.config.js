import { defineConfig } from 'vite'
import tsconfigPaths from "vite-tsconfig-paths";
//import envCompatible from "vite-plugin-env-compatible";
// import svgr from 'vite-plugin-svgr'
// import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  build: {
    outDir: 'build',
  },
  optimizeDeps: {
    entries: "src/index.tsx"
  },
  publicDir: 'public',
  assetsInclude: ['**/*.md'],
  plugins: [tsconfigPaths()],
  define: {
    'process.versions': {},
    'process.env': {}
  }
})