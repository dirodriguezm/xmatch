/**
 * Reference metadata for each catalog XWave indexes: which release is served,
 * what it covers and how to credit it. Used by the Catalogs page, citation
 * export and the data-version footer.
 */

import type { CatalogOption } from "./catalogs";

export interface CatalogReference {
  /** Short author-year label, e.g. "Gaia Collaboration, Vallenari et al. 2023". */
  label: string;
  bibcode: string;
  doi?: string;
  bibtex: string;
}

export interface CatalogMeta {
  slug: CatalogOption;
  name: string;
  /** Release served by XWave, shown in version badges. */
  release: string;
  wavelength: string;
  coverage: string;
  /** Approximate number of sources in the published release. */
  sources: string;
  /** Typical positional accuracy, for choosing a match radius. */
  astrometry: string;
  homepage: string;
  reference: CatalogReference;
  /** Acknowledgement text the survey asks users to include. */
  acknowledgement: string;
  /** Caveats worth knowing before trusting a match. */
  knownIssues: string[];
}

export const CATALOG_META: Record<CatalogOption, CatalogMeta> = {
  gaia: {
    slug: "gaia",
    name: "Gaia",
    release: "DR3",
    wavelength: "Optical (G, BP, RP; 330–1050 nm)",
    coverage: "All sky",
    sources: "1.81 billion",
    astrometry: "~0.02–0.5 mas at epoch J2016.0",
    homepage: "https://www.cosmos.esa.int/web/gaia/dr3",
    reference: {
      label: "Gaia Collaboration, Vallenari et al. 2023",
      bibcode: "2023A&A...674A...1G",
      doi: "10.1051/0004-6361/202243940",
      bibtex: `@ARTICLE{2023A&A...674A...1G,
  author = {{Gaia Collaboration} and {Vallenari}, A. and {Brown}, A.~G.~A. and others},
  title = "{Gaia Data Release 3. Summary of the content and survey properties}",
  journal = {\\aap},
  year = 2023,
  volume = {674},
  eid = {A1},
  pages = {A1},
  doi = {10.1051/0004-6361/202243940},
  adsurl = {https://ui.adsabs.harvard.edu/abs/2023A&A...674A...1G}
}`,
    },
    acknowledgement:
      "This work has made use of data from the European Space Agency (ESA) mission Gaia (https://www.cosmos.esa.int/gaia), processed by the Gaia Data Processing and Analysis Consortium (DPAC, https://www.cosmos.esa.int/web/gaia/dpac/consortium). Funding for the DPAC has been provided by national institutions, in particular the institutions participating in the Gaia Multilateral Agreement.",
    knownIssues: [
      "Positions are at epoch J2016.0; high proper-motion stars drift out of small radii when matched against older epochs.",
      "Crowded fields (Galactic plane, cluster cores) have lower completeness and spurious astrometric solutions.",
    ],
  },
  allwise: {
    slug: "allwise",
    name: "AllWISE",
    release: "AllWISE Source Catalog",
    wavelength: "Mid-infrared (W1–W4; 3.4, 4.6, 12, 22 µm)",
    coverage: "All sky",
    sources: "747 million",
    astrometry: "~0.05–0.5″ for bright sources, ~1″ near the detection limit",
    homepage: "https://wise2.ipac.caltech.edu/docs/release/allwise/",
    reference: {
      label: "Cutri et al. 2013",
      bibcode: "2013yCat.2328....0C",
      bibtex: `@MISC{2013yCat.2328....0C,
  author = {{Cutri}, R.~M. and others},
  title = "{VizieR Online Data Catalog: AllWISE Data Release (Cutri+ 2013)}",
  howpublished = {VizieR On-line Data Catalog: II/328},
  year = 2013,
  adsurl = {https://ui.adsabs.harvard.edu/abs/2013yCat.2328....0C}
}`,
    },
    acknowledgement:
      "This publication makes use of data products from the Wide-field Infrared Survey Explorer, which is a joint project of the University of California, Los Angeles, and the Jet Propulsion Laboratory/California Institute of Technology, funded by the National Aeronautics and Space Administration.",
    knownIssues: [
      "The 6″ W1 beam (12″ in W4) blends close neighbours into one source.",
      "Check cc_flags for diffraction spikes, persistence and halo artefacts near bright stars.",
    ],
  },
  erosita: {
    slug: "erosita",
    name: "eROSITA",
    release: "eRASS1 (DR1)",
    wavelength: "X-ray (0.2–2.3 keV main band)",
    coverage: "Western Galactic hemisphere (German eROSITA sky)",
    sources: "~930 thousand",
    astrometry: "~5″ typical (1σ), larger for faint sources",
    homepage: "https://erosita.mpe.mpg.de/dr1/",
    reference: {
      label: "Merloni et al. 2024",
      bibcode: "2024A&A...682A..34M",
      doi: "10.1051/0004-6361/202347165",
      bibtex: `@ARTICLE{2024A&A...682A..34M,
  author = {{Merloni}, A. and {Lamer}, G. and {Liu}, T. and others},
  title = "{The SRG/eROSITA all-sky survey. First X-ray catalogues and data release of the western Galactic hemisphere}",
  journal = {\\aap},
  year = 2024,
  volume = {682},
  eid = {A34},
  pages = {A34},
  doi = {10.1051/0004-6361/202347165},
  adsurl = {https://ui.adsabs.harvard.edu/abs/2024A&A...682A..34M}
}`,
    },
    acknowledgement:
      "This work is based on data from eROSITA, the soft X-ray instrument aboard SRG, a joint Russian-German science mission supported by the Russian Space Agency (Roskosmos), in the interests of the Russian Academy of Sciences represented by its Space Research Institute (IKI), and the Deutsches Zentrum für Luft- und Raumfahrt (DLR).",
    knownIssues: [
      "Only the western Galactic hemisphere (l > 180°) is public; searches in the eastern half return nothing.",
      "Positional errors of several arcsec mean an eROSITA source often has multiple optical candidates; use a 15–30″ radius and inspect them.",
    ],
  },
};

/** One-line data version, e.g. "Gaia DR3 · AllWISE Source Catalog · eROSITA eRASS1 (DR1)". */
export function dataVersionLine(): string {
  return Object.values(CATALOG_META)
    .map((m) =>
      m.release.startsWith(m.name) ? m.release : `${m.name} ${m.release}`
    )
    .join(" · ");
}
