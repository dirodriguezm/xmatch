"use client";

import {
  BookOutlined,
  NodeIndexOutlined,
  QuestionCircleOutlined,
} from "@ant-design/icons";
import { Tabs } from "antd";
import { useState, useSyncExternalStore } from "react";

import { PageShell } from "@/app/components/layout";
import { FaqList } from "@/app/components/learn/FaqList";
import { GlossarySection } from "@/app/components/learn/GlossarySection";
import { MethodsSection } from "@/app/components/learn/MethodsSection";
import { Section } from "@/app/components/learn/Section";

type TabKey = "faq" | "methods" | "glossary";

function subscribeHash(onChange: () => void) {
  window.addEventListener("hashchange", onChange);
  return () => window.removeEventListener("hashchange", onChange);
}

function useHash(): string {
  return useSyncExternalStore(
    subscribeHash,
    () => decodeURIComponent(window.location.hash.slice(1)),
    () => ""
  );
}

function tabFromHash(hash: string): TabKey {
  if (hash.startsWith("methods")) return "methods";
  if (hash.startsWith("glossary")) return "glossary";
  return "faq";
}

function setHash(hash: string) {
  window.history.replaceState(null, "", `#${hash}`);
  window.dispatchEvent(new HashChangeEvent("hashchange"));
}

export default function LearnPage() {
  const hash = useHash();
  const tab = tabFromHash(hash);
  const hashFaq = hash.startsWith("faq-") ? hash.slice(4) : null;

  const [open, setOpen] = useState<string[]>([]);
  const [dismissed, setDismissed] = useState<string | null>(null);
  const openIds =
    hashFaq && hashFaq !== dismissed && !open.includes(hashFaq)
      ? [...open, hashFaq]
      : open;

  return (
    <PageShell
      title="Learn"
      description="Answers to common questions, exactly how the cross-match works, and what every catalog column means."
    >
      <Tabs
        activeKey={tab}
        onChange={(key) => setHash(key)}
        size="large"
        className="mb-12"
        items={[
          {
            key: "faq",
            label: "FAQ",
            icon: <QuestionCircleOutlined />,
            children: (
              <Section id="faq" title="Frequently asked questions">
                <FaqList
                  openIds={openIds}
                  onChange={(keys) => {
                    if (hashFaq && !keys.includes(hashFaq)) {
                      setDismissed(hashFaq);
                    }
                    setOpen(keys);
                  }}
                />
              </Section>
            ),
          },
          {
            key: "methods",
            label: "How matching works",
            icon: <NodeIndexOutlined />,
            children: (
              <Section id="methods" title="How matching works">
                <MethodsSection />
              </Section>
            ),
          },
          {
            key: "glossary",
            label: "Glossary",
            icon: <BookOutlined />,
            children: (
              <Section id="glossary" title="Catalog column glossary">
                <GlossarySection />
              </Section>
            ),
          },
        ]}
      />
    </PageShell>
  );
}
