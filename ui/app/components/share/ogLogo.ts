import { readFile } from "node:fs/promises";
import { join } from "node:path";

/**
 * The site logo (public/xwave-icon.svg) as a data URI for OG images. Read
 * from process.cwd() so Next traces the file into the serverless function.
 * Resolves to undefined if it can't be read, so a card still renders.
 */
export async function loadLogoDataUri(): Promise<string | undefined> {
  try {
    const svg = await readFile(join(process.cwd(), "public", "xwave-icon.svg"));
    return `data:image/svg+xml;base64,${svg.toString("base64")}`;
  } catch {
    return undefined;
  }
}
