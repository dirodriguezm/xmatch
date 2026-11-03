import { buildChangelogRss, CHANGELOG } from "@/app/lib/constants/changelog";
import { SITE_URL } from "@/app/lib/constants/site";

export const dynamic = "force-static";

export function GET() {
  return new Response(buildChangelogRss(CHANGELOG, SITE_URL), {
    headers: {
      "Content-Type": "application/rss+xml; charset=utf-8",
      "Cache-Control": "public, max-age=3600",
    },
  });
}
