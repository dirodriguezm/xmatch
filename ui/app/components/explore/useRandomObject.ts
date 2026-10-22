"use client";

import { App } from "antd";
import { useRouter } from "next/navigation";
import { useCallback, useState } from "react";

import { buildObjectUrl } from "@/app/lib/utils/urls";

import { findRandomGaiaObject, pickRandomFeatured } from "./randomObject";

export type RandomMode = "sky" | "featured";

/** Navigate to a random object: a random sky position, or a featured target. */
export function useRandomObject() {
  const router = useRouter();
  const { message } = App.useApp();
  const [loading, setLoading] = useState(false);

  const go = useCallback(
    async (mode: RandomMode = "sky") => {
      if (mode === "featured") {
        const f = pickRandomFeatured();
        router.push(buildObjectUrl(f.objectId, f.catalog));
        return;
      }
      setLoading(true);
      const hit = await findRandomGaiaObject();
      if (hit) {
        router.push(buildObjectUrl(hit.objectId, hit.catalog));
        // Keep the spinner until the route changes.
        return;
      }
      setLoading(false);
      const f = pickRandomFeatured();
      message.info(`Empty patch of sky — here is ${f.name} instead`);
      router.push(buildObjectUrl(f.objectId, f.catalog));
    },
    [router, message]
  );

  return { go, loading };
}
