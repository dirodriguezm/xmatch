/**
 * Object basket: a small client-side collection of catalog objects the user
 * has pinned while browsing. Persisted in localStorage and exposed through a
 * subscribe/getSnapshot pair so components can read it with
 * `useSyncExternalStore` without a provider.
 */

import type { BasketItem } from "@/app/components/basket/types";
import { toCsv } from "@/app/lib/utils/csv";

export const BASKET_STORAGE_KEY = "xwave:basket";

/** sessionStorage key the bulk page reads its input from. */
export const BULK_INPUT_STORAGE_KEY = "xwave:bulk-input";

/** Hard cap so a runaway basket never bloats localStorage. */
export const BASKET_MAX_ITEMS = 500;

export function basketKey(item: Pick<BasketItem, "objectId" | "catalog">) {
  return `${item.catalog.toLowerCase()}::${item.objectId}`;
}

function isBasketItem(value: unknown): value is BasketItem {
  if (!value || typeof value !== "object") return false;
  const v = value as Record<string, unknown>;
  return (
    typeof v.objectId === "string" &&
    typeof v.catalog === "string" &&
    typeof v.ra === "number" &&
    Number.isFinite(v.ra) &&
    typeof v.dec === "number" &&
    Number.isFinite(v.dec)
  );
}

/** Parse a stored basket, dropping malformed entries and duplicates. */
export function parseBasket(raw: string | null): BasketItem[] {
  if (!raw) return [];
  try {
    const parsed: unknown = JSON.parse(raw);
    if (!Array.isArray(parsed)) return [];
    const seen = new Set<string>();
    const items: BasketItem[] = [];
    for (const entry of parsed) {
      if (!isBasketItem(entry)) continue;
      const key = basketKey(entry);
      if (seen.has(key)) continue;
      seen.add(key);
      items.push({
        objectId: entry.objectId,
        catalog: entry.catalog,
        ra: entry.ra,
        dec: entry.dec,
      });
    }
    return items.slice(0, BASKET_MAX_ITEMS);
  } catch {
    return [];
  }
}

export function hasItem(items: BasketItem[], item: BasketItem): boolean {
  const key = basketKey(item);
  return items.some((i) => basketKey(i) === key);
}

export function addItem(items: BasketItem[], item: BasketItem): BasketItem[] {
  if (hasItem(items, item) || items.length >= BASKET_MAX_ITEMS) return items;
  return [...items, item];
}

export function removeItem(
  items: BasketItem[],
  item: Pick<BasketItem, "objectId" | "catalog">
): BasketItem[] {
  const key = basketKey(item);
  return items.filter((i) => basketKey(i) !== key);
}

export function toggleItem(
  items: BasketItem[],
  item: BasketItem
): BasketItem[] {
  return hasItem(items, item) ? removeItem(items, item) : addItem(items, item);
}

export function basketToCsv(items: BasketItem[]): string {
  return toCsv(
    ["object_id", "catalog", "ra_deg", "dec_deg"],
    items.map((i) => [i.objectId, i.catalog, i.ra, i.dec])
  );
}

/**
 * "name,ra,dec" lines for the bulk cross-match input. Commas in names would
 * break the column split, so they are replaced with spaces.
 */
export function basketToBulkInput(items: BasketItem[]): string {
  return items
    .map((i) => `${i.objectId.replace(/,/g, " ")},${i.ra},${i.dec}`)
    .join("\n");
}

/** Tab-separated plain-text list for the clipboard. */
export function basketToText(items: BasketItem[]): string {
  return items
    .map((i) => `${i.objectId}\t${i.catalog}\t${i.ra}\t${i.dec}`)
    .join("\n");
}

// ---------------------------------------------------------------------------
// Store
// ---------------------------------------------------------------------------

const EMPTY: BasketItem[] = [];
let cache: BasketItem[] | null = null;
const listeners = new Set<() => void>();

function readStorage(): BasketItem[] {
  if (typeof window === "undefined") return EMPTY;
  try {
    return parseBasket(window.localStorage.getItem(BASKET_STORAGE_KEY));
  } catch {
    return EMPTY;
  }
}

function emit() {
  listeners.forEach((l) => l());
}

function write(next: BasketItem[]) {
  cache = next;
  try {
    window.localStorage.setItem(BASKET_STORAGE_KEY, JSON.stringify(next));
  } catch {
    // Storage full or blocked: the in-memory basket still works this session.
  }
  emit();
}

function onStorage(e: StorageEvent) {
  if (e.key !== null && e.key !== BASKET_STORAGE_KEY) return;
  cache = parseBasket(e.newValue);
  emit();
}

export const basketStore = {
  subscribe(listener: () => void): () => void {
    listeners.add(listener);
    if (listeners.size === 1 && typeof window !== "undefined") {
      window.addEventListener("storage", onStorage);
    }
    return () => {
      listeners.delete(listener);
      if (listeners.size === 0 && typeof window !== "undefined") {
        window.removeEventListener("storage", onStorage);
      }
    };
  },
  getSnapshot(): BasketItem[] {
    if (cache === null) cache = readStorage();
    return cache;
  },
  getServerSnapshot(): BasketItem[] {
    return EMPTY;
  },
  toggle(item: BasketItem) {
    write(toggleItem(basketStore.getSnapshot(), item));
  },
  add(item: BasketItem) {
    write(addItem(basketStore.getSnapshot(), item));
  },
  remove(item: Pick<BasketItem, "objectId" | "catalog">) {
    write(removeItem(basketStore.getSnapshot(), item));
  },
  clear() {
    write(EMPTY);
  },
};
