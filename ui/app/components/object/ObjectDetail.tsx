"use client";

import {
  CopyOutlined,
  DatabaseOutlined,
  DownloadOutlined,
  EnvironmentOutlined,
  FieldTimeOutlined,
  LineChartOutlined,
  RadarChartOutlined,
  StarOutlined,
} from "@ant-design/icons";
import {
  App,
  Button,
  Card,
  Col,
  Collapse,
  Descriptions,
  Empty,
  Flex,
  Row,
  Select,
  Space,
  Tooltip,
  Typography,
} from "antd";
import dynamic from "next/dynamic";
import Link from "next/link";
import { type ReactNode, useRef, useState } from "react";

import { ObjectActions } from "@/app/components/actions/ObjectActions";
import type { CrossmatchResult } from "@/app/components/results/ResultsTable";
import {
  type Counterpart,
  COUNTERPART_RADIUS_ARCSEC,
  type GaiaEpochParams,
  NEIGHBOR_RADIUS_ARCSEC,
  useCounterparts,
  useDesiSpectrum,
  useDesiTarget,
  useGaiaEpochPhotometry,
  useGalacticReddening,
  useLightcurve,
  useNeighbors,
  usePs1Lightcurve,
  useVizierSed,
  useZtfLightcurve,
} from "@/app/hooks/queries";
import { PHOTOMETRY_BANDS } from "@/app/lib/constants/bands";
import { describeCatalogField } from "@/app/lib/constants/catalogFields";
import {
  CATALOG_COLOR_CLASSES,
  getSearchCatalogLabel,
} from "@/app/lib/constants/catalogs";
import {
  equatorialToEcliptic,
  equatorialToGalactic,
} from "@/app/lib/utils/coordinates";
import { calculateAxisBounds } from "@/app/lib/utils/data";
import { gaiaSourceIdFromDesignation } from "@/app/lib/utils/gaiaEpoch";
import {
  detectionPointsToCsv,
  downloadCsv,
  getCatalogLabel,
  getMagSystem,
  groupDetectionsByCatalog,
} from "@/app/lib/utils/lightcurve";
import {
  buildSedPoints,
  type PhotometrySource,
  type SedPoint,
} from "@/app/lib/utils/sed";
import { buildObjectUrl } from "@/app/lib/utils/urls";
import type { AladinViewerRef } from "@/types/aladin";
import type { components } from "@/types/xwave-api";

import { AladinViewer } from "./AladinViewer";
import { LightCurveChart } from "./LightCurveChart";
import { LightCurveSkeleton } from "./LightCurveSkeleton";
import { NearbySources } from "./NearbySources";
import { ObjectArchives } from "./ObjectArchives";
import { SedChart } from "./SedChart";
import { SimbadIdentity } from "./SimbadIdentity";
import { SpectrumChart } from "./SpectrumChart";

// Split out so astronomy-engine only loads when the panel is first opened.
const ObservabilityPanel = dynamic(
  () => import("./ObservabilityPanel").then((m) => m.ObservabilityPanel),
  { ssr: false }
);

const DSS_SURVEY = "https://alasky.cds.unistra.fr/DSS/DSSColor/";

const SURVEY_OPTIONS = [
  {
    label: "DSS Optical",
    value: DSS_SURVEY,
    category: "Optical",
    description: "Digitized photographic plates (POSS/UKST), whole sky",
  },
  {
    label: "DESI DR10",
    value: "CDS/P/DESI-Legacy-Surveys/DR10/color",
    category: "Optical",
    description: "Deep Legacy Surveys g/r/z imaging, partial sky",
  },
  {
    label: "DSS2 Color",
    value: "CDS/P/DSS2/color",
    category: "Optical",
    description: "Second-epoch DSS plates, red/blue color composite",
  },
  {
    label: "2MASS",
    value: "CDS/P/2MASS/color",
    category: "Infrared",
    description: "Near-infrared J/H/Ks (1.2–2.2 µm), whole sky",
  },
  {
    label: "AllWISE",
    value: "CDS/P/allWISE/color",
    category: "Infrared",
    description: "Mid-infrared W1–W4 (3.4–22 µm), whole sky",
  },
  {
    label: "XMM-Newton",
    value: "xcatdb/P/XMM/PN/color",
    category: "X-ray",
    description: "EPIC-pn 0.2–12 keV, pointed fields only",
  },
  {
    label: "Chandra",
    value: "cxc.harvard.edu/P/cda/hips/allsky/rgb",
    category: "X-ray",
    description: "ACIS 0.5–7 keV, arcsecond resolution, pointed fields only",
  },
  {
    label: "NVSS 1.4 GHz",
    value: "CDS/P/NVSS",
    category: "Radio",
    description: "VLA 1.4 GHz, 45″ beam, dec > −40°",
  },
  {
    label: "SUMSS 843 MHz",
    value: "CDS/P/SUMSS",
    category: "Radio",
    description: "Molonglo 843 MHz, 45″ beam, dec < −30°",
  },
  {
    label: "RACS 887 MHz",
    value: "https://casda.csiro.au/hips/RACS/low/I/",
    category: "Radio",
    description: "ASKAP 887.5 MHz, ~15–25″ beam, dec < +41°",
  },
  {
    label: "RACS-mid 1.4 GHz",
    value: "https://casda.csiro.au/hips/RACSmidb_I1/",
    category: "Radio",
    description: "ASKAP 1367.5 MHz, ~10″ beam, dec < +49°",
  },
  {
    label: "VLASS 3 GHz",
    value: "https://vlass-dl.nrao.edu/vlass/HiPS/MedianStack/Quicklook/",
    category: "Radio",
    description: "VLA 2–4 GHz, 2.5″ beam, dec > −40°",
  },
];

type SurveyOption = { label: string; value: string };

const SURVEY_DESCRIPTIONS = new Map(
  SURVEY_OPTIONS.map((s) => [s.value, s.description])
);

const surveySelectOptions = Object.entries(
  SURVEY_OPTIONS.reduce<Record<string, SurveyOption[]>>((acc, s) => {
    (acc[s.category] ??= []).push({ label: s.label, value: s.value });
    return acc;
  }, {})
).map(([category, options]) => ({ label: category, options }));

const { Text, Title } = Typography;

type Allwise = components["schemas"]["repository.Allwise"];

interface ObjectDetailProps {
  object: CrossmatchResult;
  metadata?: Allwise | null;
}

function formatSigned(value: number, decimals: number): string {
  return `${value >= 0 ? "+" : "−"}${Math.abs(value).toFixed(decimals)}`;
}

function counterpartStatus(c: Counterpart): string {
  switch (c.status) {
    case "loading":
      return "checking…";
    case "error":
      return "lookup failed";
    case "none":
      return `none within ${c.radiusArcsec}″`;
    case "found":
      return c.separationArcsec != null
        ? `${c.separationArcsec.toFixed(1)}″`
        : "found";
  }
}

function formatChipTooltip(p: SedPoint | undefined, survey: string): string {
  if (!p) return `${survey} — no measurement`;
  const origin = p.isSelf
    ? "this object"
    : `${getSearchCatalogLabel(p.catalog)} counterpart at ${p.separationArcsec.toFixed(1)}″`;
  return `${survey}${p.upperLimit ? " upper limit" : ""} · ${origin}`;
}

// Convert decimal degrees to sexagesimal
function toHMS(ra: number): string {
  const hours = ra / 15;
  const h = Math.floor(hours);
  const m = Math.floor((hours - h) * 60);
  const s = ((hours - h) * 60 - m) * 60;
  return `${h.toString().padStart(2, "0")}h ${m.toString().padStart(2, "0")}m ${s.toFixed(2)}s`;
}

function toDMS(dec: number): string {
  const sign = dec >= 0 ? "+" : "-";
  const absDec = Math.abs(dec);
  const d = Math.floor(absDec);
  const m = Math.floor((absDec - d) * 60);
  const s = ((absDec - d) * 60 - m) * 60;
  return `${sign}${d}° ${m.toString().padStart(2, "0")}′ ${s.toFixed(2)}″`;
}

export function ObjectDetail({ object, metadata }: ObjectDetailProps) {
  const { message } = App.useApp();
  const aladinRef = useRef<AladinViewerRef>(null);
  const [panels, setPanels] = useState<{
    open: string[];
    autoOpened: string[];
  }>({ open: [], autoOpened: [] });
  const {
    data: lightcurveData,
    isLoading: backendLightcurveLoading,
    error: backendLightcurveError,
  } = useLightcurve({ ra: object.ra, dec: object.dec, radius: 1.5 });
  // ZTF comes straight from ALeRCE: the backend's ZTF client sends the radius
  // in the wrong unit and never returns detections (see /api/ztf-lightcurve).
  const {
    data: ztfLightcurveData,
    isLoading: ztfLightcurveLoading,
    error: ztfLightcurveError,
  } = useZtfLightcurve({ ra: object.ra, dec: object.dec, radius: 1.5 });
  const lightcurveLoading = backendLightcurveLoading || ztfLightcurveLoading;
  const lightcurveError = backendLightcurveError ?? ztfLightcurveError;
  // DESI spectrum: resolve TARGETID from coordinates (deduped with ObjectArchives'
  // identical query), then fetch the full-resolution wavelength/flux arrays.
  const { data: desiTarget } = useDesiTarget({
    ra: object.ra,
    dec: object.dec,
  });
  const desiTargetid = desiTarget?.targetid ?? null;
  const {
    data: spectrumData,
    isLoading: spectrumLoading,
    error: spectrumError,
  } = useDesiSpectrum(desiTargetid);
  const copyToClipboard = (text: string, label: string) => {
    navigator.clipboard.writeText(text);
    message.success(`${label} copied to clipboard`);
  };

  const meta = metadata as Record<string, unknown> | undefined;
  const galactic = equatorialToGalactic(object.ra, object.dec);
  const ecliptic = equatorialToEcliptic(object.ra, object.dec);

  // The detail page only loads the catalog the object came from; the other
  // catalogs' photometry comes from the nearest counterpart in each.
  const counterparts = useCounterparts({
    ra: object.ra,
    dec: object.dec,
    sourceCatalog: object.catalog,
  });
  // SIMBAD positions are good to ~1″, so never search tighter than 2″ even
  // for Gaia; for AllWISE/eROSITA use the catalog's own positional budget.
  const simbadRadiusArcsec = Math.max(
    2,
    COUNTERPART_RADIUS_ARCSEC[
      object.catalog.toLowerCase() as keyof typeof COUNTERPART_RADIUS_ARCSEC
    ] ?? 2
  );
  const nearby = useNeighbors({
    ra: object.ra,
    dec: object.dec,
    self: { id: object.objectId, catalog: object.catalog },
  });
  const photometrySources: PhotometrySource[] = [
    {
      catalog: object.catalog,
      id: object.objectId,
      record: meta ?? {},
      isSelf: true,
      separationArcsec: 0,
    },
    ...counterparts
      .filter((c) => c.status === "found" && c.record)
      .map((c) => ({
        catalog: c.catalog,
        id: c.id,
        record: c.record!,
        isSelf: false,
        separationArcsec: c.separationArcsec ?? 0,
      })),
  ];
  const sedPoints = buildSedPoints(photometrySources);

  // Gaia epoch photometry is fetched by source_id. Take it from the object
  // itself or its Gaia counterpart (from the designation: the numeric
  // source_id field has lost precision), and skip the request when Gaia says
  // there is none. With no Gaia counterpart, the route cone-searches the
  // Gaia archive by position instead.
  const gaiaCounterpart = counterparts.find((c) => c.catalog === "gaia");
  const gaiaSource =
    object.catalog.toLowerCase() === "gaia"
      ? { id: object.objectId, record: meta }
      : gaiaCounterpart?.status === "found"
        ? { id: gaiaCounterpart.id, record: gaiaCounterpart.record }
        : null;
  let gaiaEpochParams: GaiaEpochParams | null = null;
  if (gaiaSource) {
    const sourceId = gaiaSource.id
      ? gaiaSourceIdFromDesignation(gaiaSource.id)
      : null;
    const flag = gaiaSource.record?.has_epoch_photometry;
    if (sourceId && flag !== 0 && flag !== false)
      gaiaEpochParams = { sourceId };
  } else if (
    gaiaCounterpart?.status === "none" ||
    gaiaCounterpart?.status === "error"
  ) {
    gaiaEpochParams = { ra: object.ra, dec: object.dec };
  }
  const gaiaEpoch = useGaiaEpochPhotometry(gaiaEpochParams);
  const ps1Lightcurve = usePs1Lightcurve({ ra: object.ra, dec: object.dec });
  const vizierSed = useVizierSed({ ra: object.ra, dec: object.dec });
  const reddening = useGalacticReddening({ ra: object.ra, dec: object.dec });

  const photometryData = PHOTOMETRY_BANDS.map((band) => {
    // Prefer the object's own measurement over a counterpart's.
    const point =
      sedPoints.find((p) => p.band === band.band && p.isSelf) ??
      sedPoints.find((p) => p.band === band.band);
    return { band: band.band, survey: band.survey, point };
  });
  const measuredBands = photometryData.filter((p) => p.point).length;

  // Build catalog details from all metadata fields (exclude id, ra, dec already shown)
  const excludedFields = new Set(["id", "ra", "dec"]);
  const catalogDetails = meta
    ? Object.entries(meta)
        .filter(([key]) => !excludedFields.has(key))
        .map(([key, value]) => ({
          key,
          label: key,
          description: describeCatalogField(object.catalog, key),
          value: value == null ? "—" : String(value),
        }))
    : [];

  // Per-catalog light curve panels: the unified /lightcurve endpoint for every
  // survey except ZTF, which is taken from the ALeRCE proxy instead
  // Gaia and Pan-STARRS come from their own archives and appear as they
  // arrive, without holding back the panels above.
  const lightcurveByCatalog = {
    ...groupDetectionsByCatalog(lightcurveData, ["ztf"]),
    ...groupDetectionsByCatalog(ztfLightcurveData),
    ...(gaiaEpoch.data?.found ? { gaia: gaiaEpoch.data.points } : {}),
    ...(ps1Lightcurve.data?.found ? { ps1: ps1Lightcurve.data.points } : {}),
  };
  // Shared MJD range across all surveys so panels can be visually compared along the time axis
  const allLightcurveMjds = Object.values(lightcurveByCatalog)
    .flat()
    .map((p) => p.mjd)
    .filter((m): m is number => typeof m === "number");
  const sharedMjdRange =
    allLightcurveMjds.length > 0
      ? calculateAxisBounds(allLightcurveMjds, 0.05, 1)
      : undefined;
  // Stem for downloaded CSV filenames; falls back to coordinates when no objectId
  const filenameStem =
    object.objectId || `${object.ra.toFixed(5)}_${object.dec.toFixed(5)}`;
  const surveyPanelItems = Object.entries(lightcurveByCatalog)
    .sort(([a], [b]) => a.localeCompare(b))
    .map(([catalog, points]) => ({
      key: `lightcurve-${catalog}`,
      label: (
        <Space>
          <LineChartOutlined />
          <span>Light Curve ({getCatalogLabel(catalog)})</span>
          <Text type="secondary" className="text-xs">
            ({points.length} points)
          </Text>
        </Space>
      ),
      extra: (
        <Tooltip title="Download as CSV">
          <Button
            type="text"
            size="small"
            icon={<DownloadOutlined />}
            onClick={(e) => {
              e.stopPropagation();
              downloadCsv(
                `${filenameStem}_${getCatalogLabel(catalog)}_lightcurve.csv`,
                detectionPointsToCsv({ [catalog]: points })
              );
            }}
          />
        </Tooltip>
      ),
      children: (
        <LightCurveChart
          data={{
            detections: points,
            non_detections: [],
            forced_photometry: [],
          }}
          mjdRange={sharedMjdRange}
          magSystem={getMagSystem(catalog)}
        />
      ),
    }));

  // Placeholder shown while the light curve requests are in flight, or when
  // no survey returned detections (errored or empty). One source failing still
  // shows the panels from the other.
  let lightcurveStatusItem: {
    key: string;
    label: ReactNode;
    children: ReactNode;
  } | null = null;
  if (lightcurveLoading) {
    lightcurveStatusItem = {
      key: "lightcurve-loading",
      label: (
        <Space>
          <LineChartOutlined />
          <span>Light Curves</span>
          <Text type="secondary" className="text-xs">
            (loading…)
          </Text>
        </Space>
      ),
      children: <LightCurveSkeleton />,
    };
  } else if (surveyPanelItems.length === 0 && lightcurveError) {
    lightcurveStatusItem = {
      key: "lightcurve-error",
      label: (
        <Space>
          <LineChartOutlined />
          <span>Light Curves</span>
          <Text type="danger" className="text-xs">
            (failed to load)
          </Text>
        </Space>
      ),
      children: (
        <Text type="danger">
          {lightcurveError.message || "Failed to load light curves"}
        </Text>
      ),
    };
  } else if (surveyPanelItems.length === 0) {
    lightcurveStatusItem = {
      key: "lightcurve-empty",
      label: (
        <Space>
          <LineChartOutlined />
          <span>Light Curves</span>
          <Text type="secondary" className="text-xs">
            (unavailable)
          </Text>
        </Space>
      ),
      children: (
        <Empty
          image={Empty.PRESENTED_IMAGE_SIMPLE}
          description={
            <Text type="secondary">
              No light curves available for this object.
            </Text>
          }
        />
      ),
    };
  }

  // DESI spectrum panel — only present when coordinates resolve to a DESI target
  const spectrumNotFound = spectrumData?.found === false;
  const spectrumReady = spectrumData?.found === true;
  const spectrumLabel = spectrumLoading
    ? "(loading…)"
    : spectrumError
      ? "(temporarily unavailable)"
      : spectrumReady
        ? `(${spectrumData.spectype ?? "spectrum"})`
        : spectrumNotFound
          ? "(none)"
          : "";
  const spectrumItem = desiTargetid
    ? {
        key: "desi-spectrum",
        label: (
          <Space>
            <LineChartOutlined />
            <span>DESI Spectrum</span>
            <Text type="secondary" className="text-xs">
              {spectrumLabel}
            </Text>
          </Space>
        ),
        children: (
          <SpectrumChart
            wavelength={spectrumData?.wavelength}
            flux={spectrumData?.flux}
            loading={spectrumLoading}
            error={spectrumError ?? null}
            notFound={spectrumNotFound}
            targetid={desiTargetid}
            spectype={spectrumData?.spectype}
            redshift={spectrumData?.redshift}
            model={spectrumData?.model}
            ivar={spectrumData?.ivar}
          />
        ),
      }
    : null;

  const collapseItems = [
    {
      key: "photometry",
      label: (
        <Space>
          <StarOutlined />
          <span>Photometry</span>
          <Text type="secondary" className="text-xs">
            ({measuredBands} bands)
          </Text>
        </Space>
      ),
      children: (
        <Flex vertical gap={16}>
          <Flex wrap="wrap" gap={8}>
            {photometryData.map(({ band, survey, point }) => (
              <Tooltip key={band} title={formatChipTooltip(point, survey)}>
                <div className="text-center px-3 py-1 rounded border border-border bg-surface">
                  <Text
                    type={point ? "secondary" : undefined}
                    className={`text-xs ${point ? "" : "text-border"}`}
                  >
                    {band}
                  </Text>
                  <div
                    className={`font-mono text-sm ${point ? "" : "text-border"}`}
                  >
                    {point
                      ? `${point.upperLimit ? ">" : ""}${point.mag.toFixed(2)}`
                      : "—"}
                  </div>
                </div>
              </Tooltip>
            ))}
          </Flex>
          <SedChart
            ra={object.ra}
            dec={object.dec}
            points={sedPoints}
            vizier={{
              points: vizierSed.data?.points ?? [],
              rowCount: vizierSed.data?.rowCount ?? 0,
              loading: vizierSed.isPending,
              error: vizierSed.error,
            }}
            reddening={{
              value: reddening.data,
              loading: reddening.isPending,
              error: reddening.error,
            }}
            filenameStem={filenameStem}
          />
        </Flex>
      ),
    },
    {
      key: "observability",
      label: (
        <Space>
          <FieldTimeOutlined />
          <span>Observability</span>
          <Text type="secondary" className="text-xs">
            (Chilean observatories)
          </Text>
        </Space>
      ),
      children: <ObservabilityPanel ra={object.ra} dec={object.dec} />,
    },
    {
      key: "nearby",
      label: (
        <Space>
          <RadarChartOutlined />
          <span>Nearby Sources</span>
          <Text type="secondary" className="text-xs">
            {nearby.isLoading && nearby.neighbors.length === 0
              ? `(within ${NEIGHBOR_RADIUS_ARCSEC}″)`
              : `(${nearby.neighbors.length} within ${NEIGHBOR_RADIUS_ARCSEC}″)`}
          </Text>
        </Space>
      ),
      children: (
        <NearbySources
          neighbors={nearby.neighbors}
          radiusArcsec={NEIGHBOR_RADIUS_ARCSEC}
          loading={nearby.isLoading}
          failedCatalogs={nearby.failedCatalogs}
          truncatedCatalogs={nearby.truncatedCatalogs}
        />
      ),
    },
    ...(lightcurveStatusItem ? [lightcurveStatusItem] : surveyPanelItems),
    ...(spectrumItem ? [spectrumItem] : []),
    ...(catalogDetails.length > 0
      ? [
          {
            key: "catalog-details",
            label: (
              <Space>
                <DatabaseOutlined />
                <span>Catalog Details</span>
                <Text type="secondary" className="text-xs">
                  ({catalogDetails.length} fields)
                </Text>
              </Space>
            ),
            children: (
              <Descriptions
                size="small"
                column={{ xs: 1, sm: 2, md: 3 }}
                bordered
              >
                {catalogDetails.map((field) => (
                  <Descriptions.Item
                    key={field.key}
                    label={
                      <Tooltip title={field.description}>
                        <Text
                          className={`font-mono text-xs ${field.description ? "cursor-help underline decoration-dotted underline-offset-2" : ""}`}
                        >
                          {field.label}
                        </Text>
                      </Tooltip>
                    }
                  >
                    <Text className="font-mono text-xs">{field.value}</Text>
                  </Descriptions.Item>
                ))}
              </Descriptions>
            ),
          },
        ]
      : []),
  ];

  // Panels open themselves the first time their data arrives. Tracking which
  // keys were already auto-opened (instead of remounting the Collapse) keeps
  // whatever the user opened or closed in the meantime.
  const autoOpenKeys = [
    "photometry",
    ...(lightcurveStatusItem
      ? [lightcurveStatusItem.key]
      : surveyPanelItems.map((item) => item.key)),
    ...(spectrumReady ? ["desi-spectrum"] : []),
  ];
  const freshKeys = autoOpenKeys.filter((k) => !panels.autoOpened.includes(k));
  if (freshKeys.length > 0) {
    // Adjusting state during render: React re-renders before committing.
    setPanels((p) => ({
      open: [...p.open, ...freshKeys],
      autoOpened: [...p.autoOpened, ...freshKeys],
    }));
  }

  return (
    <div className="p-6 max-w-6xl mx-auto">
      {/* Main two-column layout */}
      <Row gutter={[24, 24]} className="mb-6">
        {/* Left: Object Info */}
        <Col xs={24} md={14}>
          <Card
            className="bg-surface h-full"
            size="small"
            styles={{ body: { height: "100%" } }}
          >
            <Flex vertical className="h-full" justify="space-between">
              <Flex vertical gap={16}>
                {/* Object Name */}
                <div>
                  <Title level={3} className="!m-0 !mb-2">
                    {object.objectId}
                  </Title>
                  <SimbadIdentity
                    ra={object.ra}
                    dec={object.dec}
                    radiusArcsec={simbadRadiusArcsec}
                  />
                  <div className="mt-3">
                    <ObjectActions
                      objectId={object.objectId}
                      catalog={object.catalog}
                      ra={object.ra}
                      dec={object.dec}
                    />
                  </div>
                </div>

                {/* Coordinates */}
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                  <div>
                    <Text type="secondary" className="text-xs block mb-1">
                      Right Ascension
                    </Text>
                    <Flex align="center" gap={8}>
                      <Text className="font-mono">{object.ra.toFixed(6)}°</Text>
                      <Text type="secondary" className="text-xs">
                        ({toHMS(object.ra)})
                      </Text>
                      <Button
                        type="text"
                        size="small"
                        icon={<CopyOutlined />}
                        onClick={() =>
                          copyToClipboard(object.ra.toFixed(6), "RA")
                        }
                      />
                    </Flex>
                  </div>
                  <div>
                    <Text type="secondary" className="text-xs block mb-1">
                      Declination
                    </Text>
                    <Flex align="center" gap={8}>
                      <Text className="font-mono">
                        {object.dec.toFixed(6)}°
                      </Text>
                      <Text type="secondary" className="text-xs">
                        ({toDMS(object.dec)})
                      </Text>
                      <Button
                        type="text"
                        size="small"
                        icon={<CopyOutlined />}
                        onClick={() =>
                          copyToClipboard(object.dec.toFixed(6), "Dec")
                        }
                      />
                    </Flex>
                  </div>
                  <div>
                    <Text type="secondary" className="text-xs block mb-1">
                      Galactic
                    </Text>
                    <Text className="font-mono text-sm">
                      l {galactic.l.toFixed(4)}°, b{" "}
                      {formatSigned(galactic.b, 4)}°
                    </Text>
                  </div>
                  <div>
                    <Text type="secondary" className="text-xs block mb-1">
                      Ecliptic (J2000)
                    </Text>
                    <Text className="font-mono text-sm">
                      λ {ecliptic.lambda.toFixed(4)}°, β{" "}
                      {formatSigned(ecliptic.beta, 4)}°
                    </Text>
                  </div>
                </div>

                {/* Nearest counterpart in every other catalog */}
                <div>
                  <Text type="secondary" className="text-xs block mb-1">
                    Counterparts
                  </Text>
                  <Flex vertical gap={2}>
                    {counterparts.map((c) => (
                      <Flex key={c.catalog} align="center" gap={8}>
                        <span
                          className={`inline-block w-2 h-2 rounded-full ${CATALOG_COLOR_CLASSES[c.catalog]}`}
                        />
                        <Text className="text-sm">
                          {getSearchCatalogLabel(c.catalog)}
                        </Text>
                        {c.status === "found" && c.id && (
                          <Link
                            href={buildObjectUrl(c.id, c.catalog)}
                            className="font-mono text-xs truncate"
                          >
                            {c.id}
                          </Link>
                        )}
                        <Text
                          type={c.status === "error" ? "danger" : "secondary"}
                          className="text-xs"
                        >
                          {counterpartStatus(c)}
                        </Text>
                      </Flex>
                    ))}
                  </Flex>
                </div>
              </Flex>

              <Flex justify="flex-end">
                <Tooltip
                  title={
                    surveyPanelItems.length > 0
                      ? "Download all light curves as CSV"
                      : "No light curves available for this object"
                  }
                >
                  <span>
                    <Button
                      size="small"
                      icon={<DownloadOutlined />}
                      disabled={surveyPanelItems.length === 0}
                      onClick={() =>
                        downloadCsv(
                          `${filenameStem}_lightcurve.csv`,
                          detectionPointsToCsv(lightcurveByCatalog)
                        )
                      }
                    >
                      Download All Light Curves
                    </Button>
                  </span>
                </Tooltip>
              </Flex>
            </Flex>
          </Card>
        </Col>

        {/* Right: Aladin Viewer */}
        <Col xs={24} md={10}>
          <Card
            className="bg-surface h-full flex flex-col"
            styles={{
              body: {
                padding: 0,
                flex: 1,
                display: "flex",
                flexDirection: "column",
              },
            }}
            title={
              <Space>
                <EnvironmentOutlined />
                <span>Sky View</span>
              </Space>
            }
            extra={
              <Select
                defaultValue={DSS_SURVEY}
                size="small"
                className="w-[150px]"
                options={surveySelectOptions}
                popupMatchSelectWidth={false}
                optionRender={(option) => (
                  <div className="max-w-[280px]">
                    <div>{option.label}</div>
                    <div className="text-xs text-foreground/60 whitespace-normal">
                      {SURVEY_DESCRIPTIONS.get(option.value as string)}
                    </div>
                  </div>
                )}
                onChange={(value) => aladinRef.current?.setSurvey(value)}
              />
            }
          >
            <AladinViewer
              ref={aladinRef}
              center={{ ra: object.ra, dec: object.dec }}
              fov={0.9}
              height="auto"
              className="flex-1 min-h-[200px]"
            />
          </Card>
        </Col>
      </Row>

      <div className="mb-4">
        <ObjectArchives ra={object.ra} dec={object.dec} />
      </div>

      <Collapse
        items={collapseItems}
        activeKey={panels.open}
        onChange={(keys) =>
          setPanels((p) => ({ ...p, open: ([] as string[]).concat(keys) }))
        }
        className="bg-surface"
      />
    </div>
  );
}
