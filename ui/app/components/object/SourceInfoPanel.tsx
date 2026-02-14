"use client";

import { CopyOutlined } from "@ant-design/icons";
import { App, Button, Flex, Tag, Tooltip, Typography } from "antd";

const { Title, Text } = Typography;

interface SourceInfoPanelProps {
  objectId: string;
  catalog: string;
  ra: number;
  dec: number;
}

const catalogColors: Record<string, string> = {
  "GAIA DR3": "blue",
  SIMBAD: "green",
  "2MASS": "orange",
  WISE: "purple",
  AllWISE: "purple",
};

export function SourceInfoPanel({
  objectId,
  catalog,
  ra,
  dec,
}: SourceInfoPanelProps) {
  const { message } = App.useApp();

  const copyToClipboard = (text: string, label: string) => {
    navigator.clipboard.writeText(text);
    message.success(`${label} copied`);
  };

  const formatCoordinate = (value: number) => value.toFixed(6);

  // Convert decimal degrees to sexagesimal
  const toHMS = (raVal: number) => {
    const hours = raVal / 15;
    const h = Math.floor(hours);
    const m = Math.floor((hours - h) * 60);
    const s = ((hours - h) * 60 - m) * 60;
    return `${h}h${m}m${s.toFixed(1)}s`;
  };

  const toDMS = (decVal: number) => {
    const sign = decVal >= 0 ? "+" : "-";
    const absDec = Math.abs(decVal);
    const d = Math.floor(absDec);
    const m = Math.floor((absDec - d) * 60);
    const s = ((absDec - d) * 60 - m) * 60;
    return `${sign}${d}d${m}m${s.toFixed(1)}s`;
  };

  return (
    <Flex vertical gap={16} className="h-full">
      <div>
        <Title level={4} className="!m-0 !mb-2 break-words">
          {objectId}
        </Title>
        <Tag color={catalogColors[catalog] || "default"}>{catalog}</Tag>
      </div>

      <Flex vertical gap={8}>
        <div>
          <Text type="secondary" className="text-xs block mb-1">
            RA
          </Text>
          <Flex align="center" gap={4}>
            <Tooltip title={toHMS(ra)}>
              <Text className="font-mono text-sm">{formatCoordinate(ra)}°</Text>
            </Tooltip>
            <Button
              type="text"
              size="small"
              icon={<CopyOutlined />}
              onClick={() => copyToClipboard(formatCoordinate(ra), "RA")}
              className="!p-0 !h-auto"
            />
          </Flex>
        </div>

        <div>
          <Text type="secondary" className="text-xs block mb-1">
            DEC
          </Text>
          <Flex align="center" gap={4}>
            <Tooltip title={toDMS(dec)}>
              <Text className="font-mono text-sm">
                {formatCoordinate(dec)}°
              </Text>
            </Tooltip>
            <Button
              type="text"
              size="small"
              icon={<CopyOutlined />}
              onClick={() => copyToClipboard(formatCoordinate(dec), "DEC")}
              className="!p-0 !h-auto"
            />
          </Flex>
        </div>
      </Flex>
    </Flex>
  );
}
