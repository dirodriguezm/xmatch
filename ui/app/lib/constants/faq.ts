/**
 * Frequently asked questions for /learn. Answers are plain paragraphs; links
 * are rendered after them. Keep claims in line with the Go service
 * (service/internal/search) — it does not propagate proper motion.
 */

import { MAX_RADIUS_ARCSEC } from "./search";

export interface FaqLink {
  href: string;
  label: string;
}

export interface FaqItem {
  id: string;
  question: string;
  answer: string[];
  links?: FaqLink[];
  /** Extra rich content rendered by the page. */
  extra?: "radius-table";
}

export const FAQ_ITEMS: FaqItem[] = [
  {
    id: "what-is-a-cross-match",
    question: "What is a cross-match?",
    answer: [
      "A positional cross-match finds the sources in one or more catalogs that lie close to a given sky position. XWave takes a coordinate and a radius and returns, per catalog, every indexed source within that radius together with its angular separation from your position.",
      "It is purely positional: XWave does not use colours, magnitudes or source types to decide which counterpart is “right”. Closeness is evidence, not proof — always weigh the separation against each catalog's positional accuracy.",
    ],
    links: [{ href: "#methods", label: "How matching works" }],
  },
  {
    id: "which-radius",
    question: "Which search radius should I pick for each catalog?",
    answer: [
      "Pick a radius a few times larger than the combined positional uncertainty of your input and the catalog. Too small and real counterparts are missed; too large and unrelated neighbours creep in. The defaults below are what the search form uses.",
      `The service accepts any positive radius, but this interface refuses radii above ${MAX_RADIUS_ARCSEC}″ because larger cone searches do not return in reasonable time.`,
    ],
    extra: "radius-table",
  },
  {
    id: "why-no-match",
    question: "Why did my search return no match?",
    answer: [
      "An empty result is a real answer: nothing in that catalog lies inside your radius. Common reasons are a radius smaller than the catalog's positional error, a source fainter than the survey limit, a position outside the survey footprint (eROSITA covers only half the sky), or a high proper-motion star that has moved since the catalog epoch.",
      "Also check the input: RA and Dec are decimal degrees, and the radius sent to the API is in arcseconds.",
    ],
  },
  {
    id: "epochs-proper-motion",
    question: "Does XWave account for epochs and proper motion?",
    answer: [
      "No. Positions are matched exactly as stored in each catalog, without propagating proper motion to a common epoch. Gaia DR3 positions are at epoch J2016.0, AllWISE positions are from the WISE/NEOWISE observations of 2010–2011, and eROSITA eRASS1 from 2019–2020.",
      "For most sources this does not matter at arcsecond scales, but a star moving 1″/yr will be several arcseconds away between Gaia and AllWISE. For nearby or high proper-motion stars, propagate your position with Gaia's pmra/pmdec before searching, or widen the radius and inspect the candidates.",
    ],
  },
  {
    id: "erosita-coverage",
    question: "Why does eROSITA return nothing for my target?",
    answer: [
      "Only the western Galactic hemisphere (Galactic longitude l > 180°) is public in eRASS1 (DR1). Searches in the eastern half return nothing, however bright the source.",
      "eROSITA positions are also uncertain by several arcseconds, so an X-ray source often has multiple optical candidates; use a 15–30″ radius and inspect them.",
    ],
    links: [{ href: "/catalogs", label: "Catalog coverage" }],
  },
  {
    id: "units",
    question: "What units and conventions does XWave use?",
    answer: [
      "Coordinates are ICRS (J2000) right ascension and declination in decimal degrees, RA in [0, 360] and Dec in [−90, 90]. The API radius is in arcseconds; the search form lets you type arcmin or degrees and converts for you. Separations are great-circle distances in arcseconds.",
      "Catalog columns keep the units of their source release — e.g. Gaia parallaxes in mas and proper motions in mas/yr, AllWISE and 2MASS magnitudes in the Vega system, eROSITA fluxes in erg/s/cm². See the glossary for each column.",
    ],
    links: [{ href: "#glossary", label: "Column glossary" }],
  },
  {
    id: "how-many-results",
    question: "How many objects does a search return?",
    answer: [
      "The API's nneighbor parameter caps the number of objects returned per request and defaults to 1. This interface asks for up to 100, so in practice your radius decides what comes back. In very crowded fields with a large radius you can still reach the cap — narrow the radius if the list looks truncated.",
    ],
  },
  {
    id: "rate-limits",
    question: "Are there rate limits or a maximum bulk size?",
    answer: [
      "The service does not currently enforce per-user rate limits, but it is a shared academic resource: please space out automated requests and avoid hammering it in tight loops.",
      "Bulk requests are processed server-side in chunks of positions with bounded concurrency, and there is no hard limit on list length in the service. Very large lists take proportionally longer, so split them into batches of a few thousand positions and retry failed batches rather than sending one enormous request.",
    ],
    links: [{ href: "/bulk", label: "Bulk cross-match" }],
  },
  {
    id: "how-to-cite",
    question: "How do I cite XWave and the catalogs?",
    answer: [
      "Cite XWave and the original catalog papers for every catalog you used, and include each survey's requested acknowledgement text. The Cite button on any object or results page generates BibTeX and acknowledgements for you.",
    ],
  },
  {
    id: "api-access",
    question: "Can I use XWave from code?",
    answer: [
      "Yes. Everything the interface does goes through a public HTTP API (cone search, bulk cone search, metadata and light curves). No key is needed. Try requests in the API playground, or copy ready-made curl, Python and JavaScript snippets from any result.",
    ],
    links: [
      { href: "/developers", label: "API playground" },
      {
        href: "https://xwave-astro.udp.cl/swagger/index.html",
        label: "Swagger reference",
      },
    ],
  },
  {
    id: "data-freshness",
    question: "How current is the data?",
    answer: [
      "XWave serves fixed public releases — Gaia DR3, the AllWISE Source Catalog and eROSITA eRASS1 (DR1) — indexed once. They do not change until a new release is ingested, which is announced in the changelog. Data shown on object pages from external services (SIMBAD, VizieR, ALeRCE, etc.) is fetched live or cached for up to a week.",
    ],
    links: [
      { href: "/catalogs", label: "Catalog releases" },
      { href: "/changelog", label: "Changelog" },
    ],
  },
  {
    id: "reporting-bad-matches",
    question: "I found a wrong or suspicious match. How do I report it?",
    answer: [
      "Use the Report button on the object page. Choose what looks wrong, add a note, and XWave opens a prefilled GitHub issue with the object id, catalog and coordinates for you to review and submit. Your browser remembers what you reported.",
    ],
    links: [{ href: "/contact", label: "Other ways to get in touch" }],
  },
  {
    id: "privacy",
    question: "What happens to my searches? Is there tracking?",
    answer: [
      "There are no accounts and no analytics trackers. Searches are sent to the XWave API, and object pages query third-party services (CDS, IRSA, MAST, ALeRCE, the Gaia archive, NOIRLab) mostly through this site's server. The sky viewer loads images directly from CDS in your browser. Standard server logs may record requests.",
      "Preferences such as your basket and reported matches are stored only in your browser's local storage. Reports and feedback become public GitHub issues under your GitHub account.",
    ],
  },
];
