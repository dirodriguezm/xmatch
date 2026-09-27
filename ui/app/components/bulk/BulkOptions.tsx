"use client";

import { ThunderboltOutlined } from "@ant-design/icons";
import { Button, InputNumber, Select, Tooltip, Typography } from "antd";

import { CATALOG_LABELS, CATALOG_OPTIONS } from "@/app/lib/constants/catalogs";
import { MAX_RADIUS_ARCSEC } from "@/app/lib/constants/search";

const { Text } = Typography;

export const BULK_CATALOGS = ["all", ...CATALOG_OPTIONS] as const;
export type BulkCatalog = (typeof BULK_CATALOGS)[number];
export const BULK_MAX_NNEIGHBOR = 100;

const CATALOG_SELECT = [
  { value: "all", label: "All catalogs" },
  ...CATALOG_OPTIONS.map((c) => ({ value: c, label: CATALOG_LABELS[c] })),
];

interface BulkOptionsProps {
  catalog: BulkCatalog;
  radius: number;
  nneighbor: number;
  onCatalogChange: (value: BulkCatalog) => void;
  onRadiusChange: (value: number) => void;
  onNneighborChange: (value: number) => void;
  onRun: () => void;
  running: boolean;
  disabledReason: string | null;
}

export function BulkOptions({
  catalog,
  radius,
  nneighbor,
  onCatalogChange,
  onRadiusChange,
  onNneighborChange,
  onRun,
  running,
  disabledReason,
}: BulkOptionsProps) {
  return (
    <div className="flex flex-wrap items-end gap-4">
      <label className="flex flex-col gap-1">
        <Text className="text-neutral-400 text-xs">Catalog</Text>
        <Select
          value={catalog}
          onChange={onCatalogChange}
          options={CATALOG_SELECT}
          className="w-40"
        />
      </label>
      <label className="flex flex-col gap-1">
        <Text className="text-neutral-400 text-xs">
          Radius (max {MAX_RADIUS_ARCSEC}″)
        </Text>
        <InputNumber
          value={radius}
          min={0.1}
          max={MAX_RADIUS_ARCSEC}
          step={0.5}
          onChange={(v) => v != null && onRadiusChange(v)}
          suffix="arcsec"
          className="w-40"
        />
      </label>
      <label className="flex flex-col gap-1">
        <Tooltip title="Maximum matches returned per input position (per catalog)">
          <Text className="text-neutral-400 text-xs">Nearest N</Text>
        </Tooltip>
        <InputNumber
          value={nneighbor}
          min={1}
          max={BULK_MAX_NNEIGHBOR}
          precision={0}
          onChange={(v) => v != null && onNneighborChange(v)}
          className="w-24"
        />
      </label>
      <Tooltip title={disabledReason}>
        <Button
          type="primary"
          icon={<ThunderboltOutlined />}
          onClick={onRun}
          loading={running}
          disabled={disabledReason !== null}
        >
          Cross-match
        </Button>
      </Tooltip>
    </div>
  );
}
