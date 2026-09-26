import { describe, expect, it } from "vitest";

import type { SedPoint } from "./sed";
import { sedToCsv } from "./sedCsv";
import type { VizierSedPoint } from "./vizierSed";

const own: SedPoint = {
  band: "G",
  survey: "Gaia DR3",
  catalog: "gaia",
  wavelengthUm: 0.6218,
  mag: 18.68,
  magErr: 0.002,
  upperLimit: false,
  nuFnu: 5.2e-13,
  nuFnuErr: 1e-15,
  isSelf: true,
  separationArcsec: 0,
};

const viz: VizierSedPoint = {
  filter: "SDSS:g",
  wavelengthUm: 0.482,
  fluxJy: 9e-5,
  nuFnu: 5.5e-13,
  nMeasurements: 7,
  nRejected: 2,
  tables: ["II/294/sdss7", "V/147/sdss12"],
  inconsistent: false,
};

const parse = (csv: string) => csv.split("\n").map((l) => l.split(","));

describe("sedToCsv", () => {
  it("writes both sources sorted by wavelength", () => {
    const rows = parse(sedToCsv([own], [viz], null));
    expect(rows[0][0]).toBe("source");
    expect(rows.slice(1).map((r) => r[0])).toEqual(["vizier", "xwave"]);
    expect(rows[1][11]).toBe("II/294/sdss7 V/147/sdss12");
  });

  it("leaves dereddening empty without a usable E(B−V)", () => {
    const [, row] = parse(sedToCsv([own], [], null));
    expect(row[13]).toBe("");
    expect(row[14]).toBe("");
  });

  it("adds A_λ and the dereddened flux when E(B−V) is known", () => {
    const [, row] = parse(sedToCsv([own], [], 0.1));
    const a = Number(row[13]);
    expect(a).toBeGreaterThan(0.2);
    expect(Number(row[14])).toBeCloseTo(5.2e-13 * 10 ** (0.4 * a), 20);
  });
});
