"use client";

import { AimOutlined, FilterOutlined, SearchOutlined } from "@ant-design/icons";
import { Button, Collapse, Flex, Input, Typography } from "antd";
import { useState } from "react";

import { CatalogRadiusRow } from "@/app/components/common";
import { useSearchParams } from "@/app/hooks/useSearchParamsSync";
import {
  buildDefaultCatalogConfigs,
  CATALOG_SELECT_OPTIONS,
} from "@/app/lib/constants/catalogs";
import {
  type CatalogRadiusConfig,
  decodeCatalogRadii,
  encodeCatalogRadii,
} from "@/app/lib/constants/search";
import { useCrossmatchState } from "@/app/store/crossmatch-context";

const { Title, Text } = Typography;

function initConfigs(catalogRadiiStr: string): CatalogRadiusConfig[] {
  if (!catalogRadiiStr) return buildDefaultCatalogConfigs();
  const decoded = decodeCatalogRadii(catalogRadiiStr);
  if (decoded.length === 0) return buildDefaultCatalogConfigs();
  return decoded;
}

export function SidebarSearchForm() {
  const { state, dispatch } = useCrossmatchState();
  const {
    ra,
    dec,
    catalogRadii: catalogRadiiStr,
    setSearchParams,
  } = useSearchParams();

  const [draftRa, setDraftRa] = useState(ra);
  const [draftDec, setDraftDec] = useState(dec);
  const [draftConfigs, setDraftConfigs] = useState<CatalogRadiusConfig[]>(() =>
    initConfigs(catalogRadiiStr)
  );

  const isDraftValid = Boolean(draftRa && draftDec);

  const updateConfig = (
    catalog: string,
    patch: Partial<CatalogRadiusConfig>
  ) => {
    setDraftConfigs((prev) =>
      prev.map((c) => (c.catalog === catalog ? { ...c, ...patch } : c))
    );
  };

  const handleUpdateSearch = () => {
    dispatch({ type: "SET_RESULTS_STATE", payload: "loading" });
    setSearchParams({
      ra: draftRa,
      dec: draftDec,
      catalogRadii: encodeCatalogRadii(draftConfigs),
    });
  };

  return (
    <div className="flex flex-col h-full">
      <div className="flex-1 overflow-auto">
        {/* Search Parameters Section */}
        <div className="p-5">
          <Flex vertical gap={20}>
            <Flex align="center" gap={10}>
              <AimOutlined className="text-primary text-base" />
              <Title level={5} className="m-0!">
                Search Parameters
              </Title>
            </Flex>

            <Flex vertical gap={16}>
              <Flex vertical gap={6}>
                <Text type="secondary" className="text-xs font-medium">
                  RA (°)
                </Text>
                <Input
                  placeholder="e.g., 10.6847"
                  value={draftRa}
                  onChange={(e) => setDraftRa(e.target.value)}
                  size="large"
                />
              </Flex>

              <Flex vertical gap={6}>
                <Text type="secondary" className="text-xs font-medium">
                  Dec (°)
                </Text>
                <Input
                  placeholder="e.g., 41.2689"
                  value={draftDec}
                  onChange={(e) => setDraftDec(e.target.value)}
                  size="large"
                />
              </Flex>
            </Flex>

            <Button
              type="primary"
              block
              icon={<SearchOutlined />}
              onClick={handleUpdateSearch}
              disabled={!isDraftValid}
              loading={state.resultsState === "loading"}
              size="large"
            >
              Search
            </Button>
          </Flex>
        </div>

        {/* Catalogs & Radii Section - Collapsible */}
        <Collapse
          defaultActiveKey={["catalogs"]}
          ghost
          expandIconPlacement="end"
          className="border-t border-border"
          items={[
            {
              key: "catalogs",
              label: (
                <Flex align="center" gap={10}>
                  <FilterOutlined className="text-primary text-base" />
                  <Text strong>Catalogs & Radii</Text>
                  <Text type="secondary" className="text-xs">
                    ({draftConfigs.filter((c) => c.enabled).length}/
                    {draftConfigs.length})
                  </Text>
                </Flex>
              ),
              children: (
                <Flex vertical gap={10} className="pb-2">
                  {CATALOG_SELECT_OPTIONS.map((catalog) => {
                    const config = draftConfigs.find(
                      (c) => c.catalog === catalog.value
                    )!;
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
              ),
            },
          ]}
        />
      </div>
    </div>
  );
}
