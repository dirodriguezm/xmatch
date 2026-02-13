"use client";

import { SearchOutlined } from "@ant-design/icons";
import { Button, Flex, Input, Select, Space, Typography } from "antd";

import { useCrossmatchState } from "@/app/store/crossmatch-context";

const { Title } = Typography;

const serviceOptions = [
  { value: "SIMBAD", label: "SIMBAD" },
  { value: "NED", label: "NED" },
  { value: "VizieR", label: "VizieR" },
];

interface TargetNameResolverProps {
  onResolve?: () => void;
}

export function TargetNameResolver({ onResolve }: TargetNameResolverProps) {
  const { state, dispatch } = useCrossmatchState();

  const handleResolve = () => {
    if (state.resolver.targetName && onResolve) {
      onResolve();
    }
  };

  return (
    <div className="sidebar-section">
      <Flex vertical gap="middle">
        <Space align="center">
          <SearchOutlined className="text-primary" />
          <Title level={5} className="!m-0">
            Target Name Resolver
          </Title>
        </Space>

        <Flex vertical gap="small">
          <Input.Search
            placeholder="e.g., M31, NGC 1234"
            value={state.resolver.targetName}
            onChange={(e) =>
              dispatch({
                type: "SET_RESOLVER",
                payload: { targetName: e.target.value },
              })
            }
            onSearch={handleResolve}
            enterButton={false}
          />

          <Space>
            <Select
              value={state.resolver.service}
              options={serviceOptions}
              onChange={(value) =>
                dispatch({
                  type: "SET_RESOLVER",
                  payload: { service: value },
                })
              }
              className="w-[120px]"
            />
            <Button
              onClick={handleResolve}
              loading={state.resolver.isResolving}
              disabled={!state.resolver.targetName}
            >
              Resolve
            </Button>
          </Space>
        </Flex>
      </Flex>
    </div>
  );
}
