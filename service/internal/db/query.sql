-- name: FindObjects :many
SELECT *
FROM mastercat 
WHERE ipix IN (sqlc.slice(ipix));

-- name: GetObjectsFromCatalog :many
SELECT * 
FROM mastercat 
WHERE ipix IN (sqlc.slice(ipix))
AND cat = ?;

-- name: InsertObject :one
INSERT INTO mastercat (
	id, ipix, ra, dec, cat
) VALUES (
	?, ?, ?, ?, ?
)
RETURNING *;

-- name: GetAllObjects :many
SELECT *
FROM mastercat;

-- name: GetCatalogs :many
SELECT *
FROM catalogs;

-- name: InsertCatalog :one
INSERT INTO catalogs (
	name, nside
) VALUES (
	?, ?
)
RETURNING *;

-- name: InsertAllwise :exec
INSERT INTO allwise (
	id, w1mpro, w1sigmpro, w2mpro, w2sigmpro, w3mpro, w3sigmpro, w4mpro, w4sigmpro, J_m_2mass, J_msig_2mass, H_m_2mass, H_msig_2mass, K_m_2mass, K_msig_2mass
) VALUES (
	?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
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

