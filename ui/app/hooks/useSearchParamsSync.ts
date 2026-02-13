"use client";

import { parseAsFloat, parseAsString, useQueryStates } from "nuqs";

const searchParamsParsers = {
  ra: parseAsString.withDefault(""),
  dec: parseAsString.withDefault(""),
  radius: parseAsFloat.withDefault(1),
  unit: parseAsString.withDefault("arcsec"),
  catalogs: parseAsString.withDefault(""),
};

export function useSearchParams() {
  const [params, setParams] = useQueryStates(searchParamsParsers, {
    history: "replace",
    shallow: false,
  });

  return {
    // Values
    ra: params.ra,
    dec: params.dec,
    radius: params.radius,
    unit: params.unit as "arcsec" | "arcmin" | "deg",
    catalogs: params.catalogs,

    // Setters
    setRa: (value: string) => setParams({ ra: value }),
    setDec: (value: string) => setParams({ dec: value }),
    setRadius: (value: number) => setParams({ radius: value }),
    setUnit: (value: "arcsec" | "arcmin" | "deg") => setParams({ unit: value }),
    setCatalogs: (value: string) => setParams({ catalogs: value }),

    // Bulk update
    setSearchParams: setParams,

    // Check if params are set (for enabling search button)
    isValid: Boolean(params.ra && params.dec),
  };
}
