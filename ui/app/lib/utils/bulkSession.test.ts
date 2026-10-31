import { describe, expect, it } from "vitest";

import type { BulkResultRow } from "./bulkInput";
import {
  flattenBulkMatches,
  locateInBulkNav,
  orderBulkRows,
} from "./bulkSession";

function row(
  index: number,
  name: string,
  matches: [string, string, number][]
): BulkResultRow {
  return {
    target: { index, line: index + 1, name, ra: index, dec: 0 },
    matches: matches.map(([catalog, id, separation]) => ({
      catalog,
      id,
      ra: null,
      dec: null,
      separation,
    })),
  };
}

const rows = [
  row(0, "b", [
    ["gaia", "g1", 0.2],
    ["allwise", "w1", 0.9],
  ]),
  row(1, "a", []),
  row(2, "c", [["erosita", "e1", 0.1]]),
];

describe("orderBulkRows", () => {
  it("filters before sorting", () => {
    const out = orderBulkRows(rows, {
      filter: "matched",
      sortKey: "nearest",
      sortOrder: "ascend",
    });
    expect(out.map((r) => r.target.index)).toEqual([2, 0]);
  });

  it("keeps input order without a sort", () => {
    const out = orderBulkRows(rows, {
      filter: "all",
      sortKey: null,
      sortOrder: null,
    });
    expect(out).toBe(rows);
  });

  it("reverses for descend", () => {
    const out = orderBulkRows(rows, {
      filter: "all",
      sortKey: "input",
      sortOrder: "descend",
    });
    expect(out.map((r) => r.target.name)).toEqual(["c", "b", "a"]);
  });
});

describe("bulk nav", () => {
  const items = flattenBulkMatches(rows);

  it("flattens every match in table order", () => {
    expect(items.map((i) => i.id)).toEqual(["g1", "w1", "e1"]);
    expect(items[1]).toEqual({
      id: "w1",
      catalog: "allwise",
      inputIndex: 0,
      inputName: "b",
    });
  });

  it("finds neighbours, matching the catalog case-insensitively", () => {
    const at = locateInBulkNav(items, "w1", "AllWISE");
    expect(at?.index).toBe(1);
    expect(at?.prev?.id).toBe("g1");
    expect(at?.next?.id).toBe("e1");
  });

  it("has no neighbour past the ends", () => {
    expect(locateInBulkNav(items, "g1", "gaia")?.prev).toBeNull();
    expect(locateInBulkNav(items, "e1", "erosita")?.next).toBeNull();
  });

  it("returns null for objects not in the run", () => {
    expect(locateInBulkNav(items, "g1", "allwise")).toBeNull();
  });
});
