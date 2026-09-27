// End-to-end smoke test: spawn the built server over stdio and speak raw JSON-RPC.
import { spawn } from "node:child_process";
import { fileURLToPath } from "node:url";
import { dirname, join } from "node:path";
import readline from "node:readline";

const here = dirname(fileURLToPath(import.meta.url));
const child = spawn(process.execPath, [join(here, "..", "dist", "index.js")], { stdio: ["pipe", "pipe", "inherit"] });
const rl = readline.createInterface({ input: child.stdout });
const pending = new Map();
rl.on("line", (line) => {
  const msg = JSON.parse(line);
  if (msg.id !== undefined && pending.has(msg.id)) { pending.get(msg.id)(msg); pending.delete(msg.id); }
});
let nextId = 1;
const send = (method, params) => new Promise((resolve, reject) => {
  const id = nextId++;
  const t = setTimeout(() => reject(new Error(`timeout: ${method}`)), 60000);
  pending.set(id, (m) => { clearTimeout(t); resolve(m); });
  child.stdin.write(JSON.stringify({ jsonrpc: "2.0", id, method, params }) + "\n");
});
const notify = (method, params) => child.stdin.write(JSON.stringify({ jsonrpc: "2.0", method, params }) + "\n");
const firstText = (r) => r.result?.content?.[0]?.text ?? JSON.stringify(r.error ?? r);

let failed = false;
const check = (cond, label) => { console.log(`${cond ? "PASS" : "FAIL"} ${label}`); if (!cond) failed = true; };

try {
  const init = await send("initialize", { protocolVersion: "2025-06-18", capabilities: {}, clientInfo: { name: "smoke", version: "0" } });
  check(init.result?.serverInfo?.name === "xwave", `initialize -> ${init.result?.serverInfo?.name} ${init.result?.serverInfo?.version} (protocol ${init.result?.protocolVersion})`);
  notify("notifications/initialized", {});

  const tools = await send("tools/list", {});
  const names = tools.result.tools.map((t) => t.name);
  check(names.length === 6, `tools/list -> ${names.join(", ")}`);

  const res = await send("resources/read", { uri: "xwave://catalogs" });
  check(res.result?.contents?.[0]?.text.includes("eROSITA"), "resources/read xwave://catalogs");

  const cone = await send("tools/call", { name: "cone_search", arguments: { ra: 10.6847, dec: 41.269, radius_arcsec: 5 } });
  console.log(firstText(cone));
  check(!cone.result.isError && cone.result.structuredContent.count > 0, "cone_search M31");

  const rn = await send("tools/call", { name: "resolve_name", arguments: { name: "Betelgeuse" } });
  console.log(firstText(rn));
  check(Math.abs(rn.result.structuredContent?.ra - 88.7929) < 0.01, "resolve_name Betelgeuse");

  const sbn = await send("tools/call", { name: "search_by_name", arguments: { name: "Betelgeuse", radius_arcsec: 5 } });
  console.log(firstText(sbn));
  check(!sbn.result.isError, "search_by_name Betelgeuse");

  const xm = await send("tools/call", { name: "cross_match_list", arguments: { positions: [
    { name: "M31 core", ra: 10.6847, dec: 41.269 }, { name: "Betelgeuse", ra: 88.7929, dec: 7.4071 }, { name: "empty", ra: 0.1, dec: 0.1 }] } });
  console.log(firstText(xm));
  check(xm.result.structuredContent?.matched === 2, "cross_match_list 3 positions");

  const obj = await send("tools/call", { name: "get_object", arguments: { id: "0098p408_ac51-043708", catalog: "allwise" } });
  console.log(firstText(obj));
  check(!obj.result.isError, "get_object AllWISE");

  const lc = await send("tools/call", { name: "get_lightcurve", arguments: { ra: 10.6847, dec: 41.269, radius_arcsec: 5, max_points: 3 } });
  console.log(firstText(lc));
  check(lc.result.structuredContent?.counts?.detections > 0, "get_lightcurve M31");

  const bad = await send("tools/call", { name: "cone_search", arguments: { ra: 400, dec: 0, radius_arcsec: 5 } });
  console.log(firstText(bad));
  check(bad.result?.isError === true || bad.error, "out-of-range ra rejected");
} catch (e) {
  console.error(e); failed = true;
} finally {
  child.kill();
}
process.exit(failed ? 1 : 0);
