"use client";

import { AntdRegistry } from "@ant-design/nextjs-registry";
import {
  QueryCache,
  QueryClient,
  QueryClientProvider,
} from "@tanstack/react-query";
import { App, ConfigProvider } from "antd";
import { NuqsAdapter } from "nuqs/adapters/next/app";
import { useState } from "react";

import { CommandPalette } from "@/app/components/command/CommandPalette";
import { XWaveSpinner } from "@/app/components/common/XWaveSpinner";
import {
  reportQueryError,
  reportServiceResult,
  serviceForQueryKey,
} from "@/app/lib/serviceHealth";
import { darkTheme } from "@/app/lib/theme/antd-theme";

export function Providers({ children }: { children: React.ReactNode }) {
  const [queryClient] = useState(
    () =>
      new QueryClient({
        // Feed the header's status indicator from queries already running.
        queryCache: new QueryCache({
          onError: (error, query) => reportQueryError(query.queryKey, error),
          onSuccess: (_data, query) =>
            reportServiceResult(serviceForQueryKey(query.queryKey), true),
        }),
        defaultOptions: {
          queries: {
            staleTime: 60 * 1000,
            refetchOnWindowFocus: false,
          },
        },
      })
  );

  return (
    // antd's styles go in their own CSS layer, below Tailwind's utilities, so
    // a utility class on an antd component (padding, width…) actually applies.
    <AntdRegistry layer>
      <ConfigProvider
        theme={darkTheme}
        spin={{ indicator: <XWaveSpinner size={20} /> }}
      >
        <App>
          <QueryClientProvider client={queryClient}>
            <NuqsAdapter>
              {children}
              <CommandPalette />
            </NuqsAdapter>
          </QueryClientProvider>
        </App>
      </ConfigProvider>
    </AntdRegistry>
  );
}
