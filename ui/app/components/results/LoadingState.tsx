"use client";

import { Flex } from "antd";

import { ResultsTable } from "./ResultsTable";

export function LoadingState() {
  return (
    <div className="pl-10 pr-8 py-6">
      <Flex vertical gap="large">
        <div className="h-[32px]" />
        <ResultsTable data={[]} loading={true} />
      </Flex>
    </div>
  );
}
