"use client";

import { Flex } from "antd";
import type { ReactNode } from "react";

import { BasketButton } from "@/app/components/basket/BasketButton";
import { CodeSnippetButton } from "@/app/components/code/CodeSnippetButton";
import { FlagMatchButton } from "@/app/components/feedback/FlagMatchButton";
import { ShareButton } from "@/app/components/share/ShareButton";
import {
  coneSearchRequest,
  lightcurveRequest,
  metadataRequest,
} from "@/app/lib/utils/snippets";
import { buildObjectUrl } from "@/app/lib/utils/urls";

interface ObjectActionsProps {
  objectId: string;
  catalog: string;
  ra: number;
  dec: number;
  /** Extra buttons from the page, placed before Report. */
  extra?: ReactNode;
}

/** Toolbar on the object page: code, share, basket, extras, flag. */
export function ObjectActions({
  objectId,
  catalog,
  ra,
  dec,
  extra,
}: ObjectActionsProps) {
  const path = buildObjectUrl(objectId, catalog);
  return (
    <Flex gap={8} wrap align="center">
      <CodeSnippetButton
        title={`Get ${objectId} via the API`}
        requests={[
          { label: "Metadata", request: metadataRequest(objectId, catalog) },
          {
            label: "Counterparts (5″)",
            request: coneSearchRequest({ ra, dec, radius: 5, nneighbor: 10 }),
          },
          {
            label: "Light curve",
            request: lightcurveRequest({ ra, dec, radius: 2 }),
          },
        ]}
      />
      <ShareButton
        path={path}
        title={objectId}
        embedPath={`/embed/object/${encodeURIComponent(objectId)}?catalog=${encodeURIComponent(catalog)}`}
      />
      <BasketButton
        item={{ objectId, catalog, ra, dec }}
        size="small"
        labeled
      />
      {extra}
      <FlagMatchButton
        objectId={objectId}
        catalog={catalog}
        ra={ra}
        dec={dec}
      />
    </Flex>
  );
}
