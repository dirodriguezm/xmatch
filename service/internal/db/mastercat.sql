-- name: FindObjects :many
SELECT *
FROM mastercat 
WHERE ipix IN (sqlc.slice(ipix));

-- name: InsertObject :exec
INSERT INTO mastercat (
	id, ipix, ra, dec, cat
) VALUES (
	?, ?, ?, ?, ?
);

-- name: GetAllObjects :many
SELECT *
FROM mastercat;

-- name: RemoveAllObjects :exec
DELETE FROM mastercat;
