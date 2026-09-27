"use client";

import {
  AimOutlined,
  ArrowRightOutlined,
  ExportOutlined,
  LinkOutlined,
  LoadingOutlined,
  QuestionCircleOutlined,
  SearchOutlined,
} from "@ant-design/icons";
import { App, Button, Modal, Typography } from "antd";
import { useRouter } from "next/navigation";
import {
  type KeyboardEvent as ReactKeyboardEvent,
  type ReactNode,
  useCallback,
  useEffect,
  useMemo,
  useRef,
  useState,
  useSyncExternalStore,
} from "react";

import { parseCoordinates, resolveObjectName } from "@/app/lib/api/sesame";
import { FOOTER_GROUPS, NAV_ITEMS } from "@/app/lib/constants/site";
import {
  buildNavCommands,
  buildSearchHref,
  cycleIndex,
  filterCommands,
  GO_SEQUENCE_TIMEOUT_MS,
  GO_SHORTCUTS,
  isMacPlatform,
  isTypingTarget,
  type PaletteCommand,
  paletteShortcutLabel,
} from "@/app/lib/utils/commands";

const { Text } = Typography;

/** Window event the header trigger dispatches to open the palette. */
const OPEN_EVENT = "xwave:command-palette";

/**
 * Page search inputs "/" should focus, most specific first. Pages can opt in
 * with `data-search-input` on the input (or a wrapper).
 */
const SEARCH_INPUT_SELECTORS = [
  "input[data-search-input]",
  "[data-search-input] input",
  "main .ant-input-search input",
  ".ant-input-search input",
];

function subscribeNoop() {
  return () => {};
}

/** Platform check that renders "Ctrl K" on the server and fixes up on the client. */
function useIsMac(): boolean {
  return useSyncExternalStore(
    subscribeNoop,
    () => isMacPlatform(navigator.platform || navigator.userAgent),
    () => false
  );
}

function findPageSearchInput(): HTMLInputElement | null {
  for (const selector of SEARCH_INPUT_SELECTORS) {
    const candidates = document.querySelectorAll<HTMLInputElement>(selector);
    for (const el of candidates) {
      if (!el.disabled && el.offsetParent !== null) return el;
    }
  }
  return null;
}

function anotherDialogOpen(): boolean {
  return Array.from(
    document.querySelectorAll<HTMLElement>(".ant-modal-wrap, .ant-drawer-open")
  ).some((el) => getComputedStyle(el).display !== "none");
}

function Kbd({ children }: { children: ReactNode }) {
  return (
    <kbd className="inline-flex min-w-5 items-center justify-center rounded border border-border bg-surface px-1.5 font-mono text-[11px] leading-5 text-neutral-400">
      {children}
    </kbd>
  );
}

/** Header button that opens the palette. */
export function CommandPaletteTrigger() {
  const isMac = useIsMac();
  return (
    <Button
      icon={<SearchOutlined />}
      onClick={() => window.dispatchEvent(new CustomEvent(OPEN_EVENT))}
      aria-label="Open command palette"
      aria-keyshortcuts={isMac ? "Meta+K" : "Control+K"}
      className="text-neutral-400"
    >
      <span className="hidden sm:inline text-neutral-400">Search</span>
      <span className="hidden sm:inline">
        <Kbd>{paletteShortcutLabel(isMac)}</Kbd>
      </span>
    </Button>
  );
}

const SEARCH_ID = "search:object";
const COPY_ID = "action:copy-link";
const HELP_ID = "action:shortcuts";

const ACTION_COMMANDS: PaletteCommand[] = [
  {
    id: COPY_ID,
    label: "Copy page link",
    group: "Actions",
    keywords: ["share", "url", "clipboard"],
  },
  {
    id: HELP_ID,
    label: "Keyboard shortcuts",
    group: "Actions",
    keywords: ["help", "keys", "hotkeys"],
    hint: "?",
  },
];

function commandIcon(cmd: PaletteCommand) {
  if (cmd.id === SEARCH_ID) return <AimOutlined />;
  if (cmd.id === COPY_ID) return <LinkOutlined />;
  if (cmd.id === HELP_ID) return <QuestionCircleOutlined />;
  if (cmd.external) return <ExportOutlined />;
  return <ArrowRightOutlined />;
}

/** Global ⌘K palette and keyboard shortcuts, mounted once in Providers. */
export function CommandPalette() {
  const router = useRouter();
  const { message } = App.useApp();
  const isMac = useIsMac();

  const [open, setOpen] = useState(false);
  /** Opened from the keyboard: skip the zoom animation so typing starts at once. */
  const [instant, setInstant] = useState(false);
  const [helpOpen, setHelpOpen] = useState(false);
  const [query, setQuery] = useState("");
  const [active, setActive] = useState(0);
  const [resolving, setResolving] = useState(false);

  const inputRef = useRef<HTMLInputElement>(null);
  const listRef = useRef<HTMLDivElement>(null);
  const goPendingAt = useRef(0);

  const navCommands = useMemo(
    () => buildNavCommands(NAV_ITEMS, FOOTER_GROUPS),
    []
  );

  const results = useMemo(() => {
    const trimmed = query.trim();
    const matched = filterCommands([...navCommands, ...ACTION_COMMANDS], query);
    if (!trimmed) return matched;
    const coords = parseCoordinates(trimmed);
    const search: PaletteCommand = {
      id: SEARCH_ID,
      label: coords
        ? `Search position: ${coords.ra.toFixed(5)}, ${coords.dec.toFixed(5)}`
        : `Search object: ${trimmed}`,
      group: "Search",
      hint: coords ? "RA Dec" : "Sesame",
    };
    // Coordinates are unambiguous; names only lead when nothing else matched well.
    return coords || matched.length === 0
      ? [search, ...matched]
      : [...matched.slice(0, 1), search, ...matched.slice(1)];
  }, [navCommands, query]);

  const activeIndex = Math.min(active, Math.max(results.length - 1, 0));

  const openPalette = useCallback((viaKeyboard: boolean) => {
    setInstant(viaKeyboard);
    setQuery("");
    setActive(0);
    setOpen(true);
  }, []);

  const close = useCallback(() => setOpen(false), []);

  const runSearch = useCallback(
    async (raw: string) => {
      const text = raw.trim();
      const coords = parseCoordinates(text);
      if (coords) {
        close();
        router.push(buildSearchHref(coords.ra, coords.dec));
        return;
      }
      setResolving(true);
      try {
        const resolved = await resolveObjectName(text);
        if (!resolved) {
          message.error(`Could not resolve "${text}" to coordinates`);
          return;
        }
        close();
        router.push(buildSearchHref(resolved.ra, resolved.dec));
      } finally {
        setResolving(false);
      }
    },
    [close, message, router]
  );

  const execute = useCallback(
    async (cmd: PaletteCommand | undefined) => {
      if (!cmd || resolving) return;
      if (cmd.id === SEARCH_ID) {
        await runSearch(query);
        return;
      }
      if (cmd.id === COPY_ID) {
        close();
        try {
          await navigator.clipboard.writeText(window.location.href);
          message.success("Page link copied");
        } catch {
          message.error("Could not access the clipboard");
        }
        return;
      }
      if (cmd.id === HELP_ID) {
        close();
        setHelpOpen(true);
        return;
      }
      if (cmd.href) {
        close();
        if (cmd.external) {
          window.open(cmd.href, "_blank", "noopener,noreferrer");
        } else if (cmd.href.endsWith(".txt")) {
          // Route handlers, not pages: do a full navigation.
          window.location.assign(cmd.href);
        } else {
          router.push(cmd.href);
        }
      }
    },
    [close, message, query, resolving, router, runSearch]
  );

  // Header trigger.
  useEffect(() => {
    const onOpen = () => openPalette(false);
    window.addEventListener(OPEN_EVENT, onOpen);
    return () => window.removeEventListener(OPEN_EVENT, onOpen);
  }, [openPalette]);

  // Global shortcuts.
  useEffect(() => {
    const onKeyDown = (e: KeyboardEvent) => {
      if (
        (e.metaKey || e.ctrlKey) &&
        !e.altKey &&
        e.key.toLowerCase() === "k"
      ) {
        e.preventDefault();
        if (open) setOpen(false);
        else openPalette(true);
        return;
      }
      if (open || helpOpen || e.defaultPrevented) return;
      if (e.metaKey || e.ctrlKey || e.altKey) return;
      if (isTypingTarget(e.target as HTMLElement | null)) return;
      // Another dialog (drawer, modal) owns the keyboard.
      if (anotherDialogOpen()) return;

      const key = e.key;
      const now = Date.now();
      if (
        goPendingAt.current &&
        now - goPendingAt.current < GO_SEQUENCE_TIMEOUT_MS
      ) {
        goPendingAt.current = 0;
        const target = GO_SHORTCUTS[key.toLowerCase()];
        if (target) {
          e.preventDefault();
          router.push(target.href);
          return;
        }
      }

      if (key === "/") {
        e.preventDefault();
        const input = findPageSearchInput();
        if (input) {
          input.focus();
          input.select();
        } else {
          openPalette(true);
        }
      } else if (key === "?") {
        e.preventDefault();
        setHelpOpen(true);
      } else if (key === "g" || key === "G") {
        goPendingAt.current = now;
      }
    };
    window.addEventListener("keydown", onKeyDown);
    return () => window.removeEventListener("keydown", onKeyDown);
  }, [helpOpen, open, openPalette, router]);

  // Keep the active row visible.
  useEffect(() => {
    listRef.current
      ?.querySelector<HTMLElement>(`[data-index="${activeIndex}"]`)
      ?.scrollIntoView({ block: "nearest" });
  }, [activeIndex]);

  const onInputKeyDown = (e: ReactKeyboardEvent<HTMLInputElement>) => {
    if (e.key === "ArrowDown" || (e.ctrlKey && e.key === "n")) {
      e.preventDefault();
      setActive(cycleIndex(activeIndex, 1, results.length));
    } else if (e.key === "ArrowUp" || (e.ctrlKey && e.key === "p")) {
      e.preventDefault();
      setActive(cycleIndex(activeIndex, -1, results.length));
    } else if (e.key === "Home") {
      e.preventDefault();
      setActive(0);
    } else if (e.key === "End") {
      e.preventDefault();
      setActive(Math.max(results.length - 1, 0));
    } else if (e.key === "Enter") {
      e.preventDefault();
      void execute(results[activeIndex]);
    }
  };

  let lastGroup: string | null = null;
  const activeId = results[activeIndex]
    ? `cmd-${results[activeIndex].id}`
    : undefined;

  return (
    <>
      <Modal
        open={open}
        onCancel={close}
        footer={null}
        closable={false}
        width={560}
        transitionName={instant ? "" : undefined}
        maskTransitionName={instant ? "" : undefined}
        destroyOnHidden
        afterOpenChange={(isOpen) => {
          if (isOpen) inputRef.current?.focus();
        }}
        styles={{ body: { padding: 0 } }}
        aria-label="Command palette"
      >
        <div className="flex items-center gap-3 border-b border-border px-4 py-3">
          {resolving ? (
            <LoadingOutlined className="text-neutral-400" />
          ) : (
            <SearchOutlined className="text-neutral-400" />
          )}
          <input
            ref={inputRef}
            autoFocus
            value={query}
            onChange={(e) => {
              setQuery(e.target.value);
              setActive(0);
            }}
            onKeyDown={onInputKeyDown}
            placeholder="Jump to a page, or search an object (M31, 10.68 41.27)…"
            className="min-w-0 flex-1 border-0 bg-transparent text-base text-foreground outline-none placeholder:text-neutral-500"
            role="combobox"
            aria-expanded
            aria-controls="command-palette-list"
            aria-activedescendant={activeId}
            aria-autocomplete="list"
            spellCheck={false}
          />
          <Kbd>Esc</Kbd>
        </div>

        <div
          ref={listRef}
          id="command-palette-list"
          role="listbox"
          className="max-h-[50vh] overflow-y-auto p-2"
        >
          {results.length === 0 && (
            <div className="px-3 py-6 text-center text-neutral-400">
              No matching commands
            </div>
          )}
          {results.map((cmd, i) => {
            const header = cmd.group !== lastGroup ? cmd.group : null;
            lastGroup = cmd.group;
            const selected = i === activeIndex;
            return (
              <div key={cmd.id}>
                {header && (
                  <div className="px-3 pb-1 pt-2 text-[11px] font-medium uppercase tracking-wide text-neutral-500">
                    {header}
                  </div>
                )}
                <div
                  id={`cmd-${cmd.id}`}
                  role="option"
                  aria-selected={selected}
                  data-index={i}
                  onMouseMove={() => active !== i && setActive(i)}
                  onClick={() => void execute(cmd)}
                  className={`flex cursor-pointer items-center gap-3 rounded-md px-3 py-2 ${
                    selected
                      ? "bg-surface-elevated text-foreground"
                      : "text-neutral-300"
                  }`}
                >
                  <span className="text-neutral-400">{commandIcon(cmd)}</span>
                  <span className="min-w-0 flex-1 truncate">{cmd.label}</span>
                  {cmd.hint && (
                    <span className="shrink-0 font-mono text-xs text-neutral-500">
                      {cmd.hint}
                    </span>
                  )}
                </div>
              </div>
            );
          })}
        </div>

        <div className="flex flex-wrap items-center gap-x-4 gap-y-1 border-t border-border px-4 py-2 text-xs text-neutral-500">
          <span className="flex items-center gap-1">
            <Kbd>↑</Kbd>
            <Kbd>↓</Kbd> navigate
          </span>
          <span className="flex items-center gap-1">
            <Kbd>↵</Kbd> select
          </span>
          <span className="flex items-center gap-1">
            <Kbd>?</Kbd> all shortcuts
          </span>
        </div>
      </Modal>

      <ShortcutsHelp
        open={helpOpen}
        onClose={() => setHelpOpen(false)}
        paletteKey={paletteShortcutLabel(isMac)}
      />
    </>
  );
}

function ShortcutsHelp({
  open,
  onClose,
  paletteKey,
}: {
  open: boolean;
  onClose: () => void;
  paletteKey: string;
}) {
  const rows: { keys: string[]; label: string }[] = [
    { keys: [paletteKey], label: "Open / close the command palette" },
    { keys: ["/"], label: "Focus the page search box" },
    { keys: ["?"], label: "Show this help" },
    ...Object.entries(GO_SHORTCUTS).map(([k, v]) => ({
      keys: ["G", k.toUpperCase()],
      label: `Go to ${v.label}`,
    })),
    { keys: ["Esc"], label: "Close dialogs" },
  ];
  return (
    <Modal
      open={open}
      onCancel={onClose}
      footer={null}
      title="Keyboard shortcuts"
      width={440}
      transitionName=""
      maskTransitionName=""
    >
      <div className="flex flex-col divide-y divide-border">
        {rows.map((row) => (
          <div
            key={row.label}
            className="flex items-center justify-between gap-4 py-2"
          >
            <Text className="text-neutral-300">{row.label}</Text>
            <span className="flex items-center gap-1">
              {row.keys.map((k, i) => (
                <span key={k} className="flex items-center gap-1">
                  {i > 0 && (
                    <span className="text-xs text-neutral-500">then</span>
                  )}
                  <Kbd>{k}</Kbd>
                </span>
              ))}
            </span>
          </div>
        ))}
      </div>
      <Text className="!mt-3 block text-xs text-neutral-500">
        Shortcuts are ignored while typing in a text field.
      </Text>
    </Modal>
  );
}
