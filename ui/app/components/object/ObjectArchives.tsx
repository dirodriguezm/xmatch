"use client";

import {
  EnvironmentOutlined,
  LineChartOutlined,
  TagOutlined,
} from "@ant-design/icons";
import { Button, Card, Flex, Space, Tooltip, Typography } from "antd";

import { useDesiTarget } from "@/app/hooks/queries";
import {
  buildAladinUrl,
  buildDesiSpectrumUrl,
  buildLegacySurveyViewerUrl,
  buildNedUrl,
  buildPanstarrsUrl,
  buildSdssNavigateUrl,
  buildSimbadUrl,
  buildVizierUrl,
} from "@/app/lib/utils/urls";

const { Text } = Typography;

interface ObjectArchivesProps {
  ra: number;
  dec: number;
}

// Buttons are grouped by what the destination service returns:
//   - Catalogs       → identifiers / cross-IDs / metadata "tags" (TagOutlined)
//   - Image viewers  → interactive views of the sky region      (EnvironmentOutlined)
//   - Spectra        → 1-D spectrum plot (flux vs. wavelength)  (LineChartOutlined)
// Order within each group is intentional: catalogs by CDS prominence,
// image viewers from general → optical → DESI/imaging → PS, spectra by survey.
export function ObjectArchives({ ra, dec }: ObjectArchivesProps) {
  const simbadUrl = buildSimbadUrl(ra, dec);
  const vizierUrl = buildVizierUrl(ra, dec);
  const nedUrl = buildNedUrl(ra, dec);
  const aladinUrl = buildAladinUrl(ra, dec);
  const sdssUrl = buildSdssNavigateUrl(ra, dec);
  const legacyUrl = buildLegacySurveyViewerUrl(ra, dec);
  const panstarrsUrl = buildPanstarrsUrl(ra, dec);

  // Pre-load the DESI spectrum lookup so the button reflects availability
  // (loading / enabled-with-target / disabled-no-spectrum) before the user clicks.
  const {
    data: desi,
    isLoading: desiLoading,
    isError: desiError,
  } = useDesiTarget({ ra, dec });
  const desiTargetid = desi?.targetid ?? null;
  const desiDisabled = desiLoading || desiError || !desiTargetid;
  let desiTooltip: string;
  if (desiLoading) desiTooltip = "Looking up DESI DR1 spectrum…";
  else if (desiError) desiTooltip = "DESI spectrum lookup failed";
  else if (desiTargetid)
    desiTooltip = `Open DESI DR1 spectrum (offset ${(desi?.separationArcsec ?? 0).toFixed(2)}″)`;
  else desiTooltip = "No DESI DR1 spectrum at this position";

  const groups = [
    {
      label: "Catalogs",
      icon: <TagOutlined />,
      buttons: (
        <>
          <Button
            href={simbadUrl}
            target="_blank"
            rel="noopener noreferrer"
            icon={<TagOutlined />}
            size="small"
          >
            SIMBAD
          </Button>
          <Button
            href={vizierUrl}
            target="_blank"
            rel="noopener noreferrer"
            icon={<TagOutlined />}
            size="small"
          >
            VizieR
          </Button>
          <Button
            href={nedUrl}
            target="_blank"
            rel="noopener noreferrer"
            icon={<TagOutlined />}
            size="small"
          >
            NED
          </Button>
        </>
      ),
    },
    {
      label: "Image viewers",
      icon: <EnvironmentOutlined />,
      buttons: (
        <>
          <Button
            href={aladinUrl}
            target="_blank"
            rel="noopener noreferrer"
            icon={<EnvironmentOutlined />}
            size="small"
          >
            Aladin Lite
          </Button>
          <Button
            href={sdssUrl}
            target="_blank"
            rel="noopener noreferrer"
            icon={<EnvironmentOutlined />}
            size="small"
          >
            SDSS DR19
          </Button>
          <Tooltip title="DESI Legacy Imaging Surveys viewer (imaging + DESI spectroscopic overlay available in the UI)">
            <Button
              href={legacyUrl}
              target="_blank"
              rel="noopener noreferrer"
              icon={<EnvironmentOutlined />}
              size="small"
            >
              Legacy Survey (DESI)
            </Button>
          </Tooltip>
          <Button
            href={panstarrsUrl}
            target="_blank"
            rel="noopener noreferrer"
            icon={<EnvironmentOutlined />}
            size="small"
          >
            Pan-STARRS1
          </Button>
        </>
      ),
    },
    {
      label: "Spectra",
      icon: <LineChartOutlined />,
      buttons: (
        <Tooltip title={desiTooltip}>
          {/* span wrapper lets the tooltip stay reachable while the button is disabled */}
          <span>
            <Button
              size="small"
              icon={<LineChartOutlined />}
              loading={desiLoading}
              disabled={desiDisabled}
              onClick={() =>
                desiTargetid &&
                window.open(
                  buildDesiSpectrumUrl(desiTargetid),
                  "_blank",
                  "noopener,noreferrer"
                )
              }
            >
              DESI Spectrum
            </Button>
          </span>
        </Tooltip>
      ),
    },
  ];

  return (
    <Card title="Archives" size="small" className="bg-surface">
      <Flex vertical gap={12}>
        {groups.map((g) => (
          <Flex
            key={g.label}
            align="center"
            wrap="wrap"
            gap={12}
            className="min-h-[28px]"
          >
            <Flex align="center" gap={6} className="min-w-[120px]">
              {g.icon}
              <Text type="secondary" className="text-xs">
                {g.label}
              </Text>
            </Flex>
            <Space wrap>{g.buttons}</Space>
          </Flex>
        ))}
      </Flex>
    </Card>
  );
}
