"use client";

import {
  Button,
  Empty,
  Flex,
  Segmented,
  Select,
  Skeleton,
  Tooltip,
  Typography,
} from "antd";
import { useState } from "react";

import { useSimbadRefs } from "@/app/hooks/queries";
import {
  adsAbstractUrl,
  SIMBAD_REFS_MAX,
  type SimbadRefsSort,
} from "@/app/lib/utils/simbadRefs";
import { buildAdsObjectUrl } from "@/app/lib/utils/urls";

const { Text, Link } = Typography;

/** References per page; ObjectDetail asks for the same first page to show the count. */
export const LITERATURE_PAGE = 25;

/** Year filters: `from` inclusive, `to` exclusive. */
const PERIODS = [
  { value: "all", label: "All years" },
  { value: "2020", label: "2020s", from: 2020 },
  { value: "2010", label: "2010s", from: 2010, to: 2020 },
  { value: "2000", label: "2000s", from: 2000, to: 2010 },
  { value: "1990", label: "1990s", from: 1990, to: 2000 },
  { value: "older", label: "Before 1990", to: 1990 },
] as const;
type Period = (typeof PERIODS)[number]["value"];

function simbadBibliographyUrl(mainId: string): string {
  return `https://simbad.cds.unistra.fr/simbad/sim-id?Ident=${encodeURIComponent(mainId)}#lab_bib`;
}

interface LiteraturePanelProps {
  oid: number;
  mainId: string;
}

/** Papers about the SIMBAD object, linking to ADS. */
export function LiteraturePanel({ oid, mainId }: LiteraturePanelProps) {
  const [limit, setLimit] = useState(LITERATURE_PAGE);
  const [sort, setSort] = useState<SimbadRefsSort>("newest");
  const [period, setPeriod] = useState<Period>("all");
  const range = PERIODS.find((p) => p.value === period);
  const { data, isPending, isError, isFetching } = useSimbadRefs(oid, limit, {
    sort,
    fromYear: range && "from" in range ? range.from : undefined,
    toYear: range && "to" in range ? range.to : undefined,
  });

  // A new order or period starts again from the first page.
  const changeSort = (s: SimbadRefsSort) => {
    setSort(s);
    setLimit(LITERATURE_PAGE);
  };
  const changePeriod = (p: Period) => {
    setPeriod(p);
    setLimit(LITERATURE_PAGE);
  };

  const controls = (
    <Flex justify="space-between" align="center" wrap gap={8}>
      <Segmented<SimbadRefsSort>
        size="small"
        value={sort}
        onChange={changeSort}
        options={[
          { value: "newest", label: "Newest" },
          {
            value: "focused",
            label: (
              <Tooltip title="Papers attached to the fewest objects first: these study this object, rather than list it in a survey table. SIMBAD has no citation counts.">
                <span>Most focused</span>
              </Tooltip>
            ),
          },
        ]}
      />
      <Select<Period>
        size="small"
        value={period}
        onChange={changePeriod}
        options={PERIODS.map(({ value, label }) => ({ value, label }))}
        className="w-36"
        aria-label="Publication years"
      />
    </Flex>
  );

  if (isPending) {
    return (
      <Flex vertical gap={12}>
        {controls}
        <Skeleton active paragraph={{ rows: 4 }} />
      </Flex>
    );
  }
  if (isError) {
    return (
      <Text type="secondary" className="text-sm">
        Could not load references from SIMBAD.
      </Text>
    );
  }

  const canLoadMore =
    data.references.length < data.total && limit < SIMBAD_REFS_MAX;
  const which = sort === "focused" ? "most focused" : "most recent";

  return (
    <Flex vertical gap={8}>
      {controls}
      {data.total === 0 ? (
        <Empty
          image={Empty.PRESENTED_IMAGE_SIMPLE}
          description={
            period === "all"
              ? `SIMBAD lists no papers about ${mainId}.`
              : `No papers about ${mainId} in this period.`
          }
        />
      ) : (
        <ul
          className={`m-0 list-none divide-y divide-border p-0 ${isFetching ? "opacity-60" : ""}`}
        >
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
                  {[
                    ref.year,
                    ref.journal,
                    ref.objects !== undefined &&
                      `${ref.objects.toLocaleString("en-US")} object${ref.objects === 1 ? "" : "s"}`,
                  ]
                    .filter(Boolean)
                    .join(" · ")}
                </Text>
              </Flex>
            </li>
          ))}
        </ul>
      )}
      <Flex justify="space-between" align="center" wrap gap={8}>
        <Text type="secondary" className="text-xs">
          {data.total > 0 &&
            `Showing the ${data.references.length} ${which} of ${data.total.toLocaleString("en-US")}. `}
          Full lists:{" "}
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
