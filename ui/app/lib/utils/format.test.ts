import { describe, expect, it } from "vitest";

import {
  formatFlux,
  formatMjdDate,
  formatMjdDateTime,
  formatScientific,
} from "./format";

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

describe("formatMjdDate", () => {
  it("converts MJD to the UTC calendar date and time", () => {
    expect(formatMjdDate(51544.5)).toBe("2000-01-01");
    expect(formatMjdDateTime(51544.5)).toBe("2000-01-01 12:00 UTC");
    expect(formatMjdDateTime(0)).toBe("1858-11-17 00:00 UTC");
  });
});
