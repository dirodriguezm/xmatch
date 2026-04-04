"use client";

import { PlusOutlined } from "@ant-design/icons";
import { App, Button, Flex, Input, Space, Tag, Typography } from "antd";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";

import { CatalogRadiusRow } from "@/app/components/common";
import { Logo } from "@/app/components/common";
import { parseCoordinates, resolveObjectName } from "@/app/lib/api/sesame";
import {
  buildDefaultCatalogConfigs,
  CATALOG_SELECT_OPTIONS,
} from "@/app/lib/constants/catalogs";
import {
  type CatalogRadiusConfig,
  encodeCatalogRadii,
} from "@/app/lib/constants/search";

const { Title, Text } = Typography;

interface QuickExample {
  name: string;
  query: string;
}

const ALL_EXAMPLES: QuickExample[] = [
  { name: "M31 (Andromeda)", query: "M31" },
  { name: "Sagittarius A*", query: "17:45:40.0 -29:00:28" },
  { name: "Crab Pulsar", query: "05:34:31.9 +22:00:52" },
  { name: "M87 (Virgo A)", query: "M87" },
  { name: "Eta Carinae", query: "Eta Carinae" },
  { name: "Vela Pulsar", query: "08:35:20.6 -45:10:35" },
  { name: "NGC 253 (Sculptor)", query: "NGC 253" },
  { name: "Cygnus X-1", query: "Cygnus X-1" },
  { name: "M1 (Crab Nebula)", query: "M1" },
  { name: "Centaurus A", query: "Centaurus A" },
  { name: "M42 (Orion Nebula)", query: "M42" },
  { name: "Cassiopeia A", query: "Cassiopeia A" },
  { name: "M51 (Whirlpool)", query: "M51" },
  { name: "Sirius", query: "Sirius" },
  { name: "Sombrero Galaxy", query: "M104" },
];

const EXAMPLES_TO_SHOW = 4;

function pickRandom(items: QuickExample[], n: number): QuickExample[] {
  const shuffled = [...items].sort(() => Math.random() - 0.5);
  return shuffled.slice(0, n);
}

export function LandingSearch() {
  const router = useRouter();
  const { message } = App.useApp();
  const [query, setQuery] = useState("");
  const [configs, setConfigs] = useState<CatalogRadiusConfig[]>(
    buildDefaultCatalogConfigs
  );
  const [isLoading, setIsLoading] = useState(false);
  const [quickExamples, setQuickExamples] = useState(
    ALL_EXAMPLES.slice(0, EXAMPLES_TO_SHOW)
  );

  useEffect(() => {
    // Randomize on client mount to avoid SSR hydration mismatch
    // eslint-disable-next-line react-hooks/set-state-in-effect
    setQuickExamples(pickRandom(ALL_EXAMPLES, EXAMPLES_TO_SHOW));
  }, []);

  const updateConfig = (
    catalog: string,
    patch: Partial<CatalogRadiusConfig>
  ) => {
    setConfigs((prev) =>
      prev.map((c) => (c.catalog === catalog ? { ...c, ...patch } : c))
    );
  };

  const handleSubmit = async () => {
    if (!query.trim()) {
      message.warning("Please enter coordinates or an object name");
      return;
    }

    const enabledConfigs = configs.filter((c) => c.enabled);
    if (enabledConfigs.length === 0) {
      message.warning("Please select at least one catalog");
      return;
    }

    setIsLoading(true);

    try {
      let ra: number;
      let dec: number;

      const coords = parseCoordinates(query);
      if (coords) {
        ra = coords.ra;
        dec = coords.dec;
      } else {
        const resolved = await resolveObjectName(query);
        if (!resolved) {
          message.error(`Could not resolve "${query}" to coordinates`);
          setIsLoading(false);
          return;
        }
        ra = resolved.ra;
        dec = resolved.dec;
      }

      const params = new URLSearchParams({
        ra: ra.toString(),
        dec: dec.toString(),
        catalogRadii: encodeCatalogRadii(configs),
      });

      router.push(`/search?${params.toString()}`);
    } catch {
      message.error("Failed to resolve coordinates");
      setIsLoading(false);
    }
  };

  const handleExampleClick = (example: QuickExample) => {
    setQuery(example.query);
  };

  return (
    <Flex
      vertical
      align="center"
      justify="center"
      className="min-h-[calc(100vh-64px)] py-12 px-6 bg-background"
    >
      <Flex
        vertical
        align="center"
        gap="large"
        className="max-w-[620px] w-full"
      >
        {/* Branding */}
        <Flex align="center" gap="middle">
          <Logo size="xlarge" />
          <Title level={1} className="m-0! text-5xl!">
            XWave
          </Title>
        </Flex>

        {/* Search bar */}
        <Space.Compact className="w-full">
          <Button icon={<PlusOutlined />} size="large" title="Upload file" />
          <Input.Search
            placeholder="Coordinates or name (e.g., 12:30:00 -45:00:00 or M31)"
            size="large"
            enterButton
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            onSearch={handleSubmit}
            loading={isLoading}
            className="flex-1"
          />
        </Space.Compact>

        {/* Per-catalog radius inputs */}
        <Flex vertical align="center" gap="small">
          <Text type="secondary">Catalogs:</Text>
          <Flex vertical gap={8}>
            {CATALOG_SELECT_OPTIONS.map((catalog) => {
              const config = configs.find((c) => c.catalog === catalog.value)!;
              return (
                <CatalogRadiusRow
                  key={catalog.value}
                  label={catalog.label}
                  catalog={catalog.value}
                  config={config}
                  onChange={(patch) => updateConfig(catalog.value, patch)}
                />
              );
            })}
          </Flex>
        </Flex>

        {/* Quick examples */}
        <Flex vertical align="center" gap="small">
          <Text type="secondary">Quick examples:</Text>
          <Space size="small" wrap>
            {quickExamples.map((example) => (
              <Tag
                key={example.name}
                color="default"
                className="cursor-pointer py-1 px-3"
                onClick={() => handleExampleClick(example)}
              >
                {example.name}
              </Tag>
            ))}
          </Space>
        </Flex>
      </Flex>
    </Flex>
  );
}
