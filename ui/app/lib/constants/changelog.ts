/**
 * Release notes shown on /changelog and published as RSS at
 * /changelog/rss.xml. Newest first. Entries summarise the git history.
 */

export type ChangelogTag = "New" | "Improved" | "Catalogs" | "API" | "Fixed";

export interface ChangelogEntry {
  /** Stable slug, used as the page anchor and RSS guid. */
  id: string;
  /** ISO date, YYYY-MM-DD. */
  date: string;
  title: string;
  summary: string;
  tags: ChangelogTag[];
  items: string[];
}

export const CHANGELOG: ChangelogEntry[] = [
  {
    id: "2026-09-26-new-site-sections",
    date: "2026-09-26",
    title: "New site sections",
    summary:
      "XWave grows from a search box into a full site: bulk matching, reproducible code and live status.",
    tags: ["New"],
    items: [
      "Bulk cross-match: upload a list of positions and match them against every catalog in one go.",
      "Copy as code: every search and object has ready-to-run curl, Python and JavaScript snippets.",
      "Share and embed: copyable links and an embeddable object card.",
      "Explore: browse the sky and collect objects in a basket for later.",
      "Command palette (⌘K / Ctrl K) to jump anywhere.",
      "API playground for trying endpoints from the browser.",
      "Status page, Learn (FAQ, methods, glossary), Contact and this changelog — plus a Report button to flag bad matches.",
    ],
  },
  {
    id: "2026-09-26-object-page-context",
    date: "2026-09-26",
    title: "Object pages: identity, neighbours and more light curves",
    summary:
      "The object page now tells you what SIMBAD calls a source, what else is nearby and how it varies across surveys.",
    tags: ["New", "Improved"],
    items: [
      "SIMBAD identity with redshift or radial velocity when available.",
      "Nearby sources list and a sky map of neighbours around the target.",
      "Hover descriptions for every catalog column.",
      "Gaia DR3 epoch photometry and Pan-STARRS1 DR2 light curves alongside ZTF.",
    ],
  },
  {
    id: "2026-09-26-sed-observability",
    date: "2026-09-26",
    title: "Multi-wavelength SED and observability",
    summary:
      "Counterparts across catalogs are combined into a spectral energy distribution, and you can see when a target is up from Chile.",
    tags: ["New", "Improved"],
    items: [
      "Cross-catalog counterparts on the object page, each linking to its own page.",
      "SED built from XWave photometry plus every VizieR catalog within 2″.",
      "Optional extinction correction using the IRSA dust service; SED download as CSV.",
      "Observability panel for Chilean observatories and a “tonight” summary in results.",
      "ZTF light curves are fetched from ALeRCE directly.",
    ],
  },
  {
    id: "2026-09-22-standalone-binary",
    date: "2026-09-22",
    title: "Standalone release binary",
    summary:
      "The service now ships as a single fully static executable, making deployments and rollbacks simpler.",
    tags: ["API"],
    items: [
      "HEALPix libraries and libc are linked statically; no Nix or system libraries needed on the server.",
      "Releases are built in CI and verified by digest before promotion.",
    ],
  },
  {
    id: "2026-09-12-radius-and-about",
    date: "2026-09-12",
    title: "The radius now governs results",
    summary:
      "Searches return everything inside your radius (up to a cap) instead of a single object per catalog.",
    tags: ["Improved", "Fixed"],
    items: [
      "Cone-search radius is sent in arcseconds (it was mistakenly sent in degrees).",
      "Up to 100 objects per catalog are returned; radii above 120″ are refused rather than left hanging.",
      "Photometry column and CSV export in the results table.",
      "The interface now talks to the public API at xwave-astro.udp.cl.",
      "New About page explaining how the cross-match works.",
    ],
  },
  {
    id: "2026-06-15-gaia-columns",
    date: "2026-06-15",
    title: "Richer Gaia metadata and faster cone search",
    summary:
      "Gaia DR3 rows carry more columns, and cone search uses range queries over HEALPix pixels.",
    tags: ["Catalogs", "API"],
    items: [
      "Expanded Gaia DR3 column set (astrometry quality, photometry, astrophysical parameters).",
      "Cone search queries contiguous HEALPix pixel ranges instead of pixel lists.",
      "Light-curve queries respect the requested radius.",
      "Catalog handling refactored so each catalog plugs in through an adapter.",
    ],
  },
  {
    id: "2026-05-31-archives-desi",
    date: "2026-05-31",
    title: "Archives card and DESI spectra",
    summary:
      "Object pages link out to the major archives and show a DESI DR1 spectrum when one exists.",
    tags: ["New"],
    items: [
      "Grouped external links to SIMBAD, VizieR, NED, Aladin, Legacy Survey, SDSS and Pan-STARRS.",
      "DESI DR1 target lookup and an interactive spectrum chart with line markers.",
      "Light-curve download is always visible, disabled when no data exists.",
    ],
  },
  {
    id: "2026-04-27-radio-ztf",
    date: "2026-04-27",
    title: "Radio surveys and ZTF data release light curves",
    summary:
      "The sky viewer gains radio layers, and the API can serve ZTF data-release light curves.",
    tags: ["New", "API"],
    items: [
      "VLASS 3 GHz, RACS-mid 1.4 GHz and a Gaia DR3 overlay in the sky viewer.",
      "Multi-wavelength survey selector and randomised quick examples.",
      "ZTF DR light curves and a catalog filter on the /lightcurve endpoint.",
      "Bulk cone search results carry the index of the input coordinate.",
    ],
  },
  {
    id: "2026-03-16-erosita",
    date: "2026-03-16",
    title: "eROSITA and per-catalog radii",
    summary:
      "X-ray sources from eRASS1 join Gaia and AllWISE, and each catalog gets its own search radius.",
    tags: ["Catalogs", "Improved"],
    items: [
      "eROSITA eRASS1 (DR1) catalog indexed and searchable.",
      "Per-catalog search radius in the search form.",
      "Metadata responses include coordinates; CORS enabled for the web app.",
    ],
  },
  {
    id: "2026-02-13-web-interface",
    date: "2026-02-13",
    title: "XWave web interface",
    summary:
      "A Next.js interface for cone search, results and object detail pages.",
    tags: ["New"],
    items: [
      "Cross-match search by coordinates or object name.",
      "Results table and sky plot across catalogs.",
      "Object detail viewer with Aladin Lite sky images and light curves.",
    ],
  },
];

function escapeXml(s: string): string {
  return s
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;")
    .replace(/'/g, "&apos;");
}

/** RFC 822 date at noon UTC, as RSS 2.0 requires. */
export function rssDate(isoDate: string): string {
  return new Date(`${isoDate}T12:00:00Z`).toUTCString();
}

export function buildChangelogRss(
  entries: ChangelogEntry[],
  siteUrl: string
): string {
  const base = siteUrl.replace(/\/$/, "");
  const items = entries
    .map((e) => {
      const link = `${base}/changelog#${e.id}`;
      const html = `<p>${escapeXml(e.summary)}</p><ul>${e.items
        .map((i) => `<li>${escapeXml(i)}</li>`)
        .join("")}</ul>`;
      return [
        "    <item>",
        `      <title>${escapeXml(e.title)}</title>`,
        `      <link>${link}</link>`,
        `      <guid isPermaLink="false">${escapeXml(e.id)}</guid>`,
        `      <pubDate>${rssDate(e.date)}</pubDate>`,
        ...e.tags.map((t) => `      <category>${escapeXml(t)}</category>`),
        `      <description>${escapeXml(html)}</description>`,
        "    </item>",
      ].join("\n");
    })
    .join("\n");
  const lastBuild = entries[0] ? rssDate(entries[0].date) : "";
  return `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">
  <channel>
    <title>XWave changelog</title>
    <link>${base}/changelog</link>
    <description>New features, catalogs and API changes in XWave, the astronomical cross-match service.</description>
    <language>en</language>
    <lastBuildDate>${lastBuild}</lastBuildDate>
    <atom:link href="${base}/changelog/rss.xml" rel="self" type="application/rss+xml"/>
${items}
  </channel>
</rss>
`;
}
