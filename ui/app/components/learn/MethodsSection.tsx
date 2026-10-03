"use client";

import { Card, Typography } from "antd";

import { MAX_RADIUS_ARCSEC } from "@/app/lib/constants/search";

import { Step } from "./Section";

const { Paragraph, Text, Link } = Typography;

/**
 * How the cone search works, as implemented in
 * service/internal/search/conesearch and service/internal/search/knn.
 */
export function MethodsSection() {
  return (
    <div>
      <Paragraph className="!text-neutral-300">
        Every cone search — from the search form, the object page or the API —
        runs the same pipeline in the Go service: a HEALPix pre-filter narrows
        the catalog to a few pixels, then exact separations are computed on just
        those candidates.
      </Paragraph>

      <div className="mt-6">
        <Step index={1} title="Indexing: one HEALPix pixel per source">
          <Paragraph className="!text-neutral-400 !mb-0">
            At ingestion, each source&apos;s RA/Dec is mapped to a{" "}
            <Link
              href="https://healpix.sourceforge.io/"
              target="_blank"
              rel="noopener noreferrer"
            >
              HEALPix
            </Link>{" "}
            pixel in the <Text code>NESTED</Text> scheme and stored with the
            source as <Text code>ipix</Text>. The HEALPix order is set per
            catalog and recorded alongside it; the repository&apos;s reference
            configuration uses order 18 (N<sub>side</sub> = 2<sup>18</sup> =
            262,144), i.e. pixels about 0.8″ across. Because NESTED numbering
            keeps nearby pixels numerically close, a patch of sky becomes a few
            contiguous index ranges.
          </Paragraph>
        </Step>

        <Step index={2} title="Candidate pixels: an inclusive disc query">
          <Paragraph className="!text-neutral-400 !mb-0">
            Your radius (arcseconds) is converted to radians and HEALPix&apos;s{" "}
            <Text code>query_disc_inclusive</Text> returns every pixel that
            touches the disc, as ranges of pixel indices, using an oversampling
            factor of 4. &ldquo;Inclusive&rdquo; means partially covered pixels
            are kept, and a few that barely miss the disc may be too: the
            pre-filter can over-select, but never drops a source inside the
            radius. If catalogs were indexed at different orders, the query runs
            once per order.
          </Paragraph>
        </Step>

        <Step index={3} title="Fetch candidates by pixel range">
          <Paragraph className="!text-neutral-400 !mb-0">
            Sources whose <Text code>ipix</Text> falls in those ranges are read
            with range queries. When you ask for a specific catalog (
            <Text code>catalog=gaia</Text>), the others are filtered out;{" "}
            <Text code>catalog=all</Text> keeps every catalog.
          </Paragraph>
        </Step>

        <Step index={4} title="Nearest-N ranking with a k-d tree">
          <Paragraph className="!text-neutral-400 !mb-0">
            Candidates are loaded into a 2-D k-d tree keyed on (RA, Dec) and the{" "}
            <Text code>nneighbor</Text> closest to your position are kept (the
            API default is 1; this interface asks for 100). This ranking treats
            RA and Dec as flat coordinates in degrees — it does not apply a cos
            δ factor or wrap at RA = 0°/360°.
          </Paragraph>
        </Step>

        <Step index={5} title="Great-circle separation and radius cut">
          <Paragraph className="!text-neutral-400 !mb-0">
            For each of those neighbours the service computes the true angular
            separation with the haversine formula, converts it to arcseconds and
            drops anything beyond your radius. Survivors are returned
            nearest-first with their separation — the distance column in the
            results.
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
        <ul className="text-neutral-400 space-y-2 pl-5 list-disc marker:text-neutral-600 mt-2 mb-0">
          <li>
            <Text strong className="text-foreground">
              No epoch propagation.
            </Text>{" "}
            Positions are compared as stored; proper motion is not applied, so
            fast-moving stars can fall outside small radii across catalogs of
            different epochs.
          </li>
          <li>
            <Text strong className="text-foreground">
              No probabilistic association.
            </Text>{" "}
            There is no likelihood-ratio or Bayesian counterpart probability;
            every source within the radius is returned, ranked by distance.
          </li>
          <li>
            <Text strong className="text-foreground">
              Each catalog is matched to your position, not to each other.
            </Text>{" "}
            A Gaia and an AllWISE result for the same query are both near your
            coordinate, which is not the same as being each other&apos;s
            counterpart.
          </li>
          <li>
            <Text strong className="text-foreground">
              Edge cases of the flat ranking.
            </Text>{" "}
            Close to the poles or across RA = 0°, the step-4 ranking can differ
            from true separation. With a small <Text code>nneighbor</Text> a
            genuinely closer source could be skipped; asking for more neighbours
            avoids it.
          </li>
          <li>
            <Text strong className="text-foreground">
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
