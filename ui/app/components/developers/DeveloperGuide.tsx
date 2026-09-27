"use client";

import {
  ApiOutlined,
  DashboardOutlined,
  ExportOutlined,
  LockOutlined,
  PythonOutlined,
  RobotOutlined,
  ToolOutlined,
} from "@ant-design/icons";
import { Typography } from "antd";
import type { ReactNode } from "react";

import { CodeBlock } from "@/app/components/code";
import { API_BASE_URL } from "@/app/lib/api/client";
import { MAX_RADIUS_ARCSEC } from "@/app/lib/constants/search";
import { REPO_URL, SWAGGER_URL } from "@/app/lib/constants/site";

const { Title, Paragraph, Text } = Typography;

function Section({
  id,
  icon,
  title,
  children,
}: {
  id: string;
  icon: ReactNode;
  title: string;
  children: ReactNode;
}) {
  return (
    <section id={id} className="scroll-mt-20">
      <Title level={4} className="!mb-3 flex items-center gap-2">
        <span className="text-neutral-400">{icon}</span>
        <a href={`#${id}`} className="!text-foreground no-underline">
          {title}
        </a>
      </Title>
      <div className="flex flex-col gap-3 text-neutral-300">{children}</div>
    </section>
  );
}

function Code({ children }: { children: ReactNode }) {
  return (
    <code className="rounded bg-surface-elevated px-1.5 py-0.5 font-mono text-[0.85em] text-foreground">
      {children}
    </code>
  );
}

const PYTHON_QUICKSTART = `import requests
import pandas as pd

BASE = "${API_BASE_URL}"

def conesearch(ra, dec, radius=5, catalog="all", nneighbor=10):
    """Sources within \`radius\` arcsec of (ra, dec), as a DataFrame."""
    r = requests.get(f"{BASE}/conesearch", timeout=60, params=dict(
        ra=ra, dec=dec, radius=radius, catalog=catalog, nneighbor=nneighbor))
    r.raise_for_status()
    groups = r.json() if r.status_code == 200 else []  # 204 = no match
    return pd.DataFrame([row for g in groups for row in g["data"]])

df = conesearch(10.6847, 41.269, radius=10)
print(df.sort_values("distance").head())`;

const PANDAS_ONELINER = `pd.json_normalize(requests.get("${API_BASE_URL}/conesearch", params={"ra": 10.6847, "dec": 41.269, "radius": 10, "nneighbor": 5}).json(), "data")`;

const ASTROQUERY_STYLE = `import astropy.units as u
from astropy.coordinates import SkyCoord
from astropy.table import Table
import requests

class XWave:
    """Minimal astroquery-style wrapper around the XWave REST API."""
    URL = "${API_BASE_URL}"

    @classmethod
    def query_region(cls, coord, radius=5 * u.arcsec, catalog="all", nneighbor=100):
        c = coord.icrs
        r = requests.get(f"{cls.URL}/conesearch", timeout=60, params={
            "ra": c.ra.deg, "dec": c.dec.deg,
            "radius": radius.to_value(u.arcsec),
            "catalog": catalog, "nneighbor": nneighbor,
        })
        r.raise_for_status()
        rows = [row for g in (r.json() if r.status_code == 200 else []) for row in g["data"]]
        return Table(rows=rows) if rows else Table()

    @classmethod
    def query_bulk(cls, coords, radius=2 * u.arcsec, catalog="all"):
        c = coords.icrs
        r = requests.post(f"{cls.URL}/bulk-conesearch", timeout=300, json={
            "ra": list(c.ra.deg), "dec": list(c.dec.deg),
            "radius": radius.to_value(u.arcsec), "catalog": catalog,
        })
        r.raise_for_status()
        rows = [{**row, "input_index": g["index"]}
                for g in (r.json() if r.status_code == 200 else []) for row in g["data"]]
        return Table(rows=rows) if rows else Table()

t = XWave.query_region(SkyCoord.from_name("M31"), radius=10 * u.arcsec, catalog="gaia")
t.pprint()`;

/** Reference sections under the playground. */
export function DeveloperGuide() {
  return (
    <div className="mt-14 grid gap-12 lg:grid-cols-2">
      <Section id="base-url" icon={<LockOutlined />} title="Base URL & auth">
        <CodeBlock code={API_BASE_URL} copyLabel="Base URL" />
        <Paragraph className="!mb-0 text-neutral-300">
          No authentication: no API key, no sign-up. Every endpoint returns
          JSON. Coordinates are decimal degrees (ICRS/J2000) and{" "}
          <strong>all radii are in arcseconds</strong>, bulk included. An empty
          match returns <Code>204 No Content</Code> with no body — treat it as
          an empty list, not an error.
        </Paragraph>
        <Paragraph className="!mb-0 text-neutral-300">
          Browsers can call the API directly only from XWave&apos;s own origins;
          from anywhere else, call it server-side, from a notebook or script.
        </Paragraph>
      </Section>

      <Section
        id="rate-limits"
        icon={<DashboardOutlined />}
        title="Rate limits & etiquette"
      >
        <ul className="m-0 flex list-disc flex-col gap-1.5 pl-5">
          <li>
            There is no hard rate limit — it is a shared academic service, so
            please be gentle: at most ~4 requests in parallel.
          </li>
          <li>
            Many positions? Use one <Code>POST /bulk-conesearch</Code> instead
            of a loop of cone searches.
          </li>
          <li>
            Keep radii ≤ {MAX_RADIUS_ARCSEC}″. Larger cones can hang without
            ever replying.
          </li>
          <li>
            <Code>nneighbor</Code> is a cap, not a cost — latency is flat up to
            hundreds, so ask for what you need.
          </li>
          <li>Retry 5xx with exponential backoff; don&apos;t retry 4xx.</li>
        </ul>
      </Section>

      <Section id="python" icon={<PythonOutlined />} title="Python quickstart">
        <Text className="text-neutral-400">
          <Code>pip install requests pandas</Code>
        </Text>
        <CodeBlock code={PYTHON_QUICKSTART} copyLabel="Python quickstart" />
        <Text className="text-neutral-400">Or, as a pandas one-liner:</Text>
        <CodeBlock code={PANDAS_ONELINER} copyLabel="One-liner" wrap />
      </Section>

      <Section
        id="astroquery"
        icon={<ToolOutlined />}
        title="astroquery-style usage"
      >
        <Paragraph className="!mb-0 text-neutral-300">
          XWave isn&apos;t an astroquery module (yet), but a few lines give you
          the familiar <Code>query_region(SkyCoord, radius)</Code> interface
          returning an astropy <Code>Table</Code>, with units handled for you.
        </Paragraph>
        <CodeBlock
          code={ASTROQUERY_STYLE}
          copyLabel="astroquery-style wrapper"
          size="lg"
        />
      </Section>

      <Section id="llms" icon={<RobotOutlined />} title="For LLMs & agents">
        <ul className="m-0 flex list-disc flex-col gap-1.5 pl-5">
          <li>
            <a href="/llms.txt">/llms.txt</a> — a concise, machine-readable
            index of the API following{" "}
            <a
              href="https://llmstxt.org"
              target="_blank"
              rel="noopener noreferrer"
            >
              llmstxt.org
            </a>
            .
          </li>
          <li>
            <a href="/llms-full.txt">/llms-full.txt</a> — the full reference
            (every parameter, response schema and error) in one plain-text file
            to paste into a prompt.
          </li>
        </ul>
      </Section>

      <Section id="reference" icon={<ApiOutlined />} title="Full reference">
        <Paragraph className="!mb-0 text-neutral-300">
          The OpenAPI (Swagger) reference documents every endpoint, including{" "}
          <Code>POST /bulk-metadata</Code>, with request and response schemas
          you can try in the browser.
        </Paragraph>
        <div className="flex flex-wrap gap-4">
          <a href={SWAGGER_URL} target="_blank" rel="noopener noreferrer">
            Swagger reference <ExportOutlined />
          </a>
          <a href={REPO_URL} target="_blank" rel="noopener noreferrer">
            Source on GitHub <ExportOutlined />
          </a>
        </div>
      </Section>
    </div>
  );
}
