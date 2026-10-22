"use client";

import {
  BugOutlined,
  CommentOutlined,
  DashboardOutlined,
  FlagOutlined,
  GithubOutlined,
  SendOutlined,
} from "@ant-design/icons";
import { Button, Checkbox, Form, Input, Radio, Typography } from "antd";
import Link from "next/link";
import type { ReactNode } from "react";

import {
  buildFeedbackIssue,
  buildIssueUrl,
  FEEDBACK_KINDS,
  type FeedbackKind,
  MAX_NOTE_LENGTH,
} from "@/app/components/feedback/githubIssue";
import { PageShell } from "@/app/components/layout";
import { Section } from "@/app/components/learn/Section";
import {
  DISCUSSIONS_URL,
  ISSUES_URL,
  REPO_URL,
} from "@/app/lib/constants/site";

const { Paragraph, Text } = Typography;

interface Channel {
  icon: ReactNode;
  title: string;
  body: string;
  href: string;
  cta: string;
  external?: boolean;
}

const CHANNELS: Channel[] = [
  {
    icon: <BugOutlined />,
    title: "Bugs & errors",
    body: "Something broken, wrong or slow? Open a GitHub issue so it can be tracked and fixed in the open.",
    href: ISSUES_URL,
    cta: "Open issues",
    external: true,
  },
  {
    icon: <CommentOutlined />,
    title: "Questions & ideas",
    body: "Ask how to do something, propose a catalog or share what you built in GitHub Discussions.",
    href: DISCUSSIONS_URL,
    cta: "Join discussions",
    external: true,
  },
  {
    icon: <FlagOutlined />,
    title: "A bad match",
    body: "Use the Report button on any object page; it fills in the object, catalog and coordinates for you.",
    href: "/learn#faq-reporting-bad-matches",
    cta: "How reporting works",
  },
  {
    icon: <DashboardOutlined />,
    title: "Is it down?",
    body: "Check live status of the XWave API and the external services the object page relies on.",
    href: "/status",
    cta: "View status",
  },
];

interface FormValues {
  kind: FeedbackKind;
  title: string;
  description: string;
  steps?: string;
  includeContext: boolean;
}

function ChannelCard({ channel }: { channel: Channel }) {
  const cta = (
    <span className="text-sm">
      {channel.cta} {channel.external ? "↗" : "→"}
    </span>
  );
  return (
    <div className="rounded-lg border border-border bg-surface p-4 flex flex-col gap-2">
      <div className="flex items-center gap-2 text-foreground font-medium">
        <span className="text-primary">{channel.icon}</span>
        {channel.title}
      </div>
      <Text className="!text-neutral-400 text-sm flex-1">{channel.body}</Text>
      {channel.external ? (
        <a href={channel.href} target="_blank" rel="noopener noreferrer">
          {cta}
        </a>
      ) : (
        <Link href={channel.href}>{cta}</Link>
      )}
    </div>
  );
}

export default function ContactPage() {
  const [form] = Form.useForm<FormValues>();
  const kind = Form.useWatch("kind", form);

  const submit = (values: FormValues) => {
    const issue = buildFeedbackIssue({
      kind: values.kind,
      title: values.title,
      description: values.description,
      steps: values.steps,
      userAgent: values.includeContext ? navigator.userAgent : undefined,
    });
    window.open(buildIssueUrl(issue), "_blank", "noopener,noreferrer");
  };

  return (
    <PageShell
      title="Contact"
      description="XWave is developed in the open on GitHub. Pick the channel that fits, or use the form to draft an issue."
      actions={
        <Button
          icon={<GithubOutlined />}
          href={REPO_URL}
          target="_blank"
          rel="noopener noreferrer"
        >
          Repository
        </Button>
      }
    >
      <Section title="Channels">
        <div className="grid gap-3 sm:grid-cols-2">
          {CHANNELS.map((c) => (
            <ChannelCard key={c.title} channel={c} />
          ))}
        </div>
      </Section>

      <Section title="Send feedback" icon={<SendOutlined />}>
        <Paragraph className="!text-neutral-400">
          There is no inbox behind this form: it drafts a GitHub issue and opens
          it in a new tab, where you can review and post it with your GitHub
          account. Issues are public — don&apos;t include anything private.
        </Paragraph>
        <div className="rounded-lg border border-border bg-surface p-5">
          <Form<FormValues>
            form={form}
            layout="vertical"
            requiredMark={false}
            initialValues={{ kind: "bug", includeContext: true }}
            onFinish={submit}
          >
            <Form.Item name="kind" label="What is it about?">
              <Radio.Group
                optionType="button"
                options={FEEDBACK_KINDS.map((k) => ({
                  value: k.value,
                  label: k.label,
                }))}
              />
            </Form.Item>
            <Form.Item
              name="title"
              label="Title"
              rules={[
                {
                  required: true,
                  whitespace: true,
                  message: "Add a short title",
                },
              ]}
            >
              <Input maxLength={120} placeholder="One line summary" />
            </Form.Item>
            <Form.Item
              name="description"
              label="Description"
              rules={[
                {
                  required: true,
                  whitespace: true,
                  message: "Describe what happened or what you'd like",
                },
              ]}
            >
              <Input.TextArea
                autoSize={{ minRows: 4, maxRows: 12 }}
                maxLength={MAX_NOTE_LENGTH}
                showCount
                placeholder={
                  kind === "feature"
                    ? "What would you like XWave to do, and why?"
                    : "What happened, and what did you expect?"
                }
              />
            </Form.Item>
            {kind === "bug" && (
              <Form.Item name="steps" label="Steps to reproduce (optional)">
                <Input.TextArea
                  autoSize={{ minRows: 2, maxRows: 8 }}
                  maxLength={1000}
                  placeholder={"1. Search for M31 with a 5″ radius\n2. …"}
                />
              </Form.Item>
            )}
            <Form.Item name="includeContext" valuePropName="checked">
              <Checkbox>Include my browser version</Checkbox>
            </Form.Item>
            <Button type="primary" htmlType="submit" icon={<GithubOutlined />}>
              Draft issue on GitHub
            </Button>
          </Form>
        </div>
      </Section>
    </PageShell>
  );
}
