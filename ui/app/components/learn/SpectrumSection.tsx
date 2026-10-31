"use client";

import { Typography } from "antd";
import Image from "next/image";
import Link from "next/link";
import { type KeyboardEvent, type ReactNode, useState } from "react";

import { getSearchCatalogColor } from "@/app/lib/constants/catalogs";

const { Paragraph, Title } = Typography;

// 3C 273, the first quasar identified. The cutouts in /public/learn/3c273 are
// 4′ × 4′ (360 px, 1.5 px per arcsecond) from CDS hips2fits.
const OBJECT_HREF = `/object/${encodeURIComponent("Gaia DR3 3700386905605055360")}?catalog=gaia`;
const PX_PER_ARCSEC = 1.5;

interface Band {
  file: string;
  tick: string;
  /** Static Tailwind class placing the tick on the ruler. */
  position: string;
  name: string;
  survey: string;
  meta: string[];
  /** Radio beam drawn over the image, in arcsec. */
  beam?: number;
  /** Catalog XWave cross-matches in this band, if any. */
  xwave?: string;
  text: ReactNode;
}

const BANDS: Band[] = [
  {
    file: "nvss",
    tick: "1.4 GHz",
    position: "left-[3%]",
    name: "Radio",
    survey: "NVSS · VLA · 1.4 GHz",
    meta: ["λ ≈ 21 cm", "45″ beam"],
    beam: 45,
    text: (
      <>
        In radio, 3C 273 shines from electrons spiralling at nearly the speed of
        light in magnetic fields. NVSS is very sensitive but very blurry: with a
        45″ beam the quasar and its jet merge into one round blob.
      </>
    ),
  },
  {
    file: "racs",
    tick: "1.4 GHz",
    position: "left-[14%]",
    name: "Radio",
    survey: "RACS-mid · ASKAP · 1.4 GHz",
    meta: ["λ ≈ 21 cm", "~10″ beam"],
    text: (
      <>
        Same frequency, a sharper radio telescope. Now the blob splits in two:
        the quasar core and the <strong>jet</strong> heading south-west, a
        stream of plasma some 200,000 light-years long.
      </>
    ),
  },
  {
    file: "wise",
    tick: "3–12 µm",
    position: "left-[38%]",
    name: "Mid-infrared",
    survey: "AllWISE · WISE",
    meta: ["λ = 3.4–22 µm", "~6″"],
    xwave: "AllWISE",
    text: (
      <>
        Mid-infrared shows the hot dust around the central black hole. The
        quasar is a point: WISE is not sharp enough to see the jet.
      </>
    ),
  },
  {
    file: "tmass",
    tick: "1–2 µm",
    position: "left-[50%]",
    name: "Near-infrared",
    survey: "2MASS",
    meta: ["λ = 1.2–2.2 µm", "~2–3″"],
    text: (
      <>
        In the near-infrared 3C 273 is still a bright point, next to a few stars
        of our own Galaxy that lie in the same direction.
      </>
    ),
  },
  {
    file: "optical",
    tick: "0.5 µm",
    position: "left-[66%]",
    name: "Optical",
    survey: "Pan-STARRS DR1",
    meta: ["λ ≈ 0.4–0.9 µm", "~1″"],
    xwave: "Gaia",
    text: (
      <>
        In the optical it looks like a star, hence <em>quasi-stellar</em>. A
        good image shows the jet as a faint line to the south-west. Gaia, which
        XWave cross-matches, pins its position to a fraction of a
        milliarcsecond.
      </>
    ),
  },
  {
    file: "chandra",
    tick: "0.5–7 keV",
    position: "left-[94%]",
    name: "X-rays",
    survey: "Chandra",
    meta: ["0.5–7 keV", "~0.5″"],
    xwave: "eROSITA",
    text: (
      <>
        X-rays show gas at millions of degrees near the black hole, and the jet,
        which shines here too. Chandra is exceptionally sharp; eROSITA, which
        XWave cross-matches, covers the whole sky with less detail. The diagonal
        line is a detector read-out artefact caused by how bright the source is.
      </>
    ),
  },
];

const RESULTS = [
  {
    catalog: "gaia",
    band: "Optical · Gaia DR3",
    id: "Gaia DR3 3700386905605055360",
    sep: "0.002″",
  },
  {
    catalog: "allwise",
    band: "Infrared · AllWISE",
    id: "AllWISE 1875p015_ac51-058179",
    sep: "0.04″",
  },
  {
    catalog: "erosita",
    band: "X-rays · eROSITA",
    id: "1eRASS J122906.7+020309",
    sep: "1.44″",
  },
];

/** Crosshair at the target and a 30″ scale bar, drawn over a 360 px cutout. */
function Overlay({ beam }: { beam?: number }) {
  return (
    <svg
      viewBox="0 0 360 360"
      className="absolute inset-0 h-full w-full"
      aria-hidden="true"
    >
      <path
        d="M180 166V173M180 187V194M166 180H173M187 180H194"
        stroke="#fff"
        strokeWidth={1.6}
      />
      <line
        x1={16}
        x2={16 + 30 * PX_PER_ARCSEC}
        y1={344}
        y2={344}
        stroke="#fff"
        strokeWidth={2}
      />
      <text x={16} y={336} fill="#fff" fontSize={11} fontFamily="monospace">
        30″
      </text>
      {beam && (
        <circle
          cx={180}
          cy={180}
          r={(beam * PX_PER_ARCSEC) / 2}
          fill="none"
          stroke="#fff"
          strokeOpacity={0.8}
          strokeDasharray="4 4"
          strokeWidth={1.4}
        />
      )}
    </svg>
  );
}

function Cutout({
  file,
  alt,
  priority = false,
}: {
  file: string;
  alt: string;
  priority?: boolean;
}) {
  return (
    <Image
      src={`/learn/3c273/${file}.jpg`}
      alt={alt}
      width={360}
      height={360}
      priority={priority}
      className="h-full w-full object-cover"
    />
  );
}

const P = ({ children }: { children: ReactNode }) => (
  <Paragraph className="!text-neutral-300 max-w-[70ch]">{children}</Paragraph>
);

const H = ({ children }: { children: ReactNode }) => (
  <Title level={4} className="!mt-10 !mb-3 text-foreground">
    {children}
  </Title>
);

/** "Across the spectrum": one real object, 3C 273, from radio to X-rays. */
export function SpectrumSection() {
  const [current, setCurrent] = useState(4);
  const band = BANDS[current];

  const onRulerKey = (e: KeyboardEvent<HTMLDivElement>) => {
    const next =
      e.key === "ArrowRight"
        ? current + 1
        : e.key === "ArrowLeft"
          ? current - 1
          : null;
    if (next === null || next < 0 || next >= BANDS.length) return;
    e.preventDefault();
    setCurrent(next);
    e.currentTarget.querySelectorAll("button")[next]?.focus();
  };

  return (
    <div>
      <P>
        Instead of describing wavelengths in the abstract, follow one real
        object: 3C 273, the first quasar ever identified. Every image below is
        the same 4′ patch of sky, and the matches are what XWave returns for it
        today.
      </P>

      <figure className="m-0 my-6 overflow-hidden rounded-lg border border-border bg-surface">
        <div className="grid md:grid-cols-[minmax(0,11fr)_minmax(0,10fr)]">
          <div className="relative aspect-square bg-black">
            {BANDS.map((b, i) => (
              <div
                key={b.file}
                className={`absolute inset-0 transition-opacity duration-300 motion-reduce:transition-none ${i === current ? "opacity-100" : "opacity-0"}`}
              >
                <Cutout
                  file={b.file}
                  alt={`3C 273 in ${b.name.toLowerCase()}, ${b.survey}`}
                  priority={i === current}
                />
              </div>
            ))}
            <Overlay beam={band.beam} />
          </div>
          <div className="flex flex-col gap-3 p-5" aria-live="polite">
            <span className="text-sm text-neutral-500">{band.survey}</span>
            <span className="text-xl font-semibold text-foreground">
              {band.name}
            </span>
            <div className="flex flex-wrap gap-1.5">
              {band.meta.map((m) => (
                <span
                  key={m}
                  className="rounded border border-border bg-surface-elevated px-2 py-1 font-mono text-xs text-neutral-400"
                >
                  {m}
                </span>
              ))}
              <span
                className={`rounded border px-2 py-1 font-mono text-xs ${band.xwave ? "border-[#2f5a1d] text-[#73d13d]" : "border-border text-neutral-400"}`}
              >
                {band.xwave
                  ? `XWave cross-matches ${band.xwave}`
                  : "In XWave: image only"}
              </span>
            </div>
            <p className="m-0 text-sm text-neutral-300">{band.text}</p>
          </div>
        </div>
        <div className="border-t border-border px-5 pb-4 pt-3">
          <div
            className="relative h-14"
            role="group"
            aria-label="Wavelength"
            onKeyDown={onRulerKey}
          >
            <div className="absolute inset-x-0 top-[22px] h-2.5 rounded-full bg-[linear-gradient(90deg,#ff7a45_0%,#c2410c_18%,#7c3aed_36%,#9254de_46%,#4096ff_62%,#a5d8ff_70%,#eb2f96_90%,#f9a8d4_100%)] opacity-85" />
            {BANDS.map((b, i) => (
              <button
                key={b.file}
                type="button"
                onClick={() => setCurrent(i)}
                aria-pressed={i === current}
                aria-label={`${b.name}, ${b.survey}`}
                className={`absolute top-0 -translate-x-1/2 cursor-pointer whitespace-nowrap border-0 bg-transparent p-0 font-mono text-xs ${b.position} ${i === current ? "text-foreground" : "text-neutral-500 hover:text-neutral-300"}`}
              >
                {b.tick}
                <span
                  className={`mx-auto mt-[6px] block h-3.5 w-3.5 rounded-full border-2 border-neutral-200 ${i === current ? "bg-neutral-200 ring-4 ring-white/15" : "bg-background"}`}
                />
              </button>
            ))}
          </div>
          <div className="flex justify-between font-mono text-xs text-neutral-500">
            <span>← longer waves · lower energy</span>
            <span>higher energy · shorter waves →</span>
          </div>
        </div>
      </figure>
      <p className="m-0 text-xs text-neutral-500">
        Cutouts from NVSS, RACS, AllWISE, 2MASS, Pan-STARRS and Chandra via the
        CDS hips2fits service. The cross marks the Gaia position; the bar is
        30″. With a wavelength selected, the arrow keys move along the spectrum.
      </p>

      <H>One search, three catalogs, three kinds of light</H>
      <P>
        This is what XWave returns for 3C 273&apos;s position with each
        catalog&apos;s suggested radius. Three instruments that see very
        different light find the same source within 1.5″, and XWave puts them in
        one answer.
      </P>
      <div className="mb-4 flex flex-col gap-2">
        {RESULTS.map((r) => (
          <div
            key={r.id}
            className="grid grid-cols-[1fr_auto] items-center gap-x-4 gap-y-1 rounded-lg border border-border bg-surface px-4 py-3 sm:grid-cols-[150px_1fr_auto]"
          >
            <span className="col-span-2 flex items-center gap-2 font-mono text-xs text-neutral-300 sm:col-span-1">
              <svg width={8} height={8} aria-hidden="true">
                <circle
                  cx={4}
                  cy={4}
                  r={4}
                  fill={getSearchCatalogColor(r.catalog)}
                />
              </svg>
              {r.band}
            </span>
            <code className="font-mono text-sm text-foreground">{r.id}</code>
            <span className="font-mono text-base font-semibold text-foreground">
              {r.sep}
            </span>
          </div>
        ))}
      </div>
      <P>
        eROSITA lands further away because its positions are less precise, not
        because it is a different object. That is why each catalog gets its own
        radius: 3″ for Gaia and AllWISE, 20″ for eROSITA.{" "}
        <Link href={OBJECT_HREF} className="text-primary">
          Open 3C 273 in XWave
        </Link>
        .
      </P>

      <H>What is a blob in radio has detail in the optical</H>
      <P>
        Not every kind of light is seen as sharply. NVSS sees the sky through a
        45″ beam, and 3C 273 is so bright that its blob spreads even further.
        Draw that circle on the optical image and it holds everything the
        optical resolves separately: the quasar and its jet. Here the answer is
        obvious because the source is bright and isolated. For a faint radio
        source in a field full of stars and galaxies, the same circle holds
        several candidates and the question becomes <strong>which one</strong>{" "}
        is the counterpart. That is why each instrument&apos;s sharpness sets
        the search radius.
      </P>
      <div className="mb-4 grid gap-3 sm:grid-cols-2">
        {[
          {
            file: "nvss",
            alt: "3C 273 in radio at 1.4 GHz (NVSS): a round blob",
            caption: "Radio, NVSS (1.4 GHz). The circle is the 45″ beam.",
          },
          {
            file: "optical",
            alt: "The same field in the optical (Pan-STARRS) with the 45 arcsecond circle",
            caption: "Optical, Pan-STARRS. The same 45″ circle.",
          },
        ].map((f) => (
          <figure key={f.file} className="m-0 flex flex-col gap-2">
            <div className="relative aspect-square overflow-hidden rounded-lg border border-border bg-black">
              <Cutout file={f.file} alt={f.alt} />
              <Overlay beam={45} />
            </div>
            <figcaption className="text-xs text-neutral-500">
              {f.caption}
            </figcaption>
          </figure>
        ))}
      </div>
      <P>
        XWave does not cross-match radio catalogs; the radio surveys are in each
        object&apos;s Sky View to look at, not to match against.
      </P>

      <H>XWave and other tools</H>
      <div className="overflow-hidden rounded-lg border border-border">
        {[
          [
            "I want to look at an object or a list",
            "XWave: one search across Gaia, AllWISE and eROSITA, and an object page with SED, light curves, spectra and images from radio to X-rays.",
          ],
          [
            "I want to cross-match any two tables",
            "CDS XMatch on the web, with any VizieR table or your own. Or TOPCAT on the desktop.",
          ],
          [
            "I want the probability of each counterpart",
            "Bayesian tools such as NWAY, built for X-ray sources with many candidates. This is the state of the art in association; XWave ranks by distance and does not compute probabilities.",
          ],
        ].map(([need, answer], i) => (
          <div
            key={need}
            className={`grid gap-1 bg-surface px-4 py-3 text-sm sm:grid-cols-[240px_1fr] sm:gap-4 ${i ? "border-t border-border" : ""}`}
          >
            <span className="font-semibold text-foreground">{need}</span>
            <span className="text-neutral-400">{answer}</span>
          </div>
        ))}
      </div>

      <p className="mt-8 text-xs text-neutral-500">
        Images: NVSS (NRAO), RACS (CSIRO/ASKAP), AllWISE (NASA/IPAC), 2MASS,
        Pan-STARRS DR1 and Chandra (CXC), via hips2fits (CDS, Strasbourg). XWave
        matches as of October 2026.
      </p>
    </div>
  );
}
