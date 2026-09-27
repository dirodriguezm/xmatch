"use client";

import { Flex } from "antd";
import { useMemo } from "react";

import { CiteButton } from "@/app/components/cite/CiteButton";
import { CodeSnippetButton } from "@/app/components/code/CodeSnippetButton";
import { ShareButton } from "@/app/components/share/ShareButton";
import { useSearchParams } from "@/app/hooks/useSearchParamsSync";
import { getSearchCatalogLabel } from "@/app/lib/constants/catalogs";
import {
  convertRadiusToArcsec,
  decodeCatalogRadii,
} from "@/app/lib/constants/search";
import { coneSearchRequest } from "@/app/lib/utils/snippets";

interface ResultsActionsProps {
  target: { ra: number; dec: number };
  /** Catalog slugs that returned at least one match. */
  matchedCatalogs: string[];
}

/** Toolbar above the results table: code for this search, cite, share. */
export function ResultsActions({
  target,
  matchedCatalogs,
}: ResultsActionsProps) {
  const { ra, dec, catalogRadii } = useSearchParams();

  const requests = useMemo(
    () =>
      decodeCatalogRadii(catalogRadii)
        .filter((c) => c.enabled)
        .map((c) => ({
          label: getSearchCatalogLabel(c.catalog),
          request: coneSearchRequest({
            ra: target.ra,
            dec: target.dec,
            radius: convertRadiusToArcsec(c.radius, c.unit),
            catalog: c.catalog,
          }),
        })),
    [catalogRadii, target.ra, target.dec]
  );

  const params = new URLSearchParams({ ra, dec });
  if (catalogRadii) params.set("catalogRadii", catalogRadii);

  return (
    <Flex gap={8} wrap align="center">
      <CodeSnippetButton
        title="Run this search via the API"
        requests={requests}
      />
      <CiteButton catalogs={matchedCatalogs} />
      <ShareButton
        path={`/search?${params}`}
        title={`Cone search at ${target.ra.toFixed(5)}, ${target.dec.toFixed(5)}`}
      />
    </Flex>
  );
}
