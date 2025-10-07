CREATE TABLE gaia (
    id text not null,
    phot_g_mean_flux double precision,
    phot_g_mean_flux_error double precision,
    phot_g_mean_mag double precision,
    phot_bp_mean_flux double precision,
    phot_bp_mean_flux_error double precision,
    phot_bp_mean_mag double precision,
    phot_rp_mean_flux double precision,
    phot_rp_mean_flux_error double precision,
    phot_rp_mean_mag double precision,
    PRIMARY KEY (id)
)
