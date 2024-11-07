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
