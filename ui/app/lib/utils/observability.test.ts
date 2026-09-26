import { describe, expect, it } from "vitest";

import { OBSERVATORIES } from "@/app/lib/constants/observatories";

import {
  airmass,
  computeNightVisibility,
  tonightInTimeZone,
} from "./observability";

const site = (id: string) => OBSERVATORIES.find((o) => o.id === id)!;
const NIGHT = new Date(2026, 8, 26); // 26 Sep 2026: the Sun is at RA ≈ 12h

/** Local sidereal time in hours (IAU 1982 GMST), independent of the library. */
function lstHours(time: Date, longitudeDeg: number): number {
  const jd = time.getTime() / 86_400_000 + 2440587.5;
  const gmst = 18.697374558 + 24.06570982441908 * (jd - 2451545.0);
  return (((gmst + longitudeDeg / 15) % 24) + 24) % 24;
}

describe("airmass", () => {
  it("is 1 at the zenith and about 2 at 30°", () => {
    expect(airmass(90)).toBeCloseTo(1, 3);
    expect(airmass(30)).toBeCloseTo(1.995, 2);
  });

  it("is null at or below the horizon", () => {
    expect(airmass(0)).toBeNull();
    expect(airmass(-5)).toBeNull();
  });
});

describe("computeNightVisibility", () => {
  it("orders sunset, twilight and sunrise", () => {
    const v = computeNightVisibility(0, -30, site("paranal"), NIGHT);
    expect(v.sunset < v.duskAstronomical).toBe(true);
    expect(v.duskAstronomical < v.dawnAstronomical).toBe(true);
    expect(v.dawnAstronomical < v.sunrise).toBe(true);
    // Late September in Chile the Sun sets around 20:00 local (23:00 UTC).
    const sunsetUtcHour = v.sunset.getUTCHours();
    expect(sunsetUtcHour >= 22 || sunsetUtcHour === 0).toBe(true);
  });

  it("peaks at transit, near the zenith for dec ≈ site latitude", () => {
    const pachon = site("pachon");
    const v = computeNightVisibility(0, pachon.latitude, pachon, NIGHT);

    expect(v.best).not.toBeNull();
    expect(v.best!.altitude).toBeGreaterThan(88);
    expect(v.best!.airmass!).toBeLessThan(1.01);

    // At transit, local sidereal time equals the target's RA (0h here);
    // allow one sample step plus precession of RA to the equator of date.
    const lst = lstHours(v.best!.time, pachon.longitude);
    const offsetMin = Math.min(lst, 24 - lst) * 60;
    expect(offsetMin).toBeLessThan(10);
  });

  it("keeps the observable window inside astronomical night", () => {
    const v = computeNightVisibility(0, -30, site("lasilla"), NIGHT);
    expect(v.window).not.toBeNull();
    expect(v.window!.start >= v.duskAstronomical).toBe(true);
    expect(v.window!.end <= v.dawnAstronomical).toBe(true);
    for (const s of v.samples) {
      if (s.time >= v.window!.start && s.time <= v.window!.end) {
        expect(s.altitude).toBeGreaterThanOrEqual(30);
      }
    }
  });

  it("reports no window for a target that never rises high enough", () => {
    const v = computeNightVisibility(0, 70, site("paranal"), NIGHT);
    expect(v.window).toBeNull();
    expect(v.maxPossibleAltitude).toBeLessThan(0);
  });

  it("returns a moon illumination fraction and separation", () => {
    const v = computeNightVisibility(0, -30, site("paranal"), NIGHT);
    expect(v.moon.illumination).toBeGreaterThanOrEqual(0);
    expect(v.moon.illumination).toBeLessThanOrEqual(1);
    expect(v.moon.separationDeg).toBeGreaterThanOrEqual(0);
    expect(v.moon.separationDeg).toBeLessThanOrEqual(180);
  });
});

describe("tonightInTimeZone", () => {
  const tz = "America/Santiago";

  it("is today in the evening", () => {
    // 2026-09-26 21:00 in Chile (UTC−3)
    const d = tonightInTimeZone(tz, new Date("2026-09-27T00:00:00Z"));
    expect([d.getFullYear(), d.getMonth(), d.getDate()]).toEqual([2026, 8, 26]);
  });

  it("is still last night before dawn", () => {
    // 2026-09-27 03:00 in Chile
    const d = tonightInTimeZone(tz, new Date("2026-09-27T06:00:00Z"));
    expect([d.getFullYear(), d.getMonth(), d.getDate()]).toEqual([2026, 8, 26]);
  });

  it("moves to the coming night after the morning", () => {
    // 2026-09-27 11:00 in Chile
    const d = tonightInTimeZone(tz, new Date("2026-09-27T14:00:00Z"));
    expect([d.getFullYear(), d.getMonth(), d.getDate()]).toEqual([2026, 8, 27]);
  });
});
