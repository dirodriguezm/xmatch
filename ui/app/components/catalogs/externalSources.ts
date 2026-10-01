/**
 * Public services the object page queries live (not indexed by XWave).
 * Shared by the Catalogs and About pages so the list is kept in one place.
 */

export type ExternalSourceKind =
  | "light curve"
  | "spectra"
  | "photometry"
  | "identity"
  | "context";

export interface ExternalSource {
  name: string;
  href: string;
  kind: ExternalSourceKind;
  usedFor: string;
}

export const EXTERNAL_SOURCES: ExternalSource[] = [
  {
    name: "ZTF data releases via the ALeRCE broker",
    href: "https://alerce.online/",
    kind: "light curve",
    usedFor:
      "ZTF g/r/i light curves, fetched by the XWave /lightcurve endpoint from ALeRCE's ZTF DR API",
  },
  {
    name: "NEOWISE via IRSA (NASA/IPAC)",
    href: "https://irsa.ipac.caltech.edu/Missions/wise.html",
    kind: "light curve",
    usedFor:
      "NEOWISE-R W1/W2 single-exposure photometry, fetched by the XWave /lightcurve endpoint",
  },
  {
    name: "Gaia Archive (ESA)",
    href: "https://gea.esac.esa.int/archive/",
    kind: "light curve",
    usedFor: "Gaia DR3 epoch photometry (G, BP, RP per transit)",
  },
  {
    name: "Pan-STARRS1 DR2 via MAST (STScI)",
    href: "https://catalogs.mast.stsci.edu/panstarrs/",
    kind: "light curve",
    usedFor: "per-epoch grizy detections for light curves (δ > −30°)",
  },
  {
    name: "Catalina Real-time Transient Survey (Caltech)",
    href: "http://nunuku.caltech.edu/cgi-bin/getcssconedb_release_img.cgi",
    kind: "light curve",
    usedFor:
      "Unfiltered, V-calibrated CSS/MLS/SSS photometry, 2005–2016 (Drake et al. 2009)",
  },
  {
    name: "HILIGT upper limit server (ESA)",
    href: "https://xmmuls.esac.esa.int/hiligt/",
    kind: "light curve",
    usedFor:
      "X-ray fluxes and upper limits from XMM-Newton (pointed and slew) and the ROSAT All-Sky Survey (Saxton et al. 2022; König et al. 2022)",
  },
  {
    name: "5XMM-DR15 via the XMM-Newton Science Archive (ESA)",
    href: "http://xmmssc.irap.omp.eu/Catalogue/5XMM-DR15/5XMM_DR15.html",
    kind: "context",
    usedFor:
      "Stacked X-ray source fluxes, hardness ratios and variability, compiled by the XMM-Newton Survey Science Centre and XMM2ATHENA (Webb, Traulsen et al. 2026)",
  },

  {
    name: "NOIRLab Astro Data Lab and SPARCL",
    href: "https://datalab.noirlab.edu/",
    kind: "spectra",
    usedFor: "DESI DR1 target lookup and spectra",
  },
  {
    name: "SIMBAD (CDS)",
    href: "https://simbad.cds.unistra.fr/simbad/",
    kind: "identity",
    usedFor: "the nearest known object's main identifier, redshift or velocity",
  },
  {
    name: "VizieR SED service (CDS)",
    href: "https://vizier.cds.unistra.fr/vizier/sed/",
    kind: "photometry",
    usedFor:
      "published photometry from every VizieR catalog within 2″, merged per filter in the SED",
  },
  {
    name: "IRSA Galactic Dust Reddening service (NASA/IPAC)",
    href: "https://irsa.ipac.caltech.edu/applications/DUST/",
    kind: "photometry",
    usedFor:
      "E(B−V) at the position (Schlafly & Finkbeiner 2011) for the SED's extinction correction",
  },
  {
    name: "Sesame name resolver (CDS)",
    href: "https://cds.unistra.fr/cgi-bin/Sesame",
    kind: "context",
    usedFor: "turning object names into coordinates",
  },
  {
    name: "Aladin Lite and HiPS surveys",
    href: "https://aladin.cds.unistra.fr/AladinLite/",
    kind: "context",
    usedFor:
      "the sky view (DSS, 2MASS, AllWISE, XMM, Chandra, NVSS, SUMSS, RACS, VLASS)",
  },
  {
    name: "Astronomy Engine",
    href: "https://github.com/cosinekitty/astronomy",
    kind: "context",
    usedFor:
      "Sun, Moon and target positions in the observability panel, computed in your browser",
  },
];

export const CDS_ACKNOWLEDGEMENT =
  "This research has made use of the VizieR catalogue access tool, CDS, Strasbourg, France (DOI: 10.26093/cds/vizier), and of the SIMBAD database and the Aladin sky atlas, operated at CDS, Strasbourg, France.";
