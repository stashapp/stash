import { defineConfig } from "vite";
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
    cors: false,
  },
  publicDir: "public",
  assetsInclude: ["**/*.md"],
  plugins: [
    tsconfigPaths(),
    // compress everything ahead of time (threshold: 0) as anything
    // skipped here will just be compressed on-the-fly by the server
    viteCompression({
      algorithm: "gzip",
      deleteOriginFile: true,
      threshold: 0, 
      filter: /\.(js|json|css|svg|md)$/i,
    }),
    viteCompression({
      algorithm: "brotliCompress",
      deleteOriginFile: true,
      threshold: 0,
      filter: /\.(js|json|css|svg|md)$/i,
    }),
  ],
});
