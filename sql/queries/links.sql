-- name: AddLink :one
INSERT INTO links (
    created_at, updated_at,url, item_id
)VALUES (
    current_timestamp,
    current_timestamp,
    ?,
    ?
)
RETURNING *;