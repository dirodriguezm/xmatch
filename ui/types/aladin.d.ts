/**
 * TypeScript declarations for Aladin Lite v3 API
 * @see https://aladin.cds.unistra.fr/AladinLite/doc/API/
 */

declare global {
  interface Window {
    A?: AladinGlobal;
  }
}

export interface AladinGlobal {
  /** Initialize Aladin Lite (loads WASM) */
  init: Promise<void>;
  /** Create an Aladin Lite instance */
  aladin: (
    element: string | HTMLElement,
    options?: AladinOptions
  ) => AladinInstance;
  /** Create a catalog */
  catalog: (options?: CatalogOptions) => AladinCatalog;
  /** Create a source */
  source: (
    ra: number,
    dec: number,
    data?: Record<string, unknown>
  ) => AladinSource;
  /** Create a marker */
  marker: (
    ra: number,
    dec: number,
    options?: MarkerOptions
  ) => AladinMarkerInstance;
  /** Create a circle overlay */
  circle: (
    ra: number,
    dec: number,
    radius: number,
    options?: ShapeOptions
  ) => AladinShape;
  /** Create a polyline overlay */
  polyline: (coords: [number, number][], options?: ShapeOptions) => AladinShape;
}

export interface AladinOptions {
  /** HiPS survey URL or ID */
  survey?: string;
  /** Initial field of view in degrees */
  fov?: number;
  /** Initial target (coordinates or object name) */
  target?: string;
  /** Show coordinate grid */
  showCooGrid?: boolean;
  /** Show zoom control */
  showZoomControl?: boolean;
  /** Show fullscreen control */
  showFullscreenControl?: boolean;
  /** Show settings control */
  showSettingsControl?: boolean;
  /** Show share control */
  showShareControl?: boolean;
  /** Show layers control */
  showLayersControl?: boolean;
  /** Show SIMBAD pointer control */
  showSimbadPointerControl?: boolean;
  /** Show reticle at center */
  showReticle?: boolean;
  /** Reticle color */
  reticleColor?: string;
  /** Reticle size */
  reticleSize?: number;
  /** Map projection */
  projection?: "SIN" | "AIT" | "ZEA" | "MOL" | "TAN";
  /** Coordinate frame */
  cooFrame?: "J2000" | "J2000d" | "Galactic";
  /** Background color */
  backgroundColor?: string;
  /** Show coordinate location display */
  showCooLocation?: boolean;
  /** Show field of view display */
  showFov?: boolean;
  /** Show frame selector */
  showFrame?: boolean;
}

export interface AladinInstance {
  /** Navigate to coordinates */
  gotoRaDec: (ra: number, dec: number) => void;
  /** Navigate to object by name (SIMBAD resolver) */
  gotoObject: (name: string, options?: { error?: (e: Error) => void }) => void;
  /** Set field of view */
  setFov: (fov: number) => void;
  /** Get current field of view */
  getFov: () => number[];
  /** Get current RA/Dec */
  getRaDec: () => [number, number];
  /** Set survey by URL or ID */
  setImageSurvey: (survey: string) => void;
  /** Add overlay (catalog, shapes) */
  addOverlay: (overlay: AladinCatalog | AladinShape) => void;
  /** Remove overlay */
  removeOverlay: (overlay: AladinCatalog | AladinShape) => void;
  /** Show/hide coordinate grid */
  showCooGrid: (show: boolean) => void;
  /** Show/hide reticle */
  showReticle: (show: boolean) => void;
  /** Set projection */
  setProjection: (projection: string) => void;
  /** Add event listener */
  on: (event: AladinEvent, callback: (...args: unknown[]) => void) => void;
  /** World to pixel coordinates */
  world2pix: (ra: number, dec: number) => [number, number] | null;
  /** Pixel to world coordinates */
  pix2world: (x: number, y: number) => [number, number] | null;
  /** Get view HTML element */
  getViewDiv: () => HTMLElement;
}

export type AladinEvent =
  | "positionChanged"
  | "zoomChanged"
  | "objectClicked"
  | "objectHovered"
  | "click"
  | "rightClickContextMenu";

export interface CatalogOptions {
  name?: string;
  color?: string;
  sourceSize?: number;
  shape?: "circle" | "plus" | "rhomb" | "cross" | "triangle" | "square";
  onClick?: "showTable" | "showPopup" | ((source: AladinSource) => void);
}

export interface AladinCatalog {
  addSources: (sources: AladinSource | AladinSource[]) => void;
  removeAll: () => void;
  show: () => void;
  hide: () => void;
  isShown: boolean;
  name: string;
}

export interface AladinSource {
  ra: number;
  dec: number;
  data?: Record<string, unknown>;
}

export interface MarkerOptions {
  color?: string;
  popupTitle?: string;
  popupDesc?: string;
  useMarkerDefaultIcon?: boolean;
}

export interface AladinMarkerInstance extends AladinSource {
  select: () => void;
  deselect: () => void;
}

export interface ShapeOptions {
  color?: string;
  lineWidth?: number;
  fillColor?: string;
}

export interface AladinShape {
  show: () => void;
  hide: () => void;
  isShown: boolean;
}

// Component Props and Ref Types

export interface AladinMarker {
  ra: number;
  dec: number;
  label?: string;
  color?: string;
  popup?: string;
}

export interface AladinCatalogSource {
  ra: number;
  dec: number;
  name?: string;
  data?: Record<string, unknown>;
}

export interface AladinViewerProps {
  /** Initial center coordinates */
  center?: { ra: number; dec: number };
  /** Initial field of view in degrees */
  fov?: number;
  /** HiPS survey to display */
  survey?: string;
  /** Map projection */
  projection?: "SIN" | "AIT" | "ZEA" | "MOL" | "TAN";
  /** Markers to display */
  markers?: AladinMarker[];
  /** Catalog sources to overlay */
  catalogSources?: AladinCatalogSource[];
  /** Catalog name */
  catalogName?: string;
  /** Catalog source color */
  catalogColor?: string;
  /** Container height */
  height?: string | number;
  /** Container width */
  width?: string | number;
  /** Custom CSS class */
  className?: string;
  /** Show coordinate grid */
  showCooGrid?: boolean;
  /** Show reticle at center */
  showReticle?: boolean;
  /** Show fullscreen button */
  showFullscreenControl?: boolean;
  /** Callback when viewer is ready */
  onReady?: (aladin: AladinInstance) => void;
  /** Callback when position changes */
  onPositionChange?: (ra: number, dec: number) => void;
  /** Callback when FOV changes */
  onFovChange?: (fov: number) => void;
  /** Callback when object is clicked */
  onObjectClick?: (object: AladinCatalogSource) => void;
}

export interface AladinViewerRef {
  /** Navigate to coordinates */
  goTo: (ra: number, dec: number, fov?: number) => void;
  /** Add a single marker */
  addMarker: (marker: AladinMarker) => void;
  /** Add multiple markers */
  addMarkers: (markers: AladinMarker[]) => void;
  /** Clear all markers */
  clearMarkers: () => void;
  /** Add catalog sources */
  addCatalogSources: (
    sources: AladinCatalogSource[],
    options?: { name?: string; color?: string }
  ) => void;
  /** Clear all catalog overlays */
  clearCatalogs: () => void;
  /** Change the HiPS survey */
  setSurvey: (survey: string) => void;
  /** Get current center position */
  getCenter: () => { ra: number; dec: number };
  /** Get current field of view */
  getFov: () => number;
  /** Get the raw Aladin instance */
  getAladinInstance: () => AladinInstance | null;
}

export {};
