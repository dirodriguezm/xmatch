"use client";

import { Descriptions, Flex, Tag, Tooltip, Typography } from "antd";

import type { ChandraSource, CscRange } from "@/app/lib/utils/csc";
import { formatFlux, formatFluxWithError } from "@/app/lib/utils/format";

const { Text, Link } = Typography;

const CSC_URL = "https://cxc.cfa.harvard.edu/csc/";
const CSCVIEW_URL = "https://cda.cfa.harvard.edu/cscview/";

function fluxText(v?: CscRange): string {
  if (!v) return "—";
  if (v.lo === undefined || v.hi === undefined) return formatFlux(v.value);
  // CSC gives a 68% interval; shown as its half-width.
  return formatFluxWithError(v.value, (v.hi - v.lo) / 2);
}

function hrText(v?: CscRange): string {
  if (!v) return "—";
  if (v.lo === undefined || v.hi === undefined) return v.value.toFixed(2);
  const err = (v.hi - v.lo) / 2;
  if (err < 0.001) return `${v.value.toFixed(3)} ± <0.001`;
  const digits = err < 0.005 ? 3 : 2;
  return `${v.value.toFixed(digits)} ± ${err.toFixed(digits)}`;
}

/** CSC indices grow with the evidence for variability; the scale is shown. */
function indexText(index: number | undefined, max: number): string {
  return index === undefined ? "—" : `${index} of ${max}`;
}

/** Summary of the nearest Chandra Source Catalog 2.2 master source. */
export function ChandraSourcePanel({ source: s }: { source: ChandraSource }) {
  return (
    <Flex vertical gap={12}>
      {(s.extended || s.saturated || s.confused) && (
        <Flex gap={8} wrap>
          {s.saturated && (
            <Tooltip title="Piled up in every Chandra observation: the measured flux and hardness are unreliable">
              <Tag color="warning">Piled up: flux unreliable</Tag>
            </Tooltip>
          )}
          {s.extended && (
            <Tooltip title="Larger than the Chandra PSF in at least one observation, e.g. a jet or diffuse emission">
              <Tag>Extended</Tag>
            </Tooltip>
          )}
          {s.confused && (
            <Tooltip title="Overlaps another source in at least one observation">
              <Tag color="warning">Confused with a neighbour</Tag>
            </Tooltip>
          )}
        </Flex>
      )}

      <Descriptions
        size="small"
        // Without lg/xl/xxl antd falls back to its default and packs five
        // pairs into a row on wide screens.
        column={{ xs: 1, md: 2, lg: 3, xl: 3, xxl: 3 }}
        bordered
        items={[
          {
            key: "name",
            label: "CSC source",
            children: (
              <span>
                <Text strong>{s.name}</Text>{" "}
                <Text type="secondary" className="text-xs">
                  {s.separationArcsec.toFixed(2)}″ away
                  {s.posErrArcsec !== undefined &&
                    ` · position ±${s.posErrArcsec.toFixed(2)}″ (95%)`}
                </Text>
              </span>
            ),
          },
          {
            key: "flux",
            label: "Flux 0.5–7 keV",
            children: fluxText(s.flux),
          },
          {
            key: "hm",
            label: (
              <Tooltip title="(H − M) / (H + M), hard 2–7 keV vs medium 1.2–2 keV; positive is harder">
                <span>Hardness H–M</span>
              </Tooltip>
            ),
            children: hrText(s.hardHm),
          },
          {
            key: "ms",
            label: (
              <Tooltip title="(M − S) / (M + S), medium 1.2–2 keV vs soft 0.5–1.2 keV; positive is harder">
                <span>Hardness M–S</span>
              </Tooltip>
            ),
            children: hrText(s.hardMs),
          },
          {
            key: "var-inter",
            label: (
              <Tooltip title="Likelihood-ratio index between Chandra observations, 0 (constant) to 8 (strongly variable)">
                <span>Variability between observations</span>
              </Tooltip>
            ),
            children: indexText(s.varInter, 8),
          },
          {
            key: "var-intra",
            label: (
              <Tooltip title="Gregory–Loredo index within single observations, 0 (constant) to 10">
                <span>Variability within an observation</span>
              </Tooltip>
            ),
            children: indexText(s.varIntra, 10),
          },
          {
            key: "signif",
            label: "Detection significance",
            children:
              s.significance !== undefined
                ? `${s.significance.toFixed(1)}σ`
                : "—",
          },
        ]}
      />

      <Text type="secondary" className="text-xs">
        From the{" "}
        <Link href={CSC_URL} target="_blank" rel="noopener noreferrer">
          Chandra Source Catalog 2.2
        </Link>{" "}
        master source table, via the CXC cone search ·{" "}
        <Link href={CSCVIEW_URL} target="_blank" rel="noopener noreferrer">
          CSCview
        </Link>{" "}
        for observations and light curves.
      </Text>
    </Flex>
  );
}
