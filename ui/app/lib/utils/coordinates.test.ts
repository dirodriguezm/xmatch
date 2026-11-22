import { describe, expect, it } from "vitest";

import {
  equatorialToEcliptic,
  equatorialToGalactic,
  formatDegreesParam,
  parseCoordinates,
  parseDecInput,
  parseRaInput,
  positionAngle,
} from "./coordinates";

// The Orion Nebula: 05h35m17.3s -05°23′28″ = 83.82208°, -5.39111°.
const ORION_RA = 83.822083;
const ORION_DEC = -5.391111;

describe("parseRaInput", () => {
  it.each([
    ["83.8221", 83.8221],
    ["  83.8221 ", 83.8221],
    ["0", 0],
    ["360", 360],
    ["83.8221°", 83.8221],
    ["83d", 83],
  ])("reads %j as decimal degrees", (input, deg) => {
    expect(parseRaInput(input).value).toBeCloseTo(deg, 6);
  });

  it.each([
    "05:35:17.3",
    "5:35:17.3",
    "05 35 17.3",
    "05h35m17.3s",
    "05h 35m 17.3s",
    "5h35m17.3",
  ])("reads %j as hours", (input) => {
    expect(parseRaInput(input).value).toBeCloseTo(ORION_RA, 5);
  });

  it("reads h:m with fractional minutes", () => {
    expect(parseRaInput("05:35.5").value).toBeCloseTo((5 + 35.5 / 60) * 15, 6);
  });

  it("reads a lone 'h' value as hours", () => {
    expect(parseRaInput("5h").value).toBe(75);
  });

  it("does not read h:m:s as degrees (the old parseFloat bug)", () => {
    expect(parseRaInput("05:35:17.3").value).not.toBe(5);
  });

  it.each([
    ["400", "between 0 and 360"],
    ["-1", "between 0 and 360"],
    ["24:00:00", "below 24h"],
    ["-05:35:17", "below 24h"],
  ])("rejects out-of-range %j", (input, message) => {
    const result = parseRaInput(input);
    expect(result.value).toBeNull();
    expect(result.error).toContain(message);
  });

  it.each([
    "abc",
    "M31",
    "83.82abc",
    "1e3",
    "5.5.5",
    "05:60:00",
    "05:35:60",
    "05.5:30:00",
    "05:35:17:01",
    "05:",
    "-",
  ])("rejects malformed %j with the format hint", (input) => {
    const result = parseRaInput(input);
    expect(result.value).toBeNull();
    expect(result.error).toContain("h:m:s");
  });

  it("asks for a value when empty", () => {
    expect(parseRaInput("  ").error).toBe("Enter an RA");
  });
});

describe("parseDecInput", () => {
  it.each([
    ["-5.3911", -5.3911],
    ["+41.269", 41.269],
    ["90", 90],
    ["-90", -90],
    ["\u22125.3911", -5.3911],
  ])("reads %j as decimal degrees", (input, deg) => {
    expect(parseDecInput(input).value).toBeCloseTo(deg, 6);
  });

  it.each([
    "-05:23:28",
    "-5:23:28",
    "-05 23 28",
    "-05d23m28s",
    "-5°23′28″",
    "- 05:23:28",
  ])("reads %j as d:m:s", (input) => {
    expect(parseDecInput(input).value).toBeCloseTo(ORION_DEC, 5);
  });

  it("keeps the sign of a declination just south of the equator", () => {
    expect(parseDecInput("-00:30:00").value).toBe(-0.5);
  });

  it.each(["91", "-90.5", "90:00:01", "-91:00:00"])(
    "rejects out-of-range %j",
    (input) => {
      expect(parseDecInput(input).error).toBe(
        "Dec must be between -90 and +90°"
      );
    }
  );

  it.each(["abc", "05h23m", "-05:23:60", "41.2x"])(
    "rejects malformed %j with the format hint",
    (input) => {
      expect(parseDecInput(input).error).toContain("d:m:s");
    }
  );

  it("asks for a value when empty", () => {
    expect(parseDecInput("").error).toBe("Enter a Dec");
  });
});

describe("formatDegreesParam", () => {
  it("rounds to 7 decimals and drops trailing zeros", () => {
    expect(formatDegreesParam(83.82208333333334)).toBe("83.8220833");
    expect(formatDegreesParam(10.6847)).toBe("10.6847");
    expect(formatDegreesParam(-5)).toBe("-5");
  });

  it("never writes an exponent", () => {
    expect(formatDegreesParam(-5e-7)).toBe("-0.0000005");
    expect(formatDegreesParam(1e-7)).toBe("0.0000001");
    expect(formatDegreesParam(-1e-9)).toBe("0");
    expect(formatDegreesParam(360)).toBe("360");
  });
});

describe("parseCoordinates", () => {
  it.each([
    "83.822083 -5.391111",
    "83.822083, -5.391111",
    "05:35:17.3 -05:23:28",
    "05h35m17.3s -05d23m28s",
    "05 35 17.3 -05 23 28",
  ])("reads %j", (input) => {
    const coords = parseCoordinates(input);
    expect(coords?.ra).toBeCloseTo(ORION_RA, 5);
    expect(coords?.dec).toBeCloseTo(ORION_DEC, 5);
  });

  it.each(["M31", "M 31", "NGC 1952", "3C 273", "400 10", "10 95"])(
    "leaves %j to the name resolver",
    (input) => {
      expect(parseCoordinates(input)).toBeNull();
    }
  );

  // Lines pasted from a table or a list often carry more than RA and Dec.
  it.each([
    "83.822083 -5.391111,",
    ",83.822083 -5.391111",
    "83.822083, -5.391111, 5",
    "83.822083 -5.391111 5",
    "83.822083 -5.391111 2000",
    "83.822083 -5.391111 J2000",
    "83.822083 -5.391111 b1950.0",
    "05:35:17.3 -05:23:28 J2000",
    "05h35m17.3s -05d23m28s 0.5",
  ])("reads %j, ignoring the extra column", (input) => {
    const coords = parseCoordinates(input);
    expect(coords?.ra).toBeCloseTo(ORION_RA, 5);
    expect(coords?.dec).toBeCloseTo(ORION_DEC, 5);
  });

  it("reads a comma-separated line of whole degrees", () => {
    expect(parseCoordinates("83, -5, 5")).toEqual({ ra: 83, dec: -5 });
  });

  it.each([
    // Could be the start of h m s, so not read as RA 12, Dec 30.
    "12 30 45",
    "83.822083 -5 5",
    "12 30, 45",
    // The third column must be a number or an epoch.
    "83.822083 -5.391111 extra",
    "83.822083 -5.391111 1e3",
    "83.822083 -5.391111 5.5.5",
    "M31 10.5 5",
    "1e3 10.5 5",
    "5.5.5 10.5 5",
    "83.822083 -5.391111 5 6",
  ])("still rejects %j", (input) => {
    expect(parseCoordinates(input)).toBeNull();
  });
});

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
