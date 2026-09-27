#!/usr/bin/env node
import { StdioServerTransport } from "@modelcontextprotocol/sdk/server/stdio.js";
import { configFromEnv } from "./client.js";
import { DEFAULT_WEB_URL } from "./helpers.js";
import { createServer } from "./server.js";

async function main(): Promise<void> {
  const cfg = configFromEnv();
  const webBase = process.env.XWAVE_WEB_URL || DEFAULT_WEB_URL;
  const server = createServer(cfg, webBase);
  const transport = new StdioServerTransport();
  await server.connect(transport);
  // stdout is the JSON-RPC channel: log to stderr only.
  console.error(`xwave-mcp ready (api=${cfg.apiUrl})`);

  const shutdown = async () => {
    await server.close().catch(() => {});
    process.exit(0);
  };
  process.on("SIGINT", shutdown);
  process.on("SIGTERM", shutdown);
}

main().catch((err) => {
  console.error("xwave-mcp failed to start:", err);
  process.exit(1);
});
