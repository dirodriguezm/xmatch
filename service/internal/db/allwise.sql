-- name: InsertAllwise :exec
INSERT INTO allwise (
	id, cntr, w1mpro, w1sigmpro, w2mpro, w2sigmpro, w3mpro, w3sigmpro, w4mpro, w4sigmpro, J_m_2mass, J_msig_2mass, H_m_2mass, H_msig_2mass, K_m_2mass, K_msig_2mass
) VALUES (
	?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
);

-- name: GetAllwise :one
SELECT allwise.*, mastercat.ra, mastercat.dec
FROM allwise 
JOIN mastercat ON mastercat.id = allwise.id
WHERE allwise.id = ?;

-- name: BulkGetAllwise :many
SELECT allwise.*, mastercat.ra, mastercat.dec
FROM allwise 
JOIN mastercat ON mastercat.id = allwise.id
WHERE allwise.id IN (sqlc.slice(id));
