import { describe, expect, it } from "vitest";

import { formatFlux, formatScientific } from "./format";

describe("formatScientific", () => {
  it("uses superscript exponents", () => {
    expect(formatScientific(1.4098e-10)).toBe("1.41 × 10⁻¹⁰");
    expect(formatScientific(3e20, 0)).toBe("3 × 10²⁰");
    expect(formatScientific(4.2)).toBe("4.20");
  });

  it("handles zero and non-finite values", () => {
    expect(formatScientific(0)).toBe("0");
    expect(formatScientific(NaN)).toBe("—");
  });
});

describe("formatFlux", () => {
  it("appends the flux unit", () => {
    expect(formatFlux(2.4366e-13)).toBe("2.44 × 10⁻¹³ erg s⁻¹ cm⁻²");
  });
});
