"use client";

import { FlagFilled, FlagOutlined } from "@ant-design/icons";
import {
  App,
  Button,
  Input,
  Modal,
  Radio,
  Space,
  Tooltip,
  Typography,
} from "antd";
import { useCallback, useState, useSyncExternalStore } from "react";

import {
  buildFlagIssue,
  buildIssueUrl,
  findFlag,
  FLAG_REASONS,
  FLAG_STORAGE_KEY,
  type FlagReason,
  flagReasonLabel,
  MAX_NOTE_LENGTH,
  saveFlag,
} from "./githubIssue";

const { Text, Paragraph } = Typography;

/** Fired on the window after this tab records a flag. */
const FLAG_EVENT = "xwave:flagged";

function safeStorage(): Storage | undefined {
  try {
    return typeof window === "undefined" ? undefined : window.localStorage;
  } catch {
    return undefined;
  }
}

function subscribe(onChange: () => void) {
  const onStorage = (e: StorageEvent) => {
    if (e.key === null || e.key === FLAG_STORAGE_KEY) onChange();
  };
  window.addEventListener("storage", onStorage);
  window.addEventListener(FLAG_EVENT, onChange);
  return () => {
    window.removeEventListener("storage", onStorage);
    window.removeEventListener(FLAG_EVENT, onChange);
  };
}

interface FlagMatchButtonProps {
  objectId: string;
  catalog: string;
  ra: number;
  dec: number;
}

/**
 * Report a wrong, blended or spurious match. Opens a prefilled GitHub issue in
 * a new tab and remembers the report in this browser.
 */
export function FlagMatchButton({
  objectId,
  catalog,
  ra,
  dec,
}: FlagMatchButtonProps) {
  const { message } = App.useApp();
  const [open, setOpen] = useState(false);
  const [reason, setReason] = useState<FlagReason>("wrong-counterpart");
  const [note, setNote] = useState("");

  const reportedReason = useSyncExternalStore(
    subscribe,
    () => findFlag(safeStorage(), catalog, objectId)?.reason ?? null,
    () => null
  );

  const submit = useCallback(() => {
    const url = buildIssueUrl(
      buildFlagIssue({
        objectId,
        catalog,
        ra,
        dec,
        reason,
        note,
        pageUrl: window.location.href,
      })
    );
    window.open(url, "_blank", "noopener,noreferrer");
    saveFlag(safeStorage(), catalog, objectId, reason);
    window.dispatchEvent(new Event(FLAG_EVENT));
    setOpen(false);
    setNote("");
    message.success("Thanks — finish submitting the issue on GitHub.");
  }, [catalog, dec, message, note, objectId, ra, reason]);

  const reported = reportedReason !== null;

  return (
    <>
      <Tooltip
        title={
          reported
            ? `You reported this (${flagReasonLabel(reportedReason)}). Click to report again.`
            : "Report a wrong or problematic match"
        }
      >
        <Button
          size="small"
          icon={reported ? <FlagFilled /> : <FlagOutlined />}
          onClick={() => setOpen(true)}
          className={reported ? "!text-amber-400" : undefined}
        >
          {reported ? "Reported" : "Report"}
        </Button>
      </Tooltip>
      <Modal
        open={open}
        title="Report this match"
        okText="Open GitHub issue"
        okButtonProps={{ icon: <FlagOutlined /> }}
        onOk={submit}
        onCancel={() => setOpen(false)}
        destroyOnHidden
      >
        <Paragraph className="text-neutral-400">
          Reports are filed as public GitHub issues so the fix is visible to
          everyone. You&apos;ll review the prefilled issue before posting it.
        </Paragraph>
        <div className="rounded border border-border bg-surface px-3 py-2 mb-4 font-mono text-xs text-neutral-300">
          {catalog} · {objectId} · {ra.toFixed(5)}, {dec.toFixed(5)}
        </div>
        <Text strong className="block mb-2">
          What looks wrong?
        </Text>
        <Radio.Group
          value={reason}
          onChange={(e) => setReason(e.target.value as FlagReason)}
          className="w-full"
        >
          <Space orientation="vertical" size={6} className="w-full">
            {FLAG_REASONS.map((r) => (
              <Radio key={r.value} value={r.value}>
                <span className="text-foreground">{r.label}</span>
                <span className="block text-xs text-neutral-400">{r.hint}</span>
              </Radio>
            ))}
          </Space>
        </Radio.Group>
        <Text strong className="block mt-4 mb-2">
          Details (optional)
        </Text>
        <Input.TextArea
          value={note}
          onChange={(e) => setNote(e.target.value)}
          maxLength={MAX_NOTE_LENGTH}
          showCount
          autoSize={{ minRows: 3, maxRows: 8 }}
          placeholder="e.g. the Gaia source 1.2″ away is the real counterpart; this one is a diffraction spike."
        />
      </Modal>
    </>
  );
}
