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
}

/** SIMBAD's TAP has no OFFSET, so "show more" asks for a larger TOP. */
export const SIMBAD_REFS_MAX = 200;

/** Newest first; `oid` is SIMBAD's integer object id. */
export function simbadRefsQuery(oid: number, limit: number): string {
  return (
    `SELECT TOP ${limit} r.bibcode, r.title, r."year" AS yr, r.journal, r.doi ` +
    `FROM ref AS r JOIN has_ref AS h ON h.oidbibref = r.oidbib ` +
    `WHERE h.oidref = ${oid} ORDER BY yr DESC, bibcode DESC`
  );
}

export function simbadRefsCountQuery(oid: number): string {
  return `SELECT COUNT(*) AS n FROM has_ref WHERE oidref = ${oid}`;
}

const str = (v: unknown) =>
  typeof v === "string" && v.trim() !== "" ? v.trim() : undefined;

/** Rows of {@link simbadRefsQuery} in TAP JSON format. */
export function parseSimbadRefs(data: unknown): SimbadReference[] {
  if (!Array.isArray(data)) return [];
  const refs: SimbadReference[] = [];
  for (const row of data) {
    if (!Array.isArray(row)) continue;
    const [bibcode, title, year, journal, doi] = row;
    const code = str(bibcode);
    if (!code) continue;
    refs.push({
      bibcode: code,
      // SIMBAD titles end with a full stop; drop it for display.
      title: str(title)?.replace(/\.$/, ""),
      year: typeof year === "number" ? year : undefined,
      journal: str(journal),
      doi: str(doi),
    });
  }
  return refs;
}

export function adsAbstractUrl(bibcode: string): string {
  return `https://ui.adsabs.harvard.edu/abs/${encodeURIComponent(bibcode)}/abstract`;
}
