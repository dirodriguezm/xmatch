import type { Metadata } from "next";

import { SITE_NAME, SITE_URL } from "@/app/lib/constants/site";
import {
  catalogFromObjectId,
  fetchObjectSummary,
  formatRaDec,
} from "@/app/lib/utils/share";

interface ObjectLayoutProps {
  children: React.ReactNode;
  params: Promise<{ objectId: string }>;
}

export async function generateMetadata({
  params,
}: Pick<ObjectLayoutProps, "params">): Promise<Metadata> {
  const { objectId } = await params;
  const id = decodeURIComponent(objectId);
  const catalog = catalogFromObjectId(id);

  // Layouts don't see ?catalog=, so only ids with a recognisable prefix get
  // a position lookup; it's best-effort and short so crawlers aren't stalled.
  const summary = catalog
    ? await fetchObjectSummary(id, catalog.slug, 1500)
    : null;

  const parts = [catalog ? `${catalog.label} source` : "Catalog source"];
  if (summary) parts.push(`at ${formatRaDec(summary.ra, summary.dec)}`);
  const description = `${parts.join(" ")}: position, photometry, light curves and cross-matched counterparts on ${SITE_NAME}.`;
  const title = `${id} | ${SITE_NAME}`;

  return {
    metadataBase: new URL(SITE_URL),
    title: id,
    description,
    openGraph: {
      type: "website",
      siteName: SITE_NAME,
      title,
      description,
      url: `/object/${encodeURIComponent(id)}`,
    },
    twitter: {
      card: "summary_large_image",
      title,
      description,
    },
  };
}

export default function ObjectLayout({ children }: ObjectLayoutProps) {
  return children;
}
