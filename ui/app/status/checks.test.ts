import { describe, expect, it } from "vitest";

import { classifyCheck, DEGRADED_LATENCY_MS, overallState } from "./checks";

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
