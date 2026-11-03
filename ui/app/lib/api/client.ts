import { z } from "zod";

import { validationErrorSchema } from "./schemas";

export const API_BASE_URL = "https://xwave-astro.udp.cl/v1";

export class ApiError extends Error {
  constructor(
    public status: number,
    message: string
  ) {
    super(message);
    this.name = "ApiError";
  }
}

/**
 * JSON.parse with integers of 16+ digits kept as strings. Catalogue ids
 * (Gaia source_id, AllWISE cntr) are 64-bit and lose their last digits as
 * JavaScript numbers: 3700386905605055360 would read as …055500. String
 * literals are skipped whole, and floats never match, so only bare integers
 * change.
 */
export function parseJsonKeepingBigInts(text: string): unknown {
  return JSON.parse(
    text.replace(
      // A whole string literal (left as is), or a bare big integer.
      /"(?:[^"\\]|\\.)*"|(?<![\d.eE+-])(-?\d{16,})(?=\s*[,\]}])/g,
      (match, bigInt?: string) => (bigInt ? `"${bigInt}"` : match)
    )
  );
}

/** Turns the service's `{Field, Reason}` body into "RA: Could not parse float." */
function errorMessage(text: string, fallback: string): string {
  try {
    const parsed = validationErrorSchema.safeParse(JSON.parse(text));
    if (parsed.success) return `${parsed.data.Field}: ${parsed.data.Reason}`;
  } catch {
    // Not JSON: the service answered with plain text.
  }
  return text || fallback;
}

/**
 * Calls the XWave API and checks the body against `schema`. A body that
 * doesn't match throws a 502 naming the endpoint and the offending field,
 * instead of handing the UI data with holes in it.
 */
export async function apiFetch<S extends z.ZodType>(
  endpoint: string,
  schema: S,
  options?: RequestInit
): Promise<z.output<S> | null> {
  const url = `${API_BASE_URL}${endpoint}`;
  const response = await fetch(url, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      ...options?.headers,
    },
  });

  if (response.status === 204) {
    return null;
  }

  if (!response.ok) {
    const errorText = await response.text();
    throw new ApiError(
      response.status,
      errorMessage(errorText, response.statusText)
    );
  }

  const parsed = schema.safeParse(
    parseJsonKeepingBigInts(await response.text())
  );
  if (!parsed.success) {
    const path = endpoint.split("?")[0];
    throw new ApiError(
      502,
      `Unexpected response from ${path}: ${z.prettifyError(parsed.error)}`
    );
  }
  return parsed.data;
}
