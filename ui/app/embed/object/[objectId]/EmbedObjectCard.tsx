"use client";

import { Empty, Skeleton, Tag, Typography } from "antd";

import { useMetadata } from "@/app/hooks/queries";
import { SITE_NAME } from "@/app/lib/constants/site";
import { toDMS, toHMS } from "@/app/lib/utils/coordinates";
import { catalogLabel, pickPhotometry } from "@/app/lib/utils/share";
import { buildObjectUrl } from "@/app/lib/utils/urls";

const { Text, Title } = Typography;

const TAG_COLORS: Record<string, string> = {
  gaia: "blue",
  allwise: "purple",
  erosita: "magenta",
};

interface EmbedObjectCardProps {
  objectId: string;
  catalog: string | null;
}

function formatValue(value: number, unit: string): string {
  return unit === "mag" ? value.toFixed(3) : value.toExponential(3);
}

/** Compact, self-contained object summary for `<iframe>` embeds. */
export function EmbedObjectCard({ objectId, catalog }: EmbedObjectCardProps) {
  const { data, isLoading, isError } = useMetadata(
    catalog ? { id: objectId, catalog } : null
  );
  const record = data as Record<string, unknown> | null | undefined;
  const ra = typeof record?.ra === "number" ? record.ra : null;
  const dec = typeof record?.dec === "number" ? record.dec : null;
  const photometry = catalog ? pickPhotometry(catalog, record) : [];

  const href = catalog
    ? buildObjectUrl(objectId, catalog)
    : `/object/${encodeURIComponent(objectId)}`;

  return (
    <article className="flex h-full min-h-[396px] flex-col gap-3 rounded-lg border border-border bg-surface p-4">
      <header className="flex items-start justify-between gap-3">
        <div className="min-w-0">
          {catalog && (
            <Tag color={TAG_COLORS[catalog] ?? "default"} className="!mb-1">
              {catalogLabel(catalog)}
            </Tag>
          )}
          <Title level={5} className="!m-0 break-words" title={objectId}>
            {objectId}
          </Title>
        </div>
        <a
          href={href}
          target="_blank"
          rel="noopener noreferrer"
          className="shrink-0 text-sm"
        >
          View on {SITE_NAME} ↗
        </a>
      </header>

      {!catalog ? (
        <Empty
          className="!my-auto"
          description="Add ?catalog=gaia, allwise or erosita to the embed URL"
        />
      ) : isLoading ? (
        <Skeleton active paragraph={{ rows: 5 }} />
      ) : isError || ra === null || dec === null ? (
        <Empty className="!my-auto" description="Object not found" />
      ) : (
        <>
          <section className="grid grid-cols-2 gap-3 rounded-md bg-surface-elevated p-3">
            <div className="flex flex-col">
              <Text className="!text-xs !text-neutral-400">RA (J2000)</Text>
              <Text className="font-mono">{ra.toFixed(6)}°</Text>
              <Text className="font-mono !text-xs !text-neutral-400">
                {toHMS(ra)}
              </Text>
            </div>
            <div className="flex flex-col">
              <Text className="!text-xs !text-neutral-400">Dec (J2000)</Text>
              <Text className="font-mono">
                {dec >= 0 ? "+" : ""}
                {dec.toFixed(6)}°
              </Text>
              <Text className="font-mono !text-xs !text-neutral-400">
                {toDMS(dec)}
              </Text>
            </div>
          </section>

          {photometry.length > 0 ? (
            <section className="min-h-0 flex-1 overflow-auto">
              <table className="w-full border-collapse text-sm">
                <thead>
                  <tr className="text-left text-xs text-neutral-400">
                    <th className="border-b border-border py-1 font-normal">
                      Band
                    </th>
                    <th className="border-b border-border py-1 text-right font-normal">
                      Value
                    </th>
                    <th className="border-b border-border py-1 text-right font-normal">
                      Error
                    </th>
                  </tr>
                </thead>
                <tbody>
                  {photometry.map((row) => (
                    <tr key={row.band} className="text-foreground">
                      <td className="border-b border-border py-1">
                        {row.band}
                      </td>
                      <td className="border-b border-border py-1 text-right font-mono">
                        {formatValue(row.value, row.unit)}{" "}
                        <span className="text-xs text-neutral-400">
                          {row.unit}
                        </span>
                      </td>
                      <td className="border-b border-border py-1 text-right font-mono text-neutral-400">
                        {row.error != null
                          ? formatValue(row.error, row.unit)
                          : "—"}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </section>
          ) : (
            <Text className="!text-neutral-400">
              No photometry in this catalog row.
            </Text>
          )}
        </>
      )}

      <footer className="mt-auto flex items-center justify-between text-xs text-neutral-400">
        <span>Data via the {SITE_NAME} public API</span>
        <a href="/" target="_blank" rel="noopener noreferrer">
          {SITE_NAME}
        </a>
      </footer>
    </article>
  );
}
