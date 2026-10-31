#!/usr/bin/env node
/**
 * Builds app/lib/constants/gamma-sources.json: the sources with a gamma-ray or
 * hard X-ray light curve that the object page can show, each with the
 * position of its associated counterpart so a match is a plain cone check.
 *
 * - Fermi-LAT Light Curve Repository: the 4FGL sources it monitors
 *   (variability index > 21.67) that have an association. The counterpart's
 *   position comes from NED, since 4FGL only gives the gamma-ray position.
 * - Swift-BAT Hard X-ray Transient Monitor: every listed source; the monitor
 *   table already gives its counterpart's position.
 *
 * Usage: node scripts/build-gamma-sources.cjs   (takes a few minutes)
 */

const fs = require("fs");
const path = require("path");

const LCR_LIST =
  "https://fermi.gsfc.nasa.gov/ssc/data/access/lat/LightCurveRepository/queryDB.php?typeOfRequest=SourceList&catalog=4FGL";
const LCR_MIN_VARIABILITY = 21.67;
const BAT_PAGE = "https://swift.gsfc.nasa.gov/results/transients/";
const NED_LOOKUP = "https://ned.ipac.caltech.edu/srs/ObjectLookup";
const OUT = path.join(__dirname, "../app/lib/constants/gamma-sources.json");
const CONCURRENCY = 6;

const round = (x, digits) => Number(x.toFixed(digits));

async function nedPosition(name) {
  const resp = await fetch(NED_LOOKUP, {
    method: "POST",
    body: new URLSearchParams({ json: JSON.stringify({ name: { v: name } }) }),
    signal: AbortSignal.timeout(30_000),
  });
  if (!resp.ok) return null;
  const pos = (await resp.json())?.Preferred?.Position;
  return pos && Number.isFinite(pos.RA) && Number.isFinite(pos.Dec)
    ? { ra: pos.RA, dec: pos.Dec }
    : null;
}

async function mapLimit(items, limit, fn) {
  const out = new Array(items.length);
  let next = 0;
  async function worker() {
    while (next < items.length) {
      const i = next++;
      out[i] = await fn(items[i], i);
    }
  }
  await Promise.all(Array.from({ length: limit }, worker));
  return out;
}

/** Positions already resolved by a previous run, keyed by association name. */
function previousPositions() {
  try {
    const prev = JSON.parse(fs.readFileSync(OUT, "utf8"));
    return new Map(prev.fermi.map((s) => [s.assoc, { ra: s.ra, dec: s.dec }]));
  } catch {
    return new Map();
  }
}

async function fermiSources() {
  const list = await (await fetch(LCR_LIST)).json();
  const monitored = list.filter(
    (s) =>
      Number(s.Variability_Index) > LCR_MIN_VARIABILITY && s.ASSOC1.trim()
  );
  const positions = previousPositions();
  // NED drops some requests under load; unresolved names get a few slower
  // passes before they are left out.
  for (const [pass, limit] of [CONCURRENCY, 2, 1].entries()) {
    const missing = [
      ...new Set(monitored.map((s) => s.ASSOC1.trim())),
    ].filter((name) => !positions.has(name));
    if (missing.length === 0) break;
    console.log(`  NED pass ${pass + 1}: ${missing.length} names`);
    await mapLimit(missing, limit, async (name) => {
      const pos = await nedPosition(name).catch(() => null);
      if (pos) positions.set(name, pos);
    });
  }
  const found = monitored.flatMap((s) => {
    const pos = positions.get(s.ASSOC1.trim());
    return pos
      ? [
          {
            name: s.Source_Name,
            assoc: s.ASSOC1.trim(),
            cls: s.CLASS1.trim().toLowerCase(),
            ra: round(pos.ra, 6),
            dec: round(pos.dec, 6),
          },
        ]
      : [];
  });
  console.log(
    `Fermi LCR: ${found.length} of ${monitored.length} associations resolved`
  );
  return found;
}

async function batSources() {
  const html = await (await fetch(BAT_PAGE)).text();
  const strip = (cell) => cell.replace(/<[^>]+>/g, "").trim();
  const sources = [];
  for (const row of html.split(/<tr[^>]*>/).slice(1)) {
    const cells = row.split(/<td[^>]*>/).slice(1).map((c) => c.split("</td>")[0] ?? "");
    const link = /<a href="([^"]+)">/.exec(cells[1] ?? "");
    const ra = Number(strip(cells[2] ?? ""));
    const dec = Number(strip(cells[3] ?? ""));
    if (!link?.[1] || !Number.isFinite(ra) || !Number.isFinite(dec)) continue;
    sources.push({
      name: strip(cells[1] ?? ""),
      path: link[1],
      type: strip(cells[5] ?? ""),
      ra,
      dec,
    });
  }
  console.log(`Swift-BAT monitor: ${sources.length} sources`);
  return sources;
}

async function main() {
  const [fermi, bat] = await Promise.all([fermiSources(), batSources()]);
  fs.mkdirSync(path.dirname(OUT), { recursive: true });
  fs.writeFileSync(
    OUT,
    JSON.stringify({
      generated: new Date().toISOString().slice(0, 10),
      fermi,
      bat,
    }) + "\n"
  );
  console.log(`Wrote ${path.relative(process.cwd(), OUT)}`);
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});
