/* eslint-disable react/forbid-dom-props -- Satori (next/og) only understands inline styles. */

/**
 * Markup for Open Graph cards rendered with `next/og` (Satori): flexbox only,
 * inline styles, every multi-child element needs `display: flex`.
 */

export const OG_SIZE = { width: 1200, height: 630 };

const BG = "#0a0a0a";
const ACCENT = "#1668dc";
const MUTED = "#a3a3a3";

/** Deterministic PRNG so a given object always gets the same starfield. */
function mulberry32(seed: number) {
  let a = seed >>> 0;
  return () => {
    a = (a + 0x6d2b79f5) >>> 0;
    let t = a;
    t = Math.imul(t ^ (t >>> 15), t | 1);
    t ^= t + Math.imul(t ^ (t >>> 7), t | 61);
    return ((t ^ (t >>> 14)) >>> 0) / 4294967296;
  };
}

function hashString(s: string): number {
  let h = 2166136261;
  for (let i = 0; i < s.length; i++) {
    h ^= s.charCodeAt(i);
    h = Math.imul(h, 16777619);
  }
  return h >>> 0;
}

function Starfield({ seed }: { seed: string }) {
  const rand = mulberry32(hashString(seed));
  const stars = Array.from({ length: 90 }, (_, i) => {
    const size = rand() < 0.1 ? 4 : rand() < 0.4 ? 2.5 : 1.5;
    return {
      key: i,
      left: rand() * OG_SIZE.width,
      top: rand() * OG_SIZE.height,
      size,
      opacity: 0.25 + rand() * 0.65,
    };
  });
  return (
    <div style={{ position: "absolute", inset: 0, display: "flex" }}>
      {stars.map((s) => (
        <div
          key={s.key}
          style={{
            position: "absolute",
            left: s.left,
            top: s.top,
            width: s.size,
            height: s.size,
            borderRadius: 9999,
            background: "#ffffff",
            opacity: s.opacity,
          }}
        />
      ))}
    </div>
  );
}

/** Reticle centred at (cx, cy): two rings plus four ticks. */
function Crosshair({ cx, cy }: { cx: number; cy: number }) {
  const ring = (r: number, opacity: number) => (
    <div
      style={{
        position: "absolute",
        left: cx - r,
        top: cy - r,
        width: r * 2,
        height: r * 2,
        borderRadius: 9999,
        border: `3px solid ${ACCENT}`,
        opacity,
      }}
    />
  );
  const tick = (x: number, y: number, w: number, h: number) => (
    <div
      style={{
        position: "absolute",
        left: x,
        top: y,
        width: w,
        height: h,
        background: ACCENT,
      }}
    />
  );
  return (
    <div style={{ position: "absolute", inset: 0, display: "flex" }}>
      {ring(150, 0.35)}
      {ring(70, 0.9)}
      {tick(cx - 200, cy - 1.5, 110, 3)}
      {tick(cx + 90, cy - 1.5, 110, 3)}
      {tick(cx - 1.5, cy - 200, 3, 110)}
      {tick(cx - 1.5, cy + 90, 3, 110)}
      <div
        style={{
          position: "absolute",
          left: cx - 9,
          top: cy - 9,
          width: 18,
          height: 18,
          borderRadius: 9999,
          background: "#ffffff",
          boxShadow: `0 0 40px 12px ${ACCENT}`,
        }}
      />
    </div>
  );
}

function Wordmark() {
  return (
    <div style={{ display: "flex", alignItems: "center", gap: 14 }}>
      <div
        style={{
          width: 40,
          height: 40,
          borderRadius: 10,
          background: ACCENT,
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
          color: "#fff",
          fontSize: 26,
          fontWeight: 700,
        }}
      >
        X
      </div>
      <div
        style={{
          fontSize: 34,
          fontWeight: 700,
          color: "#fff",
          letterSpacing: -0.5,
        }}
      >
        XWave
      </div>
    </div>
  );
}

interface OgCardProps {
  /** Small label above the headline, e.g. the catalog. */
  eyebrow?: string;
  headline: string;
  /** Lines under the headline (coordinates, tagline…). */
  lines?: string[];
  footer?: string;
  seed?: string;
}

export function OgCard({
  eyebrow,
  headline,
  lines = [],
  footer,
  seed = headline,
}: OgCardProps) {
  const headlineSize = headline.length > 26 ? 56 : 72;
  return (
    <div
      style={{
        width: "100%",
        height: "100%",
        display: "flex",
        position: "relative",
        background: `radial-gradient(circle at 78% 50%, #10213d 0%, ${BG} 55%)`,
        fontFamily: "sans-serif",
      }}
    >
      <Starfield seed={seed} />
      <Crosshair cx={930} cy={315} />
      <div
        style={{
          position: "relative",
          display: "flex",
          flexDirection: "column",
          justifyContent: "space-between",
          padding: 64,
          width: 760,
          height: "100%",
        }}
      >
        <Wordmark />
        <div style={{ display: "flex", flexDirection: "column", gap: 18 }}>
          {eyebrow && (
            <div style={{ display: "flex" }}>
              <div
                style={{
                  fontSize: 26,
                  color: "#fff",
                  background: "rgba(22,104,220,0.35)",
                  border: `2px solid ${ACCENT}`,
                  borderRadius: 999,
                  padding: "6px 20px",
                }}
              >
                {eyebrow}
              </div>
            </div>
          )}
          <div
            style={{
              fontSize: headlineSize,
              fontWeight: 700,
              color: "#fff",
              lineHeight: 1.1,
              letterSpacing: -1,
              wordBreak: "break-word",
            }}
          >
            {headline}
          </div>
          {lines.map((line) => (
            <div key={line} style={{ fontSize: 30, color: MUTED }}>
              {line}
            </div>
          ))}
        </div>
        <div style={{ fontSize: 24, color: MUTED }}>{footer ?? ""}</div>
      </div>
    </div>
  );
}
