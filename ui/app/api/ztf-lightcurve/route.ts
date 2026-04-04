import { NextRequest, NextResponse } from "next/server";

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

interface Detection {
  mjd: number;
  mag: number;
  magerr: number;
  fid: number;
}

function flattenLightcurves(lightcurves: AlerceLightcurve[]): Detection[] {
  const detections: Detection[] = [];

  for (const lc of lightcurves) {
    const len = Math.min(lc.hmjd.length, lc.mag.length, lc.magerr.length);
    for (let i = 0; i < len; i++) {
      detections.push({
        mjd: lc.hmjd[i],
        mag: lc.mag[i],
        magerr: lc.magerr[i],
        fid: lc.filterid,
      });
    }
  }

  return detections;
}

export async function GET(request: NextRequest) {
  const searchParams = request.nextUrl.searchParams;
  const ra = searchParams.get("ra");
  const dec = searchParams.get("dec");
  const radius = searchParams.get("radius") || "2"; // 2 arcsec default

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

    return NextResponse.json({ detections });
  } catch (error) {
    console.error("ZTF lightcurve proxy error:", error);
    return NextResponse.json(
      { error: "Failed to fetch from ZTF lightcurve service" },
      { status: 500 }
    );
  }
}
