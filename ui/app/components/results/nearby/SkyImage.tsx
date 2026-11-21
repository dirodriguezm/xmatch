"use client";

import { useMemo } from "react";

import { XWaveSpinner } from "@/app/components/common/XWaveSpinner";
import { AladinViewer } from "@/app/components/object/AladinViewer";
import type { NearbySource } from "@/app/components/results/nearby/shared";
import { useAvailableHips } from "@/app/hooks/queries";
import {
  getSearchCatalogColor,
  getSearchCatalogLabel,
} from "@/app/lib/constants/catalogs";

/**
 * Pan-STARRS (≈0.25″/px) north of −30°, SkyMapper (≈0.5″/px) south of it,
 * DSS wherever those aren't served. HiPS IDs, so Aladin can use mirrors:
 * PS1 i-r-g has one at MAST (z-zg-g is CDS-only); SkyMapper has none.
 */
const SURVEY_LABELS: Record<string, string> = {
  "CDS/P/PanSTARRS/DR1/color-i-r-g": "Pan-STARRS DR1",
  "CDS/P/Skymapper/DR4/color": "SkyMapper DR4",
  "CDS/P/DSS2/color": "DSS2",
};
const NORTH_SURVEYS = ["CDS/P/PanSTARRS/DR1/color-i-r-g", "CDS/P/DSS2/color"];
const SOUTH_SURVEYS = ["CDS/P/Skymapper/DR4/color", "CDS/P/DSS2/color"];
const PS1_MIN_DEC = -29.5;
/** Above this many circles, separation labels become clutter. */
const LABEL_MAX = 30;

interface SkyImageProps {
  target: { ra: number; dec: number };
  sources: NearbySource[];
  radii: Record<string, number>;
}

/**
 * The matches circled on a real sky image, one colour per catalog, with the
 * widest search radius as a dashed ring. Mounted only when chosen, so Aladin
 * (and its tiles) load on demand.
 */
export function SkyImage({ target, sources, radii }: SkyImageProps) {
  // Frame the catalogs that returned something: an empty 60″ eROSITA cone
  // would shrink every match to a speck.
  const outer = Math.max(
    1,
    ...[...new Set(sources.map((s) => s.slug))].map((slug) => radii[slug] ?? 0)
  );
  const north = target.dec > PS1_MIN_DEC;
  const hips = useAvailableHips(north ? NORTH_SURVEYS : SOUTH_SURVEYS);
  const layers = useMemo(
    () =>
      [...new Set(sources.map((s) => s.slug))].map((slug) => ({
        name: getSearchCatalogLabel(slug),
        color: getSearchCatalogColor(slug),
        labels: sources.length <= LABEL_MAX,
        sources: sources
          .filter((s) => s.slug === slug)
          .map((s) => ({
            ra: s.ra,
            dec: s.dec,
            name: s.objectId,
            data: {
              label: `${s.sepArcsec.toFixed(1)}″`,
              Catalog: getSearchCatalogLabel(slug),
              Separation: `${s.sepArcsec.toFixed(2)}″`,
              "Position angle": `${s.pa.toFixed(0)}° E of N`,
            },
          })),
      })),
    [sources]
  );
  // Stable object, so the viewer isn't re-initialised on every render.
  const center = useMemo(
    () => ({ ra: target.ra, dec: target.dec }),
    [target.ra, target.dec]
  );

  if (!hips.data) {
    return (
      <div className="flex h-full min-h-[240px] flex-col items-center justify-center gap-2 text-center">
        {hips.isPending ? (
          <XWaveSpinner size={24} />
        ) : (
          <>
            <p className="m-0 text-sm">Sky images unavailable</p>
            <p className="m-0 text-xs text-neutral-500">
              CDS isn&apos;t responding. The Offsets view shows the same
              matches.
            </p>
          </>
        )}
      </div>
    );
  }

  return (
    <div className="flex h-full min-h-0 flex-col gap-2">
      <AladinViewer
        center={center}
        // A little wider than the outermost ring.
        fov={(outer * 2.4) / 3600}
        survey={hips.data}
        catalogLayers={layers}
        ringArcsec={outer}
        height="auto"
        className="min-h-[240px] flex-1 overflow-hidden rounded"
      />
      <p className="m-0 text-center text-xs text-neutral-500">
        {SURVEY_LABELS[hips.data]} colour · ring: {outer}″{" · "}
        {sources.length <= LABEL_MAX ? "labels: separation · " : ""}click a
        circle for details
      </p>
    </div>
  );
}
