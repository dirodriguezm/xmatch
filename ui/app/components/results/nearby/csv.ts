import { toCsv } from "@/app/lib/utils/csv";

import type { NearbySource } from "./shared";

/** CSV of the matches, as offered by the results toolbar. */
export function nearbyCsv(sources: NearbySource[]) {
  return toCsv(
    [
      "object_id",
      "catalog",
      "ra_deg",
      "dec_deg",
      "angular_distance_arcsec",
      "position_angle_deg",
      "sigma_ratio",
      "mag",
      "band",
    ],
    sources.map((s) => [
      s.objectId,
      s.slug,
      s.ra,
      s.dec,
      s.sepArcsec,
      Number(s.pa.toFixed(1)),
      Number((s.sepArcsec / s.sigma).toFixed(2)),
      s.photometry?.mag,
      s.photometry?.band,
    ])
  );
}
