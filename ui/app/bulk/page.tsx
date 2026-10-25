"use client";

import { Alert } from "antd";
import {
  parseAsFloat,
  parseAsInteger,
  parseAsString,
  parseAsStringLiteral,
  useQueryStates,
} from "nuqs";
import { Suspense, useEffect, useMemo, useRef, useState } from "react";

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
import { BULK_INPUT_STORAGE_KEY } from "@/app/lib/utils/basket";
import { MAX_BULK_ROWS, parseBulkInput } from "@/app/lib/utils/bulkInput";
import {
  type BulkRun,
  DEFAULT_BULK_VIEW,
  loadBulkRun,
  saveBulkRun,
} from "@/app/lib/utils/bulkSession";
import { bulkConeSearchRequest } from "@/app/lib/utils/snippets";

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
  // The previous run, shown until a new one starts (e.g. back from /object).
  const [restored, setRestored] = useState<BulkRun | null>(null);
  const initialized = useRef(false);

  // Prefill once from the basket or a file picked on the landing page, then
  // drop both the stored list and the flag. Otherwise bring back the last run.
  // Deferred to a task so the state update happens outside the effect body.
  useEffect(() => {
    // Marked inside the task: strict mode's extra mount cancels the first one.
    if (initialized.current) return;
    if (params.from !== "basket" && params.from !== "file") {
      const timer = window.setTimeout(() => {
        initialized.current = true;
        const run = loadBulkRun();
        if (!run) return;
        setRestored(run);
        setText(run.text);
        void setParams({
          catalog: BULK_CATALOGS.find((c) => c === run.catalog) ?? "all",
          radius: run.radius,
          n: run.nneighbor,
        });
      }, 0);
      return () => window.clearTimeout(timer);
    }
    const timer = window.setTimeout(() => {
      initialized.current = true;
      try {
        const stored = window.sessionStorage.getItem(BULK_INPUT_STORAGE_KEY);
        if (stored) setText(stored);
        window.sessionStorage.removeItem(BULK_INPUT_STORAGE_KEY);
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
    const runText = text;
    bulk.mutate(
      { targets: parsed.targets, radius, catalog, nneighbor },
      {
        onSuccess: (rows, vars) =>
          saveBulkRun({
            ...vars,
            text: runText,
            rows,
            view: DEFAULT_BULK_VIEW,
          }),
      }
    );
  };

  // Everything below describes the run that produced the results, not the
  // form's current state — editing the form must not relabel old results.
  const showRestored = bulk.isIdle && restored !== null;
  const lastRun = showRestored ? restored : bulk.variables;
  const lastRows = showRestored ? restored.rows : bulk.data;
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

        {(bulk.isSuccess || showRestored) &&
          lastRun &&
          lastRows &&
          lastRequest && (
            <section className="flex flex-col gap-2">
              <BulkResults
                key={showRestored ? "restored" : bulk.submittedAt}
                rows={lastRows}
                radius={lastRun.radius}
                request={lastRequest}
                initialView={showRestored ? restored.view : undefined}
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
