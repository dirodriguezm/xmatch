import { describe, expect, it } from "vitest";

import {
  addItem,
  BASKET_MAX_ITEMS,
  basketToBulkInput,
  basketToCsv,
  hasItem,
  parseBasket,
  removeItem,
  toggleItem,
} from "./basket";

const a = { objectId: "Gaia DR3 1", catalog: "gaia", ra: 10, dec: -5 };
const b = {
  objectId: "0098p408_ac51-043708",
  catalog: "allwise",
  ra: 1,
  dec: 2,
};

describe("basket helpers", () => {
  it("adds without duplicates, keyed by catalog + id", () => {
    const one = addItem([], a);
    expect(addItem(one, { ...a })).toBe(one);
    expect(addItem(one, { ...a, catalog: "GAIA" })).toBe(one);
    expect(addItem(one, { ...a, catalog: "allwise" })).toHaveLength(2);
  });

  it("toggles and removes", () => {
    const one = toggleItem([], a);
    expect(hasItem(one, a)).toBe(true);
    expect(toggleItem(one, a)).toEqual([]);
    expect(removeItem([a, b], a)).toEqual([b]);
  });

  it("caps the basket size", () => {
    const full = Array.from({ length: BASKET_MAX_ITEMS }, (_, i) => ({
      ...a,
      objectId: `x${i}`,
    }));
    expect(addItem(full, b)).toBe(full);
  });

  it("parses stored JSON defensively", () => {
    expect(parseBasket(null)).toEqual([]);
    expect(parseBasket("not json")).toEqual([]);
    expect(parseBasket('{"a":1}')).toEqual([]);
    const raw = JSON.stringify([a, a, { objectId: 3 }, { ...b, extra: 1 }]);
    expect(parseBasket(raw)).toEqual([a, b]);
  });

  it("serialises to CSV and bulk input", () => {
    expect(basketToCsv([a])).toBe(
      "object_id,catalog,ra_deg,dec_deg\nGaia DR3 1,gaia,10,-5"
    );
    expect(basketToBulkInput([a, { ...b, objectId: "weird,name" }])).toBe(
      "Gaia DR3 1,10,-5\nweird name,1,2"
    );
  });
});
