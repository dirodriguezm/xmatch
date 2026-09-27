"use client";

import { ReadOutlined } from "@ant-design/icons";
import { Button, Modal, Typography } from "antd";
import { useState } from "react";

import { resolveCatalogs } from "@/app/lib/utils/citation";

import { CitationTabs } from "./CitationTabs";

const { Text } = Typography;

/** Button that opens citation export for XWave plus these catalogs. */
export function CiteButton({
  catalogs,
  objectId,
}: {
  catalogs: string[];
  objectId?: string;
}) {
  const [open, setOpen] = useState(false);
  const names = resolveCatalogs(catalogs)
    .map((m) => `${m.name} ${m.release}`)
    .join(", ");

  return (
    <>
      <Button
        size="small"
        icon={<ReadOutlined />}
        onClick={() => setOpen(true)}
      >
        Cite
      </Button>
      <Modal
        title="Cite this data"
        open={open}
        onCancel={() => setOpen(false)}
        footer={null}
        width={760}
        destroyOnHidden
      >
        <Text className="!text-sm text-neutral-400">
          If {objectId ? <>data on {objectId}</> : "these results"} helped your
          work, please cite XWave and the catalogs used ({names}).
        </Text>
        <CitationTabs catalogs={catalogs} objectId={objectId} />
      </Modal>
    </>
  );
}
