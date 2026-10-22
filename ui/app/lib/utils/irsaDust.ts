/**
 * Galactic reddening from the IRSA Galactic Dust Reddening and Extinction
 * service (https://irsa.ipac.caltech.edu/applications/DUST/).
 */

export interface GalacticReddening {
  /** E(B−V) from Schlafly & Finkbeiner (2011), the recalibrated SFD map. */
  ebvSF11: number;
  /** E(B−V) from Schlegel, Finkbeiner & Davis (1998). */
  ebvSFD: number;
}

function readValue(block: string, tag: string): number | null {
  const m = block.match(new RegExp(`<${tag}>\\s*([-\\d.eE+]+)\\s*\\(mag\\)`));
  if (!m) return null;
  const n = Number(m[1]);
  return Number.isFinite(n) ? n : null;
}

/**
 * Read the E(B−V) value at the requested position (the reference pixel, not
 * the regional mean). Throws when the service reports an error or the
 * reddening block is missing.
 */
export function parseIrsaDust(xml: string): GalacticReddening {
  if (!/<results\s+status="ok"/.test(xml)) {
    const msg = xml.match(/<message>([\s\S]*?)<\/message>/)?.[1]?.trim();
    throw new Error(`IRSA dust: ${msg ?? "service returned an error"}`);
  }
  const block = [...xml.matchAll(/<result>([\s\S]*?)<\/result>/g)]
    .map((m) => m[1])
    .find((b) => /E\(B-V\)\s+Reddening/.test(b));
  if (!block) throw new Error("IRSA dust: no E(B-V) block in response");

  const ebvSF11 = readValue(block, "refPixelValueSandF");
  const ebvSFD = readValue(block, "refPixelValueSFD");
  if (ebvSF11 == null || ebvSFD == null) {
    throw new Error("IRSA dust: E(B-V) values missing");
  }
  return { ebvSF11, ebvSFD };
}
