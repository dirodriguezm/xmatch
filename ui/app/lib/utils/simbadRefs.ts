/**
 * Papers about an object, from SIMBAD's bibliography (the `ref` and
 * `has_ref` tables of its TAP service). No API key needed, unlike ADS.
 */

export interface SimbadReference {
  bibcode: string;
  title?: string;
  year?: number;
  journal?: string;
  doi?: string;
  /** How many objects the paper is attached to in SIMBAD. */
  objects?: number;
}

/**
 * `newest` lists by year; `focused` puts papers attached to the fewest
 * objects first: those are about this object, not a survey table that
 * merely lists it. SIMBAD has no citation counts (ADS does, with a key).
 */
export type SimbadRefsSort = "newest" | "focused";

export interface SimbadRefsOptions {
  sort?: SimbadRefsSort;
  /** Publication years, from inclusive and to exclusive. */
  fromYear?: number;
  toYear?: number;
}

function refsWhere(oid: number, { fromYear, toYear }: SimbadRefsOptions) {
  return [
    `h.oidref = ${oid}`,
    ...(fromYear !== undefined ? [`r."year" >= ${fromYear}`] : []),
    ...(toYear !== undefined ? [`r."year" < ${toYear}`] : []),
  ].join(" AND ");
}

/** SIMBAD's TAP has no OFFSET, so "show more" asks for a larger TOP. */
export const SIMBAD_REFS_MAX = 200;

/** `oid` is SIMBAD's integer object id; every number here is validated. */
export function simbadRefsQuery(
  oid: number,
  limit: number,
  options: SimbadRefsOptions = {}
): string {
  const order =
    options.sort === "focused"
      ? "r.nbobject ASC, yr DESC, bibcode DESC"
      : "yr DESC, bibcode DESC";
  return (
    `SELECT TOP ${limit} r.bibcode, r.title, r."year" AS yr, r.journal, r.doi, r.nbobject ` +
    `FROM ref AS r JOIN has_ref AS h ON h.oidbibref = r.oidbib ` +
    `WHERE ${refsWhere(oid, options)} ORDER BY ${order}`
  );
}

export function simbadRefsCountQuery(
  oid: number,
  options: SimbadRefsOptions = {}
): string {
  if (options.fromYear === undefined && options.toYear === undefined) {
    return `SELECT COUNT(*) AS n FROM has_ref WHERE oidref = ${oid}`;
  }
  return (
    `SELECT COUNT(*) AS n FROM has_ref AS h JOIN ref AS r ON h.oidbibref = r.oidbib ` +
    `WHERE ${refsWhere(oid, options)}`
  );
}

const str = (v: unknown) =>
  typeof v === "string" && v.trim() !== "" ? v.trim() : undefined;

/** Rows of {@link simbadRefsQuery} in TAP JSON format. */
export function parseSimbadRefs(data: unknown): SimbadReference[] {
  if (!Array.isArray(data)) return [];
  const refs: SimbadReference[] = [];
  for (const row of data) {
    if (!Array.isArray(row)) continue;
    const [bibcode, title, year, journal, doi, objects] = row;
    const code = str(bibcode);
    if (!code) continue;
    refs.push({
      bibcode: code,
      // SIMBAD titles end with a full stop; drop it for display.
      title: str(title)?.replace(/\.$/, ""),
      year: typeof year === "number" ? year : undefined,
      journal: str(journal),
      doi: str(doi),
      objects: typeof objects === "number" ? objects : undefined,
    });
  }
  return refs;
}

export function adsAbstractUrl(bibcode: string): string {
  return `https://ui.adsabs.harvard.edu/abs/${encodeURIComponent(bibcode)}/abstract`;
}
