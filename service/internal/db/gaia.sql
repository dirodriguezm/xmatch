-- name: InsertGaia :exec
INSERT INTO gaia (
	id,
	source_id,
	ra_error,
	dec_error,
	parallax,
	parallax_error,
	pm,
	pmra,
	pmra_error,
	pmdec,
	pmdec_error,
	astrometric_excess_noise,
	astrometric_excess_noise_sig,
	ruwe,
	phot_g_n_obs,
	phot_g_mean_flux,
	phot_g_mean_flux_error,
	phot_g_mean_flux_over_error,
	phot_g_mean_mag,
	phot_bp_n_obs,
	phot_bp_mean_flux,
	phot_bp_mean_flux_error,
	phot_bp_mean_flux_over_error,
	phot_bp_mean_mag,
	phot_rp_n_obs,
	phot_rp_mean_flux,
	phot_rp_mean_flux_error,
	phot_rp_mean_flux_over_error,
	phot_rp_mean_mag,
	phot_bp_rp_excess_factor,
	phot_proc_mode,
	bp_rp,
	bp_g,
	g_rp,
	radial_velocity,
	radial_velocity_error,
	rv_method_used,
	phot_variable_flag,
	in_qso_candidates,
	in_galaxy_candidates,
	non_single_star,
	has_epoch_photometry,
	classprob_dsc_combmod_quasar,
	classprob_dsc_combmod_galaxy,
	classprob_dsc_combmod_star,
	teff_gspphot,
	teff_gspphot_lower,
	teff_gspphot_upper,
	logg_gspphot,
	logg_gspphot_lower,
	logg_gspphot_upper,
	mh_gspphot,
	mh_gspphot_lower,
	mh_gspphot_upper,
	distance_gspphot,
	distance_gspphot_lower,
	distance_gspphot_upper
) VALUES (
	?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
	?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
	?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
	?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
	?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
	?, ?, ?, ?, ?, ?, ?
);

-- name: GetGaia :one
SELECT gaia.*, mastercat.ra, mastercat.dec
FROM gaia 
JOIN mastercat ON mastercat.id = gaia.id
WHERE gaia.id = ?;

-- name: BulkGetGaia :many
SELECT gaia.*, mastercat.ra, mastercat.dec
FROM gaia 
JOIN mastercat ON mastercat.id = gaia.id
WHERE gaia.id IN (sqlc.slice(id));
