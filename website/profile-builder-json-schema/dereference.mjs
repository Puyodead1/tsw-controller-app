import fs from "node:fs";
import $RefParser from "@apidevtools/json-schema-ref-parser";

try {
  fs.writeFileSync(
    "./profile-builder-json-schema/profile.complete.schema.json",
    JSON.stringify(
      await $RefParser.dereference("./profile-builder-json-schema/profile.schema.json"),
      null,
      2,
    ),
  );
} catch (err) {
  console.error(err);
}
