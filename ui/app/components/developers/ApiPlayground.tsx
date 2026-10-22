"use client";

import {
  ClockCircleOutlined,
  CloudServerOutlined,
  ReloadOutlined,
  SendOutlined,
  StopOutlined,
} from "@ant-design/icons";
import {
  Alert,
  Button,
  Card,
  Empty,
  Input,
  InputNumber,
  Select,
  Switch,
  Tag,
  Tooltip,
  Typography,
} from "antd";
import { type ReactNode, useMemo, useRef, useState } from "react";

import { CodeBlock, CodeSnippetPanel } from "@/app/components/code";
import { CATALOG_SELECT_OPTIONS } from "@/app/lib/constants/catalogs";
import { MAX_RADIUS_ARCSEC } from "@/app/lib/constants/search";
import { requestUrl } from "@/app/lib/utils/snippets";

import { JsonViewer } from "./JsonViewer";
import {
  buildPlaygroundRequest,
  DEFAULT_STATE,
  type EndpointKey,
  ENDPOINTS,
  formatBytes,
  parsePositions,
  type PlaygroundState,
  proxyUrl,
  summarizeResult,
} from "./playground";

const { Text } = Typography;

const CONE_CATALOGS = [
  { value: "all", label: "All catalogs" },
  ...CATALOG_SELECT_OPTIONS,
];
const LIGHTCURVE_CATALOGS = [
  { value: "all", label: "All surveys" },
  { value: "ztf", label: "ZTF" },
  { value: "neowise", label: "NEOWISE" },
  { value: "allwise", label: "AllWISE" },
];
const REQUEST_TIMEOUT_MS = 60_000;

interface ResponseState {
  status: number;
  statusText: string;
  ms: number;
  bytes: number;
  via: "direct" | "proxy";
  data: unknown;
  /** Raw body when it is not JSON. */
  text?: string;
}

function Field({
  label,
  hint,
  children,
}: {
  label: string;
  hint?: ReactNode;
  children: ReactNode;
}) {
  return (
    <label className="flex flex-col gap-1">
      <span className="flex items-baseline justify-between gap-2">
        <span className="font-mono text-xs text-foreground">{label}</span>
        {hint && <span className="text-xs text-neutral-500">{hint}</span>}
      </span>
      {children}
    </label>
  );
}

function statusColor(status: number): string {
  if (status >= 500) return "red";
  if (status >= 400) return "orange";
  if (status === 204) return "blue";
  return "green";
}

export function ApiPlayground() {
  const [endpoint, setEndpoint] = useState<EndpointKey>("conesearch");
  const [state, setState] = useState<PlaygroundState>(DEFAULT_STATE);
  const [loading, setLoading] = useState(false);
  const [response, setResponse] = useState<ResponseState | null>(null);
  const [error, setError] = useState<string | null>(null);
  const abortRef = useRef<AbortController | null>(null);

  const set = <K extends keyof PlaygroundState>(
    key: K,
    value: PlaygroundState[K] | null
  ) => {
    if (value == null) return;
    setState((s) => ({ ...s, [key]: value }));
  };

  const info = ENDPOINTS.find((e) => e.key === endpoint)!;
  const request = useMemo(
    () => buildPlaygroundRequest(endpoint, state),
    [endpoint, state]
  );
  const positions = useMemo(
    () => parsePositions(state.bulkPositions),
    [state.bulkPositions]
  );

  const validation = useMemo(() => {
    if (endpoint === "metadata")
      return state.metadataId.trim()
        ? null
        : { message: "Enter a source id.", blocking: true };
    if (endpoint === "bulk") {
      if (positions.invalid.length > 0)
        return {
          message: `Could not parse line ${positions.invalid.join(", ")}.`,
          blocking: true,
        };
      if (positions.ra.length === 0)
        return { message: "Enter at least one position.", blocking: true };
    }
    const radius =
      endpoint === "lightcurve"
        ? state.lightcurveRadius
        : endpoint === "bulk"
          ? state.bulkRadius
          : state.radius;
    if (radius > MAX_RADIUS_ARCSEC)
      return {
        message: `Radius above ${MAX_RADIUS_ARCSEC}″ may never return.`,
        blocking: false,
      };
    return null;
  }, [endpoint, positions, state]);

  const send = async () => {
    abortRef.current?.abort();
    const controller = new AbortController();
    abortRef.current = controller;
    const timer = setTimeout(() => controller.abort(), REQUEST_TIMEOUT_MS);
    const init: RequestInit =
      request.method === "POST"
        ? {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(request.body),
            signal: controller.signal,
          }
        : { signal: controller.signal };

    setLoading(true);
    setError(null);
    const t0 = performance.now();
    try {
      let via: ResponseState["via"] = "direct";
      let res: Response;
      try {
        // The public API allows CORS from the XWave origins…
        res = await fetch(requestUrl(request), init);
      } catch (e) {
        if (controller.signal.aborted) throw e;
        // …anywhere else (previews, forks), fall back to the Next.js proxy.
        via = "proxy";
        res = await fetch(proxyUrl(request), init);
      }
      const text = await res.text();
      const ms = performance.now() - t0;
      let data: unknown = null;
      let raw: string | undefined;
      if (text) {
        try {
          data = JSON.parse(text);
        } catch {
          raw = text;
        }
      }
      setResponse({
        status: res.status,
        statusText: res.statusText,
        ms,
        bytes: new Blob([text]).size,
        via,
        data,
        text: raw,
      });
    } catch (e) {
      setResponse(null);
      setError(
        controller.signal.aborted
          ? "Request cancelled or timed out."
          : e instanceof Error
            ? e.message
            : "Request failed"
      );
    } finally {
      clearTimeout(timer);
      setLoading(false);
    }
  };

  const cancel = () => abortRef.current?.abort();

  const selectEndpoint = (key: EndpointKey) => {
    setEndpoint(key);
    setResponse(null);
    setError(null);
  };

  const positionFields = (
    <div className="grid grid-cols-2 gap-3">
      <Field label="ra" hint="deg">
        <InputNumber
          value={state.ra}
          min={0}
          max={360}
          step={0.001}
          onChange={(v) => set("ra", v)}
          className="!w-full"
        />
      </Field>
      <Field label="dec" hint="deg">
        <InputNumber
          value={state.dec}
          min={-90}
          max={90}
          step={0.001}
          onChange={(v) => set("dec", v)}
          className="!w-full"
        />
      </Field>
    </div>
  );

  const radiusField = (key: "radius" | "lightcurveRadius" | "bulkRadius") => (
    <Field label="radius" hint={`arcsec, ≤ ${MAX_RADIUS_ARCSEC}`}>
      <InputNumber
        value={state[key]}
        min={0.1}
        max={MAX_RADIUS_ARCSEC}
        onChange={(v) => set(key, v)}
        className="!w-full"
      />
    </Field>
  );

  const catalogField = (
    <Field label="catalog">
      <Select
        value={state.catalog}
        options={CONE_CATALOGS}
        onChange={(v) => set("catalog", v)}
      />
    </Field>
  );

  const nneighborField = (
    <Field label="nneighbor" hint="max per catalog">
      <InputNumber
        value={state.nneighbor}
        min={1}
        max={1000}
        precision={0}
        onChange={(v) => set("nneighbor", v)}
        className="!w-full"
      />
    </Field>
  );

  return (
    <div className="grid gap-6 lg:grid-cols-[minmax(0,340px)_minmax(0,1fr)]">
      {/* Left: endpoint picker + params */}
      <div className="flex flex-col gap-4">
        <nav aria-label="Endpoints" className="flex flex-col gap-1">
          {ENDPOINTS.map((e) => {
            const active = e.key === endpoint;
            return (
              <button
                key={e.key}
                type="button"
                onClick={() => selectEndpoint(e.key)}
                aria-current={active ? "true" : undefined}
                className={`flex cursor-pointer items-center gap-3 rounded-md border px-3 py-2 text-left transition-colors ${
                  active
                    ? "border-border bg-surface-elevated text-foreground"
                    : "border-transparent bg-transparent text-neutral-400 hover:bg-surface-elevated hover:text-foreground"
                }`}
              >
                <Tag
                  color={e.method === "GET" ? "green" : "blue"}
                  className="!m-0 w-12 text-center font-mono"
                >
                  {e.method}
                </Tag>
                <span className="min-w-0 flex-1">
                  <span className="block text-sm">{e.label}</span>
                  <span className="block truncate font-mono text-xs text-neutral-500">
                    {e.path}
                  </span>
                </span>
              </button>
            );
          })}
        </nav>

        <Card
          size="small"
          title="Parameters"
          extra={
            <Tooltip title="Reset to defaults (M31)">
              <Button
                size="small"
                type="text"
                icon={<ReloadOutlined />}
                onClick={() => setState(DEFAULT_STATE)}
                aria-label="Reset parameters"
              />
            </Tooltip>
          }
          className="bg-surface"
        >
          <div className="flex flex-col gap-3">
            <Text className="text-xs text-neutral-400">{info.description}</Text>

            {endpoint === "conesearch" && (
              <>
                {positionFields}
                {radiusField("radius")}
                {catalogField}
                {nneighborField}
                <div className="flex items-center justify-between">
                  <span className="font-mono text-xs">getMetadata</span>
                  <Switch
                    size="small"
                    checked={state.getMetadata}
                    onChange={(v) => set("getMetadata", v)}
                  />
                </div>
              </>
            )}

            {endpoint === "metadata" && (
              <>
                <Field label="id" hint="from a cone search">
                  <Input
                    value={state.metadataId}
                    onChange={(e) => set("metadataId", e.target.value)}
                    spellCheck={false}
                  />
                </Field>
                <Field label="catalog">
                  <Select
                    value={state.metadataCatalog}
                    options={CATALOG_SELECT_OPTIONS}
                    onChange={(v) => set("metadataCatalog", v)}
                  />
                </Field>
                <Text className="text-xs text-neutral-500">
                  Gaia ids look like{" "}
                  <button
                    type="button"
                    className="cursor-pointer border-0 bg-transparent p-0 font-mono text-neutral-300 underline"
                    onClick={() =>
                      setState((s) => ({
                        ...s,
                        metadataId: "Gaia DR3 381266999950756352",
                        metadataCatalog: "gaia",
                      }))
                    }
                  >
                    Gaia DR3 381266999950756352
                  </button>
                  .
                </Text>
              </>
            )}

            {endpoint === "lightcurve" && (
              <>
                {positionFields}
                {radiusField("lightcurveRadius")}
                <Field label="catalog">
                  <Select
                    value={state.lightcurveCatalog}
                    options={LIGHTCURVE_CATALOGS}
                    onChange={(v) => set("lightcurveCatalog", v)}
                  />
                </Field>
              </>
            )}

            {endpoint === "bulk" && (
              <>
                <Field
                  label="ra, dec"
                  hint={`${positions.ra.length} position${positions.ra.length === 1 ? "" : "s"}`}
                >
                  <Input.TextArea
                    value={state.bulkPositions}
                    onChange={(e) => set("bulkPositions", e.target.value)}
                    autoSize={{ minRows: 4, maxRows: 10 }}
                    spellCheck={false}
                    className="font-mono text-xs"
                  />
                </Field>
                <Text className="-mt-1 text-xs text-neutral-500">
                  One position per line, degrees or sexagesimal; # starts a
                  comment.
                </Text>
                {radiusField("bulkRadius")}
                {catalogField}
                {nneighborField}
              </>
            )}

            {validation && (
              <Alert type="warning" showIcon title={validation.message} />
            )}

            <div className="flex gap-2">
              <Button
                type="primary"
                icon={<SendOutlined />}
                loading={loading}
                disabled={validation?.blocking}
                onClick={send}
                block
              >
                Send request
              </Button>
              {loading && (
                <Button icon={<StopOutlined />} onClick={cancel}>
                  Cancel
                </Button>
              )}
            </div>
          </div>
        </Card>
      </div>

      {/* Right: code + response */}
      <div className="flex min-w-0 flex-col gap-4">
        <Card
          size="small"
          title={
            <span className="flex items-center gap-2">
              <Tag
                color={info.method === "GET" ? "green" : "blue"}
                className="!m-0 font-mono"
              >
                {info.method}
              </Tag>
              <span className="font-mono text-sm">/v1{info.path}</span>
            </span>
          }
          className="bg-surface"
        >
          <CodeSnippetPanel
            requests={[{ label: info.label, request }]}
            hideFooter
          />
        </Card>

        <Card
          size="small"
          title="Response"
          className="bg-surface"
          extra={
            response && (
              <span className="flex flex-wrap items-center gap-2 text-xs text-neutral-400">
                <Tag color={statusColor(response.status)} className="!m-0">
                  {response.status} {response.statusText}
                </Tag>
                <span>
                  <ClockCircleOutlined /> {Math.round(response.ms)} ms
                </span>
                <span>{formatBytes(response.bytes)}</span>
                {response.via === "proxy" && (
                  <Tooltip title="Direct call blocked (CORS); sent through this site's proxy">
                    <span>
                      <CloudServerOutlined /> via proxy
                    </span>
                  </Tooltip>
                )}
              </span>
            )
          }
        >
          {error && <Alert type="error" showIcon title={error} />}
          {!error && !response && (
            <Empty
              image={Empty.PRESENTED_IMAGE_SIMPLE}
              description={
                loading
                  ? "Waiting for the API…"
                  : "Press “Send request” to call the live API."
              }
            />
          )}
          {!error && response && (
            <div className="flex flex-col gap-2">
              {response.status === 204 ? (
                <Text className="text-neutral-400">
                  204 No Content — nothing matched. Treat as an empty list.
                </Text>
              ) : response.text !== undefined ? (
                <CodeBlock code={response.text} copyLabel="Response" wrap />
              ) : (
                <>
                  {summarizeResult(response.data) && (
                    <Text className="text-xs text-neutral-400">
                      {summarizeResult(response.data)}
                    </Text>
                  )}
                  <div className="max-h-[520px] overflow-auto rounded-md border border-border bg-background p-3">
                    <JsonViewer value={response.data} />
                  </div>
                </>
              )}
            </div>
          )}
        </Card>
      </div>
    </div>
  );
}
