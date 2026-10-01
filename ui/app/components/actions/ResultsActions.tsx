"use client";

import { Flex } from "antd";
import { type ReactNode, useMemo } from "react";

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
  /** Further buttons, e.g. the CSV export of the matches. */
  extra?: ReactNode;
}

/** Toolbar above the results table: code for this search, share. */
export function ResultsActions({ target, extra }: ResultsActionsProps) {
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
      <ShareButton
        path={`/search?${params}`}
        title={`Cone search at ${target.ra.toFixed(5)}, ${target.dec.toFixed(5)}`}
      />
      {extra}
    </Flex>
  );
}
