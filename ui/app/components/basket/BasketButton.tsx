"use client";

import { StarFilled, StarOutlined } from "@ant-design/icons";
import { Button, Tooltip } from "antd";
import type { MouseEvent } from "react";

import { hasItem } from "@/app/lib/utils/basket";

import type { BasketItem } from "./types";
import { basketStore, useBasket } from "./useBasket";

/** Adds/removes one object from the basket. */
export function BasketButton({
  item,
  size = "middle",
}: {
  item: BasketItem;
  size?: "small" | "middle";
}) {
  const items = useBasket();
  const inBasket = hasItem(items, item);
  const label = inBasket ? "Remove from basket" : "Add to basket";

  const onClick = (e: MouseEvent) => {
    // Rows in the results table navigate on click; don't let this bubble.
    e.stopPropagation();
    basketStore.toggle(item);
  };

  return (
    <Tooltip title={label}>
      <Button
        type="text"
        size={size}
        aria-label={label}
        aria-pressed={inBasket}
        onClick={onClick}
        icon={
          inBasket ? (
            <StarFilled className="text-amber-400" />
          ) : (
            <StarOutlined />
          )
        }
      />
    </Tooltip>
  );
}
