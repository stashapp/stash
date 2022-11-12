import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import tsconfigPaths from "vite-tsconfig-paths";
import viteCompression from "vite-plugin-compression";

// https://vitejs.dev/config/
export default defineConfig({
  base: "",
  build: {
    outDir: "build",
    reportCompressedSize: false,
  },
  optimizeDeps: {
    entries: "src/index.tsx",
  },
  server: {
    port: 3000,
    cors: false,
  },
  publicDir: "public",
  assetsInclude: ["**/*.md"],
  plugins: [
    react({
      babel: {
        compact: true,
      },
    }),
    tsconfigPaths(),
    viteCompression({
      algorithm: "gzip",
      disable: false,
      deleteOriginFile: true,
      filter: /\.(js|json|css|svg|md)$/i,
    }),
  ],
});
