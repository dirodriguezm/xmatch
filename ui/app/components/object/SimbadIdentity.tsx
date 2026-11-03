"use client";

import { Flex, Skeleton, Typography } from "antd";

import { useSimbad } from "@/app/hooks/queries";

const { Text } = Typography;

interface SimbadIdentityProps {
  ra: number;
  dec: number;
  radiusArcsec: number;
}

function simbadIdUrl(mainId: string): string {
  return `https://simbad.cds.unistra.fr/simbad/sim-id?Ident=${encodeURIComponent(mainId)}`;
}

/**
 * The nearest SIMBAD object's main identifier, and its redshift or radial
 * velocity when SIMBAD has one.
 */
export function SimbadIdentity({ ra, dec, radiusArcsec }: SimbadIdentityProps) {
  const { data, isPending, isError } = useSimbad({
    ra,
    dec,
    radius: radiusArcsec,
  });

  if (isPending) {
    return <Skeleton.Input active size="small" className="!w-40 !h-5" />;
  }
  if (isError) {
    return (
      <Text type="secondary" className="text-sm">
        SIMBAD lookup failed
      </Text>
    );
  }
  if (!data.found || !data.match) {
    return (
      <Text type="secondary" className="text-sm">
        Not in SIMBAD (within {radiusArcsec}″)
      </Text>
    );
  }

  const m = data.match;
  const quality = m.rvzQuality ? ` (quality ${m.rvzQuality})` : "";
  return (
    <Flex align="baseline" gap={8} wrap>
      <Text type="secondary" className="text-xs">
        SIMBAD
      </Text>
      <a
        href={simbadIdUrl(m.mainId)}
        target="_blank"
        rel="noopener noreferrer"
        className="font-medium"
      >
        {m.mainId}
      </a>
      <Text type="secondary" className="text-xs">
        {m.separationArcsec.toFixed(1)}″
      </Text>
      {m.redshift !== undefined && (
        <Text className="text-sm" title={`SIMBAD redshift${quality}`}>
          z = {m.redshift.toPrecision(4)}
        </Text>
      )}
      {m.radialVelocity !== undefined && (
        <Text className="text-sm" title={`SIMBAD radial velocity${quality}`}>
          RV = {m.radialVelocity.toFixed(1)} km/s
        </Text>
      )}
    </Flex>
  );
}
