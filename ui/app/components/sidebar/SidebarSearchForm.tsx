"use client";

import {
  AimOutlined,
  FilterOutlined,
  HomeOutlined,
  SearchOutlined,
} from "@ant-design/icons";
import {
  Button,
  Checkbox,
  Collapse,
  Flex,
  Input,
  InputNumber,
  Select,
  Space,
  Typography,
} from "antd";
import { useState } from "react";

import { useSearchParams } from "@/app/hooks/useSearchParamsSync";
import { useCrossmatchState } from "@/app/store/crossmatch-context";

const { Title, Text } = Typography;

const unitOptions = [
  { value: "arcsec", label: "arcsec" },
  { value: "arcmin", label: "arcmin" },
  { value: "deg", label: "deg" },
];

interface CatalogOption {
  value: string;
  label: string;
  colorClass: string;
}

const catalogOptions: CatalogOption[] = [
  { value: "gaia_dr3", label: "GAIA DR3", colorClass: "bg-blue-500" },
  { value: "simbad", label: "SIMBAD", colorClass: "bg-green-500" },
  { value: "2mass", label: "2MASS", colorClass: "bg-orange-500" },
  { value: "wise", label: "WISE", colorClass: "bg-purple-600" },
];

export function SidebarSearchForm() {
  const { state, dispatch } = useCrossmatchState();
  const {
    ra,
    dec,
    radius,
    unit,
    setRa,
    setDec,
    setRadius,
    setUnit,
    setSearchParams,
    isValid,
  } = useSearchParams();
  const [selectedCatalogs, setSelectedCatalogs] = useState<string[]>(
    catalogOptions.map((c) => c.value)
  );

  const handleUpdateSearch = () => {
    dispatch({ type: "SET_RESULTS_STATE", payload: "loading" });
    // Simulate API call
    setTimeout(() => {
      dispatch({ type: "SET_RESULTS_STATE", payload: "success" });
    }, 1500);
  };

  const handleNewSearch = () => {
    setSearchParams({ ra: null, dec: null, radius: null, unit: null });
    dispatch({ type: "RESET" });
  };

  const handleCatalogChange = (catalog: string, checked: boolean) => {
    if (checked) {
      setSelectedCatalogs([...selectedCatalogs, catalog]);
    } else {
      setSelectedCatalogs(selectedCatalogs.filter((c) => c !== catalog));
    }
  };

  const catalogFilterContent = (
    <Flex vertical gap={12}>
      {catalogOptions.map((catalog) => (
        <Checkbox
          key={catalog.value}
          checked={selectedCatalogs.includes(catalog.value)}
          onChange={(e) => handleCatalogChange(catalog.value, e.target.checked)}
          className="catalog-checkbox"
        >
          <Flex align="center" gap={10}>
            <span
              className={`inline-block w-3 h-3 rounded-sm ${catalog.colorClass} shadow-[0_0_0_1px_rgba(0,0,0,0.25)]`}
            />
            <span className="font-medium">{catalog.label}</span>
          </Flex>
        </Checkbox>
      ))}
    </Flex>
  );

  return (
    <div className="flex flex-col h-full">
      <div className="flex-1 overflow-auto">
        {/* Search Parameters Section */}
        <div className="p-5">
          <Flex vertical gap={20}>
            <Flex align="center" gap={10}>
              <AimOutlined className="text-primary text-base" />
              <Title level={5} className="!m-0">
                Search Parameters
              </Title>
            </Flex>

            <Flex vertical gap={16}>
              <Flex vertical gap={6}>
                <Text type="secondary" className="text-xs font-medium">
                  RA (hms)
                </Text>
                <Input
                  placeholder="e.g., 12:00:00"
                  value={ra}
                  onChange={(e) => setRa(e.target.value)}
                  size="large"
                />
              </Flex>

              <Flex vertical gap={6}>
                <Text type="secondary" className="text-xs font-medium">
                  Dec (dms)
                </Text>
                <Input
                  placeholder="e.g., -45:00:00"
                  value={dec}
                  onChange={(e) => setDec(e.target.value)}
                  size="large"
                />
              </Flex>

              <Flex vertical gap={6}>
                <Text type="secondary" className="text-xs font-medium">
                  Radius
                </Text>
                <Space.Compact className="w-full">
                  <InputNumber
                    value={radius}
                    min={0}
                    step={0.1}
                    onChange={(value) => setRadius(value ?? 1)}
                    className="flex-1"
                    size="large"
                  />
                  <Select
                    value={unit}
                    options={unitOptions}
                    onChange={(value) => setUnit(value)}
                    className="w-[100px]"
                    size="large"
                  />
                </Space.Compact>
              </Flex>
            </Flex>

            <Button
              type="primary"
              block
              icon={<SearchOutlined />}
              onClick={handleUpdateSearch}
              disabled={!isValid}
              loading={state.resultsState === "loading"}
              size="large"
            >
              Update Search
            </Button>
          </Flex>
        </div>

        {/* Active Catalogs Section - Collapsible */}
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
                  <Text strong>Filter Catalogs</Text>
                  <Text type="secondary" className="text-xs">
                    ({selectedCatalogs.length}/{catalogOptions.length})
                  </Text>
                </Flex>
              ),
              children: catalogFilterContent,
            },
          ]}
        />
      </div>

      {/* New Search Button - Fixed at bottom */}
      <div className="p-4 px-5 border-t border-border bg-surface">
        <Button
          block
          icon={<HomeOutlined />}
          onClick={handleNewSearch}
          size="large"
        >
          New Search
        </Button>
      </div>
    </div>
  );
}
