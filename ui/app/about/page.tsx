"use client";

import {
  ApiOutlined,
  BulbOutlined,
  CompassOutlined,
  DatabaseOutlined,
  GithubOutlined,
  NodeIndexOutlined,
  ReadOutlined,
  TeamOutlined,
} from "@ant-design/icons";
import { Card, Tag, Typography } from "antd";
import NextLink from "next/link";
import type { ReactNode } from "react";

import {
  CDS_ACKNOWLEDGEMENT,
  EXTERNAL_SOURCES,
  formatDefaultRadius,
} from "@/app/components/catalogs";
import { PageShell } from "@/app/components/layout";
import { API_BASE_URL } from "@/app/lib/api/client";
import { CATALOG_META } from "@/app/lib/constants/catalogMeta";
import {
  CATALOG_COLOR_CLASSES,
  CATALOG_OPTIONS,
} from "@/app/lib/constants/catalogs";
import { MAX_RADIUS_ARCSEC } from "@/app/lib/constants/search";
import { API_ORIGIN, REPO_URL, SWAGGER_URL } from "@/app/lib/constants/site";

const { Title, Paragraph, Text, Link } = Typography;

/** Pages that go deeper than this overview. */
const MORE_LINKS: { href: string; label: string; desc: string }[] = [
  {
    href: "/catalogs",
    label: "Catalogs",
    desc: "coverage, precision, columns and known issues",
  },
  {
    href: "/learn#methods",
    label: "Methods",
    desc: "choosing radii and reading match results",
  },
  { href: "/developers", label: "API playground", desc: "try every endpoint" },
  { href: "/changelog", label: "Changelog", desc: "what changed and when" },
  { href: "/contact", label: "Contact", desc: "questions, bugs, requests" },
];

/** Contributors, from the repository's git history. */
const CONTRIBUTORS = [
  "Diego Rodríguez Mancini",
  "Javier Arredondo Contreras",
  "Matías Medina",
];

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
    <PageShell
      title="About XWave"
      description="XWave is a positional cross-match service for astronomical catalogs. Give it a coordinate and a search radius, and it returns the objects each indexed catalog has at that position, along with their angular separation from your target."
    >
      <div>
        <Section title="How the cross-match works" icon={<NodeIndexOutlined />}>
          <Paragraph className="text-muted">
            Comparing a target against every row of a catalog does not scale, so
            the search runs in two stages: a cheap spatial pre-filter that
            narrows millions of objects down to a handful of candidates, then an
            exact distance computation on just those.
          </Paragraph>

          <div className="mt-6">
            <Step index={1} title="Catalogs are indexed onto a HEALPix grid">
              <Paragraph className="text-muted !mb-0">
                When a catalog is ingested, every object&apos;s sky position is
                mapped to a{" "}
                <Link
                  href="https://healpix.sourceforge.io/"
                  target="_blank"
                  rel="noopener noreferrer"
                >
                  HEALPix
                </Link>{" "}
                pixel using the <Text code>NESTED</Text> ordering scheme, and
                that pixel index is stored alongside the object. HEALPix divides
                the sphere into equal-area pixels, so &quot;which objects are
                near this point&quot; becomes an indexed lookup rather than a
                scan. That index is the <Text code>IPix</Text> column you see in
                the results table.
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
                deliberately generous: it may return objects outside the radius,
                but it will not miss any inside it.
              </Paragraph>
            </Step>

            <Step index={3} title="Candidates are ranked by proximity">
              <Paragraph className="text-muted !mb-0">
                The candidates from step 2 are loaded into a 2-D k-d tree keyed
                on right ascension and declination, and a nearest-neighbour
                query returns the closest ones to your target. This bounds how
                much work the final step has to do.
              </Paragraph>
            </Step>

            <Step index={4} title="Exact separations are computed and filtered">
              <Paragraph className="text-muted !mb-0">
                For each remaining candidate the service computes the
                great-circle separation from your target using the haversine
                formula, and discards anything beyond your radius. What survives
                is the result set, and that separation is the{" "}
                <Text code>Ang. Dist</Text> column — reported in arcseconds.
              </Paragraph>
            </Step>
          </div>

          <Card size="small" className="bg-surface border-border mt-2">
            <Text strong className="text-foreground">
              Why results are capped
            </Text>
            <Paragraph className="text-muted !mb-0 !mt-2">
              The proximity ranking in step 3 happens <em>before</em> the radius
              filter in step 4, so the neighbour limit bounds how many objects
              can come back at all. This interface requests a generous limit so
              that in practice the radius is what determines your results — but
              for a very crowded field, a large radius can still hit that
              ceiling.
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
          <div className="grid grid-cols-1 sm:grid-cols-3 gap-2 mb-4">
            {CATALOG_OPTIONS.map((catalog) => {
              const meta = CATALOG_META[catalog];
              return (
                <NextLink
                  key={catalog}
                  href={`/catalogs#${catalog}`}
                  className="block px-3 py-2.5 rounded border border-border bg-surface hover:bg-surface-elevated"
                >
                  <span className="flex items-center gap-2 text-foreground">
                    <span
                      className={`w-2 h-2 rounded-full inline-block ${CATALOG_COLOR_CLASSES[catalog] ?? "bg-gray-500"}`}
                    />
                    {meta.name}
                    <span className="text-neutral-500 text-xs">
                      {meta.release}
                    </span>
                  </span>
                  <span className="block text-xs text-neutral-400 mt-1">
                    {meta.coverage} · {formatDefaultRadius(catalog)} default
                  </span>
                </NextLink>
              );
            })}
          </div>
          <Paragraph className="text-muted !mb-0">
            Coverage maps, astrometric precision, columns and known issues for
            each one are on the{" "}
            <NextLink href="/catalogs">Catalogs page</NextLink>. An object page
            pulls in more than the positional match: photometry, light curves
            from the surveys that cover the position, and links out to SIMBAD,
            VizieR, NED, Aladin, Legacy Survey, SDSS and Pan-STARRS for the same
            coordinate.
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
              Past roughly that point the service stops answering in reasonable
              time, so the interface refuses the request rather than leaving you
              waiting on one that will not return.
            </li>
            <li>
              <Text strong className="text-foreground">
                An empty result is not an error.
              </Text>{" "}
              A catalog with nothing at your position returns no content — that
              is a real answer about the sky, not a failure.
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
            <Link href={SWAGGER_URL} target="_blank" rel="noopener noreferrer">
              Open the API documentation →
            </Link>
            <Text className="text-muted mx-2">·</Text>
            <NextLink href="/developers">
              Try it in the API playground →
            </NextLink>
          </Paragraph>
        </Section>

        <Section title="Data sources" icon={<ReadOutlined />}>
          <Paragraph className="text-muted">
            Cross-match results come from XWave&apos;s own index of{" "}
            {CATALOG_OPTIONS.map((c) => CATALOG_META[c].name).join(", ")}. The
            object page also shows data fetched live from these public services:
          </Paragraph>
          <ul className="text-muted space-y-3 pl-5 list-disc marker:text-border">
            {EXTERNAL_SOURCES.map(({ name, href, usedFor }) => (
              <li key={name}>
                <Link href={href} target="_blank" rel="noopener noreferrer">
                  {name}
                </Link>{" "}
                — {usedFor}
              </li>
            ))}
          </ul>
          <Paragraph className="text-muted !mb-0 !mt-4">
            {CDS_ACKNOWLEDGEMENT} If XWave helped your work, the Cite button on
            any object or results page gives the references and acknowledgement
            text for XWave and each catalog.
          </Paragraph>
        </Section>

        <Section title="Team & funding" icon={<TeamOutlined />}>
          <Paragraph className="text-muted">
            XWave is developed as part of the{" "}
            <Link
              href="https://alerce.online/"
              target="_blank"
              rel="noopener noreferrer"
            >
              ALeRCE
            </Link>{" "}
            ecosystem (the repository&apos;s working name is ALeRCE xmatch), and
            the public API is served by Universidad Diego Portales at{" "}
            <Text code>{API_ORIGIN.replace(/^https?:\/\//, "")}</Text>.
          </Paragraph>
          <Paragraph className="text-muted">
            Code contributors, from the git history: {CONTRIBUTORS.join(", ")},
            and{" "}
            <Link
              href={`${REPO_URL}/graphs/contributors`}
              target="_blank"
              rel="noopener noreferrer"
            >
              others on GitHub
            </Link>
            .
          </Paragraph>
          <Paragraph className="text-muted !mb-0">
            Funding acknowledgements are not listed here yet. If you need them
            for a report or a paper, please{" "}
            <NextLink href="/contact">get in touch</NextLink>.
          </Paragraph>
        </Section>

        <Section title="Open source" icon={<GithubOutlined />}>
          <Paragraph className="text-muted !mb-0">
            XWave is released under the Apache License 2.0. The indexer, the
            search service and this interface all live in{" "}
            <Link href={REPO_URL} target="_blank" rel="noopener noreferrer">
              the project repository
            </Link>
            ; see the <NextLink href="/changelog">changelog</NextLink> for what
            has shipped recently.
          </Paragraph>
        </Section>

        <Section title="Go deeper" icon={<CompassOutlined />}>
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-2">
            {MORE_LINKS.map(({ href, label, desc }) => (
              <NextLink
                key={href}
                href={href}
                className="block px-3 py-2.5 rounded border border-border bg-surface hover:bg-surface-elevated"
              >
                <span className="block text-foreground">{label} →</span>
                <span className="block text-xs text-neutral-400">{desc}</span>
              </NextLink>
            ))}
          </div>
        </Section>
      </div>
    </PageShell>
  );
}
