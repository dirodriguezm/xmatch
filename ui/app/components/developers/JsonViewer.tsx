"use client";

import { useState } from "react";

/** Array items shown before a "show more" button. */
const ARRAY_PREVIEW = 20;
/** Nodes deeper than this start collapsed. */
const OPEN_DEPTH = 3;

function Primitive({ value }: { value: unknown }) {
  if (value === null) return <span className="text-neutral-500">null</span>;
  if (typeof value === "string")
    return <span className="text-emerald-400">{JSON.stringify(value)}</span>;
  if (typeof value === "number")
    return <span className="text-sky-400">{String(value)}</span>;
  if (typeof value === "boolean")
    return <span className="text-amber-400">{String(value)}</span>;
  return <span>{String(value)}</span>;
}

function Key({ name }: { name?: string }) {
  if (name == null) return null;
  return <span className="text-purple-300">{JSON.stringify(name)}: </span>;
}

function JsonNode({
  name,
  value,
  depth,
  last,
}: {
  name?: string;
  value: unknown;
  depth: number;
  last: boolean;
}) {
  const [limit, setLimit] = useState(ARRAY_PREVIEW);
  const [expanded, setExpanded] = useState(depth < OPEN_DEPTH);
  const comma = last ? "" : ",";

  if (value === null || typeof value !== "object") {
    return (
      <div>
        <Key name={name} />
        <Primitive value={value} />
        {comma}
      </div>
    );
  }

  const isArray = Array.isArray(value);
  const entries: [string | undefined, unknown][] = isArray
    ? (value as unknown[]).map((v) => [undefined, v])
    : Object.entries(value as Record<string, unknown>);
  const [open, close] = isArray ? ["[", "]"] : ["{", "}"];

  if (entries.length === 0) {
    return (
      <div>
        <Key name={name} />
        {open}
        {close}
        {comma}
      </div>
    );
  }

  const shown = isArray ? entries.slice(0, limit) : entries;
  const hidden = entries.length - shown.length;
  const summary = isArray
    ? `${entries.length} item${entries.length === 1 ? "" : "s"}`
    : `${entries.length} key${entries.length === 1 ? "" : "s"}`;

  return (
    <details
      open={expanded}
      onToggle={(e) => setExpanded(e.currentTarget.open)}
    >
      <summary className="cursor-pointer list-none select-none marker:hidden [&::-webkit-details-marker]:hidden">
        <span className="inline-block w-3 text-neutral-500">
          {expanded ? "▾" : "▸"}
        </span>
        <Key name={name} />
        {open}
        {!expanded && (
          <span className="text-neutral-500">
            {" "}
            {summary} {close}
            {comma}
          </span>
        )}
      </summary>
      <div className="ml-4 border-l border-border pl-2">
        {shown.map(([k, v], i) => (
          <JsonNode
            key={k ?? i}
            name={k}
            value={v}
            depth={depth + 1}
            last={i === entries.length - 1}
          />
        ))}
        {hidden > 0 && (
          <button
            type="button"
            onClick={() => setLimit((l) => l + ARRAY_PREVIEW * 5)}
            className="cursor-pointer border-0 bg-transparent p-0 text-neutral-400 underline hover:text-foreground"
          >
            … {hidden} more item{hidden === 1 ? "" : "s"} (show{" "}
            {Math.min(hidden, ARRAY_PREVIEW * 5)})
          </button>
        )}
      </div>
      <div>
        <span className="inline-block w-3" />
        {close}
        {comma}
      </div>
    </details>
  );
}

/** Collapsible, lazily-expanded JSON tree; big arrays are truncated. */
export function JsonViewer({ value }: { value: unknown }) {
  return (
    <div className="font-mono text-xs leading-relaxed text-foreground">
      <JsonNode value={value} depth={0} last />
    </div>
  );
}
