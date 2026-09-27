import { ImageResponse } from "next/og";

import { OG_SIZE, OgCard } from "@/app/components/share/OgCard";
import { toDMS, toHMS } from "@/app/lib/utils/coordinates";
import {
  catalogFromObjectId,
  fetchObjectSummary,
  formatRaDec,
} from "@/app/lib/utils/share";

export const alt = "XWave object preview";
export const size = OG_SIZE;
export const contentType = "image/png";

export default async function Image({
  params,
}: {
  params: Promise<{ objectId: string }>;
}) {
  const { objectId } = await params;
  const id = decodeURIComponent(objectId);
  const catalog = catalogFromObjectId(id);
  const summary = catalog
    ? await fetchObjectSummary(id, catalog.slug, 2000)
    : null;

  const lines = summary
    ? [
        formatRaDec(summary.ra, summary.dec),
        `${toHMS(summary.ra)}  ${toDMS(summary.dec)}`,
      ]
    : ["Astronomical cross-match"];

  return new ImageResponse(
    <OgCard
      eyebrow={catalog?.label}
      headline={id}
      lines={lines}
      footer="Position · photometry · light curves · counterparts"
      seed={id}
    />,
    size
  );
}
