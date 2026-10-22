"use client";

import { InfoCircleOutlined } from "@ant-design/icons";
import { Layout, Segmented, Tooltip } from "antd";
import { type ReactNode, useMemo, useState } from "react";

import { AppSidebar } from "@/app/components/layout";
import { MapPanel } from "@/app/components/results/nearby/MapPanel";
import { TonightSummary } from "@/app/components/results/TonightSummary";
import { getSearchCatalogLabel } from "@/app/lib/constants/catalogs";

import { CatalogMark } from "./nearby/CatalogMark";
import { SIGMA_CAVEAT, summarizeCatalogs } from "./nearby/catalogSummary";
import { GroupedMatches, TOP_N } from "./nearby/GroupedMatches";
import { MatchesBody } from "./nearby/MatchesBody";
import { ResultsToolbar } from "./nearby/ResultsToolbar";
import {
  chanceWithin,
  fieldDensity,
  formatChance,
} from "./nearby/separationStats";
import { SkyImage } from "./nearby/SkyImage";
import { useNearbyResults } from "./nearby/useNearbyResults";
import type { EnrichedResult } from "./ResultsTable";

export interface SearchResultsProps {
  data: EnrichedResult[];
  errorMessage?: string;
  onRetry?: () => void;
  /** Searched position, or null before a search. */
  target: { ra: number; dec: number } | null;
  /** The search form, kept in the sidebar for further searches. */
  searchForm: ReactNode;
}

const ALL = "__all__";

const { Content } = Layout;

/**
 * Search results: the search form in the sidebar, then a balanced split. One row whose height comes from the window
 * (not the content): the map is a square of exactly that height and the
 * grouped list scrolls inside a panel of the same height, so the two panes
 * always line up — with 3 matches or 300. Below lg it stacks: a capped square
 * map, then the list flowing with the page.
 */
export function SearchResults({
  data,
  errorMessage,
  onRetry,
  target,
  searchForm,
}: SearchResultsProps) {
  const { radii, sources, highlight } = useNearbyResults(data, target);
  const [active, setActive] = useState<string>(ALL);
  const [opened, setOpened] = useState<Set<string>>(() => new Set());
  const [mapView, setMapView] = useState<"offsets" | "sky">("offsets");

  // Catalogs that returned nothing are left out entirely — no empty group,
  // tab or search circle.
  const groups = useMemo(
    () => summarizeCatalogs(sources, radii).filter((g) => g.sources.length > 0),
    [sources, radii]
  );
  const matchedRadii = useMemo(
    () =>
      Object.fromEntries(
        groups
          .filter((g) => g.radius !== undefined)
          .map((g) => [g.slug, g.radius as number])
      ),
    [groups]
  );
  const activeGroup = groups.find((g) => g.slug === active);
  const visibleGroups = activeGroup ? [activeGroup] : groups;
  const mapSources = activeGroup ? activeGroup.sources : sources;
  const mapRadii = useMemo(
    () =>
      active !== ALL && radii[active] !== undefined
        ? { [active]: radii[active] }
        : matchedRadii,
    [active, radii, matchedRadii]
  );

  // Chance of an unrelated source at least as close as each catalog's nearest
  // match; only where the field holds enough sources to estimate a density.
  const chanceBySlug = useMemo(() => {
    const out: Record<string, string> = {};
    for (const g of groups) {
      const d = fieldDensity(g.sources, g.radius);
      if (d && g.best)
        out[g.slug] = formatChance(chanceWithin(d, g.best.sepArcsec));
    }
    return out;
  }, [groups]);

  const expanded = useMemo(() => {
    const set = new Set(opened);
    if (activeGroup) set.add(activeGroup.slug);
    const sel = sources.find((s) => s.key === highlight.selectedKey);
    if (sel) {
      const g = groups.find((x) => x.slug === sel.slug);
      if (g && g.sources.indexOf(sel) >= TOP_N) set.add(sel.slug);
    }
    return set;
  }, [opened, activeGroup, sources, groups, highlight.selectedKey]);

  const toggleExpanded = (slug: string) =>
    setOpened((prev) => {
      const next = new Set(prev);
      if (next.has(slug)) next.delete(slug);
      else next.add(slug);
      return next;
    });

  return (
    // The search form stays in the sidebar so new searches are one edit
    // away; the results area alone follows the balanced recipe.
    <Layout>
      <AppSidebar>{searchForm}</AppSidebar>
      <Content className="bg-background min-h-[calc(100vh-64px)]">
        <MatchesBody
          target={target}
          matchCount={sources.length}
          errorMessage={errorMessage}
          onRetry={onRetry}
          // The form is already in the sidebar.
          searchForm={null}
        >
          {/* The window sets this row's height (64px header + toolbar +
            padding), clamped so huge screens don't get a giant map. Both
            panes take that height; the square map is centred in its pane. */}
          <div className="grid grid-cols-1 gap-4 p-4 lg:h-[clamp(380px,calc(100dvh-96px),860px)] lg:grid-cols-[minmax(0,45fr)_minmax(0,55fr)] lg:gap-6 lg:px-6">
            <section
              aria-label="Sky map of the matches"
              className="mx-auto flex aspect-square w-full max-w-[min(100%,45svh)] flex-col rounded-lg border border-border bg-surface p-3 lg:aspect-auto lg:h-full lg:max-w-none"
            >
              <Segmented
                size="small"
                block
                className="mb-2"
                value={mapView}
                onChange={(v) => setMapView(v as "offsets" | "sky")}
                aria-label="Map view"
                options={[
                  { value: "offsets", label: "Offsets" },
                  { value: "sky", label: "Sky image" },
                ]}
              />
              {mapView === "sky" && target ? (
                <SkyImage
                  key={active}
                  target={target}
                  sources={mapSources}
                  radii={mapRadii}
                />
              ) : (
                <MapPanel
                  key={active}
                  sources={mapSources}
                  radii={mapRadii}
                  highlight={highlight}
                  size={600}
                  className="h-full min-h-0"
                  mapClassName="min-h-0 w-full flex-1 h-full"
                />
              )}
            </section>

            <section
              aria-label="Matches by catalog"
              className="flex min-h-0 min-w-0 flex-col gap-3 rounded-lg border border-border bg-surface p-3"
            >
              <div className="flex flex-wrap items-center gap-x-4 gap-y-2">
                <h2 className="m-0 text-base font-semibold text-foreground">
                  {sources.length} {sources.length === 1 ? "match" : "matches"}
                </h2>
                <div className="max-w-full overflow-x-auto">
                  <Segmented
                    size="small"
                    value={active}
                    onChange={(v) => setActive(v as string)}
                    aria-label="Catalog shown on the map and in the list"
                    options={[
                      { value: ALL, label: "All" },
                      ...groups.map((g) => ({
                        value: g.slug,
                        label: (
                          <span className="inline-flex items-center gap-1.5">
                            <CatalogMark slug={g.slug} size={9} />
                            {getSearchCatalogLabel(g.slug)}
                            <span className="text-neutral-400">
                              {g.sources.length}
                              {g.capped ? "+" : ""}
                            </span>
                          </span>
                        ),
                      })),
                    ]}
                  />
                </div>
                {target && (
                  <div className="ml-auto">
                    <ResultsToolbar target={target} sources={sources} />
                  </div>
                )}
              </div>

              {/* Only this list scrolls; the map beside it never moves. */}
              <GroupedMatches
                groups={visibleGroups}
                highlight={highlight}
                expanded={expanded}
                onToggleExpanded={activeGroup ? undefined : toggleExpanded}
                chanceBySlug={chanceBySlug}
                className="lg:min-h-0 lg:flex-1 lg:overflow-y-auto border-none"
              />

              {/* Pinned footer: context that used to float elsewhere. */}
              <div className="mt-auto flex flex-col gap-1 border-t border-border pt-2">
                {target && <TonightSummary ra={target.ra} dec={target.dec} />}
                <Tooltip title={SIGMA_CAVEAT}>
                  <span className="inline-flex w-fit cursor-help items-center gap-1 text-xs text-neutral-500">
                    <InfoCircleOutlined />σ from typical survey errors; chance ≈
                    how likely an unrelated source sits this close (approximate)
                  </span>
                </Tooltip>
              </div>
            </section>
          </div>
        </MatchesBody>
      </Content>
    </Layout>
  );
}
