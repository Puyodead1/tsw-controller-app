import fs from "node:fs";
import $RefParser from "@apidevtools/json-schema-ref-parser";

try {
  fs.writeFileSync(
    "./profile_json_schema/profile.complete.schema.json",
    JSON.stringify(
      await $RefParser.dereference("./profile_json_schema/profile.schema.json"),
      null,
      2,
    ),
  );
} catch (err) {
  console.error(err);
}
