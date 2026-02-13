"use client";

import { PlusOutlined } from "@ant-design/icons";
import {
  App,
  Button,
  Checkbox,
  Flex,
  Input,
  InputNumber,
  Select,
  Space,
  Tag,
  Typography,
} from "antd";
import { useRouter } from "next/navigation";
import { useState } from "react";

import { Logo } from "@/app/components/common";
import { parseCoordinates, resolveObjectName } from "@/app/lib/api/sesame";

const { Title, Text } = Typography;

const unitOptions = [
  { value: "arcsec", label: "arcsec" },
  { value: "arcmin", label: "arcmin" },
  { value: "deg", label: "deg" },
];

interface CatalogOption {
  value: string;
  label: string;
  color: string;
}

const catalogOptions: CatalogOption[] = [
  { value: "gaia_dr3", label: "GAIA DR3", color: "#1890ff" },
  { value: "simbad", label: "SIMBAD", color: "#52c41a" },
  { value: "2mass", label: "2MASS", color: "#fa8c16" },
  { value: "wise", label: "WISE", color: "#722ed1" },
];

interface QuickExample {
  name: string;
  query: string;
}

const quickExamples: QuickExample[] = [
  { name: "M31", query: "M31" },
  { name: "Sagittarius A*", query: "17:45:40.0 -29:00:28" },
  { name: "Crab Pulsar", query: "05:34:31.9 +22:00:52" },
];

export function LandingSearch() {
  const router = useRouter();
  const { message } = App.useApp();
  const [query, setQuery] = useState("");
  const [radius, setRadius] = useState(5);
  const [unit, setUnit] = useState<"arcsec" | "arcmin" | "deg">("arcsec");
  const [selectedCatalogs, setSelectedCatalogs] = useState<string[]>(
    catalogOptions.map((c) => c.value)
  );
  const [isLoading, setIsLoading] = useState(false);

  const handleSubmit = async () => {
    if (!query.trim()) {
      message.warning("Please enter coordinates or an object name");
      return;
    }

    setIsLoading(true);

    try {
      let ra: number;
      let dec: number;

      // Try to parse as coordinates first
      const coords = parseCoordinates(query);
      if (coords) {
        ra = coords.ra;
        dec = coords.dec;
      } else {
        // Try to resolve as object name via Sesame
        const resolved = await resolveObjectName(query);
        if (!resolved) {
          message.error(`Could not resolve "${query}" to coordinates`);
          setIsLoading(false);
          return;
        }
        ra = resolved.ra;
        dec = resolved.dec;
      }

      // Build search URL and navigate
      const params = new URLSearchParams({
        ra: ra.toString(),
        dec: dec.toString(),
        radius: radius.toString(),
        unit: unit,
        catalogs: selectedCatalogs.join(","),
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

  // Map catalog colors to Tailwind classes
  const getCatalogColorClass = (color: string) => {
    const colorMap: Record<string, string> = {
      "#1890ff": "bg-blue-500",
      "#52c41a": "bg-green-500",
      "#fa8c16": "bg-orange-500",
      "#722ed1": "bg-purple-600",
    };
    return colorMap[color] || "bg-gray-500";
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
        className="max-w-[600px] w-full"
      >
        {/* Branding */}
        <Flex vertical align="center" gap="middle">
          <Flex align="center" gap="middle">
            <Logo size="xlarge" />
            <Title level={1} className="!m-0 !text-5xl">
              XWave
            </Title>
          </Flex>
        </Flex>

        {/* Unified search bar */}
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

        {/* Radius selector */}
        <Flex align="center" gap="small" justify="center">
          <Text type="secondary">Radius:</Text>
          <Space.Compact>
            <InputNumber
              value={radius}
              min={0}
              step={0.1}
              onChange={(value) => setRadius(value ?? 1)}
              className="w-20"
            />
            <Select
              value={unit}
              options={unitOptions}
              onChange={(value) => setUnit(value)}
              className="w-[90px]"
            />
          </Space.Compact>
        </Flex>

        {/* Catalog selection */}
        <Flex vertical align="center" gap="small">
          <Text type="secondary">Catalogs:</Text>
          <Checkbox.Group
            value={selectedCatalogs}
            onChange={(values) => setSelectedCatalogs(values as string[])}
          >
            <Flex gap="middle" wrap="wrap" justify="center">
              {catalogOptions.map((catalog) => (
                <Checkbox key={catalog.value} value={catalog.value}>
                  <Flex align="center" gap={4}>
                    <span
                      className={`w-2 h-2 rounded-full inline-block ${getCatalogColorClass(catalog.color)}`}
                    />
                    {catalog.label}
                  </Flex>
                </Checkbox>
              ))}
            </Flex>
          </Checkbox.Group>
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
