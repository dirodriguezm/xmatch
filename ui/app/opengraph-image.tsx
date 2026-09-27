import { ImageResponse } from "next/og";

import { OG_SIZE, OgCard } from "@/app/components/share/OgCard";

export const alt = "XWave — astronomical cross-match service";
export const size = OG_SIZE;
export const contentType = "image/png";

export default function Image() {
  return new ImageResponse(
    <OgCard
      eyebrow="Gaia DR3 · AllWISE · eROSITA"
      headline="Cross-match the sky across surveys"
      lines={["Cone search, bulk cross-match and object pages"]}
      footer="Public API · open source"
      seed="xwave"
    />,
    size
  );
}
