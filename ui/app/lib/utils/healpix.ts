/**
 * NESTED HEALPix index of the cell holding (ra, dec) at `order`, after
 * ang2pix_nest in the HEALPix C library. Used to name the HiPS tile under a
 * position (`NorderK/DirD/NpixN`), so a mirror can be checked for that patch
 * of sky rather than just for its `properties` file.
 */
export function ang2pixNest(
  order: number,
  raDeg: number,
  decDeg: number
): number {
  const nside = 2 ** order;
  const z = Math.sin((decDeg * Math.PI) / 180);
  const za = Math.abs(z);
  const phi = ((((raDeg % 360) + 360) % 360) * Math.PI) / 180;
  const tt = phi / (Math.PI / 2); // in [0, 4)

  let face: number;
  let ix: number;
  let iy: number;
  if (za <= 2 / 3) {
    // Equatorial belt.
    const t1 = nside * (0.5 + tt);
    const t2 = nside * z * 0.75;
    const jp = Math.floor(t1 - t2);
    const jm = Math.floor(t1 + t2);
    const ifp = Math.floor(jp / nside);
    const ifm = Math.floor(jm / nside);
    face = ifp === ifm ? (ifp % 4) + 4 : ifp < ifm ? ifp % 4 : (ifm % 4) + 8;
    ix = jm % nside;
    iy = nside - (jp % nside) - 1;
  } else {
    // Polar caps.
    const ntt = Math.min(Math.floor(tt), 3);
    const tp = tt - ntt;
    const tmp = nside * Math.sqrt(3 * (1 - za));
    const jp = Math.min(Math.floor(tp * tmp), nside - 1);
    const jm = Math.min(Math.floor((1 - tp) * tmp), nside - 1);
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

  // Interleave the bits: x on even positions, y on odd ones.
  let inFace = 0;
  for (let b = 0; b < order; b++) {
    inFace +=
      ((ix >> b) & 1) * 2 ** (2 * b) + ((iy >> b) & 1) * 2 ** (2 * b + 1);
  }
  return face * nside * nside + inFace;
}

/** Path of the HiPS tile holding (ra, dec) at `order`, e.g. `Norder3/Dir0/Npix334.jpg`. */
export function hipsTilePath(
  order: number,
  raDeg: number,
  decDeg: number,
  ext: string
): string {
  const npix = ang2pixNest(order, raDeg, decDeg);
  const dir = Math.floor(npix / 10_000) * 10_000;
  return `Norder${order}/Dir${dir}/Npix${npix}.${ext}`;
}
