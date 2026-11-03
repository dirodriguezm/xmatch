declare module "aladin-lite" {
  interface AladinGlobal {
    init: Promise<void>;
    aladin: (
      element: string | HTMLElement,
      options?: Record<string, unknown>
    ) => unknown;
    catalog: (options?: Record<string, unknown>) => unknown;
    source: (
      ra: number,
      dec: number,
      data?: Record<string, unknown>
    ) => unknown;
    marker: (
      ra: number,
      dec: number,
      options?: Record<string, unknown>
    ) => unknown;
    circle: (
      ra: number,
      dec: number,
      radius: number,
      options?: Record<string, unknown>
    ) => unknown;
    polyline: (
      coords: [number, number][],
      options?: Record<string, unknown>
    ) => unknown;
  }

  const A: AladinGlobal;
  export default A;
}
