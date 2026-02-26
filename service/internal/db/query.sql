-- name: FindObjects :many
SELECT *
FROM mastercat 
WHERE ipix IN (sqlc.slice(ipix));

-- name: GetObjectsFromCatalog :many
SELECT * 
FROM mastercat 
WHERE ipix IN (sqlc.slice(ipix))
AND cat = ?;

-- name: InsertObject :exec
INSERT INTO mastercat (
	id, ipix, ra, dec, cat
) VALUES (
	?, ?, ?, ?, ?
);

-- name: GetAllObjects :many
SELECT *
FROM mastercat;

-- name: GetCatalogs :many
SELECT *
FROM catalogs;

-- name: InsertCatalog :exec
INSERT INTO catalogs (
	name, nside
) VALUES (
	?, ?
);

-- name: InsertAllwise :exec
INSERT INTO allwise (
	id, cntr, w1mpro, w1sigmpro, w2mpro, w2sigmpro, w3mpro, w3sigmpro, w4mpro, w4sigmpro, J_m_2mass, J_msig_2mass, H_m_2mass, H_msig_2mass, K_m_2mass, K_msig_2mass
) VALUES (
	?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
);

-- name: GetAllwise :one
SELECT *
FROM allwise
WHERE id = ?;

-- name: BulkGetAllwise :many
SELECT *
FROM allwise
WHERE id IN (sqlc.slice(id));

-- name: RemoveAllObjects :exec
DELETE FROM mastercat;

-- name: RemoveAllAllwise :exec
DELETE FROM allwise;

-- name: RemoveAllCatalogs :exec
DELETE FROM catalogs;

-- name: GetAllwiseFromPixels :many
SELECT allwise.*, mastercat.ra, mastercat.dec
FROM allwise 
JOIN mastercat ON mastercat.id = allwise.id
WHERE mastercat.ipix IN (sqlc.slice(ipix));

-- name: InsertGaia :exec
INSERT INTO gaia (
	id, 
  phot_g_mean_flux,
  phot_g_mean_flux_error,
  phot_g_mean_mag,
  phot_bp_mean_flux,
  phot_bp_mean_flux_error,
  phot_bp_mean_mag,
  phot_rp_mean_flux,
  phot_rp_mean_flux_error,
  phot_rp_mean_mag
) VALUES (
	?, ?, ?, ?, ?, ?, ?, ?, ?, ?
);

-- name: GetGaia :one
SELECT *
FROM gaia
WHERE id = ?;

-- name: BulkGetGaia :many
SELECT *
FROM gaia
WHERE id IN (sqlc.slice(id));

-- name: RemoveAllGaia :exec
DELETE FROM gaia;

-- name: GetGaiaFromPixels :many
SELECT gaia.*, mastercat.ra, mastercat.dec
FROM gaia 
JOIN mastercat ON mastercat.id = gaia.id
WHERE mastercat.ipix IN (sqlc.slice(ipix));

-- name: InsertErosita :exec
INSERT INTO erosita (
    id, detuid, skytile, id_src, uid, uid_hard, id_cluster,
    ra, dec, ra_lowerr, ra_uperr, dec_lowerr, dec_uperr, pos_err,
    mjd, mjd_min, mjd_max, ext, ext_err, ext_like, det_like_0,
    ml_cts_1, ml_cts_err_1, ml_rate_1, ml_rate_err_1, ml_flux_1, ml_flux_err_1, ml_bkg_1, ml_exp_1, ape_bkg_1, ape_radius_1, ape_pois_1,
    det_like_p1, ml_cts_p1, ml_cts_err_p1, ml_rate_p1, ml_rate_err_p1, ml_flux_p1, ml_flux_err_p1, ml_bkg_p1, ml_exp_p1, ape_bkg_p1, ape_radius_p1, ape_pois_p1,
    det_like_p2, ml_cts_p2, ml_cts_err_p2, ml_rate_p2, ml_rate_err_p2, ml_flux_p2, ml_flux_err_p2, ml_bkg_p2, ml_exp_p2, ape_bkg_p2, ape_radius_p2, ape_pois_p2,
    det_like_p3, ml_cts_p3, ml_cts_err_p3, ml_rate_p3, ml_rate_err_p3, ml_flux_p3, ml_flux_err_p3, ml_bkg_p3, ml_exp_p3, ape_bkg_p3, ape_radius_p3, ape_pois_p3,
    det_like_p4, ml_cts_p4, ml_cts_err_p4, ml_rate_p4, ml_rate_err_p4, ml_flux_p4, ml_flux_err_p4, ml_bkg_p4, ml_exp_p4, ape_bkg_p4, ape_radius_p4, ape_pois_p4,
    det_like_p5, ml_cts_p5, ml_cts_err_p5, ml_rate_p5, ml_rate_err_p5, ml_flux_p5, ml_flux_err_p5, ml_bkg_p5, ml_exp_p5, ape_bkg_p5, ape_radius_p5, ape_pois_p5,
    det_like_p6, ml_cts_p6, ml_cts_err_p6, ml_rate_p6, ml_rate_err_p6, ml_flux_p6, ml_flux_err_p6, ml_bkg_p6, ml_exp_p6, ape_bkg_p6, ape_radius_p6, ape_pois_p6,
    flag_sp_snr, flag_sp_bps, flag_sp_scl, flag_sp_lga, flag_sp_gc_cons, flag_no_radec_err, flag_no_ext_err, flag_no_cts_err, flag_opt
) VALUES (
    $id, $detuid, $skytile, $id_src, $uid, $uid_hard, $id_cluster,
    $ra, $dec, $ra_lowerr, $ra_uperr, $dec_lowerr, $dec_uperr, $pos_err,
    $mjd, $mjd_min, $mjd_max, $ext, $ext_err, $ext_like, $det_like_0,
    $ml_cts_1, $ml_cts_err_1, $ml_rate_1, $ml_rate_err_1, $ml_flux_1, $ml_flux_err_1, $ml_bkg_1, $ml_exp_1, $ape_bkg_1, $ape_radius_1, $ape_pois_1,
    $det_like_p1, $ml_cts_p1, $ml_cts_err_p1, $ml_rate_p1, $ml_rate_err_p1, $ml_flux_p1, $ml_flux_err_p1, $ml_bkg_p1, $ml_exp_p1, $ape_bkg_p1, $ape_radius_p1, $ape_pois_p1,
    $det_like_p2, $ml_cts_p2, $ml_cts_err_p2, $ml_rate_p2, $ml_rate_err_p2, $ml_flux_p2, $ml_flux_err_p2, $ml_bkg_p2, $ml_exp_p2, $ape_bkg_p2, $ape_radius_p2, $ape_pois_p2,
    $det_like_p3, $ml_cts_p3, $ml_cts_err_p3, $ml_rate_p3, $ml_rate_err_p3, $ml_flux_p3, $ml_flux_err_p3, $ml_bkg_p3, $ml_exp_p3, $ape_bkg_p3, $ape_radius_p3, $ape_pois_p3,
    $det_like_p4, $ml_cts_p4, $ml_cts_err_p4, $ml_rate_p4, $ml_rate_err_p4, $ml_flux_p4, $ml_flux_err_p4, $ml_bkg_p4, $ml_exp_p4, $ape_bkg_p4, $ape_radius_p4, $ape_pois_p4,
    $det_like_p5, $ml_cts_p5, $ml_cts_err_p5, $ml_rate_p5, $ml_rate_err_p5, $ml_flux_p5, $ml_flux_err_p5, $ml_bkg_p5, $ml_exp_p5, $ape_bkg_p5, $ape_radius_p5, $ape_pois_p5,
    $det_like_p6, $ml_cts_p6, $ml_cts_err_p6, $ml_rate_p6, $ml_rate_err_p6, $ml_flux_p6, $ml_flux_err_p6, $ml_bkg_p6, $ml_exp_p6, $ape_bkg_p6, $ape_radius_p6, $ape_pois_p6,
    $flag_sp_snr, $flag_sp_bps, $flag_sp_scl, $flag_sp_lga, $flag_sp_gc_cons, $flag_no_radec_err, $flag_no_ext_err, $flag_no_cts_err, $flag_opt
);

-- name: GetErosita :one
SELECT * FROM erosita WHERE id = ?;

-- name: BulkGetErosita :many
SELECT * FROM erosita WHERE id IN (sqlc.slice(id));

-- name: RemoveAllErosita :exec
DELETE FROM erosita;

-- name: GetErositaFromPixels :many
SELECT erosita.*, mastercat.ra, mastercat.dec
FROM erosita 
JOIN mastercat ON mastercat.id = erosita.id
WHERE mastercat.ipix IN (sqlc.slice(ipix));
