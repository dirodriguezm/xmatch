import type { Metadata } from "next";

export const metadata: Metadata = {
  // Embeds duplicate object pages; keep them out of search results.
  robots: { index: false, follow: true },
};

/** Bare layout for iframe embeds: no site header or footer. */
export default function EmbedLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return <div className="min-h-screen bg-background p-3">{children}</div>;
}
