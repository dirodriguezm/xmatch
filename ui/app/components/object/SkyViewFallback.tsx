"use client";

import { ReloadOutlined } from "@ant-design/icons";
import { Button, Flex, Typography } from "antd";
import { useState } from "react";

import type { Neighbor } from "@/app/hooks/queries";

import { NeighborMap } from "./NeighborMap";

const { Text } = Typography;

interface SkyViewFallbackProps {
  neighbors: Neighbor[];
  radiusArcsec: number;
  retrying: boolean;
  onRetry: () => void;
}

/**
 * Stands in for Aladin when the CDS image servers don't answer: the indexed
 * sources around the object, which come from XWave and always load.
 */
export function SkyViewFallback({
  neighbors,
  radiusArcsec,
  retrying,
  onRetry,
}: SkyViewFallbackProps) {
  const [hoveredKey, setHoveredKey] = useState<string | null>(null);

  return (
    <Flex
      vertical
      align="center"
      justify="center"
      gap={16}
      className="min-h-[380px] flex-1 px-6 py-4"
    >
      <NeighborMap
        neighbors={neighbors}
        radiusArcsec={radiusArcsec}
        hoveredKey={hoveredKey}
        onHover={setHoveredKey}
      />
      <Flex vertical align="center" gap={6} className="max-w-sm text-center">
        <Text strong className="text-sm">
          Sky images unavailable
        </Text>
        <Text type="secondary" className="text-xs">
          CDS isn&apos;t responding. Showing XWave sources within {radiusArcsec}
          ″ instead.
        </Text>
        <Button
          size="small"
          icon={<ReloadOutlined />}
          loading={retrying}
          onClick={onRetry}
        >
          Try images again
        </Button>
      </Flex>
    </Flex>
  );
}
