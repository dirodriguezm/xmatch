/**
 * Coordinate conversion and formatting utilities for astronomical coordinates
 */

/**
 * Convert Right Ascension (in degrees) to Hours:Minutes:Seconds format
 * @param ra - Right Ascension in degrees (0-360)
 * @returns Formatted string like "12h 30m 45.67s"
 */
export function toHMS(ra: number): string {
  const hours = ra / 15;
  const h = Math.floor(hours);
  const m = Math.floor((hours - h) * 60);
  const s = ((hours - h) * 60 - m) * 60;
  return `${h}h ${m}m ${s.toFixed(2)}s`;
}

/**
 * Convert Declination (in degrees) to Degrees:Minutes:Seconds format
 * @param dec - Declination in degrees (-90 to +90)
 * @returns Formatted string like "+45° 30′ 15.00″"
 */
export function toDMS(dec: number): string {
  const sign = dec >= 0 ? "+" : "-";
  const absDec = Math.abs(dec);
  const d = Math.floor(absDec);
  const m = Math.floor((absDec - d) * 60);
  const s = ((absDec - d) * 60 - m) * 60;
  return `${sign}${d}° ${m}′ ${s.toFixed(2)}″`;
}

/** Result of reading one typed angle: degrees, or why it was rejected. */
export type AngleInput =
  | { value: number; error?: undefined }
  | { value: null; error: string };

/** One number, the marker after it (if any) and the spacing around it. */
const ANGLE_FIELD = /(\d+(?:\.\d*)?|\.\d+)(\s*)([hdms:°′″'"]?)(\s*)/iy;
/** Markers allowed after the 1st, 2nd and 3rd field. */
const FIELD_MARKERS = [/^[hd°:]?$/i, /^[m′':]?$/i, /^[s″"]?$/i];

interface SplitAngle {
  sign: 1 | -1;
  parts: number[];
  /** Marker after the first field: "h", "d", "°", ":" or "". */
  unit: string;
}

/**
 * Splits "05:35:17.3", "05 35 17.3", "05h35m17.3s", "-05d23m28s" or
 * "-5°23′28″" into a sign and up to three unsigned numbers; a plain "83.82"
 * comes back as one number. Anything else ("M31", "1e3", "5.5.5") is null,
 * so a typo is rejected instead of half-read the way parseFloat("05:35:17")
 * reads 5.
 */
function splitAngle(input: string): SplitAngle | null {
  let text = input.trim();
  let sign: 1 | -1 = 1;
  if (/^[+\-\u2212]/.test(text)) {
    if (text[0] !== "+") sign = -1;
    text = text.slice(1).trim();
  }

  const parts: number[] = [];
  let unit = "";
  let pos = 0;
  while (pos < text.length) {
    const allowed = FIELD_MARKERS[parts.length];
    if (!allowed) return null;
    ANGLE_FIELD.lastIndex = pos;
    const match = ANGLE_FIELD.exec(text);
    if (!match) return null;
    const [whole, num = "", gap = "", marker = "", trailing = ""] = match;
    pos += whole.length;
    const isLast = pos === text.length;
    if (!allowed.test(marker)) return null;
    // Fields need something between them, and a separator needs a field after.
    if (!marker && !gap && !trailing && !isLast) return null;
    if (marker === ":" && isLast) return null;
    if (parts.length === 0) unit = marker.toLowerCase();
    parts.push(Number(num));
  }

  if (parts.length === 0) return null;
  // Only the last field may carry a fraction: "05.5:30" is not sexagesimal.
  if (parts.slice(0, -1).some((p) => !Number.isInteger(p))) return null;
  if (parts.slice(1).some((p) => p >= 60)) return null;
  return { sign, parts, unit };
}

function sexagesimalToNumber([a = 0, b = 0, c = 0]: number[]): number {
  return a + b / 60 + c / 3600;
}

const RA_FORMAT_HELP =
  "Enter RA as degrees (0 to 360) or h:m:s, e.g. 05:35:17.3";
const DEC_FORMAT_HELP =
  "Enter Dec as degrees (-90 to 90) or d:m:s, e.g. -05:23:28";

/**
 * Reads a typed Right Ascension: decimal degrees ("83.8221") or hours as
 * h:m:s, space-separated or with h/m/s markers ("05:35:17.3", "5 35 17.3",
 * "5h35m17.3s"). Ranges are checked here so an out-of-range value is caught
 * in the form rather than by the service.
 */
export function parseRaInput(input: string): AngleInput {
  if (!input.trim()) return { value: null, error: "Enter an RA" };
  const angle = splitAngle(input);
  if (!angle) return { value: null, error: RA_FORMAT_HELP };
  // Several fields are hours unless marked as degrees ("83d49m19s").
  const inDegrees =
    angle.unit === "d" ||
    angle.unit === "°" ||
    (angle.parts.length === 1 && angle.unit !== "h" && angle.unit !== ":");
  if (!inDegrees) {
    const hours = sexagesimalToNumber(angle.parts);
    if (angle.sign < 0 || hours >= 24) {
      return { value: null, error: "RA in h:m:s must be below 24h" };
    }
    return { value: hours * 15 };
  }
  const deg = angle.sign * sexagesimalToNumber(angle.parts);
  if (deg < 0 || deg > 360) {
    return { value: null, error: "RA must be between 0 and 360°" };
  }
  return { value: deg };
}

/**
 * Reads a typed Declination: decimal degrees ("-5.3911") or d:m:s,
 * space-separated or with d/m/s markers ("-05:23:28", "-5 23 28",
 * "-5d23m28s"). The sign is read from the text, so "-00:30:00" stays south.
 */
export function parseDecInput(input: string): AngleInput {
  if (!input.trim()) return { value: null, error: "Enter a Dec" };
  const angle = splitAngle(input);
  // Hours make no sense for a declination.
  if (!angle || angle.unit === "h") {
    return { value: null, error: DEC_FORMAT_HELP };
  }
  const deg = angle.sign * sexagesimalToNumber(angle.parts);
  if (deg < -90 || deg > 90) {
    return { value: null, error: "Dec must be between -90 and +90°" };
  }
  return { value: deg };
}

/**
 * Degrees as written into the search URL. 7 decimals is ~0.4 mas, well below
 * any catalogue's astrometric error, and keeps h:m:s conversions short.
 */
export function formatDegreesParam(deg: number): string {
  // Trimmed as text, not via Number → String, which would write 5e-7 for a
  // tiny value: a form the URL parser (rightly) rejects.
  const text = deg.toFixed(7).replace(/\.?0+$/, "");
  return text === "-0" ? "0" : text;
}

/**
 * A third column after RA and Dec that pasted lists often carry and that we
 * ignore: a number (radius, magnitude, year) or an epoch ("J2000", "B1950").
 */
const IGNORED_THIRD = /^(?:[+-]?(?:\d+(?:\.\d*)?|\.\d+)|[JB]\d{4}(?:\.\d+)?)$/i;

/** A bare integer, which could be the hours or degrees of h m s / d m s. */
const BARE_INTEGER = /^[+\-\u2212]?\d+$/;

/** "a, b, c": exactly `count` comma-separated fields of one token each. */
function isCommaSeparated(input: string, count: number): boolean {
  const fields = input
    .split(",")
    .map((f) => f.trim())
    .filter(Boolean);
  return fields.length === count && fields.every((f) => !/\s/.test(f));
}

/**
 * Try to parse a string as coordinates (RA Dec format)
 * Supports: "12:30:00 -45:00:00" (HMS/DMS), "12 30 00 -45 00 00",
 * "12h30m00s -45d00m00s" or "187.5 -45.0" (degrees), with a space or comma
 * between RA and Dec. A trailing comma, or a third column such as a radius
 * or "J2000", is ignored, as in a line pasted from a table.
 * @returns Parsed coordinates in degrees, or null if not valid coordinates
 */
export function parseCoordinates(
  input: string
): { ra: number; dec: number } | null {
  // Empty tokens come from a leading or trailing comma ("10.68 41.27,").
  const parts = input.split(/[\s,]+/).filter(Boolean);

  // Two tokens are one RA and one Dec; six are space-separated h m s d m s.
  // Three are RA, Dec and an ignored column, but only when the line is
  // comma-separated or RA and Dec could not be the start of an h m s
  // ("12 30 45" stays ambiguous). Any other count is left to the name
  // resolver ("NGC 1952", "M 31").
  let raText: string;
  let decText: string;
  if (parts.length === 2) {
    [raText = "", decText = ""] = parts;
  } else if (parts.length === 6) {
    raText = parts.slice(0, 3).join(" ");
    decText = parts.slice(3).join(" ");
  } else if (
    parts.length === 3 &&
    IGNORED_THIRD.test(parts[2] ?? "") &&
    (isCommaSeparated(input, 3) ||
      parts.slice(0, 2).every((p) => !BARE_INTEGER.test(p)))
  ) {
    [raText = "", decText = ""] = parts;
  } else {
    return null;
  }

  const ra = parseRaInput(raText);
  const dec = parseDecInput(decText);
  if (ra.value === null || dec.value === null) return null;
  return { ra: ra.value, dec: dec.value };
}

const DEG = Math.PI / 180;

/**
 * ICRS → Galactic rotation (Hipparcos, ESA 1997 vol. 1 §1.5.3). Rows are the
 * Galactic x/y/z axes expressed in ICRS.
 */
const ICRS_TO_GALACTIC = [
  [-0.0548755604162154, -0.873437090234885, -0.4838350155487132],
  [0.4941094278755837, -0.4448296299600112, 0.7469822444972189],
  [-0.8676661490190047, -0.1980763734312015, 0.4559837761750669],
];

/** Mean obliquity of the ecliptic at J2000.0 (IAU 2006), degrees. */
const OBLIQUITY_J2000_DEG = 23.4392911;

function wrap360(deg: number): number {
  return ((deg % 360) + 360) % 360;
}

/**
 * Galactic longitude/latitude (degrees) for ICRS (J2000) RA/Dec in degrees.
 */
export function equatorialToGalactic(
  ra: number,
  dec: number
): { l: number; b: number } {
  const v = [
    Math.cos(dec * DEG) * Math.cos(ra * DEG),
    Math.cos(dec * DEG) * Math.sin(ra * DEG),
    Math.sin(dec * DEG),
  ];
  const [x = 0, y = 0, z = 0] = ICRS_TO_GALACTIC.map((row) =>
    row.reduce((sum, c, k) => sum + c * (v[k] ?? 0), 0)
  );
  return {
    l: wrap360(Math.atan2(y, x) / DEG),
    b: Math.asin(Math.max(-1, Math.min(1, z))) / DEG,
  };
}

/**
 * Mean ecliptic longitude/latitude of J2000 (degrees) for ICRS RA/Dec in
 * degrees.
 */
export function equatorialToEcliptic(
  ra: number,
  dec: number
): { lambda: number; beta: number } {
  const eps = OBLIQUITY_J2000_DEG * DEG;
  const a = ra * DEG;
  const d = dec * DEG;
  const sinBeta =
    Math.sin(d) * Math.cos(eps) - Math.cos(d) * Math.sin(eps) * Math.sin(a);
  const lambda = Math.atan2(
    Math.sin(a) * Math.cos(eps) + Math.tan(d) * Math.sin(eps),
    Math.cos(a)
  );
  return {
    lambda: wrap360(lambda / DEG),
    beta: Math.asin(Math.max(-1, Math.min(1, sinBeta))) / DEG,
  };
}

/**
 * Position angle (degrees, east of north, 0–360) of the point (ra2, dec2) as
 * seen from (ra1, dec1). All inputs in degrees.
 */
export function positionAngle(
  ra1: number,
  dec1: number,
  ra2: number,
  dec2: number
): number {
  const dRa = (ra2 - ra1) * DEG;
  const d1 = dec1 * DEG;
  const d2 = dec2 * DEG;
  const y = Math.sin(dRa) * Math.cos(d2);
  const x =
    Math.cos(d1) * Math.sin(d2) - Math.sin(d1) * Math.cos(d2) * Math.cos(dRa);
  return wrap360(Math.atan2(y, x) / DEG);
}
