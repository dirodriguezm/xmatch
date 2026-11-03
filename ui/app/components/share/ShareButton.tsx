"use client";

import {
  CopyOutlined,
  MailOutlined,
  ShareAltOutlined,
  XOutlined,
} from "@ant-design/icons";
import {
  App,
  Button,
  Flex,
  Input,
  Modal,
  Tabs,
  Tooltip,
  Typography,
} from "antd";
import { useMemo, useState, useSyncExternalStore } from "react";

import { SITE_URL } from "@/app/lib/constants/site";
import {
  absoluteUrl,
  buildEmbedHtml,
  buildShareTargets,
} from "@/app/lib/utils/share";

const { Text } = Typography;

interface ShareButtonProps {
  /** App-relative path including query, e.g. "/object/X?catalog=gaia". */
  path: string;
  title: string;
  /** App-relative embed path, when the view can be embedded. */
  embedPath?: string;
  /** Icon with a tooltip instead of a labelled button, for tight headers. */
  iconOnly?: boolean;
}

const noopSubscribe = () => () => {};

/** Current origin on the client, SITE_URL during SSR. */
function useOrigin(): string {
  return useSyncExternalStore(
    noopSubscribe,
    () => window.location.origin,
    () => SITE_URL
  );
}

function useCanNativeShare(): boolean {
  return useSyncExternalStore(
    noopSubscribe,
    () => typeof navigator.share === "function",
    () => false
  );
}

function BlueskyIcon() {
  return (
    <svg viewBox="0 0 600 530" width="1em" height="1em" aria-hidden>
      <path
        fill="currentColor"
        d="M135.72 44.03C202.216 93.951 273.74 195.17 300 249.49c26.262-54.316 97.782-155.54 164.28-205.46C512.26 8.009 590-19.862 590 68.825c0 17.712-10.155 148.79-16.111 170.07-20.703 73.984-96.144 92.854-163.25 81.433 117.3 19.964 147.14 86.092 82.697 152.22-122.39 125.59-175.91-31.511-189.63-71.766-2.514-7.38-3.69-10.832-3.708-7.896-.017-2.936-1.193.516-3.707 7.896-13.714 40.255-67.233 197.36-189.63 71.766-64.444-66.128-34.605-132.26 82.697-152.22-67.108 11.421-142.55-7.45-163.25-81.433C20.15 217.613 9.997 86.535 9.997 68.825c0-88.687 77.742-60.816 125.72-24.795z"
      />
    </svg>
  );
}

const TARGET_ICONS = {
  x: <XOutlined />,
  bluesky: <BlueskyIcon />,
  email: <MailOutlined />,
};

function CopyField({
  value,
  label,
  multiline = false,
}: {
  value: string;
  label: string;
  multiline?: boolean;
}) {
  const { message } = App.useApp();

  const copy = async () => {
    try {
      await navigator.clipboard.writeText(value);
      message.success(`${label} copied`);
    } catch {
      message.error("Could not access the clipboard");
    }
  };

  return (
    <Flex gap={8} align={multiline ? "start" : "center"}>
      {multiline ? (
        <Input.TextArea
          value={value}
          readOnly
          autoSize={{ minRows: 3, maxRows: 6 }}
          className="font-mono text-xs"
          onFocus={(e) => e.currentTarget.select()}
          aria-label={label}
        />
      ) : (
        <Input
          value={value}
          readOnly
          className="font-mono text-xs"
          onFocus={(e) => e.currentTarget.select()}
          aria-label={label}
        />
      )}
      <Button type="primary" icon={<CopyOutlined />} onClick={copy}>
        Copy
      </Button>
    </Flex>
  );
}

/** "Share" button: copy link, native share, social intents and an embed snippet. */
export function ShareButton({
  path,
  title,
  embedPath,
  iconOnly = false,
}: ShareButtonProps) {
  const [open, setOpen] = useState(false);
  const origin = useOrigin();
  const canNativeShare = useCanNativeShare();

  const url = absoluteUrl(path, origin);
  const targets = useMemo(() => buildShareTargets(url, title), [url, title]);
  const embedUrl = embedPath ? absoluteUrl(embedPath, origin) : null;
  const embedHtml = embedUrl ? buildEmbedHtml(embedUrl, title) : null;

  const nativeShare = async () => {
    try {
      await navigator.share({ title, url });
    } catch {
      // User dismissed the sheet or sharing failed; nothing to do.
    }
  };

  const linkTab = (
    <Flex vertical gap={16}>
      <CopyField value={url} label="Link" />
      <Flex gap={8} wrap>
        {canNativeShare && (
          <Button icon={<ShareAltOutlined />} onClick={nativeShare}>
            Share via…
          </Button>
        )}
        {targets.map((t) => (
          <Button
            key={t.id}
            icon={TARGET_ICONS[t.id]}
            href={t.href}
            target={t.id === "email" ? undefined : "_blank"}
            rel="noopener noreferrer"
          >
            {t.label}
          </Button>
        ))}
      </Flex>
    </Flex>
  );

  return (
    <>
      <Tooltip title={iconOnly ? "Share a link to this view" : undefined}>
        <Button
          size="small"
          icon={<ShareAltOutlined />}
          onClick={() => setOpen(true)}
          aria-label={iconOnly ? "Share" : undefined}
        >
          {iconOnly ? null : "Share"}
        </Button>
      </Tooltip>
      <Modal
        title={`Share ${title}`}
        open={open}
        onCancel={() => setOpen(false)}
        footer={null}
        width={600}
        destroyOnHidden
      >
        {embedHtml && embedUrl ? (
          <Tabs
            items={[
              { key: "link", label: "Link", children: linkTab },
              {
                key: "embed",
                label: "Embed",
                children: (
                  <Flex vertical gap={12}>
                    <Text className="text-neutral-400">
                      Paste this HTML into a blog post, notebook or course page
                      to show a live summary card.
                    </Text>
                    <CopyField value={embedHtml} label="Embed code" multiline />
                    <a
                      href={embedUrl}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="text-xs"
                    >
                      Preview embed ↗
                    </a>
                  </Flex>
                ),
              },
            ]}
          />
        ) : (
          linkTab
        )}
      </Modal>
    </>
  );
}
