import fs from 'node:fs';
import { dirname, resolve } from "node:path";
import { fileURLToPath } from "node:url";
import { defineConfig } from "vite";
import solid from 'vite-plugin-solid'
import tailwindcss from "@tailwindcss/vite";
import { generators } from "./vite.generate";

const __dirname = dirname(fileURLToPath(import.meta.url));

export default defineConfig(async () => {
  await generators.profile_schema()
  const docs = await generators.docs()

  return {
    build: {
      rollupOptions: {
        input: {
          main: resolve(__dirname, 'index.html'),
          'profile-builder': resolve(__dirname, 'profile-builder/index.html'),
          ...Object.fromEntries(docs.map(({ url, filepath }) => ([url, filepath])))
        }
      },
    },
    plugins: [tailwindcss(), solid()],
  }
});
