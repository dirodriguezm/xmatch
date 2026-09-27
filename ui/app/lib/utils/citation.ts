/**
 * Citation export: BibTeX, plain-text references and acknowledgement text
 * for XWave and the catalogs it serves.
 */

import {
  CATALOG_META,
  type CatalogMeta,
} from "@/app/lib/constants/catalogMeta";
import { CATALOG_OPTIONS } from "@/app/lib/constants/catalogs";
import { REPO_URL, SITE_URL } from "@/app/lib/constants/site";

export const XWAVE_CITATION = {
  key: "XWave2026",
  title: "XWave: a HEALPix cross-match service",
  author: "{ALeRCE} and {Universidad Diego Portales}",
  authorPlain: "ALeRCE / Universidad Diego Portales",
  year: 2026,
} as const;

export const XWAVE_ACKNOWLEDGEMENT =
  "This research has made use of XWave, a HEALPix cross-match service developed by the ALeRCE team at Universidad Diego Portales.";

/**
 * Resolve catalog slugs to their metadata, deduplicated and in canonical
 * order. Unknown slugs are ignored; an empty list means every catalog.
 */
export function resolveCatalogs(catalogs: string[]): CatalogMeta[] {
  const wanted = new Set(catalogs.map((c) => c.toLowerCase()));
  const all = wanted.size === 0;
  return CATALOG_OPTIONS.filter((slug) => all || wanted.has(slug)).map(
    (slug) => CATALOG_META[slug]
  );
}

/** BibTeX entry for XWave itself (no DOI yet). */
export function xwaveBibtex(siteUrl: string = SITE_URL): string {
  return `% XWave does not have a DOI yet; a Zenodo DOI is pending.
@software{${XWAVE_CITATION.key},
  author = {${XWAVE_CITATION.author}},
  title = {{${XWAVE_CITATION.title}}},
  year = ${XWAVE_CITATION.year},
  url = {${siteUrl}},
  repository = {${REPO_URL}},
  note = {Web service and REST API; accessed via ${siteUrl}}
}`;
}

/** Full BibTeX: XWave plus one entry per (deduplicated) catalog. */
export function buildBibtex(catalogs: string[], siteUrl?: string): string {
  const entries = [xwaveBibtex(siteUrl)];
  const seen = new Set<string>();
  for (const meta of resolveCatalogs(catalogs)) {
    if (seen.has(meta.reference.bibcode)) continue;
    seen.add(meta.reference.bibcode);
    entries.push(`% ${meta.name} ${meta.release}\n${meta.reference.bibtex}`);
  }
  return `${entries.join("\n\n")}\n`;
}

/** Extract the title field from a BibTeX entry, stripping braces/quotes. */
export function bibtexTitle(bibtex: string): string | undefined {
  const m = bibtex.match(
    /title\s*=\s*(?:"\{(.+?)\}"|\{\{?(.+?)\}?\})\s*,?\s*$/im
  );
  const raw = m?.[1] ?? m?.[2];
  return raw?.replace(/[{}]/g, "").trim();
}

/** Plain-text reference list, one line per work. */
export function buildPlainReferences(
  catalogs: string[],
  siteUrl: string = SITE_URL
): string[] {
  const refs = [
    `${XWAVE_CITATION.authorPlain} (${XWAVE_CITATION.year}). ${XWAVE_CITATION.title}. ${siteUrl} (source: ${REPO_URL}).`,
  ];
  const seen = new Set<string>();
  for (const meta of resolveCatalogs(catalogs)) {
    const { reference } = meta;
    if (seen.has(reference.bibcode)) continue;
    seen.add(reference.bibcode);
    const title = bibtexTitle(reference.bibtex);
    const id = reference.doi ? `doi:${reference.doi}` : reference.bibcode;
    refs.push(
      [reference.label, title, id].filter(Boolean).join(". ").concat(".")
    );
  }
  return refs;
}

/** Acknowledgement paragraph: XWave first, then each catalog's own text. */
export function buildAcknowledgement(catalogs: string[]): string {
  const parts = [XWAVE_ACKNOWLEDGEMENT];
  const seen = new Set<string>();
  for (const meta of resolveCatalogs(catalogs)) {
    if (seen.has(meta.acknowledgement)) continue;
    seen.add(meta.acknowledgement);
    parts.push(meta.acknowledgement);
  }
  return parts.join(" ");
}

/** Filename for the downloaded .bib, e.g. "xwave-gaia-dr3-12345.bib". */
export function bibFilename(objectId?: string): string {
  const slug = objectId
    ?.toLowerCase()
    .replace(/[^a-z0-9]+/g, "-")
    .replace(/^-+|-+$/g, "");
  return slug ? `xwave-${slug}.bib` : "xwave-references.bib";
}

/** Trigger a client-side download of a text file. */
export function downloadText(
  filename: string,
  text: string,
  type = "application/x-bibtex;charset=utf-8"
): void {
  if (typeof window === "undefined") return;
  const blob = new Blob([text], { type });
  const url = URL.createObjectURL(blob);
  const a = document.createElement("a");
  a.href = url;
  a.download = filename;
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
  URL.revokeObjectURL(url);
}
