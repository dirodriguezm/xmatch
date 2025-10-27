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
