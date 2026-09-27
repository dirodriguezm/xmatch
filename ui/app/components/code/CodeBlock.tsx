"use client";

import { CopyOutlined } from "@ant-design/icons";
import { App, Button, Tooltip } from "antd";
import type { ReactNode } from "react";

const HEIGHT_CLASSES = {
  sm: "max-h-32",
  md: "max-h-[360px]",
  lg: "max-h-[440px]",
} as const;

interface CodeBlockProps {
  code: string;
  /** Label used in the "copied" toast. */
  copyLabel?: string;
  /** Extra buttons rendered next to the copy button. */
  extra?: ReactNode;
  /** Maximum height before the block scrolls. */
  size?: "sm" | "md" | "lg";
  /** Wrap long lines instead of scrolling horizontally. */
  wrap?: boolean;
}

/** Monospace code block with a copy button in the corner. */
export function CodeBlock({
  code,
  copyLabel = "Code",
  extra,
  size = "md",
  wrap = false,
}: CodeBlockProps) {
  const { message } = App.useApp();

  const copy = async () => {
    try {
      await navigator.clipboard.writeText(code);
      message.success(`${copyLabel} copied`);
    } catch {
      message.error("Could not access the clipboard");
    }
  };

  return (
    <div className="relative rounded-md border border-border bg-surface">
      <div className="absolute right-2 top-2 flex gap-1">
        {extra}
        <Tooltip title="Copy">
          <Button
            size="small"
            icon={<CopyOutlined />}
            onClick={copy}
            aria-label={`Copy ${copyLabel}`}
          />
        </Tooltip>
      </div>
      <pre
        className={`m-0 overflow-auto ${HEIGHT_CLASSES[size]} p-3 pr-12 font-mono text-xs leading-relaxed text-foreground ${
          wrap ? "whitespace-pre-wrap break-words" : "whitespace-pre"
        }`}
      >
        <code>{code}</code>
      </pre>
    </div>
  );
}
