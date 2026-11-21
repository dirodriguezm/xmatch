import { useEffect, useState } from "react";

/**
 * True once `active` (e.g. a query loading) has lasted `afterMs`, so the UI
 * can explain a long wait instead of spinning silently. Resets when it ends.
 */
export function useIsSlow(active: boolean, afterMs: number): boolean {
  const [slow, setSlow] = useState(false);

  useEffect(() => {
    if (!active) return;
    const timer = setTimeout(() => setSlow(true), afterMs);
    return () => {
      clearTimeout(timer);
      setSlow(false);
    };
  }, [active, afterMs]);

  return active && slow;
}
