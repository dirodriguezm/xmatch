"use client";

import { XWaveSpinner } from "@/app/components/common/XWaveSpinner";

export function LoadingState() {
  return (
    <div className="flex h-[calc(100vh-64px)] items-center justify-center p-8">
      <XWaveSpinner variant="crossmatch" size={64} label="Cross-matching…" />
    </div>
  );
}
