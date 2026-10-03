"use client";

import { Button, Empty, Flex, Skeleton, Typography } from "antd";
import { useState } from "react";

import { useSimbadRefs } from "@/app/hooks/queries";
import { adsAbstractUrl, SIMBAD_REFS_MAX } from "@/app/lib/utils/simbadRefs";
import { buildAdsObjectUrl } from "@/app/lib/utils/urls";

const { Text, Link } = Typography;

const PAGE = 25;

function simbadBibliographyUrl(mainId: string): string {
  return `https://simbad.cds.unistra.fr/simbad/sim-id?Ident=${encodeURIComponent(mainId)}#lab_bib`;
}

interface LiteraturePanelProps {
  oid: number;
  mainId: string;
}

/** Papers about the SIMBAD object, newest first, linking to ADS. */
export function LiteraturePanel({ oid, mainId }: LiteraturePanelProps) {
  const [limit, setLimit] = useState(PAGE);
  const { data, isPending, isError, isFetching } = useSimbadRefs(oid, limit);

  if (isPending) return <Skeleton active paragraph={{ rows: 4 }} />;
  if (isError) {
    return (
      <Text type="secondary" className="text-sm">
        Could not load references from SIMBAD.
      </Text>
    );
  }
  if (data.total === 0) {
    return (
      <Empty
        image={Empty.PRESENTED_IMAGE_SIMPLE}
        description={`SIMBAD lists no papers about ${mainId}.`}
      />
    );
  }

  const canLoadMore =
    data.references.length < data.total && limit < SIMBAD_REFS_MAX;

  return (
    <Flex vertical gap={8}>
      <ul className="m-0 p-0 list-none divide-y divide-border">
        {data.references.map((ref) => (
          <li key={ref.bibcode} className="py-2">
            <Flex vertical gap={2} className="min-w-0">
              <Link
                href={adsAbstractUrl(ref.bibcode)}
                target="_blank"
                rel="noopener noreferrer"
                className="!text-foreground hover:!text-primary"
              >
                {ref.title ?? ref.bibcode}
              </Link>
              <Text type="secondary" className="text-xs">
                {[ref.year, ref.journal].filter(Boolean).join(" · ")}
                {" · "}
                <span className="font-mono">{ref.bibcode}</span>
              </Text>
            </Flex>
          </li>
        ))}
      </ul>
      <Flex justify="space-between" align="center" wrap gap={8}>
        <Text type="secondary" className="text-xs">
          Showing the {data.references.length} most recent of {data.total}{" "}
          papers SIMBAD links to {mainId}. Full lists:{" "}
          <Link
            href={buildAdsObjectUrl(mainId)}
            target="_blank"
            rel="noopener noreferrer"
          >
            NASA ADS
          </Link>{" "}
          ·{" "}
          <Link
            href={simbadBibliographyUrl(mainId)}
            target="_blank"
            rel="noopener noreferrer"
          >
            SIMBAD
          </Link>
        </Text>
        {canLoadMore && (
          <Button
            size="small"
            loading={isFetching}
            onClick={() => setLimit((l) => Math.min(l * 2, SIMBAD_REFS_MAX))}
          >
            Show more
          </Button>
        )}
      </Flex>
    </Flex>
  );
}
