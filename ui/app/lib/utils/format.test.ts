import { describe, expect, it } from "vitest";

import {
  formatCatalogValue,
  formatFlux,
  formatFluxWithError,
  formatMjdDate,
  formatMjdDateTime,
  formatMjdTick,
  formatScientific,
  mjdYear,
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

describe("formatFluxWithError", () => {
  it("puts value and error on one power of ten, unit last", () => {
    expect(formatFluxWithError(1.4168e-11, 7.3e-14)).toBe(
      "(1.4168 ± 0.0073) × 10⁻¹¹ erg s⁻¹ cm⁻²"
    );
    expect(formatFluxWithError(3.2e-13, 4.1e-14)).toBe(
      "(3.20 ± 0.41) × 10⁻¹³ erg s⁻¹ cm⁻²"
    );
  });

  it("shows small errors to four decimals, caps tinier ones", () => {
    expect(formatFluxWithError(1.41e-10, 2.8e-14)).toBe(
      "(1.4100 ± 0.0003) × 10⁻¹⁰ erg s⁻¹ cm⁻²"
    );
    expect(formatFluxWithError(1.41e-10, 2e-16)).toBe(
      "(1.4100 ± < 0.0001) × 10⁻¹⁰ erg s⁻¹ cm⁻²"
    );
  });

  it("falls back to the bare flux without a usable error", () => {
    expect(formatFluxWithError(1.41e-10, 0)).toBe("1.41 × 10⁻¹⁰ erg s⁻¹ cm⁻²");
  });
});

describe("formatMjdTick", () => {
  it("adds the year under the MJD", () => {
    expect(formatMjdTick(57000)).toBe("57000\n2014");
    expect(mjdYear(51544.5)).toBe(2000);
  });
});

describe("formatCatalogValue", () => {
  it("rounds floats to six significant figures", () => {
    expect(formatCatalogValue(0.018106548115611076)).toBe("0.0181065");
    expect(formatCatalogValue(-0.015950290834234743)).toBe("-0.0159503");
    expect(formatCatalogValue(12.84)).toBe("12.84");
  });

  it("uses scientific notation for very small or large floats", () => {
    expect(formatCatalogValue(1.4168050808297e-11)).toBe("1.4168 × 10⁻¹¹");
  });

  it("leaves integers, id strings and booleans alone, dashes nulls", () => {
    expect(formatCatalogValue(43)).toBe("43");
    expect(formatCatalogValue("3700386905605055360")).toBe(
      "3700386905605055360"
    );
    expect(formatCatalogValue(false)).toBe("false");
    expect(formatCatalogValue(null)).toBe("—");
  });
});
