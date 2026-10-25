-- name: GetCatalogs :many
SELECT *
FROM catalogs;

-- name: InsertCatalog :exec
INSERT INTO catalogs (
	name, nside
) VALUES (
	?, ?
);

-- name: RemoveAllCatalogs :exec
DELETE FROM catalogs;
