-- name: InsertSwift :exec
INSERT INTO swift (
    id
) VALUES (
    ?
);

-- name: GetSwift :one
SELECT swift.*, mastercat.ra, mastercat.dec
FROM swift
JOIN mastercat ON mastercat.id = swift.id
WHERE swift.id = ?;

-- name: BulkGetSwift :many
SELECT swift.*, mastercat.ra, mastercat.dec
FROM swift
JOIN mastercat ON mastercat.id = swift.id
WHERE swift.id IN (sqlc.slice(id));
