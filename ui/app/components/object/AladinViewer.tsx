"use client";

import { Spin } from "antd";
import dynamic from "next/dynamic";
import { forwardRef } from "react";

import type { AladinViewerProps, AladinViewerRef } from "@/types/aladin";

// Dynamically import the inner component to avoid SSR issues
const AladinViewerInner = dynamic(
  () =>
    import("./AladinViewerInner").then((mod) => ({
      default: mod.AladinViewerInner,
    })),
  {
    ssr: false,
    loading: () => (
      <div className="flex items-center justify-center bg-surface h-[200px] w-full">
        <Spin spinning>{null}</Spin>
      </div>
    ),
  }
);

/**
 * AladinViewer - A React wrapper for Aladin Lite sky visualization
 *
 * @example
 * // Basic usage with center coordinates
 * <AladinViewer
 *   center={{ ra: 10.684, dec: 41.269 }}
 *   fov={0.5}
 *   markers={[{ ra: 10.684, dec: 41.269, label: "M31" }]}
 * />
 *
 * @example
 * // With ref for imperative control
 * const aladinRef = useRef<AladinViewerRef>(null);
 * <AladinViewer ref={aladinRef} ... />
 * aladinRef.current?.goTo(180, 45, 1);
 */
export const AladinViewer = forwardRef<AladinViewerRef, AladinViewerProps>(
  function AladinViewer(props, ref) {
    return <AladinViewerInner {...props} ref={ref} />;
  }
);
