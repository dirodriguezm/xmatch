"use client";

import { ExportOutlined } from "@ant-design/icons";
import { Empty, Segmented, Tabs, Typography } from "antd";
import Link from "next/link";
import { useState, useSyncExternalStore } from "react";

import { SWAGGER_URL } from "@/app/lib/constants/site";
import {
  type ApiRequest,
  buildSnippet,
  requestUrl,
  SNIPPET_LANGUAGES,
  type SnippetLanguage,
} from "@/app/lib/utils/snippets";

import { CodeBlock } from "./CodeBlock";

const { Text } = Typography;

export interface LabeledRequest {
  label: string;
  request: ApiRequest;
}

const LANG_STORAGE_KEY = "xwave:snippet-language";

/** In-memory choice when localStorage is unavailable. */
let fallbackLanguage: SnippetLanguage = "curl";

function readStoredLanguage(): SnippetLanguage {
  try {
    const v = window.localStorage.getItem(LANG_STORAGE_KEY);
    if (SNIPPET_LANGUAGES.some((l) => l.key === v)) return v as SnippetLanguage;
  } catch {
    // storage unavailable
  }
  return fallbackLanguage;
}

const languageListeners = new Set<() => void>();

function subscribeLanguage(cb: () => void): () => void {
  languageListeners.add(cb);
  window.addEventListener("storage", cb);
  return () => {
    languageListeners.delete(cb);
    window.removeEventListener("storage", cb);
  };
}

/** Remembered snippet language, shared across panels via localStorage. */
function useSnippetLanguage(): [SnippetLanguage, (l: SnippetLanguage) => void] {
  const lang = useSyncExternalStore(
    subscribeLanguage,
    readStoredLanguage,
    () => "curl" as SnippetLanguage
  );
  const update = (l: SnippetLanguage) => {
    try {
      window.localStorage.setItem(LANG_STORAGE_KEY, l);
    } catch {
      fallbackLanguage = l;
    }
    languageListeners.forEach((cb) => cb());
  };
  return [lang, update];
}

interface CodeSnippetPanelProps {
  requests: LabeledRequest[];
  /** Hide the footer note linking to the playground and Swagger. */
  hideFooter?: boolean;
}

/** Inline copy-as-code view: request picker, language tabs, code and URL. */
export function CodeSnippetPanel({
  requests,
  hideFooter = false,
}: CodeSnippetPanelProps) {
  const [lang, setLang] = useSnippetLanguage();
  const [selected, setSelected] = useState(0);

  if (requests.length === 0) {
    return (
      <Empty
        image={Empty.PRESENTED_IMAGE_SIMPLE}
        description="No request to show — enable at least one catalog."
      />
    );
  }

  const index = Math.min(selected, requests.length - 1);
  const { request } = requests[index];
  const url = requestUrl(request);

  return (
    <div className="flex flex-col gap-3">
      {requests.length > 1 && (
        <Segmented
          block
          value={index}
          onChange={(v) => setSelected(Number(v))}
          options={requests.map((r, i) => ({ label: r.label, value: i }))}
        />
      )}

      <Tabs
        size="small"
        activeKey={lang}
        onChange={(k) => setLang(k as SnippetLanguage)}
        items={SNIPPET_LANGUAGES.map((l) => ({ key: l.key, label: l.label }))}
        className="!mb-0"
      />
      <CodeBlock
        code={buildSnippet(request, lang)}
        copyLabel={`${SNIPPET_LANGUAGES.find((l) => l.key === lang)?.label} snippet`}
      />

      <div className="flex flex-col gap-1">
        <div className="flex items-center justify-between gap-2">
          <Text className="!text-xs uppercase tracking-wide text-neutral-400">
            {request.method} request URL
          </Text>
          {request.method === "GET" && (
            <a
              href={url}
              target="_blank"
              rel="noopener noreferrer"
              className="text-xs"
            >
              Open JSON <ExportOutlined />
            </a>
          )}
        </div>
        <CodeBlock code={url} copyLabel="URL" wrap size="sm" />
      </div>

      {!hideFooter && (
        <Text className="!text-xs text-neutral-400">
          Try parameters interactively in the{" "}
          <Link href="/developers">API playground</Link>, or see every endpoint
          in the{" "}
          <a href={SWAGGER_URL} target="_blank" rel="noopener noreferrer">
            Swagger reference
          </a>
          . No API key needed.
        </Text>
      )}
    </div>
  );
}
