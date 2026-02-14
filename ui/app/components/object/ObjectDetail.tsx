"use client";

import {
  CopyOutlined,
  EnvironmentOutlined,
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
  Flex,
  Row,
  Space,
  Tooltip,
  Typography,
} from "antd";

import type { CrossmatchResult } from "@/app/components/results/ResultsTable";
import { PHOTOMETRY_BANDS } from "@/app/lib/constants/bands";
import {
  buildAladinUrl,
  buildSimbadUrl,
  buildVizierUrl,
} from "@/app/lib/utils/urls";

import { AladinViewer } from "./AladinViewer";

const { Text, Title } = Typography;

interface ObjectDetailProps {
  object: CrossmatchResult;
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

export function ObjectDetail({ object }: ObjectDetailProps) {
  const { message } = App.useApp();
  const simbadUrl = buildSimbadUrl(object.ra, object.dec);
  const vizierUrl = buildVizierUrl(object.ra, object.dec);
  const aladinUrl = buildAladinUrl(object.ra, object.dec);

  const copyToClipboard = (text: string, label: string) => {
    navigator.clipboard.writeText(text);
    message.success(`${label} copied to clipboard`);
  };

  // Build photometry data
  const photometryData = PHOTOMETRY_BANDS.map((band) => ({
    band: band.band,
    survey: band.survey,
    mag: band.band === "G" && object.gMag !== undefined ? object.gMag : null,
  }));

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
                  {object.gMag !== undefined && (
                    <>
                      <Text type="secondary">•</Text>
                      <Text>
                        G ={" "}
                        <span className="font-mono">
                          {object.gMag.toFixed(2)}
                        </span>{" "}
                        mag
                      </Text>
                    </>
                  )}
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
          >
            <AladinViewer
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
        defaultActiveKey={["photometry"]}
        className="bg-surface"
      />
    </div>
  );
}
