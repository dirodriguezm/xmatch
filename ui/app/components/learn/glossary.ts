/** Glossary of catalog columns, built from CATALOG_FIELD_DESCRIPTIONS. */

import { CATALOG_FIELD_DESCRIPTIONS } from "@/app/lib/constants/catalogFields";

export interface GlossaryEntry {
  catalog: string;
  field: string;
  description: string;
}

export function buildGlossary(
  source: Record<string, Record<string, string>> = CATALOG_FIELD_DESCRIPTIONS
): GlossaryEntry[] {
  return Object.entries(source).flatMap(([catalog, fields]) =>
    Object.entries(fields).map(([field, description]) => ({
      catalog,
      field,
      description,
    }))
  );
}

/**
 * Case-insensitive filter over field name, description and catalog. Every
 * whitespace-separated term must match; underscores and spaces are treated
 * alike so "parallax error" finds parallax_error.
 */
export function filterGlossary(
  entries: GlossaryEntry[],
  query: string
): GlossaryEntry[] {
  const norm = (s: string) => s.toLowerCase().replace(/_/g, " ");
  const terms = norm(query).split(/\s+/).filter(Boolean);
  if (terms.length === 0) return entries;
  return entries.filter((e) => {
    const hay = norm(`${e.field} ${e.catalog} ${e.description}`);
    return terms.every((t) => hay.includes(t));
  });
}

/** Group entries by catalog, keeping catalog order of first appearance. */
export function groupByCatalog(
  entries: GlossaryEntry[]
): [string, GlossaryEntry[]][] {
  const groups = new Map<string, GlossaryEntry[]>();
  for (const e of entries) {
    const list = groups.get(e.catalog) ?? [];
    list.push(e);
    groups.set(e.catalog, list);
  }
  return [...groups.entries()];
}
