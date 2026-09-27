"use client";

import { DownloadOutlined } from "@ant-design/icons";
import { Button, Tabs, Tooltip } from "antd";
import { useMemo } from "react";

import { CodeBlock } from "@/app/components/code/CodeBlock";
import {
  bibFilename,
  buildAcknowledgement,
  buildBibtex,
  buildPlainReferences,
  downloadText,
} from "@/app/lib/utils/citation";

interface CitationTabsProps {
  /** Catalog slugs; empty means every catalog. */
  catalogs: string[];
  objectId?: string;
  size?: "sm" | "md" | "lg";
}

/** BibTeX / plain-text / acknowledgement tabs, each copyable. */
export function CitationTabs({
  catalogs,
  objectId,
  size = "md",
}: CitationTabsProps) {
  const bibtex = useMemo(() => buildBibtex(catalogs), [catalogs]);
  const plain = useMemo(
    () =>
      buildPlainReferences(catalogs)
        .map((r, i) => `[${i + 1}] ${r}`)
        .join("\n\n"),
    [catalogs]
  );
  const ack = useMemo(() => buildAcknowledgement(catalogs), [catalogs]);

  const downloadButton = (
    <Tooltip title="Download .bib">
      <Button
        size="small"
        icon={<DownloadOutlined />}
        aria-label="Download .bib"
        onClick={() => downloadText(bibFilename(objectId), bibtex)}
      />
    </Tooltip>
  );

  return (
    <Tabs
      size="small"
      items={[
        {
          key: "bibtex",
          label: "BibTeX",
          children: (
            <CodeBlock
              code={bibtex}
              copyLabel="BibTeX"
              extra={downloadButton}
              size={size}
            />
          ),
        },
        {
          key: "plain",
          label: "Plain text",
          children: (
            <CodeBlock code={plain} copyLabel="References" wrap size={size} />
          ),
        },
        {
          key: "ack",
          label: "Acknowledgements",
          children: (
            <CodeBlock
              code={ack}
              copyLabel="Acknowledgement"
              wrap
              size={size}
            />
          ),
        },
      ]}
    />
  );
}
