"use client";

import { Card, Typography } from "antd";
import type { ReactNode } from "react";

import { MAX_RADIUS_ARCSEC } from "@/app/lib/constants/search";

import { ConeSearchDiagram } from "./diagrams/ConeSearchDiagram";
import { HealpixSkyDiagram } from "./diagrams/HealpixSkyDiagram";
import { HealpixZoomDiagram } from "./diagrams/HealpixZoomDiagram";
import { IndexComparison } from "./diagrams/IndexComparison";
import { IpixPipelineDiagram } from "./diagrams/IpixPipelineDiagram";
import { Step } from "./Section";

const { Paragraph, Text, Link, Title } = Typography;

function Part({
  index,
  title,
  children,
}: {
  index: number;
  title: string;
  children: ReactNode;
}) {
  return (
    <section className="mb-12">
      <div className="mb-1 font-mono text-xs uppercase tracking-wider text-primary">
        Part {index}
      </div>
      <Title level={4} className="!mt-0 !mb-3 text-foreground">
        {title}
      </Title>
      {children}
    </section>
  );
}

const P = ({ children }: { children: ReactNode }) => (
  <Paragraph className="!text-neutral-300 max-w-[70ch]">{children}</Paragraph>
);

/**
 * How the cone search works, as implemented in
 * service/internal/search/conesearch and
 * service/internal/repository/mastercat_range.go.
 */
export function MethodsSection() {
  return (
    <div>
      <P>
        From general to specific: what a cross-match is, what HEALPix is, how
        the two fit together, and how XWave uses HEALPix to answer in
        milliseconds. The diagrams are interactive.
      </P>

      <Part index={1} title="Cross-match and cone search">
        <P>
          The same object shows up in many catalogs: Gaia sees it in the
          optical, AllWISE in the infrared, eROSITA in X-rays. Each catalog
          measures its position with a different error, so the coordinates never
          agree exactly. A <strong>cross-match</strong> asks: which sources of
          each catalog lie within a radius <em>r</em> of this position? Those
          inside are candidates for being the same object, nearest first.
        </P>
        <P>
          XWave answers it with a <strong>cone search</strong>. A position and a
          radius draw a circle on the sky; seen from Earth, that circle is the
          mouth of a cone, hence the name. Each catalog gets its own cone, with
          its own radius, and a bulk cross-match is simply one cone search per
          input position.
        </P>
        <ConeSearchDiagram />
        <Card size="small" className="bg-surface border-border">
          <Text strong className="text-foreground">
            The problem is scale
          </Text>
          <Paragraph className="!text-neutral-400 !mb-0 !mt-2">
            Gaia DR3 has 1,811,709,771 sources. Measuring the distance to every
            one of them would mean almost two billion calculations per search.
            XWave answers in milliseconds because it hardly measures anything:
            first it discards the whole sky except a tiny patch around your
            position. That is what HEALPix is for.
          </Paragraph>
        </Card>
      </Part>

      <Part index={2} title="What HEALPix is">
        <P>
          <Link
            href="https://healpix.sourceforge.io/"
            target="_blank"
            rel="noopener noreferrer"
          >
            HEALPix
          </Link>{" "}
          (Hierarchical Equal Area isoLatitude Pixelation) divides the sphere of
          the sky into cells. It was created by Górski and collaborators for the
          cosmic microwave background maps of WMAP and Planck, and it is now a
          standard in astronomy: the HiPS maps in the sky viewer are built on
          it, and every Gaia <Text code>source_id</Text> carries its HEALPix
          cell number.
        </P>
        <ul className="mb-4 grid list-none gap-3 p-0 text-sm text-neutral-400 md:grid-cols-3">
          <li className="rounded-lg border border-border bg-surface p-3">
            <strong className="text-foreground">Hierarchical.</strong> It starts
            from 12 base cells. Each splits into 4, and each of those into 4, as
            many times as needed. Every split raises the order by one.
          </li>
          <li className="rounded-lg border border-border bg-surface p-3">
            <strong className="text-foreground">Equal area.</strong> At a given
            order every cell has exactly the same area, at the equator or at a
            pole.
          </li>
          <li className="rounded-lg border border-border bg-surface p-3">
            <strong className="text-foreground">On rings of latitude.</strong>{" "}
            Cell centres line up on rings of constant latitude, which makes
            spherical-harmonic work fast.
          </li>
        </ul>
        <HealpixSkyDiagram />
        <P>
          There are two ways to number the cells. <strong>RING</strong> numbers
          them ring by ring, north to south, so two neighbours on different
          rings can have very different numbers. <strong>NESTED</strong> follows
          the hierarchy: the 4 children of a cell get consecutive numbers, and
          so do its 16 grandchildren. Turn on the path above to compare them.
          With NESTED,{" "}
          <strong>cells that are close on the sky have close numbers</strong>,
          which turns a search on the sky into a search over number intervals:
          exactly what a database does fast. XWave indexes its catalogs at order
          18, with cells about 0.8″ across, smaller than almost any
          source&apos;s position error.
        </P>
      </Part>

      <Part index={3} title="How they fit together">
        <P>
          A cone search asks a geometric question: what is near this point? A
          database knows nothing about nearness on the sky, but it finds numbers
          inside an interval very quickly. HEALPix translates between the two:
          it gives every point a number so that being close on the sky almost
          always means having similar numbers. &ldquo;Sources within 3″&rdquo;
          becomes &ldquo;sources numbered between A and B&rdquo;.
        </P>
        <IndexComparison />
        <P>
          How does a point get its number? By walking down the HEALPix tree. At
          each level the cell splits into 4 and you take the child that holds
          the point: 0, 1, 2 or 3. That path, written as digits, is the
          cell&apos;s address, the way a postal address goes from country to
          city, street and number. Follow it for θ¹ Ori F, a star in the
          Trapezium of Orion.
        </P>
        <HealpixZoomDiagram />
        <P>
          Two nearby sources walk down almost the same path: they share most of
          their address and differ only in the last digits. So all the cells of
          one area have consecutive numbers, and the search circle becomes a{" "}
          <strong>few intervals of numbers</strong> that the database reads in
          one go.
        </P>
      </Part>

      <Part index={4} title="How XWave uses HEALPix: the ipix">
        <P>
          The full path of one search through the XWave service, on a real 12.8″
          × 12.8″ patch: one order-14 cell split into its 256 order-18 cells,
          numbered NESTED within the patch. It is drawn as a square to read
          better; on the sky it is a diamond like the ones above.
        </P>
        <IpixPipelineDiagram />
      </Part>

      <Title level={4} className="!mt-0 !mb-3 text-foreground">
        Technical reference
      </Title>
      <P>
        Every cone search, from the search form, the object page or the API,
        runs this pipeline in the Go service.
      </P>
      <div className="mt-4">
        <Step index={1} title="Indexing: one HEALPix pixel per source">
          <Paragraph className="!text-neutral-400 !mb-0">
            At ingestion, each source&apos;s RA/Dec is mapped to a HEALPix pixel
            in the <Text code>NESTED</Text> scheme and stored as{" "}
            <Text code>ipix</Text>. The order is set per catalog and recorded
            alongside it; the reference configuration uses order 18 (N
            <sub>side</sub> = 2<sup>18</sup> = 262,144), pixels about 0.8″
            across.
          </Paragraph>
        </Step>
        <Step index={2} title="Candidate pixels: an inclusive disc query">
          <Paragraph className="!text-neutral-400 !mb-0">
            The radius is converted to radians and HEALPix&apos;s{" "}
            <Text code>query_disc_inclusive</Text> returns every pixel touching
            the disc as ranges of indices, with an oversampling factor of 4. It
            may over-select, but never drops a source inside the radius. If
            catalogs were indexed at different orders, the query runs once per
            order.
          </Paragraph>
        </Step>
        <Step index={3} title="Fetch candidates by pixel range">
          <Paragraph className="!text-neutral-400 !mb-0">
            Sources whose <Text code>ipix</Text> falls in those ranges are read
            with range queries. With a specific catalog (
            <Text code>catalog=gaia</Text>) the others are filtered out;{" "}
            <Text code>catalog=all</Text> keeps every catalog.
          </Paragraph>
        </Step>
        <Step index={4} title="Great-circle separation and radius cut">
          <Paragraph className="!text-neutral-400 !mb-0">
            Every candidate gets its true angular separation from the haversine
            formula, in arcseconds, and anything beyond the radius is dropped.
          </Paragraph>
        </Step>
        <Step index={5} title="Rank and cut to N">
          <Paragraph className="!text-neutral-400 !mb-0">
            Survivors are sorted nearest first and cut to the requested number
            of neighbours (<Text code>nneighbor</Text>; the API default is 1,
            this interface asks for 100). The response also carries the total
            number of matches and the count per catalog, both taken before the
            cut, so you can tell when the limit hid sources.
          </Paragraph>
        </Step>
      </div>

      <Card size="small" className="bg-surface border-border mb-4">
        <Text strong className="text-foreground">
          Bulk cross-match
        </Text>
        <Paragraph className="!text-neutral-400 !mb-0 !mt-2">
          Bulk requests run the same five steps for every input position. The
          list is split into chunks (500 positions in the default configuration)
          that are processed concurrently by a small worker pool (4 by default).
          Each result carries the index of the input position it belongs to, and
          duplicates within one position are removed.
        </Paragraph>
      </Card>

      <Card size="small" className="bg-surface border-border">
        <Text strong className="text-foreground">
          What the matcher does not do
        </Text>
        <ul className="mt-2 mb-0 pl-5 text-neutral-400">
          <li>
            <Text strong className="text-neutral-300">
              No epoch propagation.
            </Text>{" "}
            Positions are compared as stored; proper motion is not applied, so
            fast-moving stars can fall outside small radii across catalogs of
            different epochs.
          </li>
          <li>
            <Text strong className="text-neutral-300">
              No probabilistic association.
            </Text>{" "}
            There is no likelihood-ratio or Bayesian counterpart probability;
            every source within the radius is returned, ranked by distance.
          </li>
          <li>
            <Text strong className="text-neutral-300">
              Each catalog is matched to your position, not to each other.
            </Text>{" "}
            A Gaia and an AllWISE result for the same query are both near your
            coordinate, which is not the same as being each other&apos;s
            counterpart.
          </li>
          <li>
            <Text strong className="text-neutral-300">
              Radius ceiling in the interface.
            </Text>{" "}
            Searches above {MAX_RADIUS_ARCSEC}″ are refused client-side because
            they do not return in reasonable time.
          </li>
        </ul>
      </Card>
    </div>
  );
}
