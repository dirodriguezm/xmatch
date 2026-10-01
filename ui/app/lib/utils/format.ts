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
  const [mantissa, exponent] = value.toExponential(digits).split("e");
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
