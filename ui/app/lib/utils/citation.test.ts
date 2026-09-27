import { describe, expect, it } from "vitest";

import {
  bibFilename,
  bibtexTitle,
  buildAcknowledgement,
  buildBibtex,
  buildPlainReferences,
  resolveCatalogs,
  XWAVE_ACKNOWLEDGEMENT,
  xwaveBibtex,
} from "./citation";

describe("citation", () => {
  it("resolves every catalog when the list is empty", () => {
    expect(resolveCatalogs([]).map((m) => m.slug)).toEqual([
      "allwise",
      "gaia",
      "erosita",
    ]);
  });

  it("dedupes, ignores unknown slugs and keeps canonical order", () => {
    expect(
      resolveCatalogs(["gaia", "GAIA", "nope", "allwise"]).map((m) => m.slug)
    ).toEqual(["allwise", "gaia"]);
  });

  it("builds an XWave software entry with a pending-DOI note", () => {
    const bib = xwaveBibtex("https://example.org");
    expect(bib).toContain("@software{XWave2026,");
    expect(bib).toContain("XWave: a HEALPix cross-match service");
    expect(bib).toContain("url = {https://example.org}");
    expect(bib).toContain("year = 2026");
    expect(bib).toMatch(/Zenodo DOI is pending/);
  });

  it("includes one entry per requested catalog", () => {
    const bib = buildBibtex(["gaia", "gaia"]);
    expect(bib.match(/@ARTICLE\{2023A&A\.\.\.674A\.\.\.1G/g)).toHaveLength(1);
    expect(bib).not.toContain("2013yCat.2328....0C");
    expect(buildBibtex([]).match(/^@/gm)).toHaveLength(4);
  });

  it("extracts BibTeX titles", () => {
    expect(
      bibtexTitle(`@X{a,\n  title = "{Gaia Data Release 3. Summary}",\n}`)
    ).toBe("Gaia Data Release 3. Summary");
    expect(bibtexTitle(`@X{a,\n  title = {{XWave: a service}},\n}`)).toBe(
      "XWave: a service"
    );
  });

  it("formats plain-text references with DOI or bibcode", () => {
    const refs = buildPlainReferences(["gaia", "allwise"], "https://x.test");
    expect(refs).toHaveLength(3);
    expect(refs[0]).toContain("XWave: a HEALPix cross-match service");
    expect(refs[0]).toContain("https://x.test");
    expect(refs[1]).toContain("Cutri et al. 2013");
    expect(refs[1]).toContain("2013yCat.2328....0C");
    expect(refs[2]).toContain("doi:10.1051/0004-6361/202243940");
  });

  it("combines acknowledgements starting with XWave", () => {
    const ack = buildAcknowledgement(["erosita"]);
    expect(ack.startsWith(XWAVE_ACKNOWLEDGEMENT)).toBe(true);
    expect(ack).toContain("eROSITA");
    expect(ack).not.toContain("Gaia");
  });

  it("builds safe .bib filenames", () => {
    expect(bibFilename("Gaia DR3 381267755864434048")).toBe(
      "xwave-gaia-dr3-381267755864434048.bib"
    );
    expect(bibFilename()).toBe("xwave-references.bib");
  });
});
