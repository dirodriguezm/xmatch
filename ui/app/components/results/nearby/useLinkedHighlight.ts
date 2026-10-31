"use client";

import { useCallback, useState } from "react";

/**
 * Hover and selection shared between a map and a list. Hover is transient;
 * selection sticks until the same source is chosen again.
 */
export function useLinkedHighlight() {
  const [hoveredKey, setHoveredKey] = useState<string | null>(null);
  const [selectedKey, setSelectedKey] = useState<string | null>(null);

  const toggleSelected = useCallback(
    (key: string | null) =>
      setSelectedKey((prev) => (key === null || prev === key ? null : key)),
    []
  );

  return {
    hoveredKey,
    selectedKey,
    /** The source to emphasise: hover wins over selection. */
    activeKey: hoveredKey ?? selectedKey,
    setHoveredKey,
    setSelectedKey,
    toggleSelected,
  };
}

export type LinkedHighlight = ReturnType<typeof useLinkedHighlight>;
