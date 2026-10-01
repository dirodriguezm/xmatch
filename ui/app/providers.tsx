"use client";

import { AntdRegistry } from "@ant-design/nextjs-registry";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { App, ConfigProvider } from "antd";
import { NuqsAdapter } from "nuqs/adapters/next/app";
import { useState } from "react";

import { CommandPalette } from "@/app/components/command/CommandPalette";
import { darkTheme } from "@/app/lib/theme/antd-theme";

export function Providers({ children }: { children: React.ReactNode }) {
  const [queryClient] = useState(
    () =>
      new QueryClient({
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
      <ConfigProvider theme={darkTheme}>
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
