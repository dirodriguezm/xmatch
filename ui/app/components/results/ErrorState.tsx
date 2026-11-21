"use client";

import { ReloadOutlined, WarningOutlined } from "@ant-design/icons";
import { Button, Typography } from "antd";

const { Text } = Typography;

export interface ErrorStateProps {
  /** Message from the failed query, when one made it through. */
  message?: string;
  onRetry?: () => void;
}

/**
 * Sibling of {@link EmptyState} for the case where a search actually failed.
 * These used to share one component, so a failed search read to the user as
 * "you haven't searched yet" — the opposite of what had happened.
 */
export function ErrorState({ message, onRetry }: ErrorStateProps) {
  return (
    <div className="flex max-w-md flex-col items-center gap-1 text-center">
      <span className="mb-3 flex size-11 items-center justify-center rounded-full border border-border bg-surface text-lg text-red-400">
        <WarningOutlined />
      </span>
      <Text strong className="text-base">
        Search failed
      </Text>
      <Text type="secondary">
        {message ||
          "The cone search could not be completed. The catalog service may be unavailable."}
      </Text>
      {onRetry && (
        <Button icon={<ReloadOutlined />} onClick={onRetry} className="mt-3">
          Retry search
        </Button>
      )}
    </div>
  );
}
