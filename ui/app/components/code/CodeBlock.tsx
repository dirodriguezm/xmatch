"use client";

import { CopyOutlined } from "@ant-design/icons";
import { App, Button, Tooltip } from "antd";
import { Highlight, Prism } from "prism-react-renderer";
import type { ReactNode } from "react";

const HEIGHT_CLASSES = {
  sm: "max-h-32",
  md: "max-h-[360px]",
  lg: "max-h-[440px]",
} as const;

export type CodeLanguage = "python" | "javascript" | "bash" | "json";

// prism-react-renderer bundles Python, JavaScript and JSON but not bash;
// this is enough for the curl snippets we generate.
Prism.languages.bash ??= {
  comment: { pattern: /(^|\s)#.*/, lookbehind: true },
  string: /"(?:\\.|[^"\\])*"|'[^']*'/,
  parameter: { pattern: /(^|\s)--?[\w-]+/, lookbehind: true },
  function: { pattern: /(^|\n)\s*[a-z][\w-]*/, lookbehind: true },
};

/**
 * Prism token type → Tailwind colour (GitHub dark palette). Tailwind rather
 * than the library's inline theme styles, since inline styles are linted out.
 */
const TOKEN_CLASSES: Record<string, string> = {
  keyword: "text-[#ff7b72]",
  operator: "text-[#ff7b72]",
  string: "text-[#a5d6ff]",
  "triple-quoted-string": "text-[#a5d6ff]",
  "template-string": "text-[#a5d6ff]",
  url: "text-[#a5d6ff]",
  comment: "text-[#8b949e] italic",
  number: "text-[#79c0ff]",
  boolean: "text-[#79c0ff]",
  builtin: "text-[#79c0ff]",
  constant: "text-[#79c0ff]",
  function: "text-[#d2a8ff]",
  "class-name": "text-[#d2a8ff]",
  decorator: "text-[#ffa657]",
  parameter: "text-[#ffa657]",
  property: "text-[#7ee787]",
};

/** Most specific mapped type wins (Prism lists outer types first). */
function tokenClass(types: string[]): string | undefined {
  for (let i = types.length - 1; i >= 0; i--) {
    const cls = TOKEN_CLASSES[types[i]];
    if (cls) return cls;
  }
  return undefined;
}

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
  /** Syntax-highlight the code; plain text when omitted. */
  language?: CodeLanguage;
}

/** Monospace code block with a copy button in the corner. */
export function CodeBlock({
  code,
  copyLabel = "Code",
  extra,
  size = "md",
  wrap = false,
  language,
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
        <code>
          {language ? (
            <Highlight code={code} language={language} prism={Prism}>
              {({ tokens }) =>
                tokens.map((line, i) => (
                  <span key={i}>
                    {line.map((token, j) => {
                      // Empty lines come back as one "\n" token; the line
                      // break below already covers them.
                      if (token.empty) return null;
                      const cls = tokenClass(token.types);
                      return cls ? (
                        <span key={j} className={cls}>
                          {token.content}
                        </span>
                      ) : (
                        token.content
                      );
                    })}
                    {i < tokens.length - 1 && "\n"}
                  </span>
                ))
              }
            </Highlight>
          ) : (
            code
          )}
        </code>
      </pre>
    </div>
  );
}
