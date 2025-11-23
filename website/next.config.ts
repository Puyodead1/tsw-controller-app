import fs from "node:fs"
import $RefParser from "@apidevtools/json-schema-ref-parser";
import type { NextConfig } from "next";

async function config(phase: string, { defaultConfig }: { defaultConfig: NextConfig }) {
  fs.writeFileSync(
    "./src/_profile-builder-json-schema/profile.complete.schema.json",
    JSON.stringify(
      await $RefParser.dereference("../profile-builder-schema/profile.schema.json"),
      null,
      2,
    ),
  );

  const nextConfig: NextConfig = {
    ...defaultConfig,
    output: 'export',
  };

  return nextConfig
}

export default config