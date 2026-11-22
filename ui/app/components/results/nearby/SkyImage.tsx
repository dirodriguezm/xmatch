"use client";

import { useCallback, useEffect, useMemo, useRef, useState } from "react";

import { XWaveSpinner } from "@/app/components/common/XWaveSpinner";
import { ServiceStatusLink } from "@/app/components/layout/ServiceStatus";
import { AladinViewer } from "@/app/components/object/AladinViewer";
import type { NearbySource } from "@/app/components/results/nearby/shared";
import { useHipsAt } from "@/app/hooks/queries";
import { useTileState } from "@/app/hooks/queries/useSkyImages";
import {
  getSearchCatalogColor,
  getSearchCatalogLabel,
} from "@/app/lib/constants/catalogs";
import type {
  AladinCatalogLayer,
  AladinGlobal,
  AladinInstance,
} from "@/types/aladin";

/**
 * Pan-STARRS (≈0.25″/px) north of −30°, SkyMapper (≈0.5″/px) south of it,
 * DSS wherever those aren't served at the target. HiPS IDs, so mirrors can
 * stand in: PS1 i-r-g has one at MAST (z-zg-g is CDS-only, and MAST's copy
 * has holes); SkyMapper has none.
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
/** Share of the view's width the widest ring spans, leaving sky around it. */
const RING_SHARE = 0.55;
/**
 * Never zoom in further than this. Narrower, a bright core (M31, the Orion
 * Trapezium) fills the frame as flat white and the survey pixels turn to blur.
 */
const MIN_FOV_ARCSEC = 60;
/**
 * How far the view may drift from the last automatic framing (as a share of
 * the field, or of the pixel scale) before it counts as the user's own.
 */
const MOVED_SHARE = 0.01;

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
  const hips = useHipsAt(
    north ? NORTH_SURVEYS : SOUTH_SURVEYS,
    target.ra,
    target.dec
  );
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

  if (!hips.data) {
    return (
      <div className="flex h-full min-h-[240px] flex-col items-center justify-center gap-2 text-center">
        {hips.isPending ? (
          <XWaveSpinner size={24} label="Loading sky image" />
        ) : (
          <Unavailable />
        )}
      </div>
    );
  }

  return (
    <div className="flex h-full min-h-0 flex-col gap-2">
      <SkyView
        // Aladin reads its centre and field only when it starts, so a new
        // search or radius gets a fresh viewer, not the old view with the
        // ring and circles drawn off-screen.
        key={`${target.ra},${target.dec},${outer},${hips.data.url}`}
        ra={target.ra}
        dec={target.dec}
        ringArcsec={outer}
        survey={hips.data.url}
        tileUrl={hips.data.tileUrl}
        layers={layers}
      />
      <p className="m-0 text-center text-xs text-neutral-500">
        {SURVEY_LABELS[hips.data.id]} colour · ring: {outer}″{" · "}
        {sources.length <= LABEL_MAX ? "labels: separation · " : ""}click a
        circle for details
      </p>
    </div>
  );
}

interface SkyViewProps {
  ra: number;
  dec: number;
  ringArcsec: number;
  survey: string;
  /** A tile under the target, probed to tell when the image is in. */
  tileUrl: string;
  layers: AladinCatalogLayer[];
}

/** One Aladin view of one position; remounted (via `key`) when that changes. */
function SkyView({
  ra,
  dec,
  ringArcsec,
  survey,
  tileUrl,
  layers,
}: SkyViewProps) {
  const [started, setStarted] = useState(false);
  const tile = useTileState(tileUrl);
  // Stable object, so the viewer isn't re-initialised on every render.
  const center = useMemo(() => ({ ra, dec }), [ra, dec]);
  // Aladin's field is the view's width.
  const fovDeg = Math.max(MIN_FOV_ARCSEC, (2 * ringArcsec) / RING_SHARE) / 3600;
  const resizeRef = useRef<ResizeObserver | null>(null);
  const frameRef = useRef(0);
  /** Pixel scale (degrees per CSS pixel) the last automatic framing left. */
  const framedScaleRef = useRef<number | null>(null);
  useEffect(
    () => () => {
      resizeRef.current?.disconnect();
      cancelAnimationFrame(frameRef.current);
    },
    []
  );

  // The ring is drawn here rather than through `ringArcsec`: a dark stroke
  // under the light dashes keeps it visible on empty sky and on a saturated
  // core alike, where a white ring alone vanishes.
  const onReady = useCallback(
    (aladin: AladinInstance) => {
      setStarted(true);
      void import("aladin-lite").then((mod) => {
        const A = mod.default as AladinGlobal;
        // Aladin applies `fov` to its default 300px canvas, then resizes the
        // canvas to the panel a frame later keeping that pixel scale, so the
        // field came out ~3x tighter than asked and a 20″ ring fell outside
        // the panel. Re-frame once the canvas has its real size, and again
        // whenever the panel resizes (e.g. the map folding away), unless the
        // user has zoomed or panned since: a resize mustn't undo that.
        const reframe = () => {
          cancelAnimationFrame(frameRef.current);
          let frames = 0;
          const settle = () => {
            const div = aladin.aladinDiv;
            const canvas = div?.querySelector("canvas");
            const want = Math.round((div?.clientWidth ?? 0) * devicePixelRatio);
            if (canvas && want > 0 && Math.abs(canvas.width - want) > 2) {
              if (++frames < 120) {
                frameRef.current = requestAnimationFrame(settle);
                return;
              }
            }
            const width = div?.clientWidth ?? 0;
            const framedScale = framedScaleRef.current;
            if (framedScale !== null && width > 0) {
              // A resize keeps Aladin's pixel scale and centre, so either
              // changing means the user moved the view.
              const fov = aladin.getFov()[0] ?? fovDeg;
              const [viewRa, viewDec] = aladin.getRaDec();
              const zoomed =
                Math.abs(fov / width / framedScale - 1) > MOVED_SHARE;
              const panned =
                separationDeg(ra, dec, viewRa, viewDec) > MOVED_SHARE * fov;
              if (zoomed || panned) return;
            }
            aladin.setFov(fovDeg);
            if (width > 0) {
              framedScaleRef.current = (aladin.getFov()[0] ?? fovDeg) / width;
            }
            aladin.view?.requestRedraw?.();
          };
          frameRef.current = requestAnimationFrame(settle);
        };
        resizeRef.current?.disconnect();
        if (aladin.aladinDiv) {
          const observer = new ResizeObserver(reframe);
          observer.observe(aladin.aladinDiv);
          resizeRef.current = observer;
        }
        reframe();
        const radiusDeg = ringArcsec / 3600;
        const halo = A.graphicOverlay({
          color: "rgba(0, 0, 0, 0.6)",
          lineWidth: 3.5,
        });
        const dashes = A.graphicOverlay({
          color: "rgba(255, 255, 255, 0.85)",
          lineWidth: 1.5,
          lineDash: [5, 4],
        });
        aladin.addOverlay(halo);
        halo.add(A.circle(ra, dec, radiusDeg));
        aladin.addOverlay(dashes);
        dashes.add(A.circle(ra, dec, radiusDeg));
        // Aladin only repaints overlays when the view moves or a tile lands,
        // so with the tiles already in the ring would wait for a zoom.
        aladin.view?.requestRedraw?.();
      });
    },
    [ra, dec, ringArcsec, fovDeg]
  );

  return (
    <div className="relative flex min-h-[240px] flex-1 flex-col overflow-hidden rounded">
      <AladinViewer
        center={center}
        fov={fovDeg}
        survey={survey}
        catalogLayers={layers}
        onReady={onReady}
        height="auto"
        // Its own stacking context, so Aladin's z-indexed controls stay
        // under the loading overlay.
        className="isolate min-h-0 flex-1"
      />
      {tile === "failed" ? (
        <div className="absolute inset-0 flex flex-col items-center justify-center gap-2 bg-surface text-center">
          <Unavailable />
        </div>
      ) : (
        (tile === "loading" ||
          tile === "slow" ||
          (tile === "loaded" && !started)) && (
          // Opaque until Aladin is up (it shows its own spinner meanwhile),
          // then see-through so the ring and circles show while tiles load.
          <div
            className={`pointer-events-none absolute inset-0 flex items-center justify-center px-6 text-center ${started ? "bg-surface/50" : "bg-surface"}`}
          >
            <XWaveSpinner
              size={24}
              label={
                tile === "slow"
                  ? "The image server is slow. The Offsets view shows the same matches."
                  : "Loading sky image"
              }
            />
          </div>
        )
      )}
    </div>
  );
}

/** Shown when no mirror serves the survey here, or its tiles fail to load. */
function Unavailable() {
  return (
    <>
      <p className="m-0 text-sm">Sky image unavailable</p>
      <p className="m-0 px-6 text-xs text-neutral-500">
        No image server has this part of the sky right now. The Offsets view
        shows the same matches. <ServiceStatusLink />
      </p>
    </>
  );
}

/**
 * Angle between two positions, in degrees; small-angle approximation, which
 * is all a drift of a fraction of the field needs.
 */
function separationDeg(
  ra1: number,
  dec1: number,
  ra2: number,
  dec2: number
): number {
  const dRa = ((((ra2 - ra1) % 360) + 540) % 360) - 180;
  const cosDec = Math.cos((((dec1 + dec2) / 2) * Math.PI) / 180);
  return Math.hypot(dRa * cosDec, dec2 - dec1);
}
