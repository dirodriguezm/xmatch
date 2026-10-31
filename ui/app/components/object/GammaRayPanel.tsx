"use client";

import { Alert, Flex, Segmented, Typography } from "antd";
import { useState } from "react";

import { XWaveSpinner } from "@/app/components/common/XWaveSpinner";
import { useBatLightcurve, useFermiLightcurve } from "@/app/hooks/queries";
import {
  BAT_BASE_URL,
  binPoints,
  GAMMA_SOURCES_DATE,
  type GammaAssociations,
  LCR_CADENCES,
  type LcrCadence,
} from "@/app/lib/utils/gamma";

import { GammaLightCurveChart } from "./GammaLightCurveChart";

const { Text, Link } = Typography;

const LCR_SOURCE_URL =
  "https://fermi.gsfc.nasa.gov/ssc/data/access/lat/LightCurveRepository/source.html?source_name=";
/** BAT days averaged per point; a faint source is noise day by day. */
const BAT_BIN_DAYS = 30;

const FERMI_COLOR = "#b37feb";
const BAT_COLOR = "#ff85c0";

function Loading() {
  return (
    <div className="flex h-56 items-center justify-center">
      <XWaveSpinner variant="lightcurve" size={48} />
    </div>
  );
}

type MjdRange = { min: number; max: number };

function FermiSection({
  name,
  mjdRange,
}: {
  name: string;
  mjdRange?: MjdRange;
}) {
  const [cadence, setCadence] = useState<LcrCadence>("weekly");
  const { data, isLoading, error } = useFermiLightcurve(name, cadence);

  return (
    <Flex vertical gap={8}>
      <Flex justify="space-between" align="center" wrap gap={8}>
        <Text strong>Fermi-LAT · 0.1–100 GeV</Text>
        <Segmented<LcrCadence>
          size="small"
          value={cadence}
          onChange={setCadence}
          options={(Object.keys(LCR_CADENCES) as LcrCadence[]).map((c) => ({
            value: c,
            label: LCR_CADENCES[c],
          }))}
        />
      </Flex>
      {isLoading ? (
        <Loading />
      ) : error ? (
        <Alert type="warning" showIcon message={error.message} />
      ) : data && data.points.length > 0 ? (
        <GammaLightCurveChart
          points={data.points}
          color={FERMI_COLOR}
          yLabel="Photon flux (ph cm⁻² s⁻¹)"
          log
          mjdRange={mjdRange}
        />
      ) : (
        <Text type="secondary">No bins in the repository yet.</Text>
      )}
      <Text type="secondary" className="text-xs">
        Bins with TS ≥ 4 as points with 1σ bars, the rest as 95% upper limits
        (▼). From the{" "}
        <Link
          href={`${LCR_SOURCE_URL}${encodeURIComponent(name.replace(/ /g, "_"))}`}
          target="_blank"
          rel="noopener noreferrer"
        >
          Fermi-LAT Light Curve Repository
        </Link>
        .
      </Text>
    </Flex>
  );
}

function BatSection({ path, mjdRange }: { path: string; mjdRange?: MjdRange }) {
  const { data, isLoading, error } = useBatLightcurve(path);
  const points = data ? binPoints(data.points, BAT_BIN_DAYS) : [];

  return (
    <Flex vertical gap={8}>
      <Text strong>Swift-BAT · 15–50 keV</Text>
      {isLoading ? (
        <Loading />
      ) : error ? (
        <Alert type="warning" showIcon message={error.message} />
      ) : points.length > 0 ? (
        <GammaLightCurveChart
          points={points}
          color={BAT_COLOR}
          yLabel="Rate (counts cm⁻² s⁻¹)"
          symbolSize={5}
          mjdRange={mjdRange}
        />
      ) : (
        <Text type="secondary">No good days in the monitor yet.</Text>
      )}
      <Text type="secondary" className="text-xs">
        {BAT_BIN_DAYS}-day weighted means of the daily rates, 1σ bars. From the{" "}
        <Link
          href={`${BAT_BASE_URL}${path}`}
          target="_blank"
          rel="noopener noreferrer"
        >
          Swift-BAT Hard X-ray Transient Monitor
        </Link>
        .
      </Text>
    </Flex>
  );
}

/**
 * Gamma-ray and hard X-ray light curves of the source this object is the
 * counterpart of. Rendered only when there is an association.
 */
export function GammaRayPanel({
  fermi,
  bat,
  mjdRange,
}: GammaAssociations & { mjdRange?: MjdRange }) {
  return (
    <Flex vertical gap={24}>
      {fermi && <FermiSection name={fermi.name} mjdRange={mjdRange} />}
      {bat && <BatSection path={bat.path} mjdRange={mjdRange} />}
      <Text type="secondary" className="text-xs">
        Matched through the counterpart, not the gamma-ray position, which is
        only good to arcminutes: this object is within{" "}
        {Math.max(
          fermi?.separationArcsec ?? 0,
          bat?.separationArcsec ?? 0
        ).toFixed(1)}
        ″ of the source associated with{" "}
        {[fermi && `${fermi.name} (4FGL)`, bat && `${bat.name} (BAT monitor)`]
          .filter(Boolean)
          .join(" and ")}
        . Association list from {GAMMA_SOURCES_DATE}.
      </Text>
    </Flex>
  );
}
