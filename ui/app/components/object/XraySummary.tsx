"use client";

import { Table, Typography } from "antd";

import type { ChandraSource } from "@/app/lib/utils/csc";
import { formatFlux, formatFluxWithError } from "@/app/lib/utils/format";
import type { XmmSource } from "@/app/lib/utils/xmm5";

const { Text } = Typography;

export interface ErositaSummary {
  id: string;
  separationArcsec?: number;
  /** eRASS1 main band 0.2–2.3 keV (erg s⁻¹ cm⁻²). */
  flux?: number;
  fluxErr?: number;
}

interface Row {
  key: string;
  mission: string;
  source: string;
  separationArcsec?: number;
  band: string;
  flux: string;
  variability: string;
}

function fluxCell(value?: number, err?: number): string {
  if (value === undefined) return "—";
  return err !== undefined
    ? formatFluxWithError(value, err)
    : formatFlux(value);
}

/**
 * One row per X-ray catalogue with a source at this position: is it an
 * X-ray source, how bright, does it vary. The panels below hold the rest.
 */
export function XraySummary({
  erosita,
  xmm,
  chandra,
}: {
  erosita?: ErositaSummary;
  xmm?: XmmSource;
  chandra?: ChandraSource;
}) {
  const rows: Row[] = [];
  if (erosita) {
    rows.push({
      key: "erosita",
      mission: "eROSITA eRASS1",
      source: erosita.id,
      separationArcsec: erosita.separationArcsec,
      band: "0.2–2.3 keV",
      flux: fluxCell(erosita.flux, erosita.fluxErr),
      variability: "single scan",
    });
  }
  if (xmm) {
    rows.push({
      key: "xmm",
      mission: "XMM-Newton 5XMM",
      source: xmm.name,
      separationArcsec: xmm.separationArcsec,
      band: "0.2–12 keV",
      flux: fluxCell(xmm.flux?.value, xmm.flux?.err),
      variability:
        xmm.variable === undefined
          ? "—"
          : xmm.variable
            ? `variable${xmm.fvar ? ` (Fvar ${xmm.fvar.value.toFixed(3)})` : ""}`
            : "not variable",
    });
  }
  if (chandra) {
    const err =
      chandra.flux?.lo !== undefined && chandra.flux.hi !== undefined
        ? (chandra.flux.hi - chandra.flux.lo) / 2
        : undefined;
    rows.push({
      key: "chandra",
      mission: "Chandra CSC 2.2",
      source: chandra.name,
      separationArcsec: chandra.separationArcsec,
      band: "0.5–7 keV",
      flux: `${fluxCell(chandra.flux?.value, err)}${chandra.saturated ? " (piled up)" : ""}`,
      variability:
        chandra.varInter === undefined
          ? "—"
          : `index ${chandra.varInter} of 8 between observations`,
    });
  }
  if (rows.length === 0) return null;

  return (
    <div className="mb-4">
      <Table<Row>
        size="small"
        pagination={false}
        dataSource={rows}
        scroll={{ x: "max-content" }}
        columns={[
          { title: "Catalogue", dataIndex: "mission" },
          {
            title: "Source",
            render: (_, r) => (
              <span>
                <span className="font-mono text-xs">{r.source}</span>
                {r.separationArcsec !== undefined && (
                  <Text type="secondary" className="ml-1.5 text-xs">
                    {r.separationArcsec.toFixed(1)}″
                  </Text>
                )}
              </span>
            ),
          },
          { title: "Band", dataIndex: "band" },
          {
            title: "Flux",
            render: (_, r) => <span className="tabular-nums">{r.flux}</span>,
          },
          { title: "Variability", dataIndex: "variability" },
        ]}
      />
      <Text type="secondary" className="mt-1.5 block text-xs">
        The bands differ, so fluxes are comparable only roughly. Details for
        each catalogue are below.
      </Text>
    </div>
  );
}
