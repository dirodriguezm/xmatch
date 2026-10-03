/**
 * Spectral energy distribution built from the photometry of an object and its
 * counterparts in the other search catalogs.
 *
 * Every magnitude here is on the Vega system, so conversion to flux density is
 * F_ν = F_0 · 10^(−m/2.5) with the per-band Vega zero point F_0. Zero points and
 * effective wavelengths come from the SVO Filter Profile Service (Gaia DR3,
 * 2MASS: Cohen+2003) and Jarrett+2011 (WISE isophotal).
 */

import { isMeasured } from "@/app/lib/constants/photometry";

/** One band a search catalog publishes, and how to read it from metadata. */
interface SedBandSpec {
  band: string;
  survey: string;
  /** Effective wavelength in µm. */
  wavelengthUm: number;
  /** Vega zero point in Jy. */
  zeroPointJy: number;
  field: string;
  /** Magnitude uncertainty field. */
  errField?: string;
  /**
   * Flux signal-to-noise field, used when the catalog publishes no magnitude
   * error (Gaia): σ_m = 2.5 / ln10 / (S/N).
   */
  fluxOverErrorField?: string;
  /**
   * The catalog marks a non-detection by publishing the magnitude with a null
   * uncertainty: the value is then a 95% upper limit (AllWISE `w?mpro` and its
   * 2MASS associations).
   */
  nullErrIsUpperLimit?: boolean;
}

export const SED_BAND_SPECS: Record<string, SedBandSpec[]> = {
  gaia: [
    {
      band: "BP",
      survey: "Gaia DR3",
      wavelengthUm: 0.511,
      zeroPointJy: 3552.01,
      field: "phot_bp_mean_mag",
      fluxOverErrorField: "phot_bp_mean_flux_over_error",
    },
    {
      band: "G",
      survey: "Gaia DR3",
      wavelengthUm: 0.6218,
      zeroPointJy: 3228.75,
      field: "phot_g_mean_mag",
      fluxOverErrorField: "phot_g_mean_flux_over_error",
    },
    {
      band: "RP",
      survey: "Gaia DR3",
      wavelengthUm: 0.7769,
      zeroPointJy: 2554.95,
      field: "phot_rp_mean_mag",
      fluxOverErrorField: "phot_rp_mean_flux_over_error",
    },
  ],
  allwise: [
    {
      band: "J",
      survey: "2MASS",
      wavelengthUm: 1.235,
      zeroPointJy: 1594,
      field: "j_m_2mass",
      errField: "j_msig_2mass",
      nullErrIsUpperLimit: true,
    },
    {
      band: "H",
      survey: "2MASS",
      wavelengthUm: 1.662,
      zeroPointJy: 1024,
      field: "h_m_2mass",
      errField: "h_msig_2mass",
      nullErrIsUpperLimit: true,
    },
    {
      band: "K",
      survey: "2MASS",
      wavelengthUm: 2.159,
      zeroPointJy: 666.7,
      field: "k_m_2mass",
      errField: "k_msig_2mass",
      nullErrIsUpperLimit: true,
    },
    {
      band: "W1",
      survey: "WISE",
      wavelengthUm: 3.3526,
      zeroPointJy: 309.54,
      field: "w1mpro",
      errField: "w1sigmpro",
      nullErrIsUpperLimit: true,
    },
    {
      band: "W2",
      survey: "WISE",
      wavelengthUm: 4.6028,
      zeroPointJy: 171.787,
      field: "w2mpro",
      errField: "w2sigmpro",
      nullErrIsUpperLimit: true,
    },
    {
      band: "W3",
      survey: "WISE",
      wavelengthUm: 11.5608,
      zeroPointJy: 31.674,
      field: "w3mpro",
      errField: "w3sigmpro",
      nullErrIsUpperLimit: true,
    },
    {
      band: "W4",
      survey: "WISE",
      wavelengthUm: 22.0883,
      zeroPointJy: 8.363,
      field: "w4mpro",
      errField: "w4sigmpro",
      nullErrIsUpperLimit: true,
    },
  ],
  // eROSITA publishes X-ray fluxes, not magnitudes; not part of this SED yet.
  erosita: [],
};

/** Where a set of photometry came from. */
export interface PhotometrySource {
  /** Search-catalog slug: `gaia`, `allwise`, `erosita`. */
  catalog: string;
  id?: string;
  record: Record<string, unknown>;
  /** True for the object itself, false for a counterpart. */
  isSelf: boolean;
  /** Separation from the object in arcsec; 0 for the object itself. */
  separationArcsec: number;
}

export interface SedPoint {
  band: string;
  survey: string;
  catalog: string;
  wavelengthUm: number;
  mag: number;
  magErr?: number;
  upperLimit: boolean;
  /** νF_ν in erg s⁻¹ cm⁻². */
  nuFnu: number;
  nuFnuErr?: number;
  isSelf: boolean;
  separationArcsec: number;
  sourceId?: string;
}

const LN10_OVER_2_5 = Math.LN10 / 2.5;
/** Speed of light in µm/s, so c/λ[µm] is ν in Hz. */
const C_UM_PER_S = 2.99792458e14;
const JY_TO_CGS = 1e-23;

function toNuFnu(mag: number, spec: SedBandSpec): number {
  const fnuJy = spec.zeroPointJy * 10 ** (-mag / 2.5);
  return fnuJy * JY_TO_CGS * (C_UM_PER_S / spec.wavelengthUm);
}

/**
 * Turn the photometry of each source into SED points, sorted by wavelength.
 * Bands a source does not measure are skipped, never emitted as zero.
 */
export function buildSedPoints(sources: PhotometrySource[]): SedPoint[] {
  const points: SedPoint[] = [];

  for (const source of sources) {
    const specs = SED_BAND_SPECS[source.catalog.toLowerCase()] ?? [];
    for (const spec of specs) {
      const mag = source.record[spec.field];
      if (!isMeasured(mag, source.catalog)) continue;

      let magErr: number | undefined;
      if (spec.errField) {
        const err = source.record[spec.errField];
        if (isMeasured(err, source.catalog)) magErr = err;
      } else if (spec.fluxOverErrorField) {
        const snr = source.record[spec.fluxOverErrorField];
        if (isMeasured(snr, source.catalog) && snr > 0)
          magErr = 1 / (LN10_OVER_2_5 * snr);
      }

      const upperLimit = Boolean(spec.nullErrIsUpperLimit) && magErr == null;
      const nuFnu = toNuFnu(mag, spec);

      points.push({
        band: spec.band,
        survey: spec.survey,
        catalog: source.catalog.toLowerCase(),
        wavelengthUm: spec.wavelengthUm,
        mag,
        magErr,
        upperLimit,
        nuFnu,
        // Linear propagation of the magnitude error: σ_F = F · ln10/2.5 · σ_m.
        nuFnuErr: magErr != null ? nuFnu * LN10_OVER_2_5 * magErr : undefined,
        isSelf: source.isSelf,
        separationArcsec: source.separationArcsec,
        sourceId: source.id,
      });
    }
  }

  return points.sort((a, b) => a.wavelengthUm - b.wavelengthUm);
}
