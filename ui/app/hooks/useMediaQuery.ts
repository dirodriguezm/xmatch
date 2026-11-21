import { useCallback, useSyncExternalStore } from "react";

/**
 * Whether a CSS media query matches, kept in sync with the window. Use it
 * when JS must agree exactly with a Tailwind breakpoint (antd's `Grid`
 * breakpoints differ: its `lg` is 992px, Tailwind's is 1024px). `serverValue`
 * is what the server and first render assume.
 */
export function useMediaQuery(query: string, serverValue = true): boolean {
  const subscribe = useCallback(
    (onChange: () => void) => {
      const mql = window.matchMedia(query);
      mql.addEventListener("change", onChange);
      return () => mql.removeEventListener("change", onChange);
    },
    [query]
  );
  return useSyncExternalStore(
    subscribe,
    () => window.matchMedia(query).matches,
    () => serverValue
  );
}
