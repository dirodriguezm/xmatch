"use client";

import { DownloadOutlined } from "@ant-design/icons";
import { Button } from "antd";

import { ResultsActions } from "@/app/components/actions/ResultsActions";
import { nearbyCsv } from "@/app/components/results/nearby";
import type { NearbySource } from "@/app/components/results/nearby/shared";
import { downloadCsv } from "@/app/lib/utils/csv";

interface ResultsToolbarProps {
  target: { ra: number; dec: number };
  sources: NearbySource[];
}

/** Code / Share / Export actions for the query bar. */
export function ResultsToolbar({ target, sources }: ResultsToolbarProps) {
  return (
    <ResultsActions
      target={target}
      extra={
        sources.length > 0 && (
          <Button
            size="small"
            icon={<DownloadOutlined />}
            onClick={() =>
              downloadCsv("xwave_crossmatch.csv", nearbyCsv(sources))
            }
          >
            Export CSV
          </Button>
        )
      }
    />
  );
}
