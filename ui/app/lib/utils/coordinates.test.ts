import { describe, expect, it } from "vitest";

import {
  equatorialToEcliptic,
  equatorialToGalactic,
  positionAngle,
} from "./coordinates";

describe("equatorialToGalactic", () => {
  it("puts Sgr A* at the Galactic centre", () => {
    const { l, b } = equatorialToGalactic(266.41683, -29.00781);
    expect(Math.min(l, 360 - l)).toBeLessThan(0.1);
    expect(Math.abs(b)).toBeLessThan(0.1);
  });

  it("maps the north Galactic pole to b = +90°", () => {
    expect(equatorialToGalactic(192.85948, 27.12825).b).toBeCloseTo(90, 3);
  });

  it("agrees with a known field (COSMOS: l ≈ 236.8°, b ≈ +42.1°)", () => {
    const { l, b } = equatorialToGalactic(150.1191, 2.2058);
    expect(l).toBeCloseTo(236.8, 0);
    expect(b).toBeCloseTo(42.1, 0);
  });
});

describe("equatorialToEcliptic", () => {
  it("keeps the vernal equinox at the origin", () => {
    const { lambda, beta } = equatorialToEcliptic(0, 0);
    expect(lambda).toBeCloseTo(0, 6);
    expect(beta).toBeCloseTo(0, 6);
  });

  it("puts the June solstice point at λ = 90°, β = 0", () => {
    const { lambda, beta } = equatorialToEcliptic(90, 23.4392911);
    expect(lambda).toBeCloseTo(90, 4);
    expect(beta).toBeCloseTo(0, 4);
  });

  it("maps the north ecliptic pole to β = +90°", () => {
    expect(equatorialToEcliptic(270, 66.5607089).beta).toBeCloseTo(90, 4);
  });
});

describe("positionAngle", () => {
  it("is 0° due north and 90° due east", () => {
    expect(positionAngle(10, 20, 10, 20.01)).toBeCloseTo(0, 6);
    expect(positionAngle(10, 20, 10.01, 20)).toBeCloseTo(90, 1);
  });

  it("is 180° due south and 270° due west", () => {
    expect(positionAngle(10, 20, 10, 19.99)).toBeCloseTo(180, 6);
    expect(positionAngle(10, 20, 9.99, 20)).toBeCloseTo(270, 1);
  });

  it("handles RA wrap-around at 0°", () => {
    expect(positionAngle(359.995, 0, 0.005, 0)).toBeCloseTo(90, 6);
  });
});
