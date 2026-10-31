/**
 * Rest-frame spectral lines for overlaying on a DESI spectrum, mirroring the set
 * shown by the Legacy Survey DESI viewer. Observed wavelength is computed at the
 * call site as `restWavelength * (1 + redshift)`.
 *
 * Wavelengths are in Ångström (vacuum), to the precision the markers need.
 */

export type SpectralLineKind = "emission" | "absorption";

export interface SpectralLine {
  /** Short display label, e.g. "Hα", "[O III]". */
  label: string;
  /** Rest-frame wavelength in Å. */
  restWavelength: number;
  kind: SpectralLineKind;
}

export const SPECTRAL_LINES: SpectralLine[] = [
  // Emission lines (typical star-forming / AGN galaxy features).
  { label: "[O II]", restWavelength: 3727, kind: "emission" },
  { label: "[Ne III]", restWavelength: 3869, kind: "emission" },
  { label: "Hδ", restWavelength: 4102, kind: "emission" },
  { label: "Hγ", restWavelength: 4340, kind: "emission" },
  { label: "Hβ", restWavelength: 4861, kind: "emission" },
  { label: "[O III]", restWavelength: 4959, kind: "emission" },
  { label: "[O III]", restWavelength: 5007, kind: "emission" },
  { label: "[N II]", restWavelength: 6548, kind: "emission" },
  { label: "Hα", restWavelength: 6563, kind: "emission" },
  { label: "[N II]", restWavelength: 6584, kind: "emission" },
  { label: "[S II]", restWavelength: 6716, kind: "emission" },
  { label: "[S II]", restWavelength: 6731, kind: "emission" },

  // Stellar absorption features (common in passive galaxies).
  { label: "Ca K", restWavelength: 3934, kind: "absorption" },
  { label: "Ca H", restWavelength: 3969, kind: "absorption" },
  { label: "G band", restWavelength: 4304, kind: "absorption" },
  { label: "Mg b", restWavelength: 5175, kind: "absorption" },
  { label: "Na D", restWavelength: 5892, kind: "absorption" },
];

/**
 * Sky-emission / telluric absorption regions (observed frame, fixed wavelengths)
 * shaded faintly behind the spectrum, as in the reference plot.
 */
export const SKY_BANDS: [number, number][] = [
  [5700, 5900],
  [7500, 7700],
];
