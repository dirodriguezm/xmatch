import { NextRequest, NextResponse } from "next/server";

/**
 * ZTF DR light curves fetched straight from ALeRCE instead of via the xwave
 * backend's /lightcurve endpoint.
 *
 * Known backend issue (service/internal/search/lightcurve/ztfdr/client.go):
 * the ZTF DR client converts the radius from arcsec to degrees before calling
 * ALeRCE, following ALeRCE's swagger ("Radius [decimal degrees]"). The live
 * endpoint actually interprets `radius` as arcsec — for NGC5252 the two ZTF
 * sources at ~0.23" appear only once radius >= 0.3 — so the backend searches
 * a ~0.0004" cone and always returns zero ZTF detections. Until that is fixed
 * the UI sources ZTF here and ignores ZTF rows from the backend.
 */
const ALERCE_API_URL = "https://api.alerce.online/ztf/dr/v1/light_curve/";

interface AlerceLightcurve {
  _id: number;
  filterid: number;
  fieldid: number;
  nepochs: number;
  objra: number;
  objdec: number;
  rcid: number;
  hmjd: number[];
  mag: number[];
  magerr: number[];
}

// Same shape as the backend /lightcurve detections so the UI's
// groupDetectionsByCatalog/expandDetection handle both sources identically.
interface Detection {
  catalog: "ztf";
  id: string;
  object_id: string;
  mjd: number;
  mag: number;
  magerr: number;
  data: { filterid: number; fieldid: number; rcid: number };
}

function flattenLightcurves(lightcurves: AlerceLightcurve[]): Detection[] {
  const detections: Detection[] = [];

  for (const lc of lightcurves) {
    const len = Math.min(lc.hmjd.length, lc.mag.length, lc.magerr.length);
    for (let i = 0; i < len; i++) {
      detections.push({
        catalog: "ztf",
        id: `${lc._id}_${i}`,
        object_id: String(lc._id),
        mjd: lc.hmjd[i],
        mag: lc.mag[i],
        magerr: lc.magerr[i],
        data: { filterid: lc.filterid, fieldid: lc.fieldid, rcid: lc.rcid },
      });
    }
  }

  return detections;
}

export async function GET(request: NextRequest) {
  const searchParams = request.nextUrl.searchParams;
  const ra = searchParams.get("ra");
  const dec = searchParams.get("dec");
  // arcsec, forwarded unconverted (see the note above ALERCE_API_URL)
  const radius = searchParams.get("radius") || "2";

  if (!ra || !dec) {
    return NextResponse.json(
      { error: "Missing required parameters: ra, dec" },
      { status: 400 }
    );
  }

  const url = `${ALERCE_API_URL}?ra=${ra}&dec=${dec}&radius=${radius}`;

  try {
    const response = await fetch(url);

    if (!response.ok) {
      return NextResponse.json(
        { error: `ALeRCE API error: ${response.status}` },
        { status: response.status }
      );
    }

    const lightcurves: AlerceLightcurve[] = await response.json();
    const detections = flattenLightcurves(lightcurves);

    return NextResponse.json({
      detections,
      non_detections: [],
      forced_photometry: [],
    });
  } catch (error) {
    console.error("ZTF lightcurve proxy error:", error);
    return NextResponse.json(
      { error: "Failed to fetch from ZTF lightcurve service" },
      { status: 500 }
    );
  }
}
