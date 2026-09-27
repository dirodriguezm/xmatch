"use client";

import { CodeOutlined } from "@ant-design/icons";
import { Button, Modal } from "antd";
import { useState } from "react";

import { CodeSnippetPanel, type LabeledRequest } from "./CodeSnippetPanel";

export type { LabeledRequest } from "./CodeSnippetPanel";

/** Button that opens copy-as-code (curl / Python / JS) for these requests. */
export function CodeSnippetButton({
  title = "Use the API",
  requests,
}: {
  title?: string;
  requests: LabeledRequest[];
}) {
  const [open, setOpen] = useState(false);

  return (
    <>
      <Button
        size="small"
        icon={<CodeOutlined />}
        onClick={() => setOpen(true)}
        disabled={requests.length === 0}
      >
        Code
      </Button>
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
