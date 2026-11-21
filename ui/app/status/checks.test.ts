import { describe, expect, it } from "vitest";

import {
  classifyCheck,
  DEGRADED_LATENCY_MS,
  mirrorInUse,
  overallState,
} from "./checks";

describe("classifyCheck", () => {
  it("treats network errors and 5xx as down", () => {
    expect(classifyCheck(null, null)).toBe("down");
    expect(classifyCheck(502, 100)).toBe("down");
    expect(classifyCheck(500, 100)).toBe("down");
  });

  it("treats fast 2xx and 404 as up", () => {
    expect(classifyCheck(200, 120)).toBe("up");
    expect(classifyCheck(404, 120)).toBe("up");
  });

  it("flags slow answers as degraded", () => {
    expect(classifyCheck(200, DEGRADED_LATENCY_MS + 1)).toBe("degraded");
  });
});

describe("overallState", () => {
  it("reports the worst state", () => {
    expect(overallState(["up", "degraded"])).toBe("degraded");
    expect(overallState(["up", "down", "degraded"])).toBe("down");
    expect(overallState([])).toBe("up");
  });
});

describe("mirrorInUse", () => {
  const check = { primaryHost: "simbad.cds.unistra.fr" };

  it("names the host when a mirror answered", () => {
    expect(mirrorInUse(check, "simbad.cfa.harvard.edu")).toBe(
      "simbad.cfa.harvard.edu"
    );
  });

  it("is null for the primary, unknown hosts or unmirrored checks", () => {
    expect(mirrorInUse(check, "simbad.cds.unistra.fr")).toBeNull();
    expect(mirrorInUse(check, null)).toBeNull();
    expect(mirrorInUse({}, "anything.test")).toBeNull();
  });
});
