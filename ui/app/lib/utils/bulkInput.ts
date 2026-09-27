/**
 * Bulk cross-match input parsing and result shaping.
 *
 * Accepts whatever people paste out of a paper, a spreadsheet or TOPCAT:
 * comma / tab / semicolon / whitespace separated, with or without a header,
 * decimal degrees or sexagesimal (colons, spaces, or h/m/s and d/m/s letters),
 * with an optional name column before or after the coordinates.
 */

import { toCsv } from "./csv";

/** Rows per bulk request. The backend has no list cap; this keeps the UI snappy. */
export const MAX_BULK_ROWS = 1000;

/** Rows per upstream POST. 1000 rows answer in ~0.6s, so 500 keeps each well under that. */
export const BULK_CHUNK_SIZE = 500;

export interface BulkTarget {
  /** 0-based position in the parsed target list (what the API calls `index`). */
  index: number;
  /** 1-based source line, for error messages and display. */
  line: number;
  name: string;
  ra: number;
  dec: number;
}

export interface BulkParseError {
  line: number;
  text: string;
  message: string;
}

export interface BulkParseResult {
  targets: BulkTarget[];
  errors: BulkParseError[];
  /** Detected header columns, or null when the input had none. */
  header: string[] | null;
}

type Delimiter = "\t" | "," | ";" | "|" | "whitespace";

const RA_NAMES = [
  "ra",
  "ra_deg",
  "radeg",
  "raj2000",
  "ra_j2000",
  "ra_icrs",
  "ra(deg)",
  "alpha",
  "right_ascension",
];
const DEC_NAMES = [
  "dec",
  "de",
  "dec_deg",
  "decdeg",
  "dej2000",
  "decj2000",
  "dec_j2000",
  "de_icrs",
  "dec_icrs",
  "dec(deg)",
  "delta",
  "declination",
];
const NAME_NAMES = [
  "name",
  "id",
  "objid",
  "object",
  "objname",
  "object_name",
  "target",
  "source",
  "source_id",
  "designation",
  "label",
];

function normalizeHeader(h: string): string {
  return h
    .trim()
    .toLowerCase()
    .replace(/^["']|["']$/g, "")
    .replace(/\s+/g, "");
}

const NUMBER_RE = /^[+-]?(\d+\.?\d*|\.\d+)(e[+-]?\d+)?$/i;

/** Split "12 30 45.6", "12:30:45.6", "12h30m45.6s", "+41°16′09″" into 3 parts. */
function sexagesimalParts(
  value: string,
  isRa: boolean
): { sign: number; a: number; b: number; c: number } | null {
  let s = value.trim();
  if (!s) return null;
  let sign = 1;
  if (s.startsWith("-") || s.startsWith("−")) {
    sign = -1;
    s = s.slice(1);
  } else if (s.startsWith("+")) {
    s = s.slice(1);
  }
  const unitRe = isRa ? /[hms:]|\s+/gi : /[dms°'′"″:]|\s+/gi;
  const parts = s.replace(unitRe, " ").trim().split(/\s+/).filter(Boolean);
  if (parts.length < 2 || parts.length > 3) return null;
  if (!parts.every((p) => /^\d+(\.\d*)?$/.test(p))) return null;
  const [a, b, c = 0] = parts.map(Number);
  // Only the last component may carry decimals.
  if (parts.slice(0, -1).some((p) => p.includes("."))) return null;
  if (b >= 60 || c >= 60) return null;
  return { sign, a, b, c };
}

function isSexagesimal(value: string): boolean {
  const s = value.trim();
  return /[:hmsd°'′"″]/i.test(s) || /\d\s+\d/.test(s);
}

/** Parse right ascension: decimal degrees or sexagesimal hours. */
export function parseRa(value: string): number | null {
  const s = value.trim();
  if (!s) return null;
  if (NUMBER_RE.test(s)) {
    const deg = Number(s);
    return deg >= 0 && deg < 360 ? deg : deg === 360 ? 0 : null;
  }
  if (/^\d+(\.\d+)?d(eg)?$/i.test(s)) return parseRa(s.replace(/d(eg)?$/i, ""));
  if (!isSexagesimal(s)) return null;
  const p = sexagesimalParts(s, true);
  if (!p || p.sign < 0 || p.a >= 24) return null;
  return (p.a + p.b / 60 + p.c / 3600) * 15;
}

/** Parse declination: decimal degrees or sexagesimal degrees. */
export function parseDec(value: string): number | null {
  const s = value.trim().replace(/^−/, "-");
  if (!s) return null;
  if (NUMBER_RE.test(s)) {
    const deg = Number(s);
    return deg >= -90 && deg <= 90 ? deg : null;
  }
  if (/^[+-]?\d+(\.\d+)?(d|deg|°)$/i.test(s))
    return parseDec(s.replace(/(d|deg|°)$/i, ""));
  if (!isSexagesimal(s)) return null;
  const p = sexagesimalParts(s, false);
  if (!p || p.a > 90) return null;
  const deg = p.sign * (p.a + p.b / 60 + p.c / 3600);
  return deg >= -90 && deg <= 90 ? deg : null;
}

function detectDelimiter(line: string): Delimiter {
  if (line.includes("\t")) return "\t";
  if (line.includes(",")) return ",";
  if (line.includes(";")) return ";";
  if (line.includes("|")) return "|";
  return "whitespace";
}

/** Split a delimited line, honouring double quotes. */
function splitDelimited(line: string, delim: string): string[] {
  const out: string[] = [];
  let cur = "";
  let quoted = false;
  for (let i = 0; i < line.length; i++) {
    const ch = line[i];
    if (ch === '"') {
      if (quoted && line[i + 1] === '"') {
        cur += '"';
        i++;
      } else {
        quoted = !quoted;
      }
    } else if (ch === delim && !quoted) {
      out.push(cur.trim());
      cur = "";
    } else {
      cur += ch;
    }
  }
  out.push(cur.trim());
  return out;
}

function splitLine(line: string, delim: Delimiter): string[] {
  return delim === "whitespace"
    ? line.trim().split(/\s+/)
    : splitDelimited(line, delim);
}

interface ColumnMap {
  ra: number;
  dec: number;
  name: number | null;
}

function headerColumns(fields: string[]): ColumnMap | null {
  const norm = fields.map(normalizeHeader);
  const find = (names: string[]) => norm.findIndex((h) => names.includes(h));
  const ra = find(RA_NAMES);
  const dec = find(DEC_NAMES);
  if (ra < 0 || dec < 0) return null;
  const name = find(NAME_NAMES);
  return { ra, dec, name: name >= 0 ? name : null };
}

type RowParse =
  | { ok: true; name: string; ra: number; dec: number }
  | { ok: false; message: string };

function coordsError(raText: string, decText: string): string {
  if (parseRa(raText) === null)
    return `Could not read RA "${raText}" (decimal degrees 0–360 or hh:mm:ss)`;
  return `Could not read Dec "${decText}" (decimal degrees ±90 or ±dd:mm:ss)`;
}

function parseByColumns(fields: string[], cols: ColumnMap): RowParse {
  const raText = fields[cols.ra] ?? "";
  const decText = fields[cols.dec] ?? "";
  const ra = parseRa(raText);
  const dec = parseDec(decText);
  if (ra === null || dec === null)
    return { ok: false, message: coordsError(raText, decText) };
  const name = cols.name !== null ? (fields[cols.name] ?? "") : "";
  return { ok: true, name, ra, dec };
}

/** No header: guess where the coordinates are. */
function parseHeuristic(fields: string[], delim: Delimiter): RowParse {
  const f = fields.filter((x) => x !== "");
  if (f.length < 2)
    return { ok: false, message: "Expected at least RA and Dec" };

  const tryPair = (raText: string, decText: string) => {
    const ra = parseRa(raText);
    const dec = parseDec(decText);
    return ra !== null && dec !== null ? { ra, dec } : null;
  };

  if (delim === "whitespace") {
    const n = f.length;
    // Space-separated sexagesimal: "[name…] hh mm ss ±dd mm ss".
    if (n >= 6) {
      const hit = tryPair(
        f.slice(n - 6, n - 3).join(":"),
        f.slice(n - 3).join(":")
      );
      if (hit && f.slice(n - 6).every((t) => /^[+\-−]?\d+(\.\d+)?$/.test(t)))
        return { ok: true, name: f.slice(0, n - 6).join(" "), ...hit };
      const lead = tryPair(f.slice(0, 3).join(":"), f.slice(3, 6).join(":"));
      if (lead && f.slice(0, 6).every((t) => /^[+\-−]?\d+(\.\d+)?$/.test(t)))
        return { ok: true, name: f.slice(6).join(" "), ...lead };
    }
    const first = tryPair(f[0], f[1]);
    if (first) return { ok: true, name: f.slice(2).join(" "), ...first };
    const last = tryPair(f[n - 2], f[n - 1]);
    if (last) return { ok: true, name: f.slice(0, n - 2).join(" "), ...last };
    return { ok: false, message: coordsError(f[0], f[1]) };
  }

  // Delimited: "ra,dec[,name…]" or "name,ra,dec[,…]".
  const first = tryPair(f[0], f[1]);
  if (first) return { ok: true, name: f[2] ?? "", ...first };
  if (f.length >= 3) {
    const second = tryPair(f[1], f[2]);
    if (second) return { ok: true, name: f[0], ...second };
    return { ok: false, message: coordsError(f[1], f[2]) };
  }
  return { ok: false, message: coordsError(f[0], f[1]) };
}

function isCommentOrBlank(line: string): boolean {
  const t = line.trim();
  return t === "" || t.startsWith("#") || t.startsWith("//");
}

/** Parse pasted or uploaded text into bulk targets plus per-line errors. */
export function parseBulkInput(text: string): BulkParseResult {
  const lines = text.replace(/^﻿/, "").split(/\r\n|\r|\n/);
  const targets: BulkTarget[] = [];
  const errors: BulkParseError[] = [];

  const firstIdx = lines.findIndex((l) => !isCommentOrBlank(l));
  if (firstIdx < 0) return { targets, errors, header: null };

  const delim = detectDelimiter(lines[firstIdx]);
  const firstFields = splitLine(lines[firstIdx], delim);
  const cols = headerColumns(firstFields);
  // A first line with no parseable coordinate and no digits at all is a header
  // even when we cannot map its columns.
  const looksLikeHeader =
    cols !== null ||
    (!parseHeuristic(firstFields, delim).ok && !/\d/.test(lines[firstIdx]));

  let header: string[] | null = null;
  if (looksLikeHeader) {
    header = firstFields;
  }

  for (
    let i = looksLikeHeader ? firstIdx + 1 : firstIdx;
    i < lines.length;
    i++
  ) {
    const raw = lines[i];
    if (isCommentOrBlank(raw)) continue;
    const fields = splitLine(raw, delim);
    const parsed = cols
      ? parseByColumns(fields, cols)
      : parseHeuristic(fields, delim);
    if (parsed.ok) {
      targets.push({
        index: targets.length,
        line: i + 1,
        name: parsed.name.trim(),
        ra: parsed.ra,
        dec: parsed.dec,
      });
    } else {
      errors.push({ line: i + 1, text: raw.trim(), message: parsed.message });
    }
  }

  return { targets, errors, header };
}

/** Well-known targets for the "Load example" button (ICRS, J2000). */
export const BULK_EXAMPLE = `name,ra,dec
M31 nucleus,10.684708,41.268750
Betelgeuse,88.792939,7.407064
Vega,279.234735,38.783689
Sirius,06:45:08.917,-16:42:58.02
Proxima Centauri,14:29:42.946,-62:40:46.17
Trapezium (theta1 Ori C),05:35:16.463,-05:23:22.85
3C 273,187.277915,2.052388
Crab Pulsar,83.633083,22.014500`;

// ── Results ────────────────────────────────────────────────────────────────

/** One element of the `/bulk-conesearch` response. */
export interface BulkApiMatch {
  catalog: string;
  data: {
    id?: string;
    ra?: number;
    dec?: number;
    cat?: string;
    distance?: number;
    ipix?: number;
  }[];
  index: number;
}

export interface BulkMatch {
  catalog: string;
  id: string;
  ra: number | null;
  dec: number | null;
  /** Separation from the input position, arcsec. */
  separation: number | null;
}

export interface BulkResultRow {
  target: BulkTarget;
  matches: BulkMatch[];
}

/** Attach matches to their input rows; matches sorted by separation. */
export function groupBulkResults(
  targets: BulkTarget[],
  response: BulkApiMatch[]
): BulkResultRow[] {
  const byIndex = new Map<number, BulkMatch[]>();
  for (const group of response) {
    for (const item of group.data ?? []) {
      if (!item.id) continue;
      const list = byIndex.get(group.index) ?? [];
      list.push({
        catalog: item.cat ?? group.catalog,
        id: item.id,
        ra: item.ra ?? null,
        dec: item.dec ?? null,
        separation: item.distance ?? null,
      });
      byIndex.set(group.index, list);
    }
  }
  return targets.map((target) => ({
    target,
    matches: (byIndex.get(target.index) ?? []).sort(
      (a, b) => (a.separation ?? Infinity) - (b.separation ?? Infinity)
    ),
  }));
}

export interface BulkSummary {
  inputs: number;
  matched: number;
  unmatched: number;
  totalMatches: number;
  /** Input rows with at least one match, per catalog. */
  perCatalog: Record<string, number>;
}

export function summarizeBulkResults(rows: BulkResultRow[]): BulkSummary {
  const perCatalog: Record<string, number> = {};
  let matched = 0;
  let totalMatches = 0;
  for (const row of rows) {
    if (row.matches.length > 0) matched++;
    totalMatches += row.matches.length;
    for (const cat of new Set(row.matches.map((m) => m.catalog))) {
      perCatalog[cat] = (perCatalog[cat] ?? 0) + 1;
    }
  }
  return {
    inputs: rows.length,
    matched,
    unmatched: rows.length - matched,
    totalMatches,
    perCatalog,
  };
}

/** One CSV line per match; unmatched inputs keep a line with empty match columns. */
export function bulkResultsToCsv(rows: BulkResultRow[]): string {
  const header = [
    "input_row",
    "input_name",
    "input_ra",
    "input_dec",
    "catalog",
    "match_id",
    "match_ra",
    "match_dec",
    "separation_arcsec",
  ];
  const out: (string | number | null)[][] = [];
  for (const { target, matches } of rows) {
    const base = [target.index + 1, target.name, target.ra, target.dec];
    if (matches.length === 0) {
      out.push([...base, null, null, null, null, null]);
      continue;
    }
    for (const m of matches) {
      out.push([
        ...base,
        m.catalog,
        m.id,
        m.ra,
        m.dec,
        m.separation === null ? null : Number(m.separation.toFixed(4)),
      ]);
    }
  }
  return toCsv(header, out);
}
