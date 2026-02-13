"use client";

import {
  CopyOutlined,
  EnvironmentOutlined,
  QuestionCircleOutlined,
  StarOutlined,
  SwapOutlined,
  TagOutlined,
} from "@ant-design/icons";
import {
  App,
  Button,
  Card,
  Col,
  Descriptions,
  Empty,
  Flex,
  Row,
  Space,
  Table,
  Tag,
  Tooltip,
  Typography,
} from "antd";

import type { CrossmatchResult } from "@/app/components/results/ResultsTable";

import { AladinViewer } from "./AladinViewer";

const { Title, Text } = Typography;

interface ObjectDetailProps {
  object: CrossmatchResult;
}

const catalogColors: Record<string, string> = {
  "GAIA DR3": "blue",
  SIMBAD: "green",
  "2MASS": "orange",
  WISE: "purple",
  AllWISE: "purple",
};

// Photometry bands configuration
const photometryBands = [
  { band: "G", survey: "Gaia", wavelength: "0.64 μm" },
  { band: "BP", survey: "Gaia", wavelength: "0.51 μm" },
  { band: "RP", survey: "Gaia", wavelength: "0.78 μm" },
  { band: "J", survey: "2MASS", wavelength: "1.24 μm" },
  { band: "H", survey: "2MASS", wavelength: "1.66 μm" },
  { band: "K", survey: "2MASS", wavelength: "2.16 μm" },
  { band: "W1", survey: "WISE", wavelength: "3.4 μm" },
  { band: "W2", survey: "WISE", wavelength: "4.6 μm" },
  { band: "W3", survey: "WISE", wavelength: "12 μm" },
  { band: "W4", survey: "WISE", wavelength: "22 μm" },
];

export function ObjectDetail({ object }: ObjectDetailProps) {
  const { message } = App.useApp();

  const copyToClipboard = (text: string, label: string) => {
    navigator.clipboard.writeText(text);
    message.success(`${label} copied to clipboard`);
  };

  const formatCoordinate = (value: number, decimals = 6) =>
    value.toFixed(decimals);

  // Convert decimal degrees to sexagesimal
  const toHMS = (ra: number) => {
    const hours = ra / 15;
    const h = Math.floor(hours);
    const m = Math.floor((hours - h) * 60);
    const s = ((hours - h) * 60 - m) * 60;
    return `${h}h ${m}m ${s.toFixed(2)}s`;
  };

  const toDMS = (dec: number) => {
    const sign = dec >= 0 ? "+" : "-";
    const absDec = Math.abs(dec);
    const d = Math.floor(absDec);
    const m = Math.floor((absDec - d) * 60);
    const s = ((absDec - d) * 60 - m) * 60;
    return `${sign}${d}° ${m}′ ${s.toFixed(2)}″`;
  };

  const simbadUrl = `https://simbad.cds.unistra.fr/simbad/sim-coo?Coord=${object.ra}+${object.dec}&Radius=2&Radius.unit=arcsec`;
  const vizierUrl = `https://vizier.cds.unistra.fr/viz-bin/VizieR?-c=${object.ra}+${object.dec}&-c.rs=2`;
  const aladinUrl = `https://aladin.cds.unistra.fr/AladinLite/?target=${object.ra}+${object.dec}&fov=0.1`;

  // Mock photometry data - only G mag available from current data
  const photometryData = photometryBands.map((band) => ({
    key: band.band,
    band: band.band,
    survey: band.survey,
    wavelength: band.wavelength,
    mag: band.band === "G" && object.gMag !== undefined ? object.gMag : null,
    error: band.band === "G" && object.gMag !== undefined ? 0.01 : null,
  }));

  const photometryColumns = [
    {
      title: "Band",
      dataIndex: "band",
      key: "band",
      width: 60,
      render: (value: string, record: { survey: string }) => (
        <Tooltip title={record.survey}>
          <Text strong>{value}</Text>
        </Tooltip>
      ),
    },
    {
      title: "λ",
      dataIndex: "wavelength",
      key: "wavelength",
      width: 80,
      render: (value: string) => (
        <Text type="secondary" className="text-xs">
          {value}
        </Text>
      ),
    },
    {
      title: "Mag",
      dataIndex: "mag",
      key: "mag",
      width: 70,
      align: "right" as const,
      render: (value: number | null) => (
        <Text className={`font-mono ${value === null ? "text-border" : ""}`}>
          {value !== null ? value.toFixed(2) : "—"}
        </Text>
      ),
    },
    {
      title: "Error",
      dataIndex: "error",
      key: "error",
      width: 70,
      align: "right" as const,
      render: (value: number | null) => (
        <Text className={`font-mono ${value === null ? "text-border" : ""}`}>
          {value !== null ? `±${value.toFixed(3)}` : "—"}
        </Text>
      ),
    },
  ];

  return (
    <div className="p-6 max-w-6xl mx-auto">
      <Flex vertical gap={24}>
        {/* Header Section with Sky Image */}
        <Row gutter={24}>
          {/* Sky Viewer */}
          <Col xs={24} sm={8} md={6}>
            <Card
              className="bg-surface h-full"
              styles={{ body: { padding: 0 } }}
            >
              <AladinViewer
                center={{ ra: object.ra, dec: object.dec }}
                fov={0.9}
                height={200}
              />
            </Card>
          </Col>

          {/* Object Info */}
          <Col xs={24} sm={16} md={18}>
            <Flex vertical gap={16} className="h-full justify-between">
              {/* Title and Tags */}
              <div>
                <Flex align="center" gap={12} wrap="wrap">
                  <Title level={2} className="!m-0 !mb-1">
                    {object.objectId}
                  </Title>
                  <Tag color={catalogColors[object.catalog] || "default"}>
                    {object.catalog}
                  </Tag>
                </Flex>
                <Flex align="center" gap={8} className="mt-2">
                  <QuestionCircleOutlined className="text-border" />
                  <Text type="secondary">
                    Unknown type (classification unavailable)
                  </Text>
                </Flex>
              </div>

              {/* Coordinates */}
              <Descriptions size="small" column={{ xs: 1, sm: 2 }}>
                <Descriptions.Item label="RA">
                  <Space>
                    <Text className="font-mono">
                      {formatCoordinate(object.ra)}°
                    </Text>
                    <Text type="secondary" className="font-mono text-xs">
                      ({toHMS(object.ra)})
                    </Text>
                    <Button
                      type="text"
                      size="small"
                      icon={<CopyOutlined />}
                      onClick={() =>
                        copyToClipboard(formatCoordinate(object.ra), "RA")
                      }
                    />
                  </Space>
                </Descriptions.Item>
                <Descriptions.Item label="Dec">
                  <Space>
                    <Text className="font-mono">
                      {formatCoordinate(object.dec)}°
                    </Text>
                    <Text type="secondary" className="font-mono text-xs">
                      ({toDMS(object.dec)})
                    </Text>
                    <Button
                      type="text"
                      size="small"
                      icon={<CopyOutlined />}
                      onClick={() =>
                        copyToClipboard(formatCoordinate(object.dec), "Dec")
                      }
                    />
                  </Space>
                </Descriptions.Item>
                <Descriptions.Item label="Angular Distance">
                  <Text className="font-mono">
                    {object.angularDistance.toFixed(3)} arcsec
                  </Text>
                </Descriptions.Item>
              </Descriptions>

              {/* External Links */}
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
            </Flex>
          </Col>
        </Row>

        {/* Photometry and Astrometry Row */}
        <Row gutter={24}>
          {/* Photometry */}
          <Col xs={24} md={14}>
            <Card
              title={
                <Space>
                  <StarOutlined />
                  <span>Photometry</span>
                </Space>
              }
              className="bg-surface"
              size="small"
            >
              <Table
                dataSource={photometryData}
                columns={photometryColumns}
                pagination={false}
                size="small"
                scroll={{ x: 280 }}
              />
            </Card>
          </Col>

          {/* Astrometry */}
          <Col xs={24} md={10}>
            <Card
              title={
                <Space>
                  <EnvironmentOutlined />
                  <span>Astrometry</span>
                </Space>
              }
              className="bg-surface h-full"
              size="small"
            >
              <Descriptions column={1} size="small">
                <Descriptions.Item label="Proper Motion (RA)">
                  <Text className="font-mono text-border">— mas/yr</Text>
                </Descriptions.Item>
                <Descriptions.Item label="Proper Motion (Dec)">
                  <Text className="font-mono text-border">— mas/yr</Text>
                </Descriptions.Item>
                <Descriptions.Item label="Parallax">
                  <Text className="font-mono text-border">— mas</Text>
                </Descriptions.Item>
                <Descriptions.Item label="Radial Velocity">
                  <Text className="font-mono text-border">— km/s</Text>
                </Descriptions.Item>
              </Descriptions>
            </Card>
          </Col>
        </Row>

        {/* Cross-Identifications */}
        <Card
          title={
            <Space>
              <SwapOutlined />
              <span>Cross-Identifications</span>
            </Space>
          }
          className="bg-surface"
          size="small"
        >
          <Empty
            image={Empty.PRESENTED_IMAGE_SIMPLE}
            description={
              <Text type="secondary">No cross-identifications available</Text>
            }
            styles={{ image: { height: 40 } }}
          />
        </Card>
      </Flex>
    </div>
  );
}
