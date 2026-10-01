"use client";

import {
  CheckCircleOutlined,
  CloseCircleOutlined,
  DownloadOutlined,
  LoadingOutlined,
  MinusCircleOutlined,
} from "@ant-design/icons";
import { Button, Empty, Flex, Segmented, Tag, Tooltip, Typography } from "antd";
import { useState } from "react";

import { HILIGT_MISSION_IDS, useHiligt } from "@/app/hooks/queries";
import { downloadCsv, toCsv } from "@/app/lib/utils/csv";
import {
  HILIGT_MISSIONS,
  HILIGT_MODEL_NOTE,
  type HiligtPoint,
} from "@/app/lib/utils/hiligt";

import {
  MISSION_DOT_CLASSES,
  XrayLightCurveChart,
} from "./XrayLightCurveChart";

const { Text, Link } = Typography;

const BANDS = [
  { value: "0.2-2", label: "0.2–2 keV" },
  { value: "2-12", label: "2–12 keV" },
  { value: "0.2-12", label: "0.2–12 keV" },
];

const HILIGT_URL = "https://xmmuls.esac.esa.int/hiligt/";

interface XrayHistoryPanelProps {
  ra: number;
  dec: number;
  /** Time axis of the optical light curves, so the panels line up. */
  mjdRange?: { min: number; max: number };
  filenameStem: string;
}

/**
 * X-ray fluxes and upper limits from ESA's HILIGT server. Mounted only when
 * its Collapse panel is opened, so the slow per-mission queries start then.
 */
export function XrayHistoryPanel({
  ra,
  dec,
  mjdRange,
  filenameStem,
}: XrayHistoryPanelProps) {
  // ROSAT only has the soft band, so it is the one every mission shares.
  const [band, setBand] = useState("0.2-2");
  const results = useHiligt({ ra, dec });

  const all: HiligtPoint[] = results.flatMap((r) => r.data?.points ?? []);
  const shown = all.filter((p) => p.band === band);
  const loading = results.some((r) => r.isLoading);

  const download = () =>
    downloadCsv(
      `${filenameStem}_xray_hiligt.csv`,
      toCsv(
        [
          "mission",
          "instrument",
          "obsid",
          "mjd",
          "band_keV",
          "flux_erg_s_cm2",
          "flux_err",
          "upper_limit_erg_s_cm2",
          "ul_sigma",
          "exposure_s",
        ],
        all.map((p) => [
          HILIGT_MISSIONS[p.mission],
          p.instrument,
          p.obsid,
          p.mjd.toFixed(5),
          p.band,
          p.flux,
          p.fluxErr,
          p.upperLimit,
          p.sigma,
          p.exposureS,
        ])
      )
    );

  return (
    <Flex vertical gap={8}>
      <Flex justify="space-between" align="center" wrap gap={8}>
        <Flex gap={6} wrap>
          {HILIGT_MISSION_IDS.map((mission, i) => {
            const r = results[i];
            const n = r.data?.points.filter((p) => p.band === band).length;
            const icon = r.isLoading ? (
              <LoadingOutlined />
            ) : r.isError ? (
              <CloseCircleOutlined />
            ) : n ? (
              <CheckCircleOutlined />
            ) : (
              <MinusCircleOutlined />
            );
            const text = r.isLoading
              ? "loading…"
              : r.isError
                ? "failed"
                : n
                  ? `${n} point${n === 1 ? "" : "s"}`
                  : "no coverage";
            return (
              <Tooltip
                key={mission}
                title={r.isError ? r.error?.message : undefined}
              >
                <Tag
                  icon={icon}
                  color={r.isError ? "error" : undefined}
                  className="!m-0"
                >
                  <span
                    className={`inline-block w-2 h-2 rounded-full mr-1 align-middle ${MISSION_DOT_CLASSES[mission]}`}
                  />
                  {HILIGT_MISSIONS[mission]}: {text}
                </Tag>
              </Tooltip>
            );
          })}
        </Flex>
        <Flex gap={8} align="center">
          <Segmented
            size="small"
            value={band}
            onChange={(v) => setBand(String(v))}
            options={BANDS}
          />
          <Tooltip title="Download all bands as CSV">
            <Button
              size="small"
              type="text"
              icon={<DownloadOutlined />}
              disabled={all.length === 0}
              onClick={download}
              aria-label="Download X-ray history as CSV"
            />
          </Tooltip>
        </Flex>
      </Flex>

      {shown.length > 0 ? (
        <XrayLightCurveChart points={shown} mjdRange={mjdRange} />
      ) : (
        <Empty
          image={Empty.PRESENTED_IMAGE_SIMPLE}
          description={
            <Text type="secondary" className="text-xs">
              {loading
                ? "Asking ESA's upper limit server — each mission takes 5–20 s…"
                : `No X-ray data in the ${band} keV band from these missions.`}
            </Text>
          }
        />
      )}

      <Text type="secondary" className="text-xs">
        ● detection · ▽ upper limit (2σ). Count rates converted to flux with an{" "}
        {HILIGT_MODEL_NOTE.charAt(0).toLowerCase() + HILIGT_MODEL_NOTE.slice(1)}
        . From{" "}
        <Link href={HILIGT_URL} target="_blank" rel="noopener noreferrer">
          HILIGT (ESA)
        </Link>
        , which also covers more missions and spectral models.
      </Text>
    </Flex>
  );
}
