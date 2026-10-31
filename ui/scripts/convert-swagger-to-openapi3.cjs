/**
 * Converts Swagger 2.0 spec to OpenAPI 3.x and generates TypeScript types.
 * Uses a temp file for the converted spec (deleted after generation).
 */
const path = require("path");
const fs = require("fs");
const os = require("os");
const { spawnSync } = require("child_process");
const swagger2openapi = require("swagger2openapi");

const SWAGGER_PATH = path.join(__dirname, "../../service/docs/swagger.json");
const TYPES_DIR = path.join(__dirname, "../types");
const OUT_DTS = path.join(TYPES_DIR, "xwave-api.d.ts");

const spec = JSON.parse(fs.readFileSync(SWAGGER_PATH, "utf8"));

swagger2openapi
  .convert(spec, { patch: true })
  .then((options) => {
    const tmpPath = path.join(os.tmpdir(), `openapi3-${Date.now()}.json`);
    fs.writeFileSync(tmpPath, JSON.stringify(options.openapi, null, 2));
    try {
      if (!fs.existsSync(TYPES_DIR)) {
        fs.mkdirSync(TYPES_DIR, { recursive: true });
      }
      const result = spawnSync(
        "npx",
        ["openapi-typescript", tmpPath, "-o", OUT_DTS],
        { stdio: "inherit", cwd: path.join(__dirname, "..") }
      );
      if (result.status !== 0) process.exit(result.status || 1);
    } finally {
      try {
        fs.unlinkSync(tmpPath);
      } catch (_) {}
    }
  })
  .catch((err) => {
    console.error("Conversion failed:", err);
    process.exit(1);
  });
