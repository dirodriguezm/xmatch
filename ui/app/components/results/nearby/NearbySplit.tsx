"use client";

import { LinkedTable } from "./LinkedTable";
import { MapPanel } from "./MapPanel";
import type { NearbySource } from "./shared";
import { useLinkedHighlight } from "./useLinkedHighlight";

interface NearbyViewProps {
  sources: NearbySource[];
  radii: Record<string, number>;
}

/** A — offset map on the left, linked match table on the right. */
export function NearbySplit({ sources, radii }: NearbyViewProps) {
  const highlight = useLinkedHighlight();
  return (
    <div className="grid grid-cols-1 gap-6 lg:grid-cols-[340px_minmax(0,1fr)]">
      <div className="rounded-lg border border-border bg-surface p-3 lg:sticky lg:top-4 lg:self-start">
        <MapPanel sources={sources} radii={radii} highlight={highlight} />
      </div>
      <LinkedTable
        sources={sources}
        radii={radii}
        highlight={highlight}
        dense
        // Use the window's height, so long lists don't scroll in a small box
        // while the map stays put beside them.
        scrollY="max(360px, calc(100vh - 300px))"
      />
    </div>
  );
}
