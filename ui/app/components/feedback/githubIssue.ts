/**
 * Pure helpers for turning user feedback into a prefilled GitHub "new issue"
 * URL, and for remembering which matches this browser has already reported.
 * There is no backend: the issue is opened in the user's own GitHub session.
 */

import { REPO_URL } from "@/app/lib/constants/site";

export interface IssueDraft {
  title: string;
  body: string;
  /** Applied only if the reporter can triage the repo; ignored otherwise. */
  labels?: string[];
}

/** GitHub rejects very long URLs; keep free text well below that. */
export const MAX_NOTE_LENGTH = 2000;

export function buildIssueUrl(
  { title, body, labels }: IssueDraft,
  repoUrl: string = REPO_URL
): string {
  const params = new URLSearchParams({ title, body });
  if (labels && labels.length > 0) params.set("labels", labels.join(","));
  return `${repoUrl}/issues/new?${params.toString()}`;
}

export const FLAG_REASONS = [
  {
    value: "wrong-counterpart",
    label: "Wrong counterpart",
    hint: "The matched source is not the same astrophysical object.",
  },
  {
    value: "blended",
    label: "Blended / confused source",
    hint: "Two or more sources merged, or the match picked the wrong neighbour.",
  },
  {
    value: "artefact",
    label: "Artefact / spurious detection",
    hint: "Diffraction spike, ghost, halo or other non-astrophysical source.",
  },
  {
    value: "bad-photometry",
    label: "Bad photometry",
    hint: "Magnitudes or fluxes look wrong for this source.",
  },
  { value: "other", label: "Other", hint: "Anything else worth a look." },
] as const;

export type FlagReason = (typeof FLAG_REASONS)[number]["value"];

export interface FlagReport {
  objectId: string;
  catalog: string;
  ra: number;
  dec: number;
  reason: FlagReason;
  note?: string;
  pageUrl?: string;
}

export function flagReasonLabel(reason: FlagReason): string {
  return FLAG_REASONS.find((r) => r.value === reason)?.label ?? reason;
}

export function buildFlagIssue(report: FlagReport): IssueDraft {
  const reasonLabel = flagReasonLabel(report.reason);
  const note = (report.note ?? "").trim().slice(0, MAX_NOTE_LENGTH);
  const lines = [
    "### Reported match",
    "",
    `- **Object:** \`${report.objectId}\``,
    `- **Catalog:** ${report.catalog}`,
    `- **RA, Dec (deg):** ${report.ra.toFixed(6)}, ${report.dec.toFixed(6)}`,
    `- **Reason:** ${reasonLabel}`,
  ];
  if (report.pageUrl) lines.push(`- **Page:** ${report.pageUrl}`);
  lines.push("", "### Details", "", note || "_No additional details._");
  lines.push("", "---", "_Reported from the XWave object page._");
  return {
    title: `[Match report] ${reasonLabel}: ${report.catalog} ${report.objectId}`,
    body: lines.join("\n"),
    labels: ["match-report"],
  };
}

export const FEEDBACK_KINDS = [
  { value: "bug", label: "Bug report", labels: ["bug"] },
  { value: "feature", label: "Feature request", labels: ["enhancement"] },
  { value: "data", label: "Data or catalog issue", labels: ["data"] },
  { value: "docs", label: "Documentation", labels: ["documentation"] },
] as const;

export type FeedbackKind = (typeof FEEDBACK_KINDS)[number]["value"];

export interface FeedbackInput {
  kind: FeedbackKind;
  title: string;
  description: string;
  steps?: string;
  pageUrl?: string;
  userAgent?: string;
}

export function buildFeedbackIssue(input: FeedbackInput): IssueDraft {
  const kind = FEEDBACK_KINDS.find((k) => k.value === input.kind);
  const lines = [
    "### Description",
    "",
    input.description.trim().slice(0, MAX_NOTE_LENGTH) || "_Not provided._",
  ];
  const steps = input.steps?.trim();
  if (steps) {
    lines.push("", "### Steps to reproduce", "", steps.slice(0, 1000));
  }
  if (input.pageUrl || input.userAgent) {
    lines.push("", "### Context", "");
    if (input.pageUrl) lines.push(`- **Page:** ${input.pageUrl}`);
    if (input.userAgent) lines.push(`- **Browser:** ${input.userAgent}`);
  }
  lines.push("", "---", "_Sent from the XWave contact page._");
  const prefix = kind ? `[${kind.label}] ` : "";
  return {
    title: `${prefix}${input.title.trim()}`,
    body: lines.join("\n"),
    labels: kind ? [...kind.labels] : [],
  };
}

/* ---------- Local memory of reported matches ---------- */

export const FLAG_STORAGE_KEY = "xwave:flagged-matches";

export interface FlagRecord {
  key: string;
  reason: FlagReason;
  at: string;
}

type StorageLike = Pick<Storage, "getItem" | "setItem">;

export function flagKey(catalog: string, objectId: string): string {
  return `${catalog.toLowerCase()}:${objectId}`;
}

export function readFlags(storage: StorageLike | undefined): FlagRecord[] {
  if (!storage) return [];
  try {
    const parsed: unknown = JSON.parse(
      storage.getItem(FLAG_STORAGE_KEY) ?? "[]"
    );
    if (!Array.isArray(parsed)) return [];
    return parsed.filter(
      (r): r is FlagRecord =>
        typeof r === "object" &&
        r !== null &&
        typeof (r as FlagRecord).key === "string"
    );
  } catch {
    return [];
  }
}

export function findFlag(
  storage: StorageLike | undefined,
  catalog: string,
  objectId: string
): FlagRecord | undefined {
  const key = flagKey(catalog, objectId);
  return readFlags(storage).find((r) => r.key === key);
}

export function saveFlag(
  storage: StorageLike | undefined,
  catalog: string,
  objectId: string,
  reason: FlagReason,
  now: Date = new Date()
): FlagRecord {
  const record: FlagRecord = {
    key: flagKey(catalog, objectId),
    reason,
    at: now.toISOString(),
  };
  if (!storage) return record;
  const rest = readFlags(storage).filter((r) => r.key !== record.key);
  try {
    storage.setItem(FLAG_STORAGE_KEY, JSON.stringify([...rest, record]));
  } catch {
    // Storage full or blocked: the report still opened on GitHub.
  }
  return record;
}
