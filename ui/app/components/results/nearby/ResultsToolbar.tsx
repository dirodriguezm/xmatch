"use client";

import { DownloadOutlined } from "@ant-design/icons";
import { Button, Tooltip } from "antd";

import { ResultsActions } from "@/app/components/actions/ResultsActions";
import { nearbyCsv } from "@/app/components/results/nearby";
import type { NearbySource } from "@/app/components/results/nearby/shared";
import { downloadCsv } from "@/app/lib/utils/csv";

interface ResultsToolbarProps {
  target: { ra: number; dec: number };
  sources: NearbySource[];
}

/** Code / Share / Export as compact icon buttons, beside the results title. */
export function ResultsToolbar({ target, sources }: ResultsToolbarProps) {
  return (
    <ResultsActions
      target={target}
      iconOnly
      extra={
        sources.length > 0 && (
          <Tooltip title="Export the matches as CSV">
            <Button
              size="small"
              icon={<DownloadOutlined />}
              aria-label="Export CSV"
              onClick={() =>
                downloadCsv("xwave_crossmatch.csv", nearbyCsv(sources))
              }
            />
          </Tooltip>
        )
      }
    />
  );
}
