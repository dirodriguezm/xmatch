import { useCallback, useSyncExternalStore } from "react";

/** Same-tab listeners; other tabs are reached through the `storage` event. */
const listeners = new Set<() => void>();
/** This page's values, for when localStorage is unavailable. */
const memory = new Map<string, boolean>();

function subscribe(onChange: () => void) {
  listeners.add(onChange);
  window.addEventListener("storage", onChange);
  return () => {
    listeners.delete(onChange);
    window.removeEventListener("storage", onChange);
  };
}

function read(storageKey: string): boolean {
  try {
    const saved = localStorage.getItem(storageKey);
    if (saved !== null) return saved === "1";
  } catch {
    // Storage unavailable: fall back to this page's value.
  }
  return memory.get(storageKey) ?? false;
}

/**
 * A boolean UI preference remembered per browser (e.g. a collapsed panel).
 * `false` on the server and during hydration, so the markup matches; the
 * saved value applies right after.
 */
export function usePersistedFlag(
  storageKey: string
): [boolean, (next: boolean) => void] {
  const value = useSyncExternalStore(
    subscribe,
    () => read(storageKey),
    () => false
  );

  const set = useCallback(
    (next: boolean) => {
      memory.set(storageKey, next);
      try {
        localStorage.setItem(storageKey, next ? "1" : "0");
      } catch {
        // Not persisted; fine.
      }
      listeners.forEach((notify) => notify());
    },
    [storageKey]
  );

  return [value, set];
}
