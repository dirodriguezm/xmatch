/**
 * Human-readable descriptions for the metadata fields each catalog returns,
 * keyed by catalog and then by the JSON field name the service emits (see
 * service/internal/repository/models.go). Fields without an entry are shown
 * with their raw name only.
 */

const ALLWISE_UPPER_LIMIT =
  "Empty when the source was not detected in this band; the magnitude is then a 95% upper limit.";
const TWOMASS_UPPER_LIMIT =
  "Empty when the magnitude is a 95% upper limit rather than a detection.";

export const CATALOG_FIELD_DESCRIPTIONS: Record<
  string,
  Record<string, string>
> = {
  allwise: {
    cntr: "Unique AllWISE source counter (row identifier in the IRSA catalog)",
    w1mpro:
      "W1 (3.4 µm) profile-fit magnitude, Vega; a 95% upper limit when w1sigmpro is empty",
    w1sigmpro: `W1 magnitude uncertainty. ${ALLWISE_UPPER_LIMIT}`,
    w2mpro:
      "W2 (4.6 µm) profile-fit magnitude, Vega; a 95% upper limit when w2sigmpro is empty",
    w2sigmpro: `W2 magnitude uncertainty. ${ALLWISE_UPPER_LIMIT}`,
    w3mpro:
      "W3 (12 µm) profile-fit magnitude, Vega; a 95% upper limit when w3sigmpro is empty",
    w3sigmpro: `W3 magnitude uncertainty. ${ALLWISE_UPPER_LIMIT}`,
    w4mpro:
      "W4 (22 µm) profile-fit magnitude, Vega; a 95% upper limit when w4sigmpro is empty",
    w4sigmpro: `W4 magnitude uncertainty. ${ALLWISE_UPPER_LIMIT}`,
    j_m_2mass:
      "2MASS J (1.24 µm) magnitude of the associated 2MASS source, Vega",
    j_msig_2mass: `2MASS J uncertainty. ${TWOMASS_UPPER_LIMIT}`,
    h_m_2mass:
      "2MASS H (1.66 µm) magnitude of the associated 2MASS source, Vega",
    h_msig_2mass: `2MASS H uncertainty. ${TWOMASS_UPPER_LIMIT}`,
    k_m_2mass:
      "2MASS Ks (2.16 µm) magnitude of the associated 2MASS source, Vega",
    k_msig_2mass: `2MASS Ks uncertainty. ${TWOMASS_UPPER_LIMIT}`,
  },
  gaia: {
    source_id: "Gaia DR3 unique source identifier",
    ra_error: "Uncertainty in right ascension (mas)",
    dec_error: "Uncertainty in declination (mas)",
    parallax:
      "Parallax (mas); distance ≈ 1000 / parallax pc for precise values",
    parallax_error: "Parallax uncertainty (mas)",
    pm: "Total proper motion (mas/yr)",
    pmra: "Proper motion in RA, × cos(dec) (mas/yr)",
    pmra_error: "Uncertainty in pmra (mas/yr)",
    pmdec: "Proper motion in declination (mas/yr)",
    pmdec_error: "Uncertainty in pmdec (mas/yr)",
    astrometric_excess_noise:
      "Extra scatter needed for the astrometric fit (mas); > 0 hints at a poor fit",
    astrometric_excess_noise_sig:
      "Significance of the excess noise; > 2 means it is significant",
    ruwe: "Renormalised unit weight error; ≳ 1.4 suggests a binary or problematic astrometry",
    phot_g_n_obs: "Number of observations contributing to G photometry",
    phot_g_mean_flux: "G-band mean flux (e⁻/s)",
    phot_g_mean_flux_error: "G-band mean flux uncertainty (e⁻/s)",
    phot_g_mean_flux_over_error: "G-band flux signal-to-noise ratio",
    phot_g_mean_mag: "G-band mean magnitude (330–1050 nm), Vega",
    phot_bp_n_obs: "Number of observations contributing to BP photometry",
    phot_bp_mean_flux: "BP (blue, 330–680 nm) mean flux (e⁻/s)",
    phot_bp_mean_flux_error: "BP mean flux uncertainty (e⁻/s)",
    phot_bp_mean_flux_over_error: "BP flux signal-to-noise ratio",
    phot_bp_mean_mag: "BP mean magnitude, Vega",
    phot_rp_n_obs: "Number of observations contributing to RP photometry",
    phot_rp_mean_flux: "RP (red, 630–1050 nm) mean flux (e⁻/s)",
    phot_rp_mean_flux_error: "RP mean flux uncertainty (e⁻/s)",
    phot_rp_mean_flux_over_error: "RP flux signal-to-noise ratio",
    phot_rp_mean_mag: "RP mean magnitude, Vega",
    phot_bp_rp_excess_factor:
      "(BP + RP flux) / G flux; high values flag blending or extended sources",
    phot_proc_mode:
      "Photometric processing mode: 0 = gold (complete colour), 1 = silver (incomplete colour), 2 = bronze (insufficient colour)",
    bp_rp: "BP − RP colour (mag); larger is redder",
    bp_g: "BP − G colour (mag)",
    g_rp: "G − RP colour (mag)",
    radial_velocity: "Radial velocity from the RVS spectrometer (km/s)",
    radial_velocity_error: "Radial velocity uncertainty (km/s)",
    rv_method_used:
      "RV method: 1 = median of epoch RVs (bright, G_RVS ≤ 12), 2 = combined cross-correlation (faint, G_RVS > 12)",
    phot_variable_flag:
      "Photometric variability flag: VARIABLE, CONSTANT or NOT_AVAILABLE",
    in_qso_candidates:
      "1 if the source is in the Gaia DR3 quasar candidates table",
    in_galaxy_candidates:
      "1 if the source is in the Gaia DR3 galaxy candidates table",
    non_single_star:
      "Non-single-star bitmask: 1 = astrometric, 2 = spectroscopic, 4 = eclipsing binary",
    has_epoch_photometry: "1 if Gaia published an epoch light curve",
    classprob_dsc_combmod_quasar:
      "DSC probability of being a quasar (combined model)",
    classprob_dsc_combmod_galaxy:
      "DSC probability of being a galaxy (combined model)",
    classprob_dsc_combmod_star:
      "DSC probability of being a single star, excluding white dwarfs (combined model)",
    teff_gspphot: "Effective temperature from GSP-Phot (K)",
    teff_gspphot_lower: "Teff 16th percentile (K)",
    teff_gspphot_upper: "Teff 84th percentile (K)",
    logg_gspphot: "Surface gravity log g from GSP-Phot (log cgs)",
    logg_gspphot_lower: "log g 16th percentile",
    logg_gspphot_upper: "log g 84th percentile",
    mh_gspphot: "Metallicity [M/H] from GSP-Phot (dex)",
    mh_gspphot_lower: "[M/H] 16th percentile (dex)",
    mh_gspphot_upper: "[M/H] 84th percentile (dex)",
    distance_gspphot: "Distance from GSP-Phot (pc)",
    distance_gspphot_lower: "Distance 16th percentile (pc)",
    distance_gspphot_upper: "Distance 84th percentile (pc)",
  },
  erosita: {
    mjd: "Modified Julian Date of the eROSITA observations",
    ml_flux_1:
      "0.2–2.3 keV flux from maximum-likelihood PSF fitting (erg/s/cm²)",
  },
};

export function describeCatalogField(
  catalog: string,
  field: string
): string | undefined {
  return CATALOG_FIELD_DESCRIPTIONS[catalog.toLowerCase()]?.[field];
}
