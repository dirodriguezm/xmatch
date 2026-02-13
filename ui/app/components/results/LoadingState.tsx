"use client";

import { Flex, Typography } from "antd";

import { ResultsTable } from "./ResultsTable";

const { Title } = Typography;

export function LoadingState() {
  return (
    <div className="p-6">
      <Flex vertical gap="large">
        <div className="flex items-center justify-between">
          <Title level={4} className="!m-0">
            Loading Results...
          </Title>
        </div>
        <ResultsTable data={[]} loading={true} />
      </Flex>
    </div>
  );
}
