import {
  CATALOG_COLOR_CLASSES,
  getSearchCatalogLabel,
} from "@/app/lib/constants/catalogs";

interface CatalogLabelProps {
  /** Search-catalog slug: `gaia`, `allwise`, `erosita`. */
  catalog: string;
  className?: string;
}

/**
 * Catalog name with its colour dot — the same mark the search form uses, so
 * a catalog looks the same everywhere instead of as a light antd Tag.
 */
export function CatalogLabel({ catalog, className = "" }: CatalogLabelProps) {
  return (
    <span
      className={`inline-flex items-center gap-1.5 whitespace-nowrap text-sm text-foreground ${className}`}
    >
      <span
        aria-hidden
        className={`inline-block w-2 h-2 rounded-full shrink-0 ${
          CATALOG_COLOR_CLASSES[catalog.toLowerCase()] ?? "bg-neutral-500"
        }`}
      />
      {getSearchCatalogLabel(catalog)}
    </span>
  );
}
