/** Number formatting shared across panels. */

const SUPERSCRIPT: Record<string, string> = {
  "-": "⁻",
  "0": "⁰",
  "1": "¹",
  "2": "²",
  "3": "³",
  "4": "⁴",
  "5": "⁵",
  "6": "⁶",
  "7": "⁷",
  "8": "⁸",
  "9": "⁹",
};

/** 1.41e-10 → "1.41 × 10⁻¹⁰". */
export function formatScientific(value: number, digits = 2): string {
  if (!Number.isFinite(value)) return "—";
  if (value === 0) return "0";
  const [mantissa = "", exponent] = value.toExponential(digits).split("e");
  const exp = String(Number(exponent));
  if (exp === "0") return mantissa;
  const sup = [...exp].map((c) => SUPERSCRIPT[c] ?? c).join("");
  return `${mantissa} × 10${sup}`;
}

export const FLUX_UNIT = "erg s⁻¹ cm⁻²";

/** X-ray flux with its unit, e.g. "1.41 × 10⁻¹⁰ erg s⁻¹ cm⁻²". */
export function formatFlux(value: number, digits = 2): string {
  return `${formatScientific(value, digits)} ${FLUX_UNIT}`;
}

/** MJD of the Unix epoch, 1970-01-01T00:00Z. */
const MJD_UNIX_EPOCH = 40587;

/** MJD → calendar date, "2012-07-16" (UTC). */
export function formatMjdDate(mjd: number): string {
  return new Date((mjd - MJD_UNIX_EPOCH) * 86_400_000)
    .toISOString()
    .slice(0, 10);
}

/** MJD → calendar date and time, "2012-07-16 22:52 UTC". */
export function formatMjdDateTime(mjd: number): string {
  const iso = new Date((mjd - MJD_UNIX_EPOCH) * 86_400_000).toISOString();
  return `${iso.slice(0, 10)} ${iso.slice(11, 16)} UTC`;
}

/**
 * Flux and its 1σ error on one power of ten, with the unit after both:
 * (1.42 ± 0.07) × 10⁻¹¹ erg s⁻¹ cm⁻². The error keeps two significant
 * figures, at most four decimals; a smaller one shows as "< 0.0001".
 */
export function formatFluxWithError(value: number, err: number): string {
  if (!Number.isFinite(value) || value === 0) return formatFlux(value);
  if (!Number.isFinite(err) || err <= 0) return formatFlux(value);
  const exponent = Math.floor(Math.log10(Math.abs(value)));
  const scale = 10 ** exponent;
  const errMantissa = err / scale;
  const decimals = Math.min(
    4,
    Math.max(2, 1 - Math.floor(Math.log10(errMantissa)))
  );
  const errText =
    errMantissa < 10 ** -decimals
      ? `< ${(10 ** -decimals).toFixed(decimals)}`
      : errMantissa.toFixed(decimals);
  const sup = [...String(exponent)].map((c) => SUPERSCRIPT[c] ?? c).join("");
  const pow = exponent === 0 ? "" : ` × 10${sup}`;
  return `(${(value / scale).toFixed(decimals)} ± ${errText})${pow} ${FLUX_UNIT}`;
}

/** Calendar year (UTC) of an MJD. */
export function mjdYear(mjd: number): number {
  return new Date((mjd - MJD_UNIX_EPOCH) * 86_400_000).getUTCFullYear();
}

/** Time-axis tick: the MJD, with its year underneath. */
export function formatMjdTick(mjd: number): string {
  return `${mjd}\n${mjdYear(mjd)}`;
}

/**
 * A catalogue cell for display: floats to six significant figures (very
 * small or large ones in scientific notation), everything else as is. The
 * catalogues store float32/64 noise past that, e.g. 0.018106548115611076.
 */
export function formatCatalogValue(value: unknown): string {
  if (value == null) return "—";
  if (typeof value !== "number" || Number.isInteger(value)) {
    return String(value);
  }
  const abs = Math.abs(value);
  if (abs !== 0 && (abs < 1e-3 || abs >= 1e7)) {
    return formatScientific(value, 4);
  }
  return String(Number(value.toPrecision(6)));
}
