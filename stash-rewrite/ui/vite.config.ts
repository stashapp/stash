import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";
import path from "path";

export default defineConfig({
  plugins: [react(), tailwindcss()],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
  server: {
    port: 5173,
    proxy: {
      "/api": {
        target: "http://localhost:9999",
        changeOrigin: true,
      },
      "/hubs": {
        target: "http://localhost:9999",
        changeOrigin: true,
        ws: true,
      },
    },
  },
  build: {
    outDir: "../src/Stash.Api/wwwroot",
    emptyOutDir: true,
    rollupOptions: {
      output: {
        manualChunks: {
          vendor: ["react", "react-dom", "@tanstack/react-query"],
          icons: ["lucide-react"],
          signalr: ["@microsoft/signalr"],
        },
      },
    },
  },
});
