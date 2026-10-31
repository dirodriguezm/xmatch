import { toCsv } from "@/app/lib/utils/csv";
import { extinctionMag } from "@/app/lib/utils/extinction";
import type { SedPoint } from "@/app/lib/utils/sed";
import type { VizierSedPoint } from "@/app/lib/utils/vizierSed";

const HEADER = [
  "source",
  "catalog",
  "band",
  "wavelength_um",
  "mag_vega",
  "mag_err",
  "upper_limit",
  "nufnu_cgs",
  "nufnu_err_cgs",
  "separation_arcsec",
  "n_measurements",
  "vizier_tables",
  "inconsistent",
  "a_lambda_mag",
  "nufnu_dereddened_cgs",
];

/**
 * Every SED point as CSV, observed fluxes always included. When `ebv` is
 * given the Galactic extinction A_λ (CCM89, R_V = 3.1) and the dereddened
 * flux are added; they stay empty when no usable E(B−V) is known.
 * νFν is in erg s⁻¹ cm⁻².
 */
export function sedToCsv(
  points: SedPoint[],
  vizier: VizierSedPoint[],
  ebv: number | null
): string {
  const dereddened = (wavelengthUm: number, nuFnu: number) => {
    const a = ebv == null ? null : extinctionMag(wavelengthUm, ebv);
    return a == null ? [null, null] : [a, nuFnu * 10 ** (0.4 * a)];
  };

  const own = points.map((p) => [
    p.isSelf ? "xwave" : "xwave_counterpart",
    p.catalog,
    `${p.survey} ${p.band}`,
    p.wavelengthUm,
    p.mag,
    p.magErr,
    p.upperLimit ? 1 : 0,
    p.nuFnu,
    p.nuFnuErr,
    p.separationArcsec,
    1,
    null,
    0,
    ...dereddened(p.wavelengthUm, p.nuFnu),
  ]);

  const fromVizier = vizier.map((p) => [
    "vizier",
    null,
    p.filter,
    p.wavelengthUm,
    null,
    null,
    0,
    p.nuFnu,
    p.nuFnuErr,
    null,
    p.nMeasurements,
    p.tables.join(" "),
    p.inconsistent ? 1 : 0,
    ...dereddened(p.wavelengthUm, p.nuFnu),
  ]);

  const rows = [...own, ...fromVizier].sort(
    (a, b) => (a[3] as number) - (b[3] as number)
  );
  return toCsv(HEADER, rows);
}
