"use client";

import { Alert } from "antd";
import {
  parseAsFloat,
  parseAsInteger,
  parseAsString,
  parseAsStringLiteral,
  useQueryStates,
} from "nuqs";
import { Suspense, useEffect, useMemo, useState } from "react";

import {
  BULK_CATALOGS,
  BULK_MAX_NNEIGHBOR,
  BulkInputPanel,
  BulkOptions,
  BulkResults,
} from "@/app/components/bulk";
import { XWaveSpinner } from "@/app/components/common/XWaveSpinner";
import { PageShell } from "@/app/components/layout";
import { useBulkConeSearch } from "@/app/hooks/queries/useBulkConeSearch";
import { MAX_RADIUS_ARCSEC } from "@/app/lib/constants/search";
import { MAX_BULK_ROWS, parseBulkInput } from "@/app/lib/utils/bulkInput";
import { bulkConeSearchRequest } from "@/app/lib/utils/snippets";

/** Written by the basket before navigating to /bulk?from=basket. */
const BASKET_STORAGE_KEY = "xwave:bulk-input";

const bulkParams = {
  catalog: parseAsStringLiteral(BULK_CATALOGS).withDefault("all"),
  radius: parseAsFloat.withDefault(2),
  n: parseAsInteger.withDefault(1),
  from: parseAsString,
};

function clamp(value: number, min: number, max: number): number {
  return Math.min(max, Math.max(min, value));
}

function BulkContent() {
  const [params, setParams] = useQueryStates(bulkParams, {
    history: "replace",
  });
  const catalog = params.catalog;
  const radius = clamp(params.radius, 0.1, MAX_RADIUS_ARCSEC);
  const nneighbor = clamp(Math.round(params.n), 1, BULK_MAX_NNEIGHBOR);

  const [text, setText] = useState("");
  const parsed = useMemo(() => parseBulkInput(text), [text]);
  const bulk = useBulkConeSearch();

  // Prefill from the basket once, then drop both the stored list and the flag.
  // Deferred to a task so the state update happens outside the effect body.
  useEffect(() => {
    if (params.from !== "basket") return;
    const timer = window.setTimeout(() => {
      try {
        const stored = window.sessionStorage.getItem(BASKET_STORAGE_KEY);
        if (stored) setText(stored);
        window.sessionStorage.removeItem(BASKET_STORAGE_KEY);
      } catch {
        // Storage blocked (private mode, sandbox): nothing to prefill.
      }
      void setParams({ from: null });
    }, 0);
    return () => window.clearTimeout(timer);
  }, [params.from, setParams]);

  const disabledReason =
    parsed.targets.length === 0
      ? "Paste or upload at least one position"
      : parsed.targets.length > MAX_BULK_ROWS
        ? `At most ${MAX_BULK_ROWS} positions per run`
        : null;

  const run = () => {
    if (disabledReason) return;
    bulk.mutate({ targets: parsed.targets, radius, catalog, nneighbor });
  };

  // Everything below describes the run that produced the results, not the
  // form's current state — editing the form must not relabel old results.
  const lastRun = bulk.variables;
  const lastRequest = useMemo(
    () =>
      lastRun
        ? bulkConeSearchRequest({
            ra: lastRun.targets.map((t) => t.ra),
            dec: lastRun.targets.map((t) => t.dec),
            radius: lastRun.radius,
            catalog: lastRun.catalog,
            nneighbor: lastRun.nneighbor,
          })
        : null,
    [lastRun]
  );

  return (
    <PageShell
      width="wide"
      title="Bulk cross-match"
      description={`Match up to ${MAX_BULK_ROWS} positions against Gaia DR3, AllWISE and eROSITA in one go. Paste a list or upload a CSV/TSV.`}
    >
      <div className="flex flex-col gap-6 pb-12">
        <section className="rounded-lg border border-border bg-surface p-4 md:p-5 flex flex-col gap-5">
          <BulkInputPanel value={text} onChange={setText} parsed={parsed} />
          <div className="border-t border-border pt-4">
            <BulkOptions
              catalog={catalog}
              radius={radius}
              nneighbor={nneighbor}
              onCatalogChange={(v) => void setParams({ catalog: v })}
              onRadiusChange={(v) => void setParams({ radius: v })}
              onNneighborChange={(v) => void setParams({ n: v })}
              onRun={run}
              running={bulk.isPending}
              disabledReason={disabledReason}
            />
          </div>
        </section>

        {bulk.isError && (
          <Alert
            type="error"
            showIcon
            title="Cross-match failed"
            description={bulk.error.message}
          />
        )}

        {bulk.isPending && (
          <div className="flex justify-center py-12">
            <XWaveSpinner
              variant="healpix"
              size={64}
              label={`Matching ${lastRun?.targets.length ?? 0} positions…`}
            />
          </div>
        )}

        {bulk.isSuccess && lastRun && lastRequest && (
          <section className="flex flex-col gap-2">
            <BulkResults
              rows={bulk.data}
              radius={lastRun.radius}
              request={lastRequest}
            />
          </section>
        )}
      </div>
    </PageShell>
  );
}

export default function BulkPage() {
  return (
    <Suspense
      fallback={
        <div className="flex min-h-screen items-center justify-center">
          <XWaveSpinner variant="xwave" size={64} />
        </div>
      }
    >
      <BulkContent />
    </Suspense>
  );
}
