"use client";

import type { ReactNode } from "react";

import type { CatalogOption } from "@/app/lib/constants/catalogs";

/** Deterministic Park–Miller generator, so example fields look the same on every render. */
export function seededRandom(seed: number) {
  let s = seed;
  return () => {
    s = (s * 16807) % 2147483647;
    return (s - 1) / 2147483646;
  };
}

/** Example catalog mix: mostly Gaia, some AllWISE, a few eROSITA. */
export function pickCatalog(u: number): CatalogOption {
  return u < 0.55 ? "gaia" : u < 0.82 ? "allwise" : "erosita";
}

export const HIT = "#7ea2ff";
export const MUTED = "#868a96";

export function formatNumber(n: number, digits = 0) {
  return n.toLocaleString("en-US", {
    minimumFractionDigits: digits,
    maximumFractionDigits: digits,
  });
}

/** Figure shell shared by the Learn diagrams: drawing, side panel, controls. */
export function DiagramFrame({
  stage,
  side,
  controls,
  wide = false,
}: {
  stage: ReactNode;
  side?: ReactNode;
  controls?: ReactNode;
  wide?: boolean;
}) {
  return (
    <figure className="m-0 my-6 overflow-hidden rounded-lg border border-border bg-surface">
      <div
        className={
          wide || !side
            ? "grid"
            : "grid md:grid-cols-[minmax(0,5fr)_minmax(0,4fr)]"
        }
      >
        <div
          className={`flex min-w-0 flex-col items-center justify-center gap-3 bg-background p-4 ${
            wide || !side
              ? ""
              : "border-b border-border md:border-b-0 md:border-r"
          }`}
        >
          {stage}
        </div>
        {side && (
          <div
            className={`flex min-w-0 flex-col gap-4 p-4 ${wide ? "border-t border-border" : ""}`}
          >
            {side}
          </div>
        )}
      </div>
      {controls && (
        <div className="flex flex-wrap items-center gap-x-6 gap-y-3 border-t border-border px-4 py-3 text-sm text-neutral-400">
          {controls}
        </div>
      )}
    </figure>
  );
}

/** Big number with a caption, as used in the diagrams' side panels. */
export function Stat({ value, label }: { value: ReactNode; label: string }) {
  return (
    <div className="flex flex-col">
      <span className="font-mono text-lg font-semibold tabular-nums text-foreground">
        {value}
      </span>
      <span className="text-xs text-neutral-400">{label}</span>
    </div>
  );
}
