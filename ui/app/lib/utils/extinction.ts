/**
 * Galactic extinction curve: Cardelli, Clayton & Mathis (1989) with the
 * O'Donnell (1994) optical coefficients, the standard Milky Way law.
 *
 * Valid for 0.1–3.33 µm (x = 1/λ from 0.3 to 10 µm⁻¹). Longward of 3.33 µm
 * the CCM infrared power law is extrapolated; there A_λ is a few per cent of
 * A_V, so for typical high-latitude reddening the correction is negligible
 * either way. Shortward of 0.1 µm no correction is returned.
 */

/** Ratio of total to selective extinction for the diffuse ISM. */
export const DEFAULT_RV = 3.1;

/** A_λ / A_V at wavelength λ (µm), or null outside the curve's range. */
export function extinctionRatio(
  wavelengthUm: number,
  rv = DEFAULT_RV
): number | null {
  const x = 1 / wavelengthUm;
  let a: number;
  let b: number;

  if (x > 10) {
    return null;
  } else if (x < 1.1) {
    // Infrared (extrapolated below x = 0.3)
    const p = Math.pow(x, 1.61);
    a = 0.574 * p;
    b = -0.527 * p;
  } else if (x < 3.3) {
    // Optical / near-IR, O'Donnell (1994)
    const y = x - 1.82;
    a =
      1 +
      0.104 * y -
      0.609 * y ** 2 +
      0.701 * y ** 3 +
      1.137 * y ** 4 -
      1.718 * y ** 5 -
      0.827 * y ** 6 +
      1.647 * y ** 7 -
      0.505 * y ** 8;
    b =
      1.952 * y +
      2.908 * y ** 2 -
      3.989 * y ** 3 -
      7.985 * y ** 4 +
      11.102 * y ** 5 +
      5.491 * y ** 6 -
      10.805 * y ** 7 +
      3.347 * y ** 8;
  } else if (x < 8) {
    // Ultraviolet
    let fa = 0;
    let fb = 0;
    if (x >= 5.9) {
      const d = x - 5.9;
      fa = -0.04473 * d ** 2 - 0.009779 * d ** 3;
      fb = 0.213 * d ** 2 + 0.1207 * d ** 3;
    }
    a = 1.752 - 0.316 * x - 0.104 / ((x - 4.67) ** 2 + 0.341) + fa;
    b = -3.09 + 1.825 * x + 1.206 / ((x - 4.62) ** 2 + 0.263) + fb;
  } else {
    // Far ultraviolet
    const d = x - 8;
    a = -1.073 - 0.628 * d + 0.137 * d ** 2 - 0.07 * d ** 3;
    b = 13.67 + 4.257 * d - 0.42 * d ** 2 + 0.374 * d ** 3;
  }

  return a + b / rv;
}

/** Extinction A_λ in magnitudes for a given E(B−V), or null out of range. */
export function extinctionMag(
  wavelengthUm: number,
  ebv: number,
  rv = DEFAULT_RV
): number | null {
  const ratio = extinctionRatio(wavelengthUm, rv);
  return ratio == null ? null : ratio * rv * ebv;
}

/**
 * Multiplicative factor that removes the extinction from a flux:
 * F_intrinsic = F_observed · 10^(0.4 A_λ). 1 when out of range.
 */
export function dereddenFactor(
  wavelengthUm: number,
  ebv: number,
  rv = DEFAULT_RV
): number {
  const a = extinctionMag(wavelengthUm, ebv, rv);
  return a == null ? 1 : 10 ** (0.4 * a);
}

/**
 * The SFD/SF11 maps give total line-of-sight dust, which is meaningless for
 * anything inside the Galactic plane and wildly large toward its centre
 * (E(B−V) ≈ 85 at Sgr A*). Correction is only offered outside these limits.
 */
const MIN_ABS_GALACTIC_LATITUDE = 5;
const MAX_EBV = 1;

/**
 * Why a map E(B−V) should not be used to deredden a source at Galactic
 * latitude `b` (degrees), or null when it can be.
 */
export function dereddeningBlocker(ebv: number, b: number): string | null {
  if (Math.abs(b) < MIN_ABS_GALACTIC_LATITUDE)
    return `reddening maps are unreliable in the Galactic plane (b = ${b.toFixed(1)}°)`;
  if (ebv > MAX_EBV)
    return `E(B−V) = ${ebv.toFixed(2)} is total line-of-sight dust, too high to correct reliably`;
  return null;
}
