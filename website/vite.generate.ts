import fs from "node:fs";
import $RefParser from "@apidevtools/json-schema-ref-parser";
import path from "node:path";
import { marked } from "marked";

const navHtml = `<nav class="navbar bg-base-100 shadow-sm">
  <div class="flex-1">
    <ul class="menu menu-horizontal px-1">
      <li><a class="link link-hover" href="/">Home</a></li>
      <li><a class="link link-hover" href="/profile-builder/index.html">Profile Builder</a></li>
      <li><a class="link link-hover" href="/docs/index.html">Documentation</a></li>
    </ul>
  </div>
</nav>`

export const generators = {
  profile_schema: async () => {
    console.info('Generating profile schema...')
    fs.writeFileSync(
      "./profile-builder-json-schema/profile.complete.schema.json",
      JSON.stringify(
        await $RefParser.dereference("./profile-builder-json-schema/profile.schema.json"),
        null,
        2,
      ),
    );
  },
  docs: async () => {
    const pages = [
      { url: 'index', title: 'Index', filepath: path.resolve("../DOCS.md") },
      { url: 'creating-a-profile-from-scratch', title: 'Creating a profile from scratch', filepath: path.resolve("../CREATING_PROFILE_QUICKSTART.md") },
      { url: 'profile-explainer', title: 'Profile explainer', filepath: path.resolve("../PROFILE_EXPLAINER.md") },
      { url: 'steam-input-setup', title: 'Setting up steam input', filepath: path.resolve("../STEAM_INPUT_SETUP.md") }
    ]
    for (const page of pages) {
      const contentHtml = await marked.parse(fs.readFileSync(page.filepath, { encoding: 'utf8' }), {
        async: true,
        walkTokens(token) {
          if (token.type === 'image' && token.href.startsWith('./images')) {
            token.href = `https://raw.githubusercontent.com/LiamMartens/tsw-controller-app/refs/heads/feat/tsw-api/${token.href.replace('./', '')}`
          }
        },
      })
      fs.writeFileSync(path.resolve(`./docs/${page.url}.html`), `<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>TSW Controller App - Documentation - ${page.title}</title>
    <link href="../style.css" rel="stylesheet" />
  </head>
  <body>
    ${navHtml}
    <main class="prose mx-auto max-w-4xl px-8 my-8">${contentHtml}</main>
  </body>
</html>`)
    }
    return pages
  }
}

