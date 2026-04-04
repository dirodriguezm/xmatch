"use client";

import {
  CopyOutlined,
  DatabaseOutlined,
  EnvironmentOutlined,
  LineChartOutlined,
  QuestionCircleOutlined,
  StarOutlined,
  TagOutlined,
} from "@ant-design/icons";
import {
  App,
  Button,
  Card,
  Col,
  Collapse,
  Descriptions,
  Flex,
  Row,
  Select,
  Space,
  Tooltip,
  Typography,
} from "antd";
import { useRef } from "react";

import type { CrossmatchResult } from "@/app/components/results/ResultsTable";
import { useZtfLightcurve } from "@/app/hooks/queries";
import { PHOTOMETRY_BANDS } from "@/app/lib/constants/bands";
import {
  buildAladinUrl,
  buildSimbadUrl,
  buildVizierUrl,
} from "@/app/lib/utils/urls";
import type { AladinViewerRef } from "@/types/aladin";
import type { components } from "@/types/xwave-api";

import { AladinViewer } from "./AladinViewer";
import { LightCurveChart } from "./LightCurveChart";

const DSS_SURVEY = "https://alasky.cds.unistra.fr/DSS/DSSColor/";

const SURVEY_OPTIONS = [
  { label: "DSS Optical", value: DSS_SURVEY, category: "Optical" },
  {
    label: "DESI DR10",
    value: "CDS/P/DESI-Legacy-Surveys/DR10/color",
    category: "Optical",
  },
  { label: "DSS2 Color", value: "CDS/P/DSS2/color", category: "Optical" },
  { label: "2MASS", value: "CDS/P/2MASS/color", category: "Infrared" },
  { label: "AllWISE", value: "CDS/P/allWISE/color", category: "Infrared" },
  { label: "XMM-Newton", value: "xcatdb/P/XMM/PN/color", category: "X-ray" },
  {
    label: "Chandra",
    value: "cxc.harvard.edu/P/cda/hips/allsky/rgb",
    category: "X-ray",
  },
  { label: "NVSS 1.4 GHz", value: "CDS/P/NVSS", category: "Radio" },
  { label: "SUMSS 843 MHz", value: "CDS/P/SUMSS", category: "Radio" },
  {
    label: "RACS 887 MHz",
    value: "https://casda.csiro.au/hips/RACS/low/I/",
    category: "Radio",
  },
];

const surveySelectOptions = Object.entries(
  SURVEY_OPTIONS.reduce<Record<string, { label: string; value: string }[]>>(
    (acc, s) => {
      (acc[s.category] ??= []).push({ label: s.label, value: s.value });
      return acc;
    },
    {}
  )
).map(([category, options]) => ({ label: category, options }));

const { Text, Title } = Typography;

type Allwise = components["schemas"]["repository.Allwise"];

interface ObjectDetailProps {
  object: CrossmatchResult;
  metadata?: Allwise | null;
}

function mag(field?: number): number | null {
  if (field == null) return null;
  return field;
}

// Convert decimal degrees to sexagesimal
function toHMS(ra: number): string {
  const hours = ra / 15;
  const h = Math.floor(hours);
  const m = Math.floor((hours - h) * 60);
  const s = ((hours - h) * 60 - m) * 60;
  return `${h.toString().padStart(2, "0")}h ${m.toString().padStart(2, "0")}m ${s.toFixed(2)}s`;
}

function toDMS(dec: number): string {
  const sign = dec >= 0 ? "+" : "-";
  const absDec = Math.abs(dec);
  const d = Math.floor(absDec);
  const m = Math.floor((absDec - d) * 60);
  const s = ((absDec - d) * 60 - m) * 60;
  return `${sign}${d}° ${m.toString().padStart(2, "0")}′ ${s.toFixed(2)}″`;
}

export function ObjectDetail({ object, metadata }: ObjectDetailProps) {
  const { message } = App.useApp();
  const aladinRef = useRef<AladinViewerRef>(null);
  const {
    data: ztfData,
    isLoading: ztfLoading,
    error: ztfError,
  } = useZtfLightcurve({ ra: object.ra, dec: object.dec });
  const simbadUrl = buildSimbadUrl(object.ra, object.dec);
  const vizierUrl = buildVizierUrl(object.ra, object.dec);
  const aladinUrl = buildAladinUrl(object.ra, object.dec);

  const copyToClipboard = (text: string, label: string) => {
    navigator.clipboard.writeText(text);
    message.success(`${label} copied to clipboard`);
  };

  // Map metadata fields to photometry bands (supports multiple catalogs)
  const meta = metadata as Record<string, unknown> | undefined;
  const bandMagMap: Record<string, number | null> = {
    G: mag(meta?.phot_g_mean_mag as number | undefined),
    BP: mag(meta?.phot_bp_mean_mag as number | undefined),
    RP: mag(meta?.phot_rp_mean_mag as number | undefined),
    J: mag(metadata?.j_m_2mass),
    H: mag(metadata?.h_m_2mass),
    K: mag(metadata?.k_m_2mass),
    W1: mag(metadata?.w1mpro),
    W2: mag(metadata?.w2mpro),
    W3: mag(metadata?.w3mpro),
    W4: mag(metadata?.w4mpro),
  };

  const photometryData = PHOTOMETRY_BANDS.map((band) => ({
    band: band.band,
    survey: band.survey,
    mag: bandMagMap[band.band] ?? null,
  }));

  // Build catalog details from all metadata fields (exclude id, ra, dec already shown)
  const excludedFields = new Set(["id", "ra", "dec"]);
  const catalogDetails = meta
    ? Object.entries(meta)
        .filter(([key]) => !excludedFields.has(key))
        .map(([key, value]) => ({
          key,
          label: key,
          value: value == null ? "—" : String(value),
        }))
    : [];

  const collapseItems = [
    {
      key: "photometry",
      label: (
        <Space>
          <StarOutlined />
          <span>Photometry</span>
          <Text type="secondary" className="text-xs">
            ({photometryData.filter((p) => p.mag !== null).length} bands)
          </Text>
        </Space>
      ),
      children: (
        <Flex wrap="wrap" gap={8}>
          {photometryData.map((p) => (
            <Tooltip key={p.band} title={p.survey}>
              <div
                className={`text-center px-3 py-2 rounded border ${
                  p.mag !== null
                    ? "border-primary bg-primary/10"
                    : "border-border bg-surface"
                }`}
              >
                <Text
                  strong
                  className={p.mag === null ? "text-border" : undefined}
                >
                  {p.band}
                </Text>
                <div
                  className={`font-mono text-sm ${p.mag === null ? "text-border" : ""}`}
                >
                  {p.mag !== null ? p.mag.toFixed(2) : "—"}
                </div>
              </div>
            </Tooltip>
          ))}
        </Flex>
      ),
    },
    {
      key: "ztf-lightcurve",
      label: (
        <Space>
          <LineChartOutlined />
          <span>Light Curve (ZTF)</span>
          {ztfData?.detections && ztfData.detections.length > 0 && (
            <Text type="secondary" className="text-xs">
              ({ztfData.detections.length} points)
            </Text>
          )}
        </Space>
      ),
      children: (
        <LightCurveChart
          data={ztfData}
          loading={ztfLoading}
          error={ztfError ?? null}
        />
      ),
    },
    ...(catalogDetails.length > 0
      ? [
          {
            key: "catalog-details",
            label: (
              <Space>
                <DatabaseOutlined />
                <span>Catalog Details</span>
                <Text type="secondary" className="text-xs">
                  ({catalogDetails.length} fields)
                </Text>
              </Space>
            ),
            children: (
              <Descriptions
                size="small"
                column={{ xs: 1, sm: 2, md: 3 }}
                bordered
              >
                {catalogDetails.map((field) => (
                  <Descriptions.Item
                    key={field.key}
                    label={
                      <Text className="font-mono text-xs">{field.label}</Text>
                    }
                  >
                    <Text className="font-mono text-xs">{field.value}</Text>
                  </Descriptions.Item>
                ))}
              </Descriptions>
            ),
          },
        ]
      : []),
    {
      key: "external",
      label: (
        <Space>
          <TagOutlined />
          <span>External Catalogs</span>
        </Space>
      ),
      children: (
        <Space wrap>
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
            href={aladinUrl}
            target="_blank"
            rel="noopener noreferrer"
            icon={<EnvironmentOutlined />}
            size="small"
          >
            Aladin Lite
          </Button>
        </Space>
      ),
    },
  ];

  return (
    <div className="p-6 max-w-6xl mx-auto">
      {/* Main two-column layout */}
      <Row gutter={[24, 24]} className="mb-6">
        {/* Left: Object Info */}
        <Col xs={24} md={14}>
          <Card className="bg-surface h-full" size="small">
            <Flex vertical gap={16}>
              {/* Object Name */}
              <div>
                <Title level={3} className="!m-0 !mb-2">
                  {object.objectId}
                </Title>
                <Flex align="center" gap={8}>
                  <QuestionCircleOutlined className="text-border" />
                  <Text type="secondary">Unknown type</Text>
                </Flex>
              </div>

              {/* Coordinates */}
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div>
                  <Text type="secondary" className="text-xs block mb-1">
                    Right Ascension
                  </Text>
                  <Flex align="center" gap={8}>
                    <Text className="font-mono">{object.ra.toFixed(6)}°</Text>
                    <Text type="secondary" className="text-xs">
                      ({toHMS(object.ra)})
                    </Text>
                    <Button
                      type="text"
                      size="small"
                      icon={<CopyOutlined />}
                      onClick={() =>
                        copyToClipboard(object.ra.toFixed(6), "RA")
                      }
                    />
                  </Flex>
                </div>
                <div>
                  <Text type="secondary" className="text-xs block mb-1">
                    Declination
                  </Text>
                  <Flex align="center" gap={8}>
                    <Text className="font-mono">{object.dec.toFixed(6)}°</Text>
                    <Text type="secondary" className="text-xs">
                      ({toDMS(object.dec)})
                    </Text>
                    <Button
                      type="text"
                      size="small"
                      icon={<CopyOutlined />}
                      onClick={() =>
                        copyToClipboard(object.dec.toFixed(6), "Dec")
                      }
                    />
                  </Flex>
                </div>
              </div>

              {/* Quick Links */}
              <Space wrap>
                <Button
                  href={simbadUrl}
                  target="_blank"
                  rel="noopener noreferrer"
                  size="small"
                >
                  SIMBAD
                </Button>
                <Button
                  href={vizierUrl}
                  target="_blank"
                  rel="noopener noreferrer"
                  size="small"
                >
                  VizieR
                </Button>
                <Button
                  href={aladinUrl}
                  target="_blank"
                  rel="noopener noreferrer"
                  size="small"
                >
                  Aladin
                </Button>
              </Space>
            </Flex>
          </Card>
        </Col>

        {/* Right: Aladin Viewer */}
        <Col xs={24} md={10}>
          <Card
            className="bg-surface h-full"
            styles={{ body: { padding: 0 } }}
            title={
              <Space>
                <EnvironmentOutlined />
                <span>Sky View</span>
              </Space>
            }
            extra={
              <Select
                defaultValue={DSS_SURVEY}
                size="small"
                className="w-[150px]"
                options={surveySelectOptions}
                onChange={(value) => aladinRef.current?.setSurvey(value)}
              />
            }
          >
            <AladinViewer
              ref={aladinRef}
              center={{ ra: object.ra, dec: object.dec }}
              fov={0.9}
              height={280}
            />
          </Card>
        </Col>
      </Row>

      {/* Collapsible sections */}
      <Collapse
        items={collapseItems}
        defaultActiveKey={["photometry", "ztf-lightcurve"]}
        className="bg-surface"
      />
    </div>
  );
}
