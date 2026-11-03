"use client";

import { LeftOutlined, RightOutlined } from "@ant-design/icons";
import { Button, Tooltip } from "antd";
import { useRouter } from "next/navigation";
import { useEffect, useMemo, useState } from "react";

import { getSearchCatalogLabel } from "@/app/lib/constants/catalogs";
import {
  type BulkNavItem,
  loadBulkNav,
  locateInBulkNav,
} from "@/app/lib/utils/bulkSession";
import { buildObjectUrl } from "@/app/lib/utils/urls";

interface ResultNavigatorProps {
  objectId: string;
  catalog: string;
}

/** Typing in a field must not flip to another object. */
function isEditable(target: EventTarget | null): boolean {
  if (!(target instanceof HTMLElement)) return false;
  return (
    target.isContentEditable ||
    ["INPUT", "TEXTAREA", "SELECT"].includes(target.tagName)
  );
}

/**
 * ‹ 12 / 340 › through the matches of the last bulk run, in the order the
 * results table showed them. Arrow keys and j/k step too. Navigates with
 * `replace` so "Back to results" is always one step away.
 */
export function ResultNavigator({ objectId, catalog }: ResultNavigatorProps) {
  const router = useRouter();
  const [items, setItems] = useState<BulkNavItem[]>([]);

  // sessionStorage only exists in the browser; read it after mount.
  useEffect(() => {
    const timer = window.setTimeout(() => setItems(loadBulkNav()), 0);
    return () => window.clearTimeout(timer);
  }, []);

  const at = useMemo(
    () => locateInBulkNav(items, objectId, catalog),
    [items, objectId, catalog]
  );

  useEffect(() => {
    if (!at) return;
    const onKey = (e: KeyboardEvent) => {
      if (e.metaKey || e.ctrlKey || e.altKey || isEditable(e.target)) return;
      const step =
        e.key === "ArrowLeft" || e.key === "k"
          ? at.prev
          : e.key === "ArrowRight" || e.key === "j"
            ? at.next
            : null;
      if (!step) return;
      e.preventDefault();
      router.replace(buildObjectUrl(step.id, step.catalog, "bulk"));
    };
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  }, [at, router]);

  if (!at) return null;

  const go = (item: BulkNavItem | null) =>
    item && router.replace(buildObjectUrl(item.id, item.catalog, "bulk"));

  return (
    <div className="flex items-center gap-1 min-w-0">
      <Tooltip title="Previous match (← or k)">
        <Button
          type="text"
          size="small"
          icon={<LeftOutlined />}
          disabled={!at.prev}
          onClick={() => go(at.prev)}
          aria-label="Previous match"
        />
      </Tooltip>
      <span className="text-sm tabular-nums text-foreground whitespace-nowrap">
        {at.index + 1} / {items.length}
      </span>
      <Tooltip title="Next match (→ or j)">
        <Button
          type="text"
          size="small"
          icon={<RightOutlined />}
          disabled={!at.next}
          onClick={() => go(at.next)}
          aria-label="Next match"
        />
      </Tooltip>
      <span className="hidden xl:inline text-xs text-neutral-400 truncate">
        input #{at.item.inputIndex + 1}
        {at.item.inputName ? ` · ${at.item.inputName}` : ""} ·{" "}
        {getSearchCatalogLabel(at.item.catalog)}
      </span>
    </div>
  );
}
