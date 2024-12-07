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
