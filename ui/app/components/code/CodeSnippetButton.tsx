"use client";

import { CodeOutlined } from "@ant-design/icons";
import { Button, Modal, Tooltip } from "antd";
import { useState } from "react";

import { CodeSnippetPanel, type LabeledRequest } from "./CodeSnippetPanel";

export type { LabeledRequest } from "./CodeSnippetPanel";

/** Button that opens copy-as-code (curl / Python / JS) for these requests. */
export function CodeSnippetButton({
  title = "Use the API",
  requests,
  iconOnly = false,
}: {
  title?: string;
  requests: LabeledRequest[];
  /** Icon with a tooltip instead of a labelled button, for tight headers. */
  iconOnly?: boolean;
}) {
  const [open, setOpen] = useState(false);

  return (
    <>
      <Tooltip title={iconOnly ? "Code: run this via the API" : undefined}>
        <Button
          size="small"
          icon={<CodeOutlined />}
          onClick={() => setOpen(true)}
          disabled={requests.length === 0}
          aria-label={iconOnly ? "Code" : undefined}
        >
          {iconOnly ? null : "Code"}
        </Button>
      </Tooltip>
      <Modal
        title={title}
        open={open}
        onCancel={() => setOpen(false)}
        footer={null}
        width={760}
        destroyOnHidden
      >
        <CodeSnippetPanel requests={requests} />
      </Modal>
    </>
  );
}
