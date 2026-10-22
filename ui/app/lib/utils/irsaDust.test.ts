import { describe, expect, it } from "vitest";

import { parseIrsaDust } from "./irsaDust";

/** Trimmed real response for 150.0991255 +2.2003759 (COSMOS). */
const OK = `<?xml version="1.0"?>
<results status="ok">
  <input><objname> 150.099130 2.200380 equ J2000 </objname></input>
  <result>
    <desc>
      E(B-V) Reddening
    </desc>
    <statistics>
      <refPixelValueSandF>
             0.0163 (mag)
      </refPixelValueSandF>
      <refPixelValueSFD>
             0.0189 (mag)
      </refPixelValueSFD>
      <meanValueSandF>
           0.0158 (mag)
      </meanValueSandF>
    </statistics>
  </result>
  <result>
    <desc>
      100 Micron Emission
    </desc>
    <statistics>
      <refPixelValueSFD>
             0.6400 (MJy/sr)
      </refPixelValueSFD>
    </statistics>
  </result>
</results>`;

describe("parseIrsaDust", () => {
  it("reads the reference-pixel E(B−V), not the regional mean", () => {
    expect(parseIrsaDust(OK)).toEqual({ ebvSF11: 0.0163, ebvSFD: 0.0189 });
  });

  it("surfaces service errors", () => {
    const err = `<results status="error"><message>Invalid location</message></results>`;
    expect(() => parseIrsaDust(err)).toThrow(/Invalid location/);
  });

  it("rejects a response without the reddening block", () => {
    const noEbv = OK.replace("E(B-V) Reddening", "Something else");
    expect(() => parseIrsaDust(noEbv)).toThrow(/no E\(B-V\) block/);
  });
});
