"use client";

import {
  CopyOutlined,
  DeleteOutlined,
  DownloadOutlined,
  NodeIndexOutlined,
  ShoppingOutlined,
} from "@ant-design/icons";
import {
  App,
  Badge,
  Button,
  Drawer,
  Empty,
  Flex,
  Popconfirm,
  Tag,
  Tooltip,
  Typography,
} from "antd";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useState } from "react";

import {
  getSearchCatalogColor,
  getSearchCatalogLabel,
} from "@/app/lib/constants/catalogs";
import {
  basketKey,
  basketToBulkInput,
  basketToCsv,
  basketToText,
  BULK_INPUT_STORAGE_KEY,
} from "@/app/lib/utils/basket";
import { downloadCsv } from "@/app/lib/utils/csv";
import { buildObjectUrl } from "@/app/lib/utils/urls";

import { basketStore, useBasket } from "./useBasket";

const { Text } = Typography;

/** Header basket badge + drawer. */
export function HeaderBasket() {
  const items = useBasket();
  const [open, setOpen] = useState(false);
  const router = useRouter();
  const { message } = App.useApp();

  const crossMatchAll = () => {
    try {
      window.sessionStorage.setItem(
        BULK_INPUT_STORAGE_KEY,
        basketToBulkInput(items)
      );
    } catch {
      message.error("Could not hand the basket to the bulk page");
      return;
    }
    setOpen(false);
    router.push("/bulk?from=basket");
  };

  const copyList = async () => {
    try {
      await navigator.clipboard.writeText(basketToText(items));
      message.success(`Copied ${items.length} objects`);
    } catch {
      message.error("Clipboard unavailable");
    }
  };

  const empty = items.length === 0;

  return (
    <>
      <Tooltip title="Object basket">
        <Badge count={items.length} size="small" offset={[-4, 4]}>
          <Button
            type="text"
            aria-label={`Object basket (${items.length})`}
            icon={<ShoppingOutlined />}
            onClick={() => setOpen(true)}
          />
        </Badge>
      </Tooltip>
      <Drawer
        title={`Basket (${items.length})`}
        open={open}
        onClose={() => setOpen(false)}
        placement="right"
        footer={
          <Flex wrap gap={8}>
            <Button
              type="primary"
              icon={<NodeIndexOutlined />}
              disabled={empty}
              onClick={crossMatchAll}
            >
              Cross-match all
            </Button>
            <Button
              icon={<DownloadOutlined />}
              disabled={empty}
              onClick={() =>
                downloadCsv("xwave_basket.csv", basketToCsv(items))
              }
            >
              CSV
            </Button>
            <Button icon={<CopyOutlined />} disabled={empty} onClick={copyList}>
              Copy
            </Button>
            <Popconfirm
              title="Empty the basket?"
              onConfirm={() => basketStore.clear()}
              okText="Clear"
              disabled={empty}
            >
              <Button danger type="text" disabled={empty}>
                Clear
              </Button>
            </Popconfirm>
          </Flex>
        }
      >
        {empty ? (
          <Empty
            image={Empty.PRESENTED_IMAGE_SIMPLE}
            description={
              <span className="text-neutral-400">
                Star objects in search results or on an object page to collect
                them here.
              </span>
            }
          />
        ) : (
          <ul className="m-0 p-0 list-none flex flex-col gap-2">
            {items.map((item) => (
              <li
                key={basketKey(item)}
                className="flex items-center gap-2 rounded-md border border-border bg-surface px-3 py-2"
              >
                <div className="min-w-0 flex-1">
                  <Link
                    href={buildObjectUrl(item.objectId, item.catalog)}
                    onClick={() => setOpen(false)}
                    className="block truncate text-foreground hover:underline"
                  >
                    {item.objectId}
                  </Link>
                  <Flex align="center" gap={6} className="mt-1">
                    <Tag
                      color={getSearchCatalogColor(item.catalog)}
                      className="m-0!"
                    >
                      {getSearchCatalogLabel(item.catalog)}
                    </Tag>
                    <Text type="secondary" className="font-mono text-xs">
                      {item.ra.toFixed(5)}, {item.dec.toFixed(5)}
                    </Text>
                  </Flex>
                </div>
                <Tooltip title="Remove">
                  <Button
                    type="text"
                    size="small"
                    aria-label={`Remove ${item.objectId}`}
                    icon={<DeleteOutlined />}
                    onClick={() => basketStore.remove(item)}
                  />
                </Tooltip>
              </li>
            ))}
          </ul>
        )}
      </Drawer>
    </>
  );
}
