"use client";

import {
  ApiOutlined,
  BulbOutlined,
  DatabaseOutlined,
  GithubOutlined,
  NodeIndexOutlined,
} from "@ant-design/icons";
import { Card, Layout, Tag, Typography } from "antd";
import type { ReactNode } from "react";

import { AppHeader } from "@/app/components/layout";
import { API_BASE_URL } from "@/app/lib/api/client";
import {
  CATALOG_COLOR_CLASSES,
  CATALOG_LABELS,
  CATALOG_OPTIONS,
} from "@/app/lib/constants/catalogs";
import { MAX_RADIUS_ARCSEC } from "@/app/lib/constants/search";

const { Content } = Layout;
const { Title, Paragraph, Text, Link } = Typography;

const SWAGGER_URL = "https://xwave-astro.udp.cl/swagger/index.html";
const REPO_URL = "https://github.com/dirodriguezm/xmatch";

interface StepProps {
  index: number;
  title: string;
  children: ReactNode;
}

function Step({ index, title, children }: StepProps) {
  return (
    <div className="flex gap-4">
      <div className="shrink-0 w-8 h-8 rounded-full bg-surface-elevated border border-border flex items-center justify-center font-mono text-sm text-foreground">
        {index}
      </div>
      <div className="min-w-0 flex-1 pb-6">
        <Title level={5} className="!mt-0 !mb-2 text-foreground">
          {title}
        </Title>
        <div className="text-muted">{children}</div>
      </div>
    </div>
  );
}

function Section({
  title,
  icon,
  children,
}: {
  title: string;
  icon?: ReactNode;
  children: ReactNode;
}) {
  return (
    <section className="mb-10">
      <Title
        level={3}
        className="!mb-4 text-foreground flex items-center gap-2"
      >
        {icon}
        {title}
      </Title>
      {children}
    </section>
  );
}

export default function AboutPage() {
  return (
    <Layout className="min-h-screen">
      <AppHeader />
      <Content className="bg-background">
        <div className="max-w-3xl mx-auto px-6 py-12">
          <Title level={2} className="!mb-2 text-foreground">
            About XWave
          </Title>
          <Paragraph className="text-muted !text-base">
            XWave is a positional cross-match service for astronomical catalogs.
            Give it a coordinate and a search radius, and it returns the objects
            each indexed catalog has at that position, along with their angular
            separation from your target.
          </Paragraph>

          <Section
            title="How the cross-match works"
            icon={<NodeIndexOutlined />}
          >
            <Paragraph className="text-muted">
              Comparing a target against every row of a catalog does not scale,
              so the search runs in two stages: a cheap spatial pre-filter that
              narrows millions of objects down to a handful of candidates, then
              an exact distance computation on just those.
            </Paragraph>

            <div className="mt-6">
              <Step index={1} title="Catalogs are indexed onto a HEALPix grid">
                <Paragraph className="text-muted !mb-0">
                  When a catalog is ingested, every object&apos;s sky position
                  is mapped to a{" "}
                  <Link
                    href="https://healpix.sourceforge.io/"
                    target="_blank"
                    rel="noopener noreferrer"
                  >
                    HEALPix
                  </Link>{" "}
                  pixel using the <Text code>NESTED</Text> ordering scheme, and
                  that pixel index is stored alongside the object. HEALPix
                  divides the sphere into equal-area pixels, so &quot;which
                  objects are near this point&quot; becomes an indexed lookup
                  rather than a scan. That index is the <Text code>IPix</Text>{" "}
                  column you see in the results table.
                </Paragraph>
              </Step>

              <Step
                index={2}
                title="The search disc is resolved to a set of pixels"
              >
                <Paragraph className="text-muted !mb-0">
                  Your coordinate and radius define a disc on the sphere. The
                  service asks HEALPix for every pixel that the disc touches —
                  inclusively, so pixels that only partially overlap are kept —
                  and fetches the objects stored in those pixels. This step is
                  deliberately generous: it may return objects outside the
                  radius, but it will not miss any inside it.
                </Paragraph>
              </Step>

              <Step index={3} title="Candidates are ranked by proximity">
                <Paragraph className="text-muted !mb-0">
                  The candidates from step 2 are loaded into a 2-D k-d tree
                  keyed on right ascension and declination, and a
                  nearest-neighbour query returns the closest ones to your
                  target. This bounds how much work the final step has to do.
                </Paragraph>
              </Step>

              <Step
                index={4}
                title="Exact separations are computed and filtered"
              >
                <Paragraph className="text-muted !mb-0">
                  For each remaining candidate the service computes the
                  great-circle separation from your target using the haversine
                  formula, and discards anything beyond your radius. What
                  survives is the result set, and that separation is the{" "}
                  <Text code>Ang. Dist</Text> column — reported in arcseconds.
                </Paragraph>
              </Step>
            </div>

            <Card size="small" className="bg-surface border-border mt-2">
              <Text strong className="text-foreground">
                Why results are capped
              </Text>
              <Paragraph className="text-muted !mb-0 !mt-2">
                The proximity ranking in step 3 happens <em>before</em> the
                radius filter in step 4, so the neighbour limit bounds how many
                objects can come back at all. This interface requests a generous
                limit so that in practice the radius is what determines your
                results — but for a very crowded field, a large radius can still
                hit that ceiling.
              </Paragraph>
            </Card>
          </Section>

          <Section title="Catalogs" icon={<DatabaseOutlined />}>
            <Paragraph className="text-muted">
              Each catalog is indexed independently and searched with its own
              radius, since the positional uncertainty that makes sense for an
              infrared source is not the one that makes sense for an X-ray
              detection.
            </Paragraph>
            <div className="flex flex-wrap gap-2 mb-4">
              {CATALOG_OPTIONS.map((catalog) => (
                <span
                  key={catalog}
                  className="inline-flex items-center gap-2 px-3 py-1.5 rounded border border-border bg-surface"
                >
                  <span
                    className={`w-2 h-2 rounded-full inline-block ${CATALOG_COLOR_CLASSES[catalog] ?? "bg-gray-500"}`}
                  />
                  <Text className="text-foreground">
                    {CATALOG_LABELS[catalog]}
                  </Text>
                </span>
              ))}
            </div>
            <Paragraph className="text-muted !mb-0">
              An object page pulls in more than the positional match: available
              photometry, time-series photometry from the surveys that cover the
              position, and direct links out to SIMBAD, VizieR, NED, Aladin,
              Legacy Survey, SDSS and Pan-STARRS for the same coordinate.
            </Paragraph>
          </Section>

          <Section title="Practical notes" icon={<BulbOutlined />}>
            <ul className="text-muted space-y-3 pl-5 list-disc marker:text-border">
              <li>
                <Text strong className="text-foreground">
                  Radii are in arcseconds.
                </Text>{" "}
                The form lets you enter arcmin or degrees and converts for you;
                the API itself always takes arcseconds.
              </li>
              <li>
                <Text strong className="text-foreground">
                  Searches are capped at {MAX_RADIUS_ARCSEC}
                  &Prime;.
                </Text>{" "}
                Past roughly that point the service stops answering in
                reasonable time, so the interface refuses the request rather
                than leaving you waiting on one that will not return.
              </li>
              <li>
                <Text strong className="text-foreground">
                  An empty result is not an error.
                </Text>{" "}
                A catalog with nothing at your position returns no content —
                that is a real answer about the sky, not a failure.
              </li>
              <li>
                <Text strong className="text-foreground">
                  Separations are great-circle distances,
                </Text>{" "}
                not projected ones, so they stay correct near the poles.
              </li>
            </ul>
          </Section>

          <Section title="API" icon={<ApiOutlined />}>
            <Paragraph className="text-muted">
              Everything this interface does is available over HTTP. The service
              is at <Text code>{API_BASE_URL}</Text>, and the full schema is
              browsable:
            </Paragraph>
            <div className="flex flex-wrap gap-2 mb-4">
              {[
                { path: "/conesearch", desc: "single position" },
                { path: "/bulk-conesearch", desc: "many positions at once" },
                { path: "/metadata", desc: "one object" },
                { path: "/bulk-metadata", desc: "many objects" },
                { path: "/lightcurve", desc: "time-series photometry" },
              ].map(({ path, desc }) => (
                <Tag key={path} className="!mr-0 !bg-surface !border-border">
                  <Text code className="!text-foreground">
                    {path}
                  </Text>
                  <Text className="text-muted ml-2 text-xs">{desc}</Text>
                </Tag>
              ))}
            </div>
            <Paragraph className="!mb-0">
              <Link
                href={SWAGGER_URL}
                target="_blank"
                rel="noopener noreferrer"
              >
                Open the API documentation →
              </Link>
            </Paragraph>
          </Section>

          <Section title="Open source" icon={<GithubOutlined />}>
            <Paragraph className="text-muted !mb-0">
              XWave is released under the Apache License 2.0. The indexer, the
              search service and this interface all live in{" "}
              <Link href={REPO_URL} target="_blank" rel="noopener noreferrer">
                the project repository
              </Link>
              .
            </Paragraph>
          </Section>
        </div>
      </Content>
    </Layout>
  );
}
