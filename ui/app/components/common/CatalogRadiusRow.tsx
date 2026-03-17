"use client";

import { Checkbox, Flex, InputNumber, Select, Space, Typography } from "antd";

import { CATALOG_COLOR_CLASSES } from "@/app/lib/constants/catalogs";
import {
  type CatalogRadiusConfig,
  RADIUS_UNIT_OPTIONS,
  type RadiusUnit,
} from "@/app/lib/constants/search";

const { Text } = Typography;

interface CatalogRadiusRowProps {
  label: string;
  catalog: string;
  config: CatalogRadiusConfig;
  onChange: (patch: Partial<CatalogRadiusConfig>) => void;
}

export function CatalogRadiusRow({
  label,
  catalog,
  config,
  onChange,
}: CatalogRadiusRowProps) {
  return (
    <Flex align="center" gap={10}>
      <Checkbox
        checked={config.enabled}
        onChange={(e) => onChange({ enabled: e.target.checked })}
      />
      <Flex align="center" gap={6} className="w-[80px] shrink-0">
        <span
          className={`w-2 h-2 rounded-full inline-block shrink-0 ${CATALOG_COLOR_CLASSES[catalog] ?? "bg-gray-500"}`}
        />
        <Text>{label}</Text>
      </Flex>
      <Space.Compact>
        <InputNumber
          value={config.radius}
          min={0}
          step={0.1}
          disabled={!config.enabled}
          onChange={(value) => onChange({ radius: value ?? 1 })}
          className="w-[72px]"
        />
        <Select
          value={config.unit}
          options={RADIUS_UNIT_OPTIONS}
          disabled={!config.enabled}
          onChange={(value) => onChange({ unit: value as RadiusUnit })}
          className="w-[85px]"
        />
      </Space.Compact>
    </Flex>
  );
}
