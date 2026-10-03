import { ImageResponse } from "next/og";

import { OG_SIZE, OgCard } from "@/app/components/share/OgCard";
import { loadLogoDataUri } from "@/app/components/share/ogLogo";

export const alt = "XWave — astronomical cross-match service";
export const size = OG_SIZE;
export const contentType = "image/png";

export default async function Image() {
  const logoSrc = await loadLogoDataUri();
  return new ImageResponse(
    <OgCard
      eyebrow="Gaia DR3 · AllWISE · eROSITA"
      headline="Cross-match the sky across surveys"
      lines={["Cone search, bulk cross-match and object pages"]}
      footer="Public API · open source"
      seed="xwave"
      logoSrc={logoSrc}
    />,
    size
  );
}
