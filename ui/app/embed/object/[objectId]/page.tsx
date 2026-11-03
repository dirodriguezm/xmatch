import type { Metadata } from "next";

import { catalogFromObjectId } from "@/app/lib/utils/share";

import { EmbedObjectCard } from "./EmbedObjectCard";

interface EmbedObjectPageProps {
  params: Promise<{ objectId: string }>;
  searchParams: Promise<{ catalog?: string | string[] }>;
}

export async function generateMetadata({
  params,
}: Pick<EmbedObjectPageProps, "params">): Promise<Metadata> {
  const { objectId } = await params;
  return { title: decodeURIComponent(objectId) };
}

export default async function EmbedObjectPage({
  params,
  searchParams,
}: EmbedObjectPageProps) {
  const { objectId } = await params;
  const { catalog: catalogParam } = await searchParams;
  const id = decodeURIComponent(objectId);
  const catalog =
    (Array.isArray(catalogParam) ? catalogParam[0] : catalogParam) ??
    catalogFromObjectId(id)?.slug ??
    null;

  return <EmbedObjectCard objectId={id} catalog={catalog} />;
}
