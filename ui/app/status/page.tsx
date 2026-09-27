"use client";

import {
  ApiOutlined,
  CloudServerOutlined,
  DatabaseOutlined,
  ReloadOutlined,
} from "@ant-design/icons";
import { useQueries } from "@tanstack/react-query";
import { Badge, Button, Table, Tag, Tooltip, Typography } from "antd";
import Link from "next/link";

import { useApiHealth } from "@/app/components/feedback/useApiHealth";
import { PageShell } from "@/app/components/layout";
import { Section } from "@/app/components/learn/Section";
import { CATALOG_META } from "@/app/lib/constants/catalogMeta";
import { API_ORIGIN } from "@/app/lib/constants/site";

import {
  type CheckState,
  classifyCheck,
  DEGRADED_LATENCY_MS,
  overallState,
  UPSTREAM_CHECKS,
  type UpstreamCheck,
} from "./checks";

const { Text, Paragraph } = Typography;

/** Upstream probes are heavier than /ping; re-run them less often. */
const UPSTREAM_POLL_MS = 5 * 60_000;
const PROBE_TIMEOUT_MS = 20_000;

interface ProbeResult {
  state: CheckState;
  httpStatus: number | null;
  latencyMs: number | null;
  checkedAt: string;
}

async function probe(check: UpstreamCheck): Promise<ProbeResult> {
  const started = performance.now();
  try {
    const resp = await fetch(check.path, {
      cache: "no-store",
      signal: AbortSignal.timeout(PROBE_TIMEOUT_MS),
    });
    await resp.arrayBuffer();
    const latencyMs = Math.round(performance.now() - started);
    return {
      state: classifyCheck(resp.status, latencyMs),
      httpStatus: resp.status,
      latencyMs,
      checkedAt: new Date().toISOString(),
    };
  } catch {
    return {
      state: "down",
      httpStatus: null,
      latencyMs: null,
      checkedAt: new Date().toISOString(),
    };
  }
}

const STATE_BADGE: Record<
  CheckState,
  { status: "success" | "warning" | "error"; label: string; color: string }
> = {
  up: { status: "success", label: "Operational", color: "green" },
  degraded: { status: "warning", label: "Degraded", color: "gold" },
  down: { status: "error", label: "Down", color: "red" },
};

function StateTag({ state }: { state: CheckState | "checking" }) {
  if (state === "checking") {
    return (
      <Tag>
        <Badge status="processing" /> Checking
      </Tag>
    );
  }
  const b = STATE_BADGE[state];
  return (
    <Tag color={b.color}>
      <Badge status={b.status} /> {b.label}
    </Tag>
  );
}

function timeLabel(iso: string | undefined) {
  return iso ? new Date(iso).toLocaleTimeString() : "—";
}

const CATALOG_ROWS = Object.values(CATALOG_META).map((m) => ({
  key: m.slug,
  name: m.name,
  release: m.release,
  sources: m.sources,
  coverage: m.coverage,
}));

export default function StatusPage() {
  const health = useApiHealth();
  const upstream = useQueries({
    queries: UPSTREAM_CHECKS.map((check) => ({
      queryKey: ["status-probe", check.id],
      queryFn: () => probe(check),
      refetchInterval: UPSTREAM_POLL_MS,
      staleTime: 0,
      retry: false,
    })),
  });

  const apiState: CheckState | "checking" = health.isPending
    ? "checking"
    : health.data?.status === "up"
      ? classifyCheck(200, health.data.latencyMs)
      : "down";
  const settled = [
    ...(apiState === "checking" ? [] : [apiState]),
    ...upstream.flatMap((q) => (q.data ? [q.data.state] : [])),
  ];
  const allChecked =
    apiState !== "checking" && upstream.every((q) => q.data !== undefined);
  const overall = overallState(settled);
  const anyFetching = health.isFetching || upstream.some((q) => q.isFetching);

  const rerun = () => {
    void health.refetch();
    upstream.forEach((q) => void q.refetch());
  };

  const bannerTone = !allChecked
    ? "border-border"
    : overall === "up"
      ? "border-green-700/60"
      : overall === "degraded"
        ? "border-amber-600/60"
        : "border-red-700/60";
  const bannerText = !allChecked
    ? "Running checks…"
    : overall === "up"
      ? "All systems operational"
      : apiState === "down"
        ? "The XWave API is not responding"
        : overall === "down"
          ? "Some external services are unavailable"
          : "Some services are slow";

  return (
    <PageShell
      title="Status"
      description="Live health of the XWave API and of the external services the object page uses. Checks run from your browser through this site."
      actions={
        <Button icon={<ReloadOutlined />} onClick={rerun} loading={anyFetching}>
          Re-run checks
        </Button>
      }
    >
      <div
        className={`rounded-lg border bg-surface px-5 py-4 mb-10 flex items-center gap-3 ${bannerTone}`}
        role="status"
        aria-live="polite"
      >
        <Badge
          status={!allChecked ? "processing" : STATE_BADGE[overall].status}
        />
        <Text className="text-foreground text-base font-medium">
          {bannerText}
        </Text>
      </div>

      <Section title="XWave API" icon={<ApiOutlined />}>
        <div className="rounded-lg border border-border bg-surface p-4 grid gap-4 sm:grid-cols-4">
          <div>
            <Text className="text-neutral-500 text-xs block">State</Text>
            <StateTag state={apiState} />
          </div>
          <div>
            <Text className="text-neutral-500 text-xs block">Latency</Text>
            <Text className="text-foreground font-mono">
              {health.data?.latencyMs != null
                ? `${health.data.latencyMs} ms`
                : "—"}
            </Text>
          </div>
          <div>
            <Text className="text-neutral-500 text-xs block">Last checked</Text>
            <Text className="text-foreground font-mono">
              {timeLabel(health.data?.checkedAt)}
            </Text>
          </div>
          <div>
            <Text className="text-neutral-500 text-xs block">Endpoint</Text>
            <Text code className="text-xs">
              {API_ORIGIN.replace("https://", "")}/ping
            </Text>
          </div>
          {health.data?.error && (
            <Text className="!text-red-400 text-sm sm:col-span-4">
              {health.data.error}
            </Text>
          )}
        </div>
        <Paragraph className="!text-neutral-500 text-xs !mt-2">
          Pinged every 30 seconds via <Text code>/api/health</Text>.
        </Paragraph>
      </Section>

      <Section title="External services" icon={<CloudServerOutlined />}>
        <Table
          size="small"
          pagination={false}
          rowKey="id"
          scroll={{ x: 600 }}
          dataSource={UPSTREAM_CHECKS.map((check, i) => ({
            ...check,
            result: upstream[i].data,
          }))}
          columns={[
            {
              title: "Service",
              dataIndex: "name",
              render: (_: unknown, row) => (
                <div>
                  <Text className="text-foreground">{row.name}</Text>
                  <Text className="text-neutral-500 text-xs block">
                    {row.provider} · {row.usedFor}
                  </Text>
                </div>
              ),
            },
            {
              title: "State",
              width: 140,
              render: (_: unknown, row) => (
                <StateTag state={row.result?.state ?? "checking"} />
              ),
            },
            {
              title: "Latency",
              width: 110,
              render: (_: unknown, row) => (
                <span className="font-mono text-sm">
                  {row.result?.latencyMs != null ? (
                    row.cached ? (
                      <Tooltip title="This proxy caches upstream answers; latency may reflect the cache.">
                        {row.result.latencyMs} ms*
                      </Tooltip>
                    ) : (
                      `${row.result.latencyMs} ms`
                    )
                  ) : (
                    "—"
                  )}
                </span>
              ),
            },
            {
              title: "Checked",
              width: 110,
              render: (_: unknown, row) => (
                <span className="font-mono text-sm text-neutral-400">
                  {timeLabel(row.result?.checkedAt)}
                </span>
              ),
            },
          ]}
        />
        <Paragraph className="!text-neutral-500 text-xs !mt-2">
          Each service is probed with a query for M31 every 5 minutes. Answers
          slower than {DEGRADED_LATENCY_MS / 1000} s count as degraded. * Cached
          server-side, so latency may not reflect the upstream service.
        </Paragraph>
      </Section>

      <Section title="Catalog data" icon={<DatabaseOutlined />}>
        <Paragraph className="!text-neutral-400">
          Releases currently indexed. These are static; changes are announced in
          the <Link href="/changelog">changelog</Link>.
        </Paragraph>
        <Table
          size="small"
          pagination={false}
          scroll={{ x: 600 }}
          dataSource={CATALOG_ROWS}
          columns={[
            { title: "Catalog", dataIndex: "name" },
            {
              title: "Release",
              dataIndex: "release",
              render: (r: string) => <Tag>{r}</Tag>,
            },
            { title: "Sources", dataIndex: "sources" },
            { title: "Coverage", dataIndex: "coverage" },
          ]}
        />
      </Section>
    </PageShell>
  );
}
