"use client";

import { ExportOutlined } from "@ant-design/icons";
import { Button } from "antd";

import { ApiPlayground, DeveloperGuide } from "@/app/components/developers";
import { PageShell } from "@/app/components/layout";
import { SWAGGER_URL } from "@/app/lib/constants/site";

export default function DevelopersPage() {
  return (
    <PageShell
      width="wide"
      title="API playground"
      description="Build a request, copy it as curl, Python or JavaScript, and run it against the live XWave API. No key required."
      actions={
        <>
          <Button href="/llms.txt" target="_blank">
            llms.txt
          </Button>
          <Button
            href={SWAGGER_URL}
            target="_blank"
            rel="noopener noreferrer"
            icon={<ExportOutlined />}
            iconPlacement="end"
          >
            Swagger
          </Button>
        </>
      }
    >
      <ApiPlayground />
      <DeveloperGuide />
      <div className="h-16" />
    </PageShell>
  );
}
