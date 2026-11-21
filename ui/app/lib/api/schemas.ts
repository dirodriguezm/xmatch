/**
 * The XWave API contract: request/response shapes of the endpoints the
 * interface uses. Each proxy route under `app/api/` checks its upstream
 * response against one of these schemas before returning it, and every type
 * the UI works with is inferred from them, so there is no separate generated
 * type set to keep in sync.
 *
 * They follow what the service returns, not what the Swagger spec says: the
 * spec declares `/conesearch` and `/bulk-conesearch` as `MastercatResult[]`,
 * although the handlers return results grouped per catalog, and it types every
 * catalog's metadata as `Allwise` when each has its own columns.
 *
 * Objects are loose: unknown fields pass through untouched, so a new column
 * in the service never breaks the UI. Only the fields the UI reads are
 * checked.
 */

import { z } from "zod";

/** The service sends `null` for an empty list. */
const listOf = <T extends z.ZodType>(item: T) =>
  z
    .array(item)
    .nullish()
    .transform((d) => d ?? []);

/** One source in a cone search group. */
const coneMatchSchema = z.looseObject({
  id: z.string(),
  ra: z.number(),
  dec: z.number(),
  cat: z.string().optional(),
  /** Arcsec from the search position. */
  distance: z.number().optional(),
  ipix: z.number().optional(),
});

/** Matches for one catalog and one input position (`index`). */
const coneGroupSchema = z.looseObject({
  catalog: z.string(),
  data: listOf(coneMatchSchema),
  index: z.number().optional(),
});

/** /conesearch. */
export const coneSearchResponseSchema = z.array(coneGroupSchema);

/** /bulk-conesearch: `index` is the input row the group belongs to. */
const bulkConeGroupSchema = coneGroupSchema.extend({ index: z.number() });
export const bulkConeSearchResponseSchema = z.array(bulkConeGroupSchema);

/**
 * A catalog row. Each catalog has its own columns (gaia → `phot_*_mean_mag`,
 * allwise → `w1mpro`…), so only the id and the position are guaranteed.
 */
export const metadataRecordSchema = z.looseObject({
  id: z.string(),
  ra: z.number().optional(),
  dec: z.number().optional(),
});

export const bulkMetadataResponseSchema = z.array(metadataRecordSchema);

const detectionSchema = z.looseObject({
  catalog: z.string(),
  mjd: z.number(),
  mag: z.number(),
  magerr: z.number().nullish(),
});

/** Non-detections carry an upper limit instead of a magnitude. */
const epochSchema = z.looseObject({ mjd: z.number() });

/** /lightcurve. */
export const lightcurveResponseSchema = z.looseObject({
  detections: listOf(detectionSchema),
  non_detections: listOf(epochSchema),
  forced_photometry: listOf(epochSchema),
});

/** The service's 400 body, e.g. `{Field: "RA", Reason: "Could not parse float."}`. */
export const validationErrorSchema = z.object({
  Field: z.string(),
  Reason: z.string(),
  ErrValue: z.string().optional(),
});

export type ConeGroup = z.infer<typeof coneGroupSchema>;
export type BulkConeGroup = z.infer<typeof bulkConeGroupSchema>;
export type MetadataRecord = z.infer<typeof metadataRecordSchema>;
export type Lightcurve = z.infer<typeof lightcurveResponseSchema>;
