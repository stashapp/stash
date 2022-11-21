import { defineConfig } from 'vite'
import tsconfigPaths from "vite-tsconfig-paths";
import viteCompression from 'vite-plugin-compression';
import { VitePWA } from "vite-plugin-pwa";

// https://vitejs.dev/config/
export default defineConfig({
  base: "",
  build: {
    outDir: 'build'
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
      filter: /\.(js|json|css|svg|md)$/i,
    }),
    VitePWA({
      includeAssets: ["favicon.ico", "apple-touch-icon.png"],
      registerType: "autoUpdate",
      devOptions: {
        enabled: true,
        type: "module",
      },
      strategies: 'injectManifest',
      srcDir: 'src',
      filename: 'serviceWorker.ts',
      injectManifest: {
        globIgnores: ['assets/**/*.gz'],
        globPatterns: ['assets/**'],
      },
      manifest: {
        name: "Stash: Porn Organizer",
        short_name: "Stash",
        display: "standalone",
        background_color: "#FFFFFF",
        theme_color: "#000000",
        icons: [
          {
            src: "pwa-192x192.png",
            sizes: "192x192",
            type: "image/png",
          },
          {
            src: "pwa-512x512.png",
            sizes: "512x512",
            type: "image/png",
          },
          {
            src: 'pwa-512x512.png',
            sizes: '512x512',
            type: 'image/png',
            purpose: 'any maskable'
          }
        ],
      },
    }),
  ],
})
