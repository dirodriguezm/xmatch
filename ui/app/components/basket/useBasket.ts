"use client";

import { useSyncExternalStore } from "react";

import { basketStore } from "@/app/lib/utils/basket";

import type { BasketItem } from "./types";

/** Live basket contents; empty during SSR and before hydration. */
export function useBasket(): BasketItem[] {
  return useSyncExternalStore(
    basketStore.subscribe,
    basketStore.getSnapshot,
    basketStore.getServerSnapshot
  );
}

export { basketStore };
