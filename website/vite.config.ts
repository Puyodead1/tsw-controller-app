import { dirname, resolve } from "node:path";
import { fileURLToPath } from "node:url";
import { defineConfig } from "vite";
import solid from 'vite-plugin-solid'
import tailwindcss from "@tailwindcss/vite";

const __dirname = dirname(fileURLToPath(import.meta.url));

export default defineConfig({
  build: {
    rollupOptions: {
      input: {
        main: resolve(__dirname, "index.html"),
        "profile-builder": resolve(__dirname, "profile-builder/index.html"),
      },
    },
  },
  plugins: [tailwindcss(), solid()],
});
