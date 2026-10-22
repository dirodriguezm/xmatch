"use client";

import { Descriptions, Flex, Tooltip, Typography } from "antd";

import {
  formatFlux,
  formatMjdDate,
  formatScientific,
} from "@/app/lib/utils/format";
import type { XmmSource, XmmValue } from "@/app/lib/utils/xmm5";

const { Text, Link } = Typography;

const XSA_URL = "https://nxsa.esac.esa.int/nxsa-web/#search";

function fluxText(v?: XmmValue): string {
  if (!v) return "—";
  return v.err
    ? `${formatFlux(v.value)} ± ${formatScientific(v.err, 1)}`
    : formatFlux(v.value);
}

function hrText(v?: XmmValue): string {
  if (!v) return "—";
  if (v.err === undefined) return v.value.toFixed(2);
  // Bright sources have tiny errors; show them rather than "± 0.00".
  if (v.err < 0.001) return `${v.value.toFixed(3)} ± <0.001`;
  const digits = v.err < 0.005 ? 3 : 2;
  return `${v.value.toFixed(digits)} ± ${v.err.toFixed(digits)}`;
}

function probability(p?: number): string {
  return p === undefined ? "" : ` (match probability ${(p * 100).toFixed(0)}%)`;
}

/** Summary of the nearest 5XMM-DR15 source (stacked over all observations). */
export function XmmSourcePanel({ source: s }: { source: XmmSource }) {
  let variability: string;
  if (s.variable === undefined) variability = "—";
  else if (s.variable)
    variability = `Variable${s.fvar ? ` · Fvar = ${s.fvar.value.toFixed(3)}${s.fvar.err ? ` ± ${s.fvar.err.toFixed(3)}` : ""}` : ""}`;
  else variability = "Not significantly variable";

  return (
    <Flex vertical gap={12}>
      <Descriptions
        size="small"
        column={{ xs: 1, md: 2 }}
        bordered
        items={[
          {
            key: "name",
            label: "5XMM source",
            children: (
              <span>
                <Text strong>{s.name}</Text>{" "}
                <Text type="secondary" className="text-xs">
                  {s.separationArcsec.toFixed(1)}″ away
                </Text>
              </span>
            ),
          },
          {
            key: "flux",
            label: "Flux 0.2–12 keV",
            children: fluxText(s.flux),
          },
          {
            key: "obs",
            label: "Observations",
            children:
              s.observations !== undefined
                ? `${s.observations}${
                    s.mjdFirst !== undefined && s.mjdLast !== undefined
                      ? ` · ${formatMjdDate(s.mjdFirst)} → ${formatMjdDate(s.mjdLast)}`
                      : ""
                  }`
                : "—",
          },
          {
            key: "var",
            label: (
              <Tooltip title="Long-term variability between XMM-Newton observations, from the stacked catalogue">
                <span>Variability</span>
              </Tooltip>
            ),
            children: variability,
          },
          {
            key: "spec",
            label: (
              <Tooltip title="Simple absorbed power-law fit from the catalogue">
                <span>Spectrum</span>
              </Tooltip>
            ),
            children:
              s.gamma !== undefined
                ? `Γ = ${s.gamma.toFixed(2)}${s.nh !== undefined ? ` · NH = ${formatScientific(s.nh, 1)} cm⁻²` : ""}`
                : "—",
          },
          {
            key: "extent",
            label: "Extent",
            children:
              s.extentArcsec === undefined
                ? "—"
                : s.extentArcsec > 0
                  ? `Extended (${s.extentArcsec.toFixed(1)}″)`
                  : "Point-like",
          },
          {
            key: "gaia",
            label: "Gaia DR3 counterpart",
            children: s.gaia
              ? `${s.gaia.sourceId}${probability(s.gaia.probability)}`
              : "—",
          },
          {
            key: "wise",
            label: "WISE counterpart",
            children: s.wise
              ? `${s.wise.name}${probability(s.wise.probability)}`
              : "—",
          },
        ]}
      />

      <Descriptions
        size="small"
        column={{ xs: 1, sm: 2, md: 3 }}
        title={
          <Text className="text-xs text-neutral-400">EPIC band fluxes</Text>
        }
        items={s.bandFluxes.map((b) => ({
          key: b.label,
          label: b.label,
          children: (
            <Text className="font-mono text-xs">{fluxText(b.flux)}</Text>
          ),
        }))}
      />

      <Descriptions
        size="small"
        column={{ xs: 2, md: 4 }}
        title={
          <Tooltip title="HRn = (B(n+1) − Bn) / (B(n+1) + Bn) between adjacent EPIC bands; positive is harder">
            <Text className="text-xs text-neutral-400">Hardness ratios</Text>
          </Tooltip>
        }
        items={s.hardnessRatios.map((h) => ({
          key: h.label,
          label: h.label,
          children: <Text className="font-mono text-xs">{hrText(h.hr)}</Text>,
        }))}
      />

      <Text type="secondary" className="text-xs">
        From the 5XMM-DR15 stacked catalogue (Webb, Traulsen et al. 2026) via
        the{" "}
        <Link href={XSA_URL} target="_blank" rel="noopener noreferrer">
          XMM-Newton Science Archive
        </Link>
        {s.infoUrl && (
          <>
            {" "}
            ·{" "}
            <Link href={s.infoUrl} target="_blank" rel="noopener noreferrer">
              source page with all counterparts (XMM SSC)
            </Link>
          </>
        )}
        .
      </Text>
    </Flex>
  );
}
