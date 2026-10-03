"use client";

import { ReloadOutlined, WarningOutlined } from "@ant-design/icons";
import { Button, Empty, Typography } from "antd";

const { Title, Text } = Typography;

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
    <div className="flex flex-col items-center justify-center h-full">
      <Empty
        image={<WarningOutlined className="text-[64px] text-red-500/70" />}
        description={
          <div className="text-center max-w-md">
            <Title level={4} className="text-foreground mb-2">
              Search failed
            </Title>
            <Text type="secondary">
              {message ||
                "The cone search could not be completed. The catalog service may be unavailable."}
            </Text>
          </div>
        }
      >
        {onRetry && (
          <Button icon={<ReloadOutlined />} onClick={onRetry}>
            Retry search
          </Button>
        )}
      </Empty>
    </div>
  );
}
