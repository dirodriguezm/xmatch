"use client";

import {
  forwardRef,
  useCallback,
  useEffect,
  useImperativeHandle,
  useRef,
  useState,
} from "react";

import type {
  AladinCatalog,
  AladinCatalogSource,
  AladinGlobal,
  AladinInstance,
  AladinMarker,
  AladinViewerProps,
  AladinViewerRef,
} from "@/types/aladin";

const DEFAULT_OPTIONS = {
  // Use direct CDS URL to avoid CORS issues with mirror servers
  survey: "https://alasky.cds.unistra.fr/DSS/DSSColor/",
  fov: 0.5,
  projection: "SIN" as const,
  showCooGrid: false,
  showReticle: true,
  showFullscreenControl: false,
};

export const AladinViewerInner = forwardRef<AladinViewerRef, AladinViewerProps>(
  function AladinViewerInner(props, ref) {
    const {
      center,
      fov = DEFAULT_OPTIONS.fov,
      survey = DEFAULT_OPTIONS.survey,
      projection = DEFAULT_OPTIONS.projection,
      markers = [],
      catalogSources = [],
      catalogName = "Sources",
      catalogColor = "#1677ff",
      height = 400,
      width = "100%",
      className = "",
      showCooGrid = DEFAULT_OPTIONS.showCooGrid,
      showReticle = DEFAULT_OPTIONS.showReticle,
      showFullscreenControl = DEFAULT_OPTIONS.showFullscreenControl,
      onReady,
      onPositionChange,
      onFovChange,
      onObjectClick,
    } = props;

    const containerRef = useRef<HTMLDivElement>(null);
    const aladinRef = useRef<AladinInstance | null>(null);
    const aladinGlobalRef = useRef<AladinGlobal | null>(null);
    const catalogRef = useRef<AladinCatalog | null>(null);
    const markerCatalogRef = useRef<AladinCatalog | null>(null);
    const [isInitialized, setIsInitialized] = useState(false);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    // Initialize Aladin
    useEffect(() => {
      if (!containerRef.current || isInitialized) {
        return;
      }

      let mounted = true;

      const initAladin = async () => {
        try {
          // Dynamic import of aladin-lite
          const A = (await import("aladin-lite")).default as AladinGlobal;

          if (!mounted || !containerRef.current) return;

          aladinGlobalRef.current = A;

          // Wait for WASM initialization
          await A.init;

          if (!mounted || !containerRef.current) return;

          // Build target string from center coordinates
          const target = center
            ? `${center.ra.toFixed(6)} ${center.dec >= 0 ? "+" : ""}${center.dec.toFixed(6)}`
            : "0 +0";

          const aladin = A.aladin(containerRef.current, {
            survey,
            fov,
            target,
            projection,
            showCooGrid,
            showReticle,
            showFullscreenControl: false,
            showZoomControl: false,
            showSettingsControl: false,
            showShareControl: false,
            showLayersControl: false,
            showSimbadPointerControl: false,
            showCooLocation: false,
            showFov: false,
            showFrame: false,
            reticleColor: "#1677ff",
            backgroundColor: "#0a0a0a",
          });

          if (!mounted) return;

          aladinRef.current = aladin;
          setIsInitialized(true);
          setIsLoading(false);

          // Set up event listeners
          if (onPositionChange) {
            aladin.on("positionChanged", ((...args: unknown[]) => {
              const position = args[0] as { ra: number; dec: number };
              onPositionChange(position.ra, position.dec);
            }) as (...args: unknown[]) => void);
          }

          if (onFovChange) {
            aladin.on("zoomChanged", ((...args: unknown[]) => {
              const fovValues = args[0] as number[];
              onFovChange(fovValues[0]);
            }) as (...args: unknown[]) => void);
          }

          if (onObjectClick) {
            aladin.on("objectClicked", ((...args: unknown[]) => {
              const object = args[0] as {
                ra: number;
                dec: number;
                data?: Record<string, unknown>;
              } | null;
              if (object) {
                onObjectClick({
                  ra: object.ra,
                  dec: object.dec,
                  data: object.data,
                });
              }
            }) as (...args: unknown[]) => void);
          }

          if (onReady) {
            onReady(aladin);
          }
        } catch (err) {
          if (mounted) {
            console.error("Failed to initialize Aladin:", err);
            setError("Failed to load sky viewer");
            setIsLoading(false);
          }
        }
      };

      initAladin();

      return () => {
        mounted = false;
      };
    }, [
      isInitialized,
      center,
      fov,
      survey,
      projection,
      showCooGrid,
      showReticle,
      showFullscreenControl,
      onReady,
      onPositionChange,
      onFovChange,
      onObjectClick,
    ]);

    // Handle markers updates
    useEffect(() => {
      if (!isInitialized || !aladinRef.current || !aladinGlobalRef.current)
        return;

      const A = aladinGlobalRef.current;

      // Clear existing marker catalog
      if (markerCatalogRef.current) {
        aladinRef.current.removeOverlay(markerCatalogRef.current);
        markerCatalogRef.current = null;
      }

      if (markers.length > 0) {
        const markerCatalog = A.catalog({
          name: "Markers",
          color: "#ff4d4f",
          sourceSize: 15,
          shape: "circle",
        });

        const aladinMarkers = markers.map((m) =>
          A.marker(m.ra, m.dec, {
            color: m.color || "#ff4d4f",
            popupTitle: m.label || "",
            popupDesc: m.popup || "",
          })
        );

        markerCatalog.addSources(aladinMarkers);
        aladinRef.current.addOverlay(markerCatalog);
        markerCatalogRef.current = markerCatalog;
      }
    }, [markers, isInitialized]);

    // Handle catalog sources updates
    useEffect(() => {
      if (!isInitialized || !aladinRef.current || !aladinGlobalRef.current)
        return;

      const A = aladinGlobalRef.current;

      // Clear existing catalog
      if (catalogRef.current) {
        aladinRef.current.removeOverlay(catalogRef.current);
        catalogRef.current = null;
      }

      if (catalogSources.length > 0) {
        const catalog = A.catalog({
          name: catalogName,
          color: catalogColor,
          sourceSize: 12,
          shape: "plus",
          onClick: "showPopup",
        });

        const sources = catalogSources.map((s) =>
          A.source(s.ra, s.dec, {
            name: s.name,
            ...s.data,
          })
        );

        catalog.addSources(sources);
        aladinRef.current.addOverlay(catalog);
        catalogRef.current = catalog;
      }
    }, [catalogSources, catalogName, catalogColor, isInitialized]);

    // Imperative handle methods
    const goTo = useCallback((ra: number, dec: number, newFov?: number) => {
      if (aladinRef.current) {
        aladinRef.current.gotoRaDec(ra, dec);
        if (newFov !== undefined) {
          aladinRef.current.setFov(newFov);
        }
      }
    }, []);

    const addMarker = useCallback((marker: AladinMarker) => {
      const A = aladinGlobalRef.current;
      if (!A || !aladinRef.current || !markerCatalogRef.current) return;

      const aladinMarker = A.marker(marker.ra, marker.dec, {
        color: marker.color || "#ff4d4f",
        popupTitle: marker.label || "",
        popupDesc: marker.popup || "",
      });
      markerCatalogRef.current.addSources(aladinMarker);
    }, []);

    const addMarkers = useCallback((newMarkers: AladinMarker[]) => {
      const A = aladinGlobalRef.current;
      if (!A || !aladinRef.current || !markerCatalogRef.current) return;

      const aladinMarkers = newMarkers.map((m) =>
        A.marker(m.ra, m.dec, {
          color: m.color || "#ff4d4f",
          popupTitle: m.label || "",
          popupDesc: m.popup || "",
        })
      );
      markerCatalogRef.current.addSources(aladinMarkers);
    }, []);

    const clearMarkers = useCallback(() => {
      if (markerCatalogRef.current) {
        markerCatalogRef.current.removeAll();
      }
    }, []);

    const addCatalogSources = useCallback(
      (
        sources: AladinCatalogSource[],
        options?: { name?: string; color?: string }
      ) => {
        const A = aladinGlobalRef.current;
        if (!A || !aladinRef.current) return;

        const catalog = A.catalog({
          name: options?.name || "Custom",
          color: options?.color || "#1677ff",
          sourceSize: 12,
          shape: "plus",
        });

        const aladinSources = sources.map((s) =>
          A.source(s.ra, s.dec, { name: s.name, ...s.data })
        );
        catalog.addSources(aladinSources);
        aladinRef.current.addOverlay(catalog);
      },
      []
    );

    const clearCatalogs = useCallback(() => {
      if (aladinRef.current && catalogRef.current) {
        aladinRef.current.removeOverlay(catalogRef.current);
        catalogRef.current = null;
      }
    }, []);

    const setSurvey = useCallback((newSurvey: string) => {
      if (aladinRef.current) {
        aladinRef.current.setImageSurvey(newSurvey);
      }
    }, []);

    const getCenter = useCallback(() => {
      if (aladinRef.current) {
        const [ra, dec] = aladinRef.current.getRaDec();
        return { ra, dec };
      }
      return { ra: 0, dec: 0 };
    }, []);

    const getFov = useCallback(() => {
      if (aladinRef.current) {
        const fovValues = aladinRef.current.getFov();
        return fovValues[0];
      }
      return 0;
    }, []);

    const getAladinInstance = useCallback(() => aladinRef.current, []);

    useImperativeHandle(ref, () => ({
      goTo,
      addMarker,
      addMarkers,
      clearMarkers,
      addCatalogSources,
      clearCatalogs,
      setSurvey,
      getCenter,
      getFov,
      getAladinInstance,
    }));

    // Cleanup on unmount
    useEffect(() => {
      return () => {
        aladinRef.current = null;
        aladinGlobalRef.current = null;
        catalogRef.current = null;
        markerCatalogRef.current = null;
      };
    }, []);

    // Always render the container so Aladin can initialize
    // Aladin requires explicit inline pixel dimensions to work correctly
    return (
      <div
        className={`relative ${className}`}
        // eslint-disable-next-line react/forbid-dom-props
        style={{
          height: typeof height === "number" ? `${height}px` : height,
          width: typeof width === "number" ? `${width}px` : width,
        }}
      >
        {/* Container for Aladin - needs explicit inline dimensions */}
        <div
          ref={containerRef}
          className="aladin-container"
          // eslint-disable-next-line react/forbid-dom-props
          style={{
            position: "absolute",
            top: 0,
            left: 0,
            right: 0,
            bottom: 0,
            visibility: isLoading || error ? "hidden" : "visible",
          }}
        />

        {/* Loading overlay */}
        {isLoading && (
          <div className="absolute inset-0 flex items-center justify-center bg-surface">
            <div className="text-center">
              <div className="animate-spin h-6 w-6 border-2 border-primary border-t-transparent rounded-full mx-auto mb-2" />
              <span className="text-xs text-foreground/60">Loading...</span>
            </div>
          </div>
        )}

        {/* Error state */}
        {error && (
          <div className="absolute inset-0 flex items-center justify-center bg-surface">
            <span className="text-xs text-red-500">{error}</span>
          </div>
        )}
      </div>
    );
  }
);
