import { describe, expect, it } from "vitest";

import { aggregateVizierSed, parseVizierSedVotable } from "./vizierSed";

/** Trimmed copy of a real response for 150.0991255 +2.2003759 (COSMOS). */
const VOTABLE = `<?xml version="1.0" encoding="UTF-8"?>
<VOTABLE version="1.4" xmlns="http://www.ivoa.net/xml/VOTable/v1.3">
<RESOURCE ID="VizieR_S843769039" name="VizieR(2026-09-26T20:17:19)">
  <TABLE ID="VizieR_0" name="allVizieR">
    <FIELD name="_RAJ2000" ucd="pos.eq.ra" datatype="double" unit="deg"></FIELD>
    <FIELD name="_DEJ2000" ucd="pos.eq.dec" datatype="double" unit="deg"></FIELD>
    <FIELD name="_tabname" ucd="meta.table" datatype="char" arraysize="32*"></FIELD>
    <FIELD name="_ID" ucd="meta.id" datatype="char" arraysize="128*"></FIELD>
    <FIELD name="_time" ucd="time.epoch" datatype="double" unit="d"></FIELD>
    <FIELD name="_etime" ucd="stat.error;time.epoch" datatype="double"></FIELD>
    <FIELD ID="sed_freq" name="_sed_freq" ucd="em.freq" unit="GHz" datatype="double"></FIELD>
    <FIELD ID="sed_flux" name="_sed_flux" ucd="phot.flux.density" unit="Jy" datatype="float"></FIELD>
    <FIELD ID="sed_eflux" name="_sed_eflux" ucd="stat.error;phot.flux.density" unit="Jy" datatype="float"></FIELD>
    <FIELD ID="sed_filter" name="_sed_filter" ucd="meta.id;instr.filter" datatype="char" arraysize="32*"></FIELD>
<DATA><TABLEDATA>
<TR><TD>150.09914400</TD><TD>+02.20041700</TD><TD>I/305/out</TD><TD>GSC2.3===N6SS001314&amp;-c=150.099144 +02.200417</TD><TD></TD><TD></TD><TD>382.54e+3</TD><TD>106.e-6</TD><TD>46.e-6</TD><TD>POSS-II:i</TD></TR>
<TR><TD>150.0991233900598</TD><TD>+02.2004256669209</TD><TD>I/353/gsc242</TD><TD>-c=150.09908858924 +02.20039814228</TD><TD></TD><TD></TD><TD>89.490e+3</TD><TD>45.0e-6</TD><TD>5.5e-6</TD><TD>WISE:W1</TD></TR>
<TR><TD>150.0991</TD><TD>+02.2004</TD><TD>II/328/allwise</TD><TD>x</TD><TD></TD><TD></TD><TD>89.490e+3</TD><TD>47.0e-6</TD><TD>5.0e-6</TD><TD>WISE:W1</TD></TR>
<TR><TD>150.0991</TD><TD>+02.2004</TD><TD>II/311/wise</TD><TD>x</TD><TD></TD><TD></TD><TD>89.490e+3</TD><TD>46.0e-6</TD><TD/><TD>WISE:W1</TD></TR>
<TR><TD>150.0991</TD><TD>+02.2004</TD><TD>J/ApJS/1</TD><TD>x</TD><TD></TD><TD></TD><TD>89.490e+3</TD><TD>460.e-6</TD><TD></TD><TD>WISE:W1</TD></TR>
<TR><TD>150.0991</TD><TD>+02.2004</TD><TD>J/ApJS/2</TD><TD>x</TD><TD></TD><TD></TD><TD>621.98e+3</TD><TD>-1.e-6</TD><TD></TD><TD>SDSS:g</TD></TR>
</TABLEDATA></DATA>
</TABLE>
</RESOURCE>
</VOTABLE>`;

describe("parseVizierSedVotable", () => {
  it("reads rows by column name, decoding entities and empty cells", () => {
    const rows = parseVizierSedVotable(VOTABLE);
    expect(rows).toHaveLength(6);
    expect(rows[0]).toEqual({
      table: "I/305/out",
      frequencyGHz: 382540,
      fluxJy: 106e-6,
      fluxErrJy: 46e-6,
      filter: "POSS-II:i",
    });
    // <TD/> and <TD></TD> both mean "no error"
    expect(rows[3].fluxErrJy).toBeUndefined();
    expect(rows[4].fluxErrJy).toBeUndefined();
  });

  it("rejects payloads that are not a VOTable", () => {
    expect(() => parseVizierSedVotable("<html>Service down</html>")).toThrow(
      /not a VOTable/
    );
  });

  it("returns no rows for an empty table", () => {
    const empty = VOTABLE.replace(/<TR>[\s\S]*<\/TR>\n/, "");
    expect(parseVizierSedVotable(empty)).toEqual([]);
  });
});

describe("aggregateVizierSed", () => {
  const points = aggregateVizierSed(parseVizierSedVotable(VOTABLE));
  const w1 = points.find((p) => p.filter === "WISE:W1")!;

  it("merges repeated filters and rejects the outlier", () => {
    expect(w1.nMeasurements).toBe(3);
    expect(w1.nRejected).toBe(1); // the 460 µJy row: a neighbour, 10× brighter
    expect(w1.fluxJy).toBeCloseTo(46e-6, 12);
    expect(w1.tables).toEqual([
      "I/353/gsc242",
      "II/311/wise",
      "II/328/allwise",
    ]);
  });

  it("converts frequency to wavelength and flux to νFν", () => {
    expect(w1.wavelengthUm).toBeCloseTo(3.35, 2);
    expect(w1.nuFnu).toBeCloseTo(46e-6 * 1e-23 * 89.49e12, 20);
  });

  it("drops non-positive fluxes and sorts by wavelength", () => {
    expect(points.map((p) => p.filter)).toEqual(["POSS-II:i", "WISE:W1"]);
  });
});

describe("inconsistent filters", () => {
  const row = (filter: string, frequencyGHz: number, fluxJy: number) => ({
    table: "T",
    filter,
    frequencyGHz,
    fluxJy,
  });
  // Five filters around 0.5–0.7 µm at ~100 µJy, one of them 500× too faint.
  const points = aggregateVizierSed([
    row("A", 600e3, 100e-6),
    row("B", 560e3, 110e-6),
    row("C", 520e3, 95e-6),
    row("D", 480e3, 105e-6),
    row("BAD", 540e3, 0.2e-6),
    // Isolated far-IR point: nothing to compare against, never flagged.
    row("FIR", 3e3, 1e-9),
  ]);
  const flagged = points.filter((p) => p.inconsistent).map((p) => p.filter);

  it("flags a filter far off its spectral neighbours", () => {
    expect(flagged).toEqual(["BAD"]);
  });

  it("keeps modest real features and isolated points", () => {
    // A 0.8 dex step (a star's UV drop) is below the threshold.
    const uv = aggregateVizierSed([
      row("A", 600e3, 100e-6),
      row("B", 560e3, 110e-6),
      row("C", 520e3, 95e-6),
      row("UV", 640e3, 16e-6),
    ]);
    expect(uv.some((p) => p.inconsistent)).toBe(false);
    expect(points.find((p) => p.filter === "FIR")!.inconsistent).toBe(false);
  });
});
