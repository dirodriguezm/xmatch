"use client";

import { useMemo } from "react";

import { AladinViewer } from "@/app/components/object/AladinViewer";
import type { NearbySource } from "@/app/components/results/nearby/shared";
import {
  getSearchCatalogColor,
  getSearchCatalogLabel,
} from "@/app/lib/constants/catalogs";

/** Pan-STARRS (≈0.25″/px) north of −30°, SkyMapper (≈0.5″/px) south of it. */
const PS1_SURVEY = "CDS/P/PanSTARRS/DR1/color-z-zg-g";
const SKYMAPPER_SURVEY = "CDS/P/Skymapper/DR4/color";
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

  return (
    <div className="flex h-full min-h-0 flex-col gap-2">
      <AladinViewer
        center={center}
        // A little wider than the outermost ring.
        fov={(outer * 2.4) / 3600}
        survey={north ? PS1_SURVEY : SKYMAPPER_SURVEY}
        catalogLayers={layers}
        ringArcsec={outer}
        height="auto"
        className="min-h-[240px] flex-1 overflow-hidden rounded"
      />
      <p className="m-0 text-center text-xs text-neutral-500">
        {north ? "Pan-STARRS DR1" : "SkyMapper DR4"} colour · ring: {outer}″
        {" · "}
        {sources.length <= LABEL_MAX ? "labels: separation · " : ""}click a
        circle for details
      </p>
    </div>
  );
}
