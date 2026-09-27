"use client";

import { Flex } from "antd";

import { BasketButton } from "@/app/components/basket/BasketButton";
import { CiteButton } from "@/app/components/cite/CiteButton";
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
}

/** Toolbar on the object page: code, cite, share, basket, flag. */
export function ObjectActions({
  objectId,
  catalog,
  ra,
  dec,
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
      <CiteButton catalogs={[catalog]} objectId={objectId} />
      <ShareButton
        path={path}
        title={objectId}
        embedPath={`/embed/object/${encodeURIComponent(objectId)}?catalog=${encodeURIComponent(catalog)}`}
      />
      <BasketButton item={{ objectId, catalog, ra, dec }} />
      <FlagMatchButton
        objectId={objectId}
        catalog={catalog}
        ra={ra}
        dec={dec}
      />
    </Flex>
  );
}
