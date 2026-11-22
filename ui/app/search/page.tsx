"use client";

import { Layout } from "antd";
import { Suspense, useCallback, useLayoutEffect, useMemo } from "react";

import { XWaveSpinner } from "@/app/components/common/XWaveSpinner";
import { AppHeader } from "@/app/components/layout";
import { SearchResults } from "@/app/components/results/SearchResults";
import { SidebarSearchForm } from "@/app/components/sidebar";
import { useParallelConeSearch } from "@/app/hooks/queries";
import { useBulkMetadata } from "@/app/hooks/queries/useBulkMetadata";
import { isRetryableConeError } from "@/app/hooks/queries/useConeSearch";
import { useSearchParams } from "@/app/hooks/useSearchParamsSync";
import { decodeCatalogRadii } from "@/app/lib/constants/search";
import { parseDecInput, parseRaInput } from "@/app/lib/utils/coordinates";
import { mapConeSearchResults } from "@/app/lib/utils/mapConeSearchResults";
import {
  buildBulkMetadataRequests,
  buildMagIndex,
  enrichWithPhotometry,
} from "@/app/lib/utils/mergePhotometry";
import {
  CrossmatchProvider,
  useCrossmatchState,
} from "@/app/store/crossmatch-context";

const { Content } = Layout;

function SearchContent() {
  const { dispatch } = useCrossmatchState();
  const { ra, dec, catalogRadii: catalogRadiiStr, isValid } = useSearchParams();

  // Same rules as the form: the form writes degrees, but an older or
  // hand-written link may carry h:m:s, which parseFloat would read as 5°.
  const raInput = parseRaInput(ra);
  const decInput = parseDecInput(dec);
  const raDeg = raInput.value;
  const decDeg = decInput.value;
  // A bad position in the URL (?ra=400) is for the user to fix. It is never
  // sent, so the service does not answer it with a 400 that reads as a failure.
  const urlError = isValid ? (raInput.error ?? decInput.error) : undefined;

  const base = useMemo(
    () =>
      raDeg === null || decDeg === null ? null : { ra: raDeg, dec: decDeg },
    [raDeg, decDeg]
  );

  const catalogConfigs = useMemo(
    () => decodeCatalogRadii(catalogRadiiStr),
    [catalogRadiiStr]
  );

  const queryResults = useParallelConeSearch(base, catalogConfigs);

  // queryResults[i] pairs with enabledCatalogs[i]: useParallelConeSearch builds
  // its queries from this same `enabled` filter, in this same order.
  const enabledCatalogs = useMemo(
    () => catalogConfigs.filter((c) => c.enabled).map((c) => c.catalog),
    [catalogConfigs]
  );

  const failed = queryResults.filter((r) => r.isError);
  const isLoading = queryResults.some((r) => r.isLoading);
  // One failing catalog no longer takes the others' results down with it: the
  // page errors only when every catalog failed, and lists the rest otherwise.
  const allFailed =
    queryResults.length > 0 && failed.length === queryResults.length;
  const isSettled =
    queryResults.length > 0 &&
    queryResults.every((r) => r.isSuccess || r.isError);
  // Retrying a failed catalog puts only that query back to loading. Keep the
  // other catalogs' results on screen instead of the full-page spinner.
  const isRetryingInPlace =
    queryResults.some((r) => r.isSuccess) &&
    queryResults.every((r) => !r.isLoading || r.errorUpdateCount > 0);

  const errorMessage = urlError ?? failed[0]?.error?.message;
  // Retrying cannot fix a request the service rejected (4xx: bad position,
  // radius over the limit), only a network error or a 5xx.
  const isRetryable =
    !urlError && failed.some((r) => isRetryableConeError(r.error));

  const mappedResults = useMemo(
    () =>
      queryResults.flatMap((r, i) =>
        r.data ? mapConeSearchResults(r.data, enabledCatalogs[i]) : []
      ),
    [queryResults, enabledCatalogs]
  );

  // Photometry enrichment. Deliberately kept out of the effect below: the table
  // renders from mappedResults immediately and the Mag column fills in later,
  // so a slow or failed metadata request must never blank the results.
  const bulkRequests = useMemo(
    () => buildBulkMetadataRequests(mappedResults),
    [mappedResults]
  );
  const { groups: metadataGroups, isFetching: photometryPending } =
    useBulkMetadata(bulkRequests);

  const enrichedResults = useMemo(
    () => enrichWithPhotometry(mappedResults, buildMagIndex(metadataGroups)),
    [mappedResults, metadataGroups]
  );

  const retryFailed = useCallback(() => {
    queryResults.forEach((r) => {
      // A rejected request would only fail again, and flicker while it does.
      if (r.isError && isRetryableConeError(r.error)) r.refetch();
    });
  }, [queryResults]);
  const handleRetry = isRetryable ? retryFailed : undefined;

  // NOTE: photometry state must never enter this effect or its dependencies —
  // SearchResults switches on resultsState to choose between the results and
  // the empty/error states.
  // This effect is the only writer of resultsState: the form just changes the
  // URL, and the queries decide whether that means loading. A layout effect,
  // so a new search goes straight to the spinner without painting one frame
  // of the previous state.
  useLayoutEffect(() => {
    if (urlError) {
      dispatch({ type: "SET_RESULTS_STATE", payload: "error" });
    } else if (!base || queryResults.length === 0) {
      // Also covers a valid target with every catalog unchecked, which
      // previously dispatched nothing and left the panel on its last state.
      dispatch({ type: "SET_RESULTS_STATE", payload: "empty" });
    } else if (isLoading && !isRetryingInPlace) {
      dispatch({ type: "SET_RESULTS_STATE", payload: "loading" });
    } else if (allFailed) {
      dispatch({ type: "SET_RESULTS_STATE", payload: "error" });
    } else if (isSettled || isRetryingInPlace) {
      dispatch({ type: "SET_RESULTS_STATE", payload: "success" });
    }
  }, [
    urlError,
    base,
    queryResults.length,
    isLoading,
    isRetryingInPlace,
    allFailed,
    isSettled,
    dispatch,
  ]);

  return (
    <Layout className="min-h-dvh">
      <AppHeader />
      <SearchResults
        data={enrichedResults}
        errorMessage={errorMessage}
        onRetry={handleRetry}
        target={base}
        searchForm={<SidebarSearchForm />}
        photometryPending={photometryPending}
        // Lets the list say which catalog failed rather than drop it.
        catalogErrors={Object.fromEntries(
          queryResults.flatMap((r, i) =>
            r.isError && enabledCatalogs[i]
              ? [[enabledCatalogs[i].toLowerCase(), r.error?.message ?? ""]]
              : []
          )
        )}
        // A catalog retrying (or re-fetching after an earlier failure) has no
        // data and no error yet; without this it would read as "0 within".
        catalogPending={queryResults.flatMap((r, i) =>
          r.isPending && enabledCatalogs[i]
            ? [enabledCatalogs[i].toLowerCase()]
            : []
        )}
      />
    </Layout>
  );
}

function LoadingFallback() {
  return (
    <Layout className="min-h-dvh">
      <AppHeader />
      <Content className="bg-background min-h-[calc(100dvh-64px)] flex items-center justify-center">
        <XWaveSpinner variant="xwave" size={64} />
      </Content>
    </Layout>
  );
}

export default function SearchPage() {
  return (
    <CrossmatchProvider>
      <Suspense fallback={<LoadingFallback />}>
        <SearchContent />
      </Suspense>
    </CrossmatchProvider>
  );
}
