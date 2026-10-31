/**
 * Minimal HEALPix geometry for the Learn page diagrams (Górski et al. 2005).
 *
 * Positions are in the HEALPix projection: x = φ ∈ [0, 2π), y ∈ [−π/2, π/2],
 * where the 12 base pixels are congruent diamonds with half-diagonal π/4.
 * Inside a face, (u, v) ∈ [0, 1]² run from the face's south vertex: u along
 * ix (north-east) and v along iy (north-west).
 */

export const QUARTER = Math.PI / 4;

/** Centre of each base pixel (face 0–11) in the HEALPix projection. */
export const FACE_CENTERS: ReadonlyArray<readonly [number, number]> =
  Array.from({ length: 12 }, (_, f) => {
    const row = Math.floor(f / 4);
    const k = f % 4;
    return row === 1
      ? ([k * 2 * QUARTER, 0] as const)
      : ([QUARTER + k * 2 * QUARTER, row === 0 ? QUARTER : -QUARTER] as const);
  });

/** Projection position of face coordinates (u, v). */
export function facePosition(face: number, u: number, v: number) {
  const [cx, cy] = FACE_CENTERS[face] ?? [0, 0];
  return [cx + (u - v) * QUARTER, cy - QUARTER + (u + v) * QUARTER] as const;
}

/**
 * Interleave the bits of (ix, iy) into a Z-order index: ix on even bits, iy on
 * odd ones. Arithmetic rather than bit operators, so it stays exact past 32
 * bits (order 18 indices reach ~8×10¹¹).
 */
export function interleave(ix: number, iy: number, bits: number): number {
  let z = 0;
  for (let b = bits - 1; b >= 0; b--) {
    z =
      z * 4 + (Math.floor(ix / 2 ** b) % 2) + 2 * (Math.floor(iy / 2 ** b) % 2);
  }
  return z;
}

/** NESTED index of pixel (ix, iy) of a face at the given order. */
export function nestIndex(face: number, ix: number, iy: number, order: number) {
  return face * 4 ** order + interleave(ix, iy, order);
}

const JRLL = [2, 2, 2, 2, 3, 3, 3, 3, 4, 4, 4, 4];
const JPLL = [1, 3, 5, 7, 0, 2, 4, 6, 1, 3, 5, 7];

/** RING index of pixel (ix, iy) of a face (healpix_base xyf2ring). */
export function ringIndex(face: number, ix: number, iy: number, order: number) {
  const nside = 2 ** order;
  const nl4 = 4 * nside;
  const npix = 12 * nside * nside;
  const ncap = 2 * nside * (nside - 1);
  const jr = (JRLL[face] ?? 0) * nside - ix - iy - 1;
  let nr: number;
  let before: number;
  let shift = 0;
  if (jr < nside) {
    nr = jr;
    before = 2 * nr * (nr - 1);
  } else if (jr > 3 * nside) {
    nr = nl4 - jr;
    before = npix - 2 * (nr + 1) * nr;
  } else {
    nr = nside;
    before = ncap + (jr - nside) * nl4;
    shift = (jr - nside) & 1;
  }
  let jp = Math.trunc(((JPLL[face] ?? 0) * nr + ix - iy + 1 + shift) / 2);
  if (jp > nl4) jp -= nl4;
  if (jp < 1) jp += nl4;
  return before + jp - 1;
}

export interface FacePixel {
  face: number;
  ix: number;
  iy: number;
  /** NESTED index at the requested order. */
  ipix: number;
}

/** Pixel containing (ra, dec) in degrees, NESTED scheme (ang2pix). */
export function angToNest(ra: number, dec: number, order: number): FacePixel {
  const nside = 2 ** order;
  const z = Math.sin((dec * Math.PI) / 180);
  const tt = ((((ra * Math.PI) / 180 / (Math.PI / 2)) % 4) + 4) % 4;
  let face: number;
  let ix: number;
  let iy: number;
  if (Math.abs(z) <= 2 / 3) {
    const t1 = nside * (0.5 + tt);
    const t2 = nside * z * 0.75;
    const jp = Math.floor(t1 - t2);
    const jm = Math.floor(t1 + t2);
    const ifp = Math.floor(jp / nside);
    const ifm = Math.floor(jm / nside);
    face = ifp === ifm ? ifp | 4 : ifp < ifm ? ifp : ifm + 8;
    ix = jm % nside;
    iy = nside - (jp % nside) - 1;
  } else {
    const ntt = Math.min(3, Math.floor(tt));
    const tp = tt - ntt;
    const tmp = nside * Math.sqrt(3 * (1 - Math.abs(z)));
    const jp = Math.min(nside - 1, Math.floor(tp * tmp));
    const jm = Math.min(nside - 1, Math.floor((1 - tp) * tmp));
    if (z >= 0) {
      face = ntt;
      ix = nside - jm - 1;
      iy = nside - jp - 1;
    } else {
      face = ntt + 8;
      ix = jp;
      iy = jm;
    }
  }
  return { face, ix, iy, ipix: nestIndex(face, ix, iy, order) };
}

/** Which of its parent's 4 children a pixel is, at each order 1…order. */
export function childDigits(pixel: FacePixel, order: number): number[] {
  return Array.from({ length: order }, (_, i) => {
    const shift = 2 ** (order - 1 - i);
    return (
      (Math.floor(pixel.ix / shift) % 2) +
      2 * (Math.floor(pixel.iy / shift) % 2)
    );
  });
}

/** Collapse sorted integers into inclusive [start, end] runs. */
export function toRanges(sorted: number[]): Array<[number, number]> {
  const ranges: Array<[number, number]> = [];
  for (const n of sorted) {
    const last = ranges[ranges.length - 1];
    if (last && n === last[1] + 1) last[1] = n;
    else ranges.push([n, n]);
  }
  return ranges;
}
