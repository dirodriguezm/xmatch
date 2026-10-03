"use client";

import { parseAsString, useQueryStates } from "nuqs";

const searchParamsParsers = {
  ra: parseAsString.withDefault(""),
  dec: parseAsString.withDefault(""),
  catalogRadii: parseAsString.withDefault(""),
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
    catalogRadii: params.catalogRadii,

    // Setters
    setRa: (value: string) => setParams({ ra: value }),
    setDec: (value: string) => setParams({ dec: value }),
    setCatalogRadii: (value: string) => setParams({ catalogRadii: value }),

    // Bulk update
    setSearchParams: setParams,

    // Check if params are set (for enabling search button)
    isValid: Boolean(params.ra && params.dec),
  };
}
