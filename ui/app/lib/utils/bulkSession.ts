import type { BulkResultRow, BulkTarget } from "@/app/lib/utils/bulkInput";

/**
 * The last bulk run, kept in sessionStorage so leaving /bulk for an object
 * page and coming back shows the same results instead of an empty form.
 */
export const BULK_RUN_STORAGE_KEY = "xwave:bulk-last-run";

/** Matches in the order the results table showed them, for ‹ › on /object. */
export const BULK_NAV_STORAGE_KEY = "xwave:bulk-nav";

export type BulkFilter = "all" | "matched" | "unmatched";
export type BulkSortKey = "row" | "input" | "matches" | "nearest";
export type BulkSortOrder = "ascend" | "descend";

/** How the results table was left: filter, sort and page. */
export interface BulkView {
  filter: BulkFilter;
  sortKey: BulkSortKey | null;
  sortOrder: BulkSortOrder | null;
  page: number;
  pageSize: number;
}

export const DEFAULT_BULK_VIEW: BulkView = {
  filter: "all",
  sortKey: null,
  sortOrder: null,
  page: 1,
  pageSize: 50,
};

export interface BulkRun {
  /** Input text as it was when the run started. */
  text: string;
  targets: BulkTarget[];
  radius: number;
  catalog: string;
  nneighbor: number;
  rows: BulkResultRow[];
  view: BulkView;
}

export interface BulkNavItem {
  id: string;
  catalog: string;
  /** 0-based input index, shown 1-based like the table's "#" column. */
  inputIndex: number;
  inputName: string;
}

const SORTERS: Record<
  BulkSortKey,
  (a: BulkResultRow, b: BulkResultRow) => number
> = {
  row: (a, b) => a.target.index - b.target.index,
  input: (a, b) => a.target.name.localeCompare(b.target.name),
  matches: (a, b) => a.matches.length - b.matches.length,
  nearest: (a, b) =>
    (a.matches[0]?.separation ?? Infinity) -
    (b.matches[0]?.separation ?? Infinity),
};

/** Rows as the table shows them: filtered, then sorted. */
export function orderBulkRows(
  rows: BulkResultRow[],
  view: Pick<BulkView, "filter" | "sortKey" | "sortOrder">
): BulkResultRow[] {
  const filtered =
    view.filter === "all"
      ? rows
      : rows.filter((r) =>
          view.filter === "matched"
            ? r.matches.length > 0
            : r.matches.length === 0
        );
  if (!view.sortKey || !view.sortOrder) return filtered;
  const cmp = SORTERS[view.sortKey];
  const sign = view.sortOrder === "ascend" ? 1 : -1;
  return [...filtered].sort((a, b) => sign * cmp(a, b));
}

/** Every match of every row, in table order (nearest first within a row). */
export function flattenBulkMatches(rows: BulkResultRow[]): BulkNavItem[] {
  return rows.flatMap((r) =>
    r.matches.map((m) => ({
      id: m.id,
      catalog: m.catalog,
      inputIndex: r.target.index,
      inputName: r.target.name,
    }))
  );
}

/** Position of an object in the nav list, with its neighbours. */
export function locateInBulkNav(
  items: BulkNavItem[],
  id: string,
  catalog: string
): {
  index: number;
  item: BulkNavItem;
  prev: BulkNavItem | null;
  next: BulkNavItem | null;
} | null {
  const cat = catalog.toLowerCase();
  const index = items.findIndex(
    (i) => i.id === id && i.catalog.toLowerCase() === cat
  );
  const item = items[index];
  if (!item) return null;
  return {
    index,
    item,
    prev: items[index - 1] ?? null,
    next: items[index + 1] ?? null,
  };
}

function read<T>(key: string): T | null {
  try {
    const raw = window.sessionStorage.getItem(key);
    return raw ? (JSON.parse(raw) as T) : null;
  } catch {
    return null;
  }
}

function write(key: string, value: unknown) {
  try {
    window.sessionStorage.setItem(key, JSON.stringify(value));
  } catch {
    // Storage full or blocked: coming back just shows an empty form.
  }
}

export function loadBulkRun(): BulkRun | null {
  const run = read<BulkRun>(BULK_RUN_STORAGE_KEY);
  return run && Array.isArray(run.rows) && Array.isArray(run.targets)
    ? { ...run, view: { ...DEFAULT_BULK_VIEW, ...run.view } }
    : null;
}

export function saveBulkRun(run: BulkRun) {
  write(BULK_RUN_STORAGE_KEY, run);
}

/** Update only the table view of the stored run. */
export function saveBulkView(view: BulkView) {
  const run = loadBulkRun();
  if (run) saveBulkRun({ ...run, view });
}

export function loadBulkNav(): BulkNavItem[] {
  const items = read<BulkNavItem[]>(BULK_NAV_STORAGE_KEY);
  return Array.isArray(items) ? items : [];
}

export function saveBulkNav(items: BulkNavItem[]) {
  write(BULK_NAV_STORAGE_KEY, items);
}
