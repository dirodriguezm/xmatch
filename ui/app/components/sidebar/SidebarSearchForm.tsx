"use client";

import { AimOutlined, FilterOutlined, SearchOutlined } from "@ant-design/icons";
import { Button, Collapse, Flex, Input, Typography } from "antd";
import { type FormEvent, useId, useState } from "react";

import { CatalogRadiusRow } from "@/app/components/common";
import { useSearchParams } from "@/app/hooks/useSearchParamsSync";
import {
  buildDefaultCatalogConfigs,
  CATALOG_SELECT_OPTIONS,
} from "@/app/lib/constants/catalogs";
import {
  type CatalogRadiusConfig,
  convertRadiusToArcsec,
  decodeCatalogRadii,
  encodeCatalogRadii,
  MAX_RADIUS_ARCSEC,
} from "@/app/lib/constants/search";
import {
  type AngleInput,
  formatDegreesParam,
  parseDecInput,
  parseRaInput,
} from "@/app/lib/utils/coordinates";
import { useCrossmatchState } from "@/app/store/crossmatch-context";

const { Title, Text } = Typography;

/**
 * Always returns one config per known catalog.
 *
 * A URL may name only a subset (`?catalogRadii=allwise:5:deg:1`), but the form
 * renders a row for every catalog and looks its config up by name — so a
 * partial list used to yield `undefined` and crash the row. Catalogs the URL
 * omits fall back to their defaults, unchecked, since they were not searched.
 */
function initConfigs(catalogRadiiStr: string): CatalogRadiusConfig[] {
  const defaults = buildDefaultCatalogConfigs();
  if (!catalogRadiiStr) return defaults;

  const decoded = decodeCatalogRadii(catalogRadiiStr);
  if (decoded.length === 0) return defaults;

  const byCatalog = new Map(decoded.map((c) => [c.catalog, c]));
  const merged = defaults.map(
    (d) => byCatalog.get(d.catalog) ?? { ...d, enabled: false }
  );
  // Keep any catalog the URL names that we do not know about.
  const known = new Set(defaults.map((d) => d.catalog));
  return [...merged, ...decoded.filter((c) => !known.has(c.catalog))];
}

/**
 * Whether the typed text already means the URL's value. After a search the
 * URL holds degrees while the input may still read "05:35:17.3"; that must
 * not count as a change, or the user's text would be swapped for degrees.
 */
function sameAngle(
  draft: string,
  urlValue: string,
  parse: (input: string) => AngleInput
): boolean {
  if (draft === urlValue) return true;
  const a = parse(draft).value;
  const b = parse(urlValue).value;
  return a !== null && b !== null && Math.abs(a - b) < 1e-6;
}

/** Which URL values are present but unreadable, to flag them on arrival. */
function invalidInUrl(ra: string, dec: string) {
  return {
    ra: Boolean(ra) && parseRaInput(ra).value === null,
    dec: Boolean(dec) && parseDecInput(dec).value === null,
  };
}

export function SidebarSearchForm() {
  const { state } = useCrossmatchState();
  const {
    ra,
    dec,
    catalogRadii: catalogRadiiStr,
    setSearchParams,
  } = useSearchParams();

  const [draftRa, setDraftRa] = useState(ra);
  const [draftDec, setDraftDec] = useState(dec);
  const [draftConfigs, setDraftConfigs] = useState<CatalogRadiusConfig[]>(() =>
    initConfigs(catalogRadiiStr)
  );
  // Errors show once a field is left or the form submitted, not mid-typing.
  // A bad value arriving in the URL (?ra=400) shows its error straight away.
  const [touched, setTouched] = useState(() => invalidInUrl(ra, dec));

  // Follow the URL when it changes under the form: "try M31", the landing
  // page, the command palette, back/forward. Adjusted during render rather
  // than in an effect so the inputs never paint one frame stale. Fields that
  // already mean the new URL value keep the user's text.
  const urlKey = `${ra}|${dec}|${catalogRadiiStr}`;
  const [syncedUrlKey, setSyncedUrlKey] = useState(urlKey);
  if (urlKey !== syncedUrlKey) {
    setSyncedUrlKey(urlKey);
    if (!sameAngle(draftRa, ra, parseRaInput)) setDraftRa(ra);
    if (!sameAngle(draftDec, dec, parseDecInput)) setDraftDec(dec);
    if (catalogRadiiStr !== encodeCatalogRadii(draftConfigs)) {
      setDraftConfigs(initConfigs(catalogRadiiStr));
    }
    setTouched(invalidInUrl(ra, dec));
  }

  const raInput = parseRaInput(draftRa);
  const decInput = parseDecInput(draftDec);
  const raError = touched.ra ? raInput.error : undefined;
  const decError = touched.dec ? decInput.error : undefined;
  const isDraftValid = Boolean(draftRa.trim() && draftDec.trim());
  // The row already flags it; searching anyway would only be refused.
  const radiusOverLimit = draftConfigs.some(
    (c) =>
      c.enabled && convertRadiusToArcsec(c.radius, c.unit) > MAX_RADIUS_ARCSEC
  );

  const raId = useId();
  const decId = useId();
  const radiusHintId = useId();

  const updateConfig = (
    catalog: string,
    patch: Partial<CatalogRadiusConfig>
  ) => {
    setDraftConfigs((prev) =>
      prev.map((c) => (c.catalog === catalog ? { ...c, ...patch } : c))
    );
  };

  const handleSubmit = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (raInput.value === null || decInput.value === null) {
      setTouched({ ra: true, dec: true });
      return;
    }
    if (radiusOverLimit) return;
    // The URL always carries degrees, whatever format was typed, so shared
    // links and the page's own parsing see one canonical form. No "loading"
    // here: the page derives it from the queries, so a search whose queries
    // are already cached (or unchanged) can never get stuck on a spinner.
    setSearchParams({
      ra: formatDegreesParam(raInput.value),
      dec: formatDegreesParam(decInput.value),
      catalogRadii: encodeCatalogRadii(draftConfigs),
    });
  };

  return (
    <div className="flex flex-col h-full">
      <div className="flex-1 overflow-auto">
        {/* Search Parameters Section */}
        <div className="p-5">
          <Flex vertical gap={20}>
            <Flex align="center" gap={10}>
              <AimOutlined className="text-primary text-base" />
              <Title level={5} className="m-0!">
                Search parameters
              </Title>
            </Flex>

            {/* A real form so Enter in either field searches. */}
            <form onSubmit={handleSubmit} noValidate>
              <Flex vertical gap={20}>
                <Flex vertical gap={16}>
                  <Flex vertical gap={6}>
                    <label htmlFor={raId}>
                      <Text type="secondary" className="text-xs font-medium">
                        RA (° or h:m:s)
                      </Text>
                    </label>
                    <Input
                      id={raId}
                      placeholder="e.g., 10.6847 or 00:42:44.3"
                      value={draftRa}
                      onChange={(e) => setDraftRa(e.target.value)}
                      onBlur={() =>
                        setTouched((t) => ({ ...t, ra: Boolean(draftRa) }))
                      }
                      status={raError ? "error" : undefined}
                      aria-invalid={raError ? true : undefined}
                      aria-describedby={raError ? `${raId}-error` : undefined}
                      size="large"
                    />
                    {raError && (
                      <Text
                        id={`${raId}-error`}
                        type="danger"
                        className="text-xs"
                      >
                        {raError}
                      </Text>
                    )}
                  </Flex>

                  <Flex vertical gap={6}>
                    <label htmlFor={decId}>
                      <Text type="secondary" className="text-xs font-medium">
                        Dec (° or d:m:s)
                      </Text>
                    </label>
                    <Input
                      id={decId}
                      placeholder="e.g., 41.269 or +41:16:08"
                      value={draftDec}
                      onChange={(e) => setDraftDec(e.target.value)}
                      onBlur={() =>
                        setTouched((t) => ({ ...t, dec: Boolean(draftDec) }))
                      }
                      status={decError ? "error" : undefined}
                      aria-invalid={decError ? true : undefined}
                      aria-describedby={decError ? `${decId}-error` : undefined}
                      size="large"
                    />
                    {decError && (
                      <Text
                        id={`${decId}-error`}
                        type="danger"
                        className="text-xs"
                      >
                        {decError}
                      </Text>
                    )}
                  </Flex>
                </Flex>

                {/* The hint sits with the button it explains. */}
                <Flex vertical gap={6}>
                  <Button
                    type="primary"
                    htmlType="submit"
                    block
                    icon={<SearchOutlined />}
                    disabled={!isDraftValid || radiusOverLimit}
                    loading={state.resultsState === "loading"}
                    size="large"
                    aria-describedby={
                      radiusOverLimit ? radiusHintId : undefined
                    }
                  >
                    Search
                  </Button>
                  {radiusOverLimit && (
                    <Text id={radiusHintId} type="danger" className="text-xs">
                      A radius is over the {MAX_RADIUS_ARCSEC}″ limit. Shorten
                      it under Catalogs &amp; Radii.
                    </Text>
                  )}
                </Flex>
              </Flex>
            </form>
          </Flex>
        </div>

        {/* Catalogs & Radii Section - Collapsible */}
        <Collapse
          defaultActiveKey={["catalogs"]}
          ghost
          expandIconPlacement="end"
          className="border-t border-border"
          items={[
            {
              key: "catalogs",
              label: (
                <Flex align="center" gap={10}>
                  <FilterOutlined className="text-primary text-base" />
                  <Text strong>Catalogs & Radii</Text>
                  <Text type="secondary" className="text-xs">
                    ({draftConfigs.filter((c) => c.enabled).length}/
                    {draftConfigs.length})
                  </Text>
                </Flex>
              ),
              children: (
                <Flex vertical gap={10} className="pb-2">
                  {CATALOG_SELECT_OPTIONS.map((catalog) => {
                    const config = draftConfigs.find(
                      (c) => c.catalog === catalog.value
                    );
                    // initConfigs guarantees one per catalog; skip rather than
                    // crash the sidebar if that invariant ever breaks.
                    if (!config) return null;
                    return (
                      <CatalogRadiusRow
                        key={catalog.value}
                        label={catalog.label}
                        catalog={catalog.value}
                        config={config}
                        onChange={(patch) => updateConfig(catalog.value, patch)}
                      />
                    );
                  })}
                </Flex>
              ),
            },
          ]}
        />
      </div>
    </div>
  );
}
