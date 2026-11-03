import { useQuery } from "@tanstack/react-query";

/** How long the sky view waits for the image server before falling back. */
const PROBE_TIMEOUT_MS = 12_000;

/**
 * Whether the HiPS server behind `surveyUrl` answers. Aladin shows a spinner
 * forever when it doesn't, so the page asks first. `no-cors` resolves on any
 * answer and rejects only when the server is unreachable or times out.
 */
export function useSkyImagesReachable(surveyUrl: string) {
  return useQuery({
    queryKey: ["sky-images-reachable", surveyUrl],
    queryFn: async () => {
      await fetch(`${surveyUrl}properties`, {
        mode: "no-cors",
        signal: AbortSignal.timeout(PROBE_TIMEOUT_MS),
      });
      return true;
    },
    staleTime: 5 * 60 * 1000,
    retry: false,
  });
}
