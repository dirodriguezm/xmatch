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
SELECT
    erosita.id, erosita.detuid, erosita.skytile, erosita.id_src, erosita.uid, erosita.uid_hard, erosita.id_cluster,
    erosita.ra_lowerr, erosita.ra_uperr, erosita.dec_lowerr, erosita.dec_uperr, erosita.pos_err,
    erosita.mjd, erosita.mjd_min, erosita.mjd_max, erosita.ext, erosita.ext_err, erosita.ext_like, erosita.det_like_0,
    erosita.ml_cts_1, erosita.ml_cts_err_1, erosita.ml_rate_1, erosita.ml_rate_err_1, erosita.ml_flux_1, erosita.ml_flux_err_1, erosita.ml_bkg_1, erosita.ml_exp_1, erosita.ape_bkg_1, erosita.ape_radius_1, erosita.ape_pois_1,
    erosita.det_like_p1, erosita.ml_cts_p1, erosita.ml_cts_err_p1, erosita.ml_rate_p1, erosita.ml_rate_err_p1, erosita.ml_flux_p1, erosita.ml_flux_err_p1, erosita.ml_bkg_p1, erosita.ml_exp_p1, erosita.ape_bkg_p1, erosita.ape_radius_p1, erosita.ape_pois_p1,
    erosita.det_like_p2, erosita.ml_cts_p2, erosita.ml_cts_err_p2, erosita.ml_rate_p2, erosita.ml_rate_err_p2, erosita.ml_flux_p2, erosita.ml_flux_err_p2, erosita.ml_bkg_p2, erosita.ml_exp_p2, erosita.ape_bkg_p2, erosita.ape_radius_p2, erosita.ape_pois_p2,
    erosita.det_like_p3, erosita.ml_cts_p3, erosita.ml_cts_err_p3, erosita.ml_rate_p3, erosita.ml_rate_err_p3, erosita.ml_flux_p3, erosita.ml_flux_err_p3, erosita.ml_bkg_p3, erosita.ml_exp_p3, erosita.ape_bkg_p3, erosita.ape_radius_p3, erosita.ape_pois_p3,
    erosita.det_like_p4, erosita.ml_cts_p4, erosita.ml_cts_err_p4, erosita.ml_rate_p4, erosita.ml_rate_err_p4, erosita.ml_flux_p4, erosita.ml_flux_err_p4, erosita.ml_bkg_p4, erosita.ml_exp_p4, erosita.ape_bkg_p4, erosita.ape_radius_p4, erosita.ape_pois_p4,
    erosita.det_like_p5, erosita.ml_cts_p5, erosita.ml_cts_err_p5, erosita.ml_rate_p5, erosita.ml_rate_err_p5, erosita.ml_flux_p5, erosita.ml_flux_err_p5, erosita.ml_bkg_p5, erosita.ml_exp_p5, erosita.ape_bkg_p5, erosita.ape_radius_p5, erosita.ape_pois_p5,
    erosita.det_like_p6, erosita.ml_cts_p6, erosita.ml_cts_err_p6, erosita.ml_rate_p6, erosita.ml_rate_err_p6, erosita.ml_flux_p6, erosita.ml_flux_err_p6, erosita.ml_bkg_p6, erosita.ml_exp_p6, erosita.ape_bkg_p6, erosita.ape_radius_p6, erosita.ape_pois_p6,
    erosita.flag_sp_snr, erosita.flag_sp_bps, erosita.flag_sp_scl, erosita.flag_sp_lga, erosita.flag_sp_gc_cons, erosita.flag_no_radec_err, erosita.flag_no_ext_err, erosita.flag_no_cts_err, erosita.flag_opt,
    mastercat.ra, mastercat.dec
FROM erosita 
JOIN mastercat ON mastercat.id = erosita.id
WHERE erosita.id = ?;

-- name: BulkGetErosita :many
SELECT
    erosita.id, erosita.detuid, erosita.skytile, erosita.id_src, erosita.uid, erosita.uid_hard, erosita.id_cluster,
    erosita.ra_lowerr, erosita.ra_uperr, erosita.dec_lowerr, erosita.dec_uperr, erosita.pos_err,
    erosita.mjd, erosita.mjd_min, erosita.mjd_max, erosita.ext, erosita.ext_err, erosita.ext_like, erosita.det_like_0,
    erosita.ml_cts_1, erosita.ml_cts_err_1, erosita.ml_rate_1, erosita.ml_rate_err_1, erosita.ml_flux_1, erosita.ml_flux_err_1, erosita.ml_bkg_1, erosita.ml_exp_1, erosita.ape_bkg_1, erosita.ape_radius_1, erosita.ape_pois_1,
    erosita.det_like_p1, erosita.ml_cts_p1, erosita.ml_cts_err_p1, erosita.ml_rate_p1, erosita.ml_rate_err_p1, erosita.ml_flux_p1, erosita.ml_flux_err_p1, erosita.ml_bkg_p1, erosita.ml_exp_p1, erosita.ape_bkg_p1, erosita.ape_radius_p1, erosita.ape_pois_p1,
    erosita.det_like_p2, erosita.ml_cts_p2, erosita.ml_cts_err_p2, erosita.ml_rate_p2, erosita.ml_rate_err_p2, erosita.ml_flux_p2, erosita.ml_flux_err_p2, erosita.ml_bkg_p2, erosita.ml_exp_p2, erosita.ape_bkg_p2, erosita.ape_radius_p2, erosita.ape_pois_p2,
    erosita.det_like_p3, erosita.ml_cts_p3, erosita.ml_cts_err_p3, erosita.ml_rate_p3, erosita.ml_rate_err_p3, erosita.ml_flux_p3, erosita.ml_flux_err_p3, erosita.ml_bkg_p3, erosita.ml_exp_p3, erosita.ape_bkg_p3, erosita.ape_radius_p3, erosita.ape_pois_p3,
    erosita.det_like_p4, erosita.ml_cts_p4, erosita.ml_cts_err_p4, erosita.ml_rate_p4, erosita.ml_rate_err_p4, erosita.ml_flux_p4, erosita.ml_flux_err_p4, erosita.ml_bkg_p4, erosita.ml_exp_p4, erosita.ape_bkg_p4, erosita.ape_radius_p4, erosita.ape_pois_p4,
    erosita.det_like_p5, erosita.ml_cts_p5, erosita.ml_cts_err_p5, erosita.ml_rate_p5, erosita.ml_rate_err_p5, erosita.ml_flux_p5, erosita.ml_flux_err_p5, erosita.ml_bkg_p5, erosita.ml_exp_p5, erosita.ape_bkg_p5, erosita.ape_radius_p5, erosita.ape_pois_p5,
    erosita.det_like_p6, erosita.ml_cts_p6, erosita.ml_cts_err_p6, erosita.ml_rate_p6, erosita.ml_rate_err_p6, erosita.ml_flux_p6, erosita.ml_flux_err_p6, erosita.ml_bkg_p6, erosita.ml_exp_p6, erosita.ape_bkg_p6, erosita.ape_radius_p6, erosita.ape_pois_p6,
    erosita.flag_sp_snr, erosita.flag_sp_bps, erosita.flag_sp_scl, erosita.flag_sp_lga, erosita.flag_sp_gc_cons, erosita.flag_no_radec_err, erosita.flag_no_ext_err, erosita.flag_no_cts_err, erosita.flag_opt,
    mastercat.ra, mastercat.dec
FROM erosita 
JOIN mastercat ON mastercat.id = erosita.id
WHERE erosita.id IN (sqlc.slice(id));
